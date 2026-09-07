package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"app/src/config"
	"app/src/model"
	"app/src/satws"

	nibussatws "github.com/InsaneTreset/nibus-sat-ws"
	"github.com/InsaneTreset/nibus-sat-ws/backfill"
	"github.com/InsaneTreset/nibus-sat-ws/satpackage"
	nservice "github.com/InsaneTreset/nibus-sat-ws/service"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var meses = [...]string{
	"", "enero", "febrero", "marzo", "abril", "mayo", "junio",
	"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre",
}

func canStartSync(cfdiCount int64, job *model.SatMasivaJob) bool {
	if job != nil && (job.Status == "queued" || job.Status == "running") {
		return false
	}
	if job != nil && job.Status == "done" {
		return false
	}
	if cfdiCount > 0 && (job == nil || (job.Status != "failed" && job.Status != "aborted")) {
		return false
	}
	return true
}

func (s *satMasivaService) latestJob(userID uuid.UUID) (*model.SatMasivaJob, error) {
	var job model.SatMasivaJob
	err := s.DB.Where("user_id = ?", userID).Order("created_at DESC").First(&job).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &job, err
}

func (s *satMasivaService) cfdiCount(userID uuid.UUID) (int64, error) {
	var n int64
	err := s.DB.Model(&model.CFDI{}).Where("user_id = ?", userID).Count(&n).Error
	return n, err
}

func (s *satMasivaService) SyncStatus(userID uuid.UUID) (map[string]any, error) {
	n, err := s.cfdiCount(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to count CFDIs")
	}
	job, err := s.latestJob(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to load sync job")
	}
	return map[string]any{
		"cfdi_count": n,
		"can_sync":   canStartSync(n, job),
		"job":        job,
	}, nil
}

func (s *satMasivaService) StartSync(c *fiber.Ctx, userID uuid.UUID) (map[string]any, error) {
	n, err := s.cfdiCount(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to count CFDIs")
	}
	prev, err := s.latestJob(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to load sync job")
	}
	if !canStartSync(n, prev) {
		return nil, fiber.NewError(fiber.StatusConflict, "Ya hay CFDIs o un job en curso")
	}
	cer, key, pass, err := s.DatosFiscalesService.LoadFIEL(c, userID)
	if err != nil {
		return nil, err
	}
	cer = append([]byte(nil), cer...)
	key = append([]byte(nil), key...)
	pass = strings.TrimSpace(pass)
	if _, err := satws.NewClient(cer, key, pass, "cfdi", config.SATProxy); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No se pudo abrir la e.firma (contraseña o archivo .key). Vuelve a cargarla en Mi FIEL.")
	}
	year := time.Now().Year()
	job := &model.SatMasivaJob{
		UUID:    uuid.New(),
		UserID:  userID,
		Year:    year,
		Status:  "queued",
		Message: "En cola",
		Logs:    []string{time.Now().Format("15:04:05") + " Job en cola"},
	}
	if err := s.DB.Create(job).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "Ya hay un job en curso")
	}
	ctx, cancel := context.WithCancel(context.Background())
	s.setCancel(job.UUID, cancel)
	go s.runYearSync(ctx, job.UUID, userID, year, cer, key, pass)
	return s.SyncStatus(userID)
}

func (s *satMasivaService) setCancel(id uuid.UUID, cancel context.CancelFunc) {
	s.cancelMu.Lock()
	s.cancels[id] = cancel
	s.cancelMu.Unlock()
}

func (s *satMasivaService) popCancel(id uuid.UUID) {
	s.cancelMu.Lock()
	if c, ok := s.cancels[id]; ok {
		c()
		delete(s.cancels, id)
	}
	s.cancelMu.Unlock()
}

func (s *satMasivaService) AbortSync(userID uuid.UUID) (map[string]any, error) {
	job, err := s.latestJob(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to load sync job")
	}
	if job == nil || (job.Status != "queued" && job.Status != "running") {
		return nil, fiber.NewError(fiber.StatusConflict, "No hay un job en curso")
	}
	s.popCancel(job.UUID)
	now := time.Now()
	job.Status = "aborted"
	job.Error = "Abortado por el usuario"
	job.FinishedAt = &now
	s.logJob(job, job.Month, "Abortado por el usuario")
	s.DB.Model(job).Select("Status", "Error", "FinishedAt", "UpdatedAt").Updates(job)
	return s.SyncStatus(userID)
}

func (s *satMasivaService) failStuckJobs() {
	msg := "El servidor se reinició durante la extracción"
	s.DB.Model(&model.SatMasivaJob{}).
		Where("status IN ?", []string{"queued", "running"}).
		Updates(map[string]any{"status": "failed", "error": msg, "finished_at": time.Now(), "updated_at": time.Now()})
}

func (s *satMasivaService) logJob(job *model.SatMasivaJob, month, msg string) {
	job.Month = month
	job.Message = msg
	job.Logs = append(job.Logs, time.Now().Format("15:04:05")+" "+msg)
	job.UpdatedAt = time.Now()
	s.DB.Model(job).Select("Month", "Message", "Logs", "UpdatedAt").Updates(job)
}

func (s *satMasivaService) runYearSync(ctx context.Context, jobID, userID uuid.UUID, year int, cer, key []byte, pass string) {
	// ponytail: in-process goroutine dies with the API; queue if you run >1 instance.
	defer s.popCancel(jobID)
	defer func() {
		if rec := recover(); rec != nil {
			var job model.SatMasivaJob
			if s.DB.Where("uuid = ?", jobID).First(&job).Error == nil && job.Status != "aborted" {
				s.failJob(&job, fmt.Errorf("%v", rec))
			}
		}
	}()

	var job model.SatMasivaJob
	if err := s.DB.Where("uuid = ?", jobID).First(&job).Error; err != nil {
		return
	}
	if ctx.Err() != nil {
		return
	}
	job.Status = "running"
	s.DB.Model(&job).Update("status", "running")
	s.logJob(&job, "", fmt.Sprintf("Extrayendo CFDIs %d", year))

	cli, err := satws.NewClient(cer, key, pass, "cfdi", config.SATProxy)
	if err != nil {
		s.failJob(&job, err)
		return
	}

	now := time.Now()
	loc := now.Location()
	for m := 1; m <= int(now.Month()); m++ {
		if ctx.Err() != nil {
			return
		}
		if s.DB.Where("uuid = ?", jobID).First(&job).Error != nil {
			return
		}
		if job.Status == "aborted" {
			return
		}
		monthName := fmt.Sprintf("%s %d", meses[m], year)
		from := time.Date(year, time.Month(m), 1, 0, 0, 0, 0, loc)
		to := from.AddDate(0, 1, 0)
		if to.After(now) {
			to = now
		}
		if !from.Before(to) {
			continue
		}
		s.logJob(&job, monthName, "Extrayendo facturas del mes "+monthName+" ....")
		for _, dl := range []nservice.DownloadType{nservice.Issued, nservice.Received} {
			if ctx.Err() != nil {
				return
			}
			kind := "emitidas"
			if dl == nservice.Received {
				kind = "recibidas"
			}
			s.logJob(&job, monthName, "SAT "+kind+" "+monthName)
			err := s.backfillMonth(ctx, cli, &job, userID, from, to, dl)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				s.failJob(&job, err)
				return
			}
		}
	}
	if ctx.Err() != nil {
		return
	}
	now = time.Now()
	job.Status = "done"
	job.FinishedAt = &now
	s.logJob(&job, "", "Extracción completada")
	s.DB.Model(&job).Select("Status", "FinishedAt", "UpdatedAt").Updates(&job)
}

func (s *satMasivaService) failJob(job *model.SatMasivaJob, err error) {
	if job.Status == "aborted" {
		return
	}
	job.Status = "failed"
	job.Error = err.Error()
	now := time.Now()
	job.FinishedAt = &now
	s.logJob(job, job.Month, "ERROR: "+err.Error())
	s.DB.Model(job).Select("Status", "Error", "FinishedAt", "UpdatedAt").Updates(job)
}

func (s *satMasivaService) backfillMonth(ctx context.Context, cli *nservice.Client, job *model.SatMasivaJob, userID uuid.UUID, from, to time.Time, dl nservice.DownloadType) error {
	opts := backfill.Options{
		From:     from,
		To:       to,
		Download: dl,
		Request:  nservice.RequestTypeMetadata,
		Chunk:    32 * 24 * time.Hour,
		OnProgress: func(msg string) {
			if ctx.Err() != nil {
				return
			}
			s.logJob(job, job.Month, msg)
		},
	}
	_, err := backfill.New(cli).Run(ctx, opts, func(zip []byte) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		recs, err := satpackage.ReadMetadata(zip)
		if err != nil {
			return err
		}
		return s.insertMetadata(userID, recs)
	})
	if err != nil && !nibussatws.IsNoInfo(err) {
		return err
	}
	return nil
}

func (s *satMasivaService) insertMetadata(userID uuid.UUID, recs []satpackage.MetadataRecord) error {
	if len(recs) == 0 {
		return nil
	}
	rows := make([]model.CFDI, 0, len(recs))
	for _, r := range recs {
		if r.UUID == "" {
			continue
		}
		total := parseMonto(r.Monto)
		fe := parseSATTime(r.FechaEmision)
		fc := parseSATTime(r.FechaCertificacionSat)
		rows = append(rows, model.CFDI{
			UUID:               uuid.New(),
			UserID:             userID,
			FolioFiscal:        r.UUID,
			RFCEmisor:          nz(r.RfcEmisor),
			NombreEmisor:       nz(r.NombreEmisor),
			RFCReceptor:        nz(r.RfcReceptor),
			NombreReceptor:     nz(r.NombreReceptor),
			FechaEmision:       fe,
			FechaCertificacion: fc,
			PACCertifico:       r.PacCertifico,
			Total:              total,
			EfectoComprobante:  r.EfectoComprobante,
			EstatusCancelacion: r.FechaCancelacion,
			EstadoComprobante:  r.Estatus,
		})
	}
	if len(rows) == 0 {
		return nil
	}
	return s.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "folio_fiscal"}},
		DoNothing: true,
	}).Create(&rows).Error
}

func nz(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func parseMonto(s string) *float64 {
	s = strings.TrimSpace(strings.ReplaceAll(s, ",", ""))
	if s == "" {
		return nil
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil
	}
	return &f
}

func parseSATTime(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}
