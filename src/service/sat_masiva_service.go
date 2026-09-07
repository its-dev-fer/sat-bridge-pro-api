package service

import (
	"app/src/config"
	"app/src/satws"
	"app/src/utils"
	"app/src/validation"
	"context"
	"sync"
	"time"

	nibussatws "github.com/InsaneTreset/nibus-sat-ws"
	"github.com/InsaneTreset/nibus-sat-ws/backfill"
	"github.com/InsaneTreset/nibus-sat-ws/satpackage"
	nservice "github.com/InsaneTreset/nibus-sat-ws/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SatMasivaService interface {
	FielStatus(c *fiber.Ctx, userID uuid.UUID) (map[string]any, error)
	Authenticate(c *fiber.Ctx, userID uuid.UUID, req *validation.SatServiceRequest) (map[string]any, error)
	Query(c *fiber.Ctx, userID uuid.UUID, req *validation.SatQueryRequest) (map[string]any, error)
	Verify(c *fiber.Ctx, userID uuid.UUID, req *validation.SatVerifyRequest) (map[string]any, error)
	Download(c *fiber.Ctx, userID uuid.UUID, req *validation.SatDownloadRequest) ([]byte, any, error)
	Backfill(c *fiber.Ctx, userID uuid.UUID, req *validation.SatBackfillRequest) (map[string]any, error)
	SyncStatus(userID uuid.UUID) (map[string]any, error)
	StartSync(c *fiber.Ctx, userID uuid.UUID) (map[string]any, error)
	AbortSync(userID uuid.UUID) (map[string]any, error)
}

type satMasivaService struct {
	Log                  *logrus.Logger
	DB                   *gorm.DB
	Validate             *validator.Validate
	DatosFiscalesService DatosFiscalesService
	cancels              map[uuid.UUID]context.CancelFunc
	cancelMu             sync.Mutex
}

func NewSatMasivaService(db *gorm.DB, validate *validator.Validate, datos DatosFiscalesService) SatMasivaService {
	s := &satMasivaService{
		Log:                  utils.Log,
		DB:                   db,
		Validate:             validate,
		DatosFiscalesService: datos,
		cancels:              map[uuid.UUID]context.CancelFunc{},
	}
	s.failStuckJobs()
	return s
}

func (s *satMasivaService) client(c *fiber.Ctx, userID uuid.UUID, kind string) (*nservice.Client, error) {
	cer, key, pass, err := s.DatosFiscalesService.LoadFIEL(c, userID)
	if err != nil {
		return nil, err
	}
	cli, err := satws.NewClient(cer, key, pass, kind, config.SATProxy)
	if err != nil {
		s.Log.Errorf("SAT client: %+v", err)
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid or expired e.firma")
	}
	return cli, nil
}

func (s *satMasivaService) FielStatus(c *fiber.Ctx, userID uuid.UUID) (map[string]any, error) {
	cli, err := s.client(c, userID, "cfdi")
	if err != nil {
		return nil, err
	}
	st := cli.Credential().Status()
	return map[string]any{
		"rfc":           st.RFC,
		"legal_name":    st.LegalName,
		"serial_number": st.SerialNumber,
		"is_fiel":       st.IsFIEL,
		"valid_from":    st.ValidFrom,
		"expires_at":    st.ExpiresAt,
		"expired":       st.Expired,
		"not_yet_valid": st.NotYetValid,
		"problems":      st.Problems,
	}, nil
}

func (s *satMasivaService) Authenticate(c *fiber.Ctx, userID uuid.UUID, req *validation.SatServiceRequest) (map[string]any, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	cli, err := s.client(c, userID, req.Service)
	if err != nil {
		return nil, err
	}
	tok, err := cli.Authenticate(c.Context())
	if err != nil {
		return nil, mapSATError(err)
	}
	return map[string]any{
		"expires_at": tok.Expires,
		"created_at": tok.Created,
	}, nil
}

func (s *satMasivaService) Query(c *fiber.Ctx, userID uuid.UUID, req *validation.SatQueryRequest) (map[string]any, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	cli, err := s.client(c, userID, req.Service)
	if err != nil {
		return nil, err
	}
	from, to, err := parseRange(req.From, req.To)
	if err != nil {
		return nil, err
	}
	dl, err := satws.ParseDownloadType(req.Download)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	rt, err := satws.ParseRequestType(req.Request)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	q := nservice.NewQuery(from, to)
	q.Download = dl
	q.Request = rt
	q.RfcEmisor = req.RfcEmisor
	q.RfcReceptores = req.RfcReceptores
	q.TipoComprobante = req.TipoComprobante
	q.EstadoComprobante = req.EstadoComprobante
	q.RfcACuentaTerceros = req.RfcACuentaTerceros
	q.Complemento = req.Complemento
	res, err := cli.Query(c.Context(), q)
	if err != nil {
		return nil, mapSATError(err)
	}
	return map[string]any{
		"request_id": res.RequestID,
		"code":       res.Code,
		"message":    res.Message,
	}, nil
}

func (s *satMasivaService) Verify(c *fiber.Ctx, userID uuid.UUID, req *validation.SatVerifyRequest) (map[string]any, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	cli, err := s.client(c, userID, req.Service)
	if err != nil {
		return nil, err
	}
	res, err := cli.Verify(c.Context(), req.RequestID)
	if err != nil {
		return nil, mapSATError(err)
	}
	return map[string]any{
		"state":        res.State.String(),
		"status_code":  res.StatusCode,
		"message":      res.Message,
		"number_cfdis": res.NumberCFDIs,
		"package_ids":  res.PackageIDs,
	}, nil
}

func (s *satMasivaService) Download(c *fiber.Ctx, userID uuid.UUID, req *validation.SatDownloadRequest) ([]byte, any, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, nil, err
	}
	cli, err := s.client(c, userID, req.Service)
	if err != nil {
		return nil, nil, err
	}
	zip, err := cli.Download(c.Context(), req.PackageID)
	if err != nil {
		return nil, nil, mapSATError(err)
	}
	if req.Parse == "" {
		return zip, nil, nil
	}
	parsed, err := parsePackage(zip, req.Parse)
	if err != nil {
		return nil, nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	return nil, parsed, nil
}

func (s *satMasivaService) Backfill(c *fiber.Ctx, userID uuid.UUID, req *validation.SatBackfillRequest) (map[string]any, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}
	cli, err := s.client(c, userID, req.Service)
	if err != nil {
		return nil, err
	}
	from, to, err := parseRange(req.From, req.To)
	if err != nil {
		return nil, err
	}
	dl, err := satws.ParseDownloadType(req.Download)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	rt, err := satws.ParseRequestType(req.Request)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	opts := backfill.Options{
		From:              from,
		To:                to,
		Download:          dl,
		Request:           rt,
		RfcEmisor:         req.RfcEmisor,
		TipoComprobante:   req.TipoComprobante,
		EstadoComprobante: req.EstadoComprobante,
		Complemento:       req.Complemento,
	}
	if req.ChunkDays > 0 {
		opts.Chunk = time.Duration(req.ChunkDays) * 24 * time.Hour
	}
	if req.PollSeconds > 0 {
		opts.PollEvery = time.Duration(req.PollSeconds) * time.Second
	}
	parseKind := "cfdi"
	if rt == nservice.RequestTypeMetadata {
		parseKind = "metadata"
	}
	var uuids []string
	// ponytail: holds the HTTP request for the whole SAT poll; queue/worker if this times out.
	stats, err := backfill.New(cli).Run(c.Context(), opts, func(zip []byte) error {
		parsed, err := parsePackage(zip, parseKind)
		if err != nil {
			return err
		}
		uuids = append(uuids, parsed.UUIDs...)
		return nil
	})
	if err != nil {
		return nil, mapSATError(err)
	}
	return map[string]any{
		"chunks":   stats.Chunks,
		"requests": stats.Requests,
		"packages": stats.Packages,
		"cfdis":    stats.CFDIs,
		"bytes":    stats.Bytes,
		"uuids":    uuids,
	}, nil
}

type parsedPackage struct {
	UUIDs    []string                    `json:"uuids"`
	Metadata []satpackage.MetadataRecord `json:"metadata,omitempty"`
	CFDIs    []map[string]string         `json:"cfdis,omitempty"`
}

func parsePackage(zip []byte, kind string) (*parsedPackage, error) {
	out := &parsedPackage{}
	if kind == "metadata" {
		recs, err := satpackage.ReadMetadata(zip)
		if err != nil {
			return nil, err
		}
		out.Metadata = recs
		for _, r := range recs {
			out.UUIDs = append(out.UUIDs, r.UUID)
		}
		return out, nil
	}
	recs, err := satpackage.ReadCFDIs(zip)
	if err != nil {
		return nil, err
	}
	for _, r := range recs {
		out.UUIDs = append(out.UUIDs, r.UUID)
		out.CFDIs = append(out.CFDIs, map[string]string{"uuid": r.UUID, "filename": r.Filename})
	}
	return out, nil
}

func parseRange(from, to string) (time.Time, time.Time, error) {
	f, err := time.Parse(time.RFC3339, from)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "from must be RFC3339")
	}
	t, err := time.Parse(time.RFC3339, to)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "to must be RFC3339")
	}
	if !f.Before(t) {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "from must be before to")
	}
	return f, t, nil
}

func mapSATError(err error) error {
	switch {
	case nibussatws.IsCredentialRejected(err):
		return fiber.NewError(fiber.StatusUnauthorized, err.Error())
	case nibussatws.IsNoInfo(err):
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	case nibussatws.IsDuplicate(err):
		return fiber.NewError(fiber.StatusConflict, err.Error())
	case nibussatws.IsLimitExceeded(err), nibussatws.IsExhausted(err):
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	default:
		return fiber.NewError(fiber.StatusBadGateway, err.Error())
	}
}
