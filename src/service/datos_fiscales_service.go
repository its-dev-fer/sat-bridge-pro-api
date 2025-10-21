package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DatosFiscalesService interface {
	CreateDatosFiscales(c *fiber.Ctx, userID uuid.UUID, req *validation.DatosFiscalesRequest, cerFile, keyFile *multipart.FileHeader) error
	GetDatosFiscalesByUserID(c *fiber.Ctx, userID uuid.UUID) (*model.DatosFiscalesSAT, error)
	UpdateDatosFiscales(c *fiber.Ctx, userID uuid.UUID, req *validation.DatosFiscalesRequest) error
	DeleteDatosFiscales(c *fiber.Ctx, userID uuid.UUID) error
}

type datosFiscalesService struct {
	Log           *logrus.Logger
	DB            *gorm.DB
	Validate      *validator.Validate
	EncryptionKey string
}

func NewDatosFiscalesService(db *gorm.DB, validate *validator.Validate, encryptionKey string) DatosFiscalesService {
	return &datosFiscalesService{
		Log:           utils.Log,
		DB:            db,
		Validate:      validate,
		EncryptionKey: encryptionKey,
	}
}

func (s *datosFiscalesService) CreateDatosFiscales(c *fiber.Ctx, userID uuid.UUID, req *validation.DatosFiscalesRequest, cerFile, keyFile *multipart.FileHeader) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	var existingData model.DatosFiscalesSAT
	result := s.DB.WithContext(c.Context()).Where("user_id = ?", userID).First(&existingData)
	if result.Error == nil {
		return fiber.NewError(fiber.StatusConflict, "Fiscal data already exists for this user")
	}

	if err := s.validateFileExtensions(cerFile, keyFile); err != nil {
		return err
	}

	cerB64, err := s.fileToBase64(cerFile)
	if err != nil {
		s.Log.Errorf("Error processing .cer file: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error processing certificate file")
	}

	keyB64, err := s.fileToBase64(keyFile)
	if err != nil {
		s.Log.Errorf("Error processing .key file: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "Error processing key file")
	}

	cerEncrypted, err := s.encrypt(cerB64)
	if err != nil {
		s.Log.Errorf("Error encrypting .cer file: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error processing certificate")
	}

	keyEncrypted, err := s.encrypt(keyB64)
	if err != nil {
		s.Log.Errorf("Error encrypting .key file: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error processing key")
	}

	passwordEncrypted, err := s.encrypt(req.Password)
	if err != nil {
		s.Log.Errorf("Error encrypting password: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error processing password")
	}

	datosFiscales := &model.DatosFiscalesSAT{
		UserID:               userID,
		RFC:                  req.RFC,
		CerB64Encriptado:     cerEncrypted,
		KeyB64Encriptado:     keyEncrypted,
		PasswordEfirmaEncrip: passwordEncrypted,
		CreatedBy:            &userID,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}

	result = s.DB.WithContext(c.Context()).Create(datosFiscales)
	if result.Error != nil {
		s.Log.Errorf("Failed to create fiscal data: %+v", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to save fiscal data")
	}

	return nil
}

func (s *datosFiscalesService) GetDatosFiscalesByUserID(c *fiber.Ctx, userID uuid.UUID) (*model.DatosFiscalesSAT, error) {
	datosFiscales := new(model.DatosFiscalesSAT)

	result := s.DB.WithContext(c.Context()).Where("user_id = ?", userID).First(datosFiscales)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Fiscal data not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to get fiscal data: %+v", result.Error)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Failed to retrieve fiscal data")
	}

	return datosFiscales, nil
}

func (s *datosFiscalesService) UpdateDatosFiscales(c *fiber.Ctx, userID uuid.UUID, req *validation.DatosFiscalesRequest) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	datosFiscales, err := s.GetDatosFiscalesByUserID(c, userID)
	if err != nil {
		return err
	}

	if req.Password != "" {
		passwordEncrypted, err := s.encrypt(req.Password)
		if err != nil {
			s.Log.Errorf("Error encrypting password: %+v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Error processing password")
		}
		datosFiscales.PasswordEfirmaEncrip = passwordEncrypted
	}

	if req.RFC != "" {
		datosFiscales.RFC = req.RFC
	}

	datosFiscales.UpdatedAt = time.Now()
	datosFiscales.UpdatedBy = &userID

	result := s.DB.WithContext(c.Context()).Save(datosFiscales)
	if result.Error != nil {
		s.Log.Errorf("Failed to update fiscal data: %+v", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update fiscal data")
	}

	return nil
}

func (s *datosFiscalesService) DeleteDatosFiscales(c *fiber.Ctx, userID uuid.UUID) error {
	result := s.DB.WithContext(c.Context()).Where("user_id = ?", userID).Delete(&model.DatosFiscalesSAT{})

	if result.Error != nil {
		s.Log.Errorf("Failed to delete fiscal data: %+v", result.Error)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete fiscal data")
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Fiscal data not found")
	}

	return nil
}

func (s *datosFiscalesService) validateFileExtensions(cerFile, keyFile *multipart.FileHeader) error {
	if cerFile == nil || keyFile == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Both .cer and .key files are required")
	}

	maxSize := int64(5 * 1024 * 1024)
	if cerFile.Size > maxSize || keyFile.Size > maxSize {
		return fiber.NewError(fiber.StatusBadRequest, "File size must be less than 5MB")
	}

	if len(cerFile.Filename) < 4 || cerFile.Filename[len(cerFile.Filename)-4:] != ".cer" {
		return fiber.NewError(fiber.StatusBadRequest, "Certificate file must have .cer extension")
	}

	if len(keyFile.Filename) < 4 || keyFile.Filename[len(keyFile.Filename)-4:] != ".key" {
		return fiber.NewError(fiber.StatusBadRequest, "Key file must have .key extension")
	}

	return nil
}

func (s *datosFiscalesService) fileToBase64(fileHeader *multipart.FileHeader) (string, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(fileBytes), nil
}

func (s *datosFiscalesService) encrypt(plaintext string) (string, error) {
	hash := sha256.Sum256([]byte(s.EncryptionKey))
	key := hash[:]

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Crear GCM
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