package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CIECService interface {
	SaveCIEC(c *fiber.Ctx, req *validation.CreateCIEC, userID string) error
	UpdateCIEC(c *fiber.Ctx, req *validation.UpdateCIEC, userID string) error
	GetDecryptedCIEC(c *fiber.Ctx, userID string) (string, error)
	DeleteCIEC(c *fiber.Ctx, userID string) error
	HasCIEC(c *fiber.Ctx, userID string) (bool, error)
}

type ciecService struct {
	Log          *logrus.Logger
	DB           *gorm.DB
	Validate     *validator.Validate
	FirmaService FirmaElectronicaService
}

func NewCIECService(db *gorm.DB, validate *validator.Validate, firmaService FirmaElectronicaService) CIECService {
	return &ciecService{
		Log:          utils.Log,
		DB:           db,
		Validate:     validate,
		FirmaService: firmaService,
	}
}

func (s *ciecService) SaveCIEC(c *fiber.Ctx, req *validation.CreateCIEC, userID string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	// Encrypt CIEC using the same encryption method as FIEL
	encryptedCIEC, err := s.FirmaService.(*firmaElectronicaService).encryptData(req.CIEC)
	if err != nil {
		s.Log.Errorf("Failed to encrypt CIEC: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to encrypt CIEC")
	}

	// Update user with RFC and encrypted CIEC
	result := s.DB.WithContext(c.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"rfc":  req.RFC,
			"ciec": encryptedCIEC,
		})

	if result.Error != nil {
		s.Log.Errorf("Failed to save CIEC: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}

func (s *ciecService) UpdateCIEC(c *fiber.Ctx, req *validation.UpdateCIEC, userID string) error {
	if err := s.Validate.Struct(req); err != nil {
		return err
	}

	encryptedCIEC, err := s.FirmaService.(*firmaElectronicaService).encryptData(req.CIEC)
	if err != nil {
		s.Log.Errorf("Failed to encrypt CIEC: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to encrypt CIEC")
	}

	result := s.DB.WithContext(c.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		Update("ciec", encryptedCIEC)

	if result.Error != nil {
		s.Log.Errorf("Failed to update CIEC: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}

func (s *ciecService) GetDecryptedCIEC(c *fiber.Ctx, userID string) (string, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fiber.NewError(fiber.StatusNotFound, "User not found")
		}
		s.Log.Errorf("Failed to get user: %+v", err)
		return "", err
	}

	if user.CIEC == "" {
		return "", fiber.NewError(fiber.StatusNotFound, "CIEC not found for this user")
	}

	decryptedCIEC, err := s.FirmaService.DecryptCIEC(user.CIEC)
	if err != nil {
		s.Log.Errorf("Failed to decrypt CIEC: %+v", err)
		return "", fiber.NewError(fiber.StatusInternalServerError, "Failed to decrypt CIEC")
	}

	return decryptedCIEC, nil
}

func (s *ciecService) DeleteCIEC(c *fiber.Ctx, userID string) error {
	result := s.DB.WithContext(c.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"ciec": "",
		})

	if result.Error != nil {
		s.Log.Errorf("Failed to delete CIEC: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}

func (s *ciecService) HasCIEC(c *fiber.Ctx, userID string) (bool, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).Select("ciec").First(&user, "id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		s.Log.Errorf("Failed to check CIEC: %+v", err)
		return false, err
	}

	return user.CIEC != "", nil
}

