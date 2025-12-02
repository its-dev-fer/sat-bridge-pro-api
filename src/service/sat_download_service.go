package service

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SatDownloadService interface {
	DownloadCFDIs(c *fiber.Ctx, req *DownloadCFDIRequest, userID string) (*DownloadCFDIResponse, error)
	QueryMetadata(c *fiber.Ctx, req *DownloadCFDIRequest, userID string) (*MetadataQueryResponse, error)
	DownloadByUUID(c *fiber.Ctx, req *DownloadByUUIDRequest, userID string) (*SingleCFDIResponse, error)
}

type satDownloadService struct {
	Log          *logrus.Logger
	DB           *gorm.DB
	LimitService DownloadLimitService
	FirmaService FirmaElectronicaService
	CIECService  CIECService
	PHPApiURL    string
	PHPApiKey    string
}

type DownloadCFDIRequest struct {
	AuthType      string `json:"auth_type" validate:"required,oneof=ciec fiel"`
	TipoCFDI      string `json:"tipo_cfdi" validate:"required,oneof=emitidos recibidos todos"`
	FechaInicio   string `json:"fecha_inicio" validate:"required"`
	FechaFin      string `json:"fecha_fin" validate:"required"`
	RFCEmisor     string `json:"rfc_emisor,omitempty"`
	RFCReceptor   string `json:"rfc_receptor,omitempty"`
	Estado        string `json:"estado,omitempty"`
	Complemento   string `json:"complemento,omitempty"`
	SaveToDatabase bool  `json:"save_to_database"`
}

type DownloadByUUIDRequest struct {
	AuthType    string `json:"auth_type" validate:"required,oneof=ciec fiel"`
	UUID        string `json:"uuid" validate:"required,uuid4"`
	TipoCFDI    string `json:"tipo_cfdi" validate:"required"`
	FechaInicio string `json:"fecha_inicio" validate:"required"`
	FechaFin    string `json:"fecha_fin" validate:"required"`
}

type DownloadCFDIResponse struct {
	Success         bool                  `json:"success"`
	TotalFound      int                   `json:"total_found"`
	TotalDownloaded int                   `json:"total_downloaded"`
	CFDIs           []CFDIData            `json:"cfdis"`
	Errors          []map[string]string   `json:"errors"`
	SolicitudID     string                `json:"solicitud_id"`
}

type MetadataQueryResponse struct {
	Success    bool       `json:"success"`
	TotalFound int        `json:"total_found"`
	Metadata   []CFDIData `json:"metadata"`
}

type SingleCFDIResponse struct {
	Success bool     `json:"success"`
	CFDI    CFDIData `json:"cfdi"`
}

type CFDIData struct {
	UUID         string  `json:"uuid"`
	RFCEmisor    string  `json:"rfc_emisor"`
	RFCReceptor  string  `json:"rfc_receptor"`
	FechaEmision string  `json:"fecha_emision"`
	MontoTotal   float64 `json:"monto_total"`
	StatusSAT    string  `json:"status_sat"`
	TipoCFDI     string  `json:"tipo_cfdi"`
	ArchivoXML   string  `json:"archivo_xml,omitempty"`
	XMLContent   string  `json:"xml_content,omitempty"`
}

func NewSatDownloadService(db *gorm.DB, limitService DownloadLimitService, 
	firmaService FirmaElectronicaService, ciecService CIECService) SatDownloadService {
	return &satDownloadService{
		Log:          utils.Log,
		DB:           db,
		LimitService: limitService,
		FirmaService: firmaService,
		CIECService:  ciecService,
		PHPApiURL:    config.PHPSatScraperURL,
		PHPApiKey:    config.PHPSatScraperAPIKey,
	}
}

func (s *satDownloadService) DownloadCFDIs(c *fiber.Ctx, req *DownloadCFDIRequest, userID string) (*DownloadCFDIResponse, error) {
	// Check download limits
	canDownload, message, err := s.LimitService.CanDownload(c, userID)
	if err != nil {
		s.Log.Errorf("Failed to check download limits: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to check download limits")
	}

	if !canDownload {
		return nil, fiber.NewError(fiber.StatusForbidden, message)
	}

	// Get authentication credentials
	phpRequest, err := s.buildPHPRequest(c, req, userID)
	if err != nil {
		return nil, err
	}

	// Call PHP microservice
	phpResponse, err := s.callPHPService("/api/download", phpRequest)
	if err != nil {
		s.Log.Errorf("Failed to call PHP service: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to download CFDIs from SAT")
	}

	var result DownloadCFDIResponse
	if err := json.Unmarshal(phpResponse, &result); err != nil {
		s.Log.Errorf("Failed to parse PHP response: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse download response")
	}

	if !result.Success {
		return &result, nil
	}

	// Create solicitud descarga
	solicitudID := uuid.New().String()
	solicitud := &model.SolicitudDescarga{
		ID:             solicitudID,
		UsuarioID:      userID,
		TipoCFDI:       req.TipoCFDI,
		RFCSolicitante: phpRequest["rfc"].(string),
		FechaSolicitud: time.Now(),
		Status:         "activo",
	}

	if err := s.DB.WithContext(c.Context()).Create(solicitud).Error; err != nil {
		s.Log.Errorf("Failed to create solicitud descarga: %+v", err)
	} else {
		result.SolicitudID = solicitudID
	}

	// Save CFDIs to database if requested
	if req.SaveToDatabase && result.Success {
		if err := s.saveCFDIsToDatabase(c, result.CFDIs, solicitudID); err != nil {
			s.Log.Errorf("Failed to save CFDIs to database: %+v", err)
		}
	}

	// Increment download counter
	if err := s.LimitService.IncrementDownloadCount(c, userID); err != nil {
		s.Log.Errorf("Failed to increment download count: %+v", err)
	}

	return &result, nil
}

func (s *satDownloadService) QueryMetadata(c *fiber.Ctx, req *DownloadCFDIRequest, userID string) (*MetadataQueryResponse, error) {
	// Get authentication credentials
	phpRequest, err := s.buildPHPRequest(c, req, userID)
	if err != nil {
		return nil, err
	}

	// Call PHP microservice
	phpResponse, err := s.callPHPService("/api/query-metadata", phpRequest)
	if err != nil {
		s.Log.Errorf("Failed to call PHP service: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to query metadata from SAT")
	}

	var result MetadataQueryResponse
	if err := json.Unmarshal(phpResponse, &result); err != nil {
		s.Log.Errorf("Failed to parse PHP response: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse metadata response")
	}

	return &result, nil
}

func (s *satDownloadService) DownloadByUUID(c *fiber.Ctx, req *DownloadByUUIDRequest, userID string) (*SingleCFDIResponse, error) {
	// Check download limits
	canDownload, message, err := s.LimitService.CanDownload(c, userID)
	if err != nil {
		return nil, err
	}

	if !canDownload {
		return nil, fiber.NewError(fiber.StatusForbidden, message)
	}

	// Build PHP request
	phpRequest := make(map[string]interface{})
	phpRequest["uuid"] = req.UUID
	phpRequest["tipo_cfdi"] = req.TipoCFDI
	phpRequest["fecha_inicio"] = req.FechaInicio
	phpRequest["fecha_fin"] = req.FechaFin

	// Add authentication
	if req.AuthType == "ciec" {
		var user model.User
		if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		if user.RFC == "" || user.CIEC == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User has no CIEC configured")
		}

		decryptedCIEC, err := s.CIECService.GetDecryptedCIEC(c, userID)
		if err != nil {
			return nil, err
		}

		phpRequest["auth_type"] = "ciec"
		phpRequest["rfc"] = user.RFC
		phpRequest["ciec"] = decryptedCIEC
	} else {
		firma, err := s.FirmaService.GetFirmaByUserID(c, userID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User has no FIEL configured")
		}

		decryptedPassword, err := s.FirmaService.DecryptPassword(firma.PasswordKEY)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to decrypt FIEL password")
		}

		phpRequest["auth_type"] = "fiel"
		phpRequest["rfc"] = firma.RFC
		phpRequest["certificado_cer"] = firma.CertificadoCER
		phpRequest["clave_privada_key"] = firma.ClavePrivadaKEY
		phpRequest["password_key"] = decryptedPassword
	}

	// Call PHP microservice
	phpResponse, err := s.callPHPService("/api/download-by-uuid", phpRequest)
	if err != nil {
		s.Log.Errorf("Failed to call PHP service: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to download CFDI from SAT")
	}

	var result SingleCFDIResponse
	if err := json.Unmarshal(phpResponse, &result); err != nil {
		s.Log.Errorf("Failed to parse PHP response: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to parse download response")
	}

	// Increment download counter
	if result.Success {
		if err := s.LimitService.IncrementDownloadCount(c, userID); err != nil {
			s.Log.Errorf("Failed to increment download count: %+v", err)
		}
	}

	return &result, nil
}

func (s *satDownloadService) buildPHPRequest(c *fiber.Ctx, req *DownloadCFDIRequest, userID string) (map[string]interface{}, error) {
	phpRequest := make(map[string]interface{})
	
	phpRequest["tipo_cfdi"] = req.TipoCFDI
	phpRequest["fecha_inicio"] = req.FechaInicio
	phpRequest["fecha_fin"] = req.FechaFin

	if req.RFCEmisor != "" {
		phpRequest["rfc_emisor"] = req.RFCEmisor
	}
	if req.RFCReceptor != "" {
		phpRequest["rfc_receptor"] = req.RFCReceptor
	}
	if req.Estado != "" {
		phpRequest["estado"] = req.Estado
	}
	if req.Complemento != "" {
		phpRequest["complemento"] = req.Complemento
	}

	// Add authentication based on type
	if req.AuthType == "ciec" {
		var user model.User
		if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
			return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
		}

		if user.RFC == "" || user.CIEC == "" {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User has no CIEC configured")
		}

		decryptedCIEC, err := s.CIECService.GetDecryptedCIEC(c, userID)
		if err != nil {
			return nil, err
		}

		phpRequest["auth_type"] = "ciec"
		phpRequest["rfc"] = user.RFC
		phpRequest["ciec"] = decryptedCIEC
	} else {
		firma, err := s.FirmaService.GetFirmaByUserID(c, userID)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusBadRequest, "User has no FIEL configured")
		}

		if !firma.IsValid() {
			return nil, fiber.NewError(fiber.StatusBadRequest, "FIEL is not valid or has expired")
		}

		decryptedPassword, err := s.FirmaService.DecryptPassword(firma.PasswordKEY)
		if err != nil {
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to decrypt FIEL password")
		}

		phpRequest["auth_type"] = "fiel"
		phpRequest["rfc"] = firma.RFC
		phpRequest["certificado_cer"] = firma.CertificadoCER
		phpRequest["clave_privada_key"] = firma.ClavePrivadaKEY
		phpRequest["password_key"] = decryptedPassword
	}

	return phpRequest, nil
}

func (s *satDownloadService) callPHPService(endpoint string, payload map[string]interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := s.PHPApiURL + endpoint
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", s.PHPApiKey)

	client := &http.Client{
		Timeout: 5 * time.Minute, // Long timeout for SAT downloads
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		s.Log.Errorf("PHP service returned error: %s", string(body))
		return nil, errors.New("PHP service returned error: " + string(body))
	}

	return body, nil
}

func (s *satDownloadService) saveCFDIsToDatabase(c *fiber.Ctx, cfdis []CFDIData, solicitudID string) error {
	for _, cfdi := range cfdis {
		fechaEmision, _ := time.Parse("2006-01-02", cfdi.FechaEmision)
		
		cfdiModel := &model.CfdiDescargado{
			ID:           uuid.New().String(),
			SolicitudID:  solicitudID,
			TipoCFDI:     cfdi.TipoCFDI,
			CfdiUUID:     cfdi.UUID,
			RFCEmisor:    cfdi.RFCEmisor,
			RFCReceptor:  cfdi.RFCReceptor,
			StatusSAT:    cfdi.StatusSAT,
			FechaEmision: fechaEmision,
			MontoTotal:   cfdi.MontoTotal,
			ArchivoXML:   cfdi.XMLContent, // Store base64 encoded XML
		}

		if err := s.DB.WithContext(c.Context()).Create(cfdiModel).Error; err != nil {
			s.Log.Errorf("Failed to save CFDI %s: %+v", cfdi.UUID, err)
			continue
		}
	}

	return nil
}

