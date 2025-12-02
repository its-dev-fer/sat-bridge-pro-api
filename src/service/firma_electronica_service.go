package service

import (
	"app/src/config"
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FirmaElectronicaService interface {
	CreateFirmaElectronica(c *fiber.Ctx, req *validation.CreateFirmaElectronica, userID string) (*model.FirmaElectronica, error)
	GetFirmaByUserID(c *fiber.Ctx, userID string) (*model.FirmaElectronica, error)
	GetFirmaByID(c *fiber.Ctx, id string) (*model.FirmaElectronica, error)
	UpdateFirmaElectronica(c *fiber.Ctx, req *validation.UpdateFirmaElectronica, id string, userID string) (*model.FirmaElectronica, error)
	DeleteFirmaElectronica(c *fiber.Ctx, id string, userID string) error
	DecryptPassword(encryptedPassword string) (string, error)
	DecryptCIEC(encryptedCIEC string) (string, error)
	GetAllFirmasByUser(c *fiber.Ctx, userID string) ([]model.FirmaElectronica, error)
}

type firmaElectronicaService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewFirmaElectronicaService(db *gorm.DB, validate *validator.Validate) FirmaElectronicaService {
	return &firmaElectronicaService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

// EncryptData encrypts data using AES-256-GCM
func (s *firmaElectronicaService) encryptData(plaintext string) (string, error) {
	// Use a secret key from config (should be 32 bytes for AES-256)
	key := []byte(config.JWTSecret) // In production, use a separate encryption key
	if len(key) < 32 {
		// Pad the key to 32 bytes if necessary
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	} else if len(key) > 32 {
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptData decrypts data using AES-256-GCM
func (s *firmaElectronicaService) decryptData(encryptedText string) (string, error) {
	key := []byte(config.JWTSecret)
	if len(key) < 32 {
		paddedKey := make([]byte, 32)
		copy(paddedKey, key)
		key = paddedKey
	} else if len(key) > 32 {
		key = key[:32]
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func (s *firmaElectronicaService) DecryptPassword(encryptedPassword string) (string, error) {
	return s.decryptData(encryptedPassword)
}

func (s *firmaElectronicaService) DecryptCIEC(encryptedCIEC string) (string, error) {
	return s.decryptData(encryptedCIEC)
}

func (s *firmaElectronicaService) CreateFirmaElectronica(c *fiber.Ctx, req *validation.CreateFirmaElectronica, userID string) (*model.FirmaElectronica, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Parse user UUID
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid user ID")
	}

	// Encrypt the password
	encryptedPassword, err := s.encryptData(req.PasswordKEY)
	if err != nil {
		s.Log.Errorf("Failed to encrypt password: %+v", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to encrypt password")
	}

	firma := &model.FirmaElectronica{
		UsuarioID:           userUUID,
		RFC:                 req.RFC,
		CertificadoCER:      req.CertificadoCER,
		ClavePrivadaKEY:     req.ClavePrivadaKEY,
		PasswordKEY:         encryptedPassword,
		NombreCertificado:   req.NombreCertificado,
		FechaVigenciaInicio: req.FechaVigenciaInicio,
		FechaVigenciaFin:    req.FechaVigenciaFin,
		Activo:              true,
	}

	result := s.DB.WithContext(c.Context()).Create(firma)
	if result.Error != nil {
		s.Log.Errorf("Failed to create firma electronica: %+v", result.Error)
		return nil, result.Error
	}

	return firma, nil
}

func (s *firmaElectronicaService) GetFirmaByUserID(c *fiber.Ctx, userID string) (*model.FirmaElectronica, error) {
	firma := new(model.FirmaElectronica)

	result := s.DB.WithContext(c.Context()).
		Where("usuario_id = ? AND activo = ?", userID, true).
		Order("created_at DESC").
		First(firma)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Firma electronica not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to get firma electronica by user id: %+v", result.Error)
		return nil, result.Error
	}

	return firma, nil
}

func (s *firmaElectronicaService) GetFirmaByID(c *fiber.Ctx, id string) (*model.FirmaElectronica, error) {
	firma := new(model.FirmaElectronica)

	result := s.DB.WithContext(c.Context()).First(firma, "id = ?", id)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Firma electronica not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to get firma electronica by id: %+v", result.Error)
		return nil, result.Error
	}

	return firma, nil
}

func (s *firmaElectronicaService) GetAllFirmasByUser(c *fiber.Ctx, userID string) ([]model.FirmaElectronica, error) {
	var firmas []model.FirmaElectronica

	result := s.DB.WithContext(c.Context()).
		Where("usuario_id = ?", userID).
		Order("created_at DESC").
		Find(&firmas)

	if result.Error != nil {
		s.Log.Errorf("Failed to get all firmas by user: %+v", result.Error)
		return nil, result.Error
	}

	return firmas, nil
}

func (s *firmaElectronicaService) UpdateFirmaElectronica(c *fiber.Ctx, req *validation.UpdateFirmaElectronica, id string, userID string) (*model.FirmaElectronica, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Get existing firma
	firma, err := s.GetFirmaByID(c, id)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	if firma.UsuarioID.String() != userID {
		return nil, fiber.NewError(fiber.StatusForbidden, "You don't have permission to update this firma")
	}

	// Update fields
	if req.CertificadoCER != "" {
		firma.CertificadoCER = req.CertificadoCER
	}
	if req.ClavePrivadaKEY != "" {
		firma.ClavePrivadaKEY = req.ClavePrivadaKEY
	}
	if req.PasswordKEY != "" {
		encryptedPassword, err := s.encryptData(req.PasswordKEY)
		if err != nil {
			s.Log.Errorf("Failed to encrypt password: %+v", err)
			return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to encrypt password")
		}
		firma.PasswordKEY = encryptedPassword
	}
	if req.NombreCertificado != "" {
		firma.NombreCertificado = req.NombreCertificado
	}
	if !req.FechaVigenciaInicio.IsZero() {
		firma.FechaVigenciaInicio = req.FechaVigenciaInicio
	}
	if !req.FechaVigenciaFin.IsZero() {
		firma.FechaVigenciaFin = req.FechaVigenciaFin
	}
	if req.Activo != nil {
		firma.Activo = *req.Activo
	}

	result := s.DB.WithContext(c.Context()).Save(firma)
	if result.Error != nil {
		s.Log.Errorf("Failed to update firma electronica: %+v", result.Error)
		return nil, result.Error
	}

	return firma, nil
}

func (s *firmaElectronicaService) DeleteFirmaElectronica(c *fiber.Ctx, id string, userID string) error {
	firma, err := s.GetFirmaByID(c, id)
	if err != nil {
		return err
	}

	// Verify ownership
	if firma.UsuarioID.String() != userID {
		return fiber.NewError(fiber.StatusForbidden, "You don't have permission to delete this firma")
	}

	result := s.DB.WithContext(c.Context()).Delete(&model.FirmaElectronica{}, "id = ?", id)
	if result.Error != nil {
		s.Log.Errorf("Failed to delete firma electronica: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Firma electronica not found")
	}

	return nil
}

