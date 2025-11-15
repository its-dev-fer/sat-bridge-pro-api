package service

import (
    "app/src/model"
    "app/src/utils"
    "app/src/validation"
    "errors"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/google/uuid"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
)

type CfdiDescargadoService interface {
    GetAll(c *fiber.Ctx) ([]model.CfdiDescargado, error)
    GetByID(c *fiber.Ctx, id string) (*model.CfdiDescargado, error)
    Create(c *fiber.Ctx, req *validation.CreateCfdiDescargado) (*model.CfdiDescargado, error)
    Update(c *fiber.Ctx, req *validation.UpdateCfdiDescargado, id string) (*model.CfdiDescargado, error)
    Delete(c *fiber.Ctx, id string) error
}

type cfdiDescargadoService struct {
    Log      *logrus.Logger
    DB       *gorm.DB
    Validate *validator.Validate
}

func NewCfdiDescargadoService(db *gorm.DB, validate *validator.Validate) CfdiDescargadoService {
    return &cfdiDescargadoService{
        Log:      utils.Log,
        DB:       db,
        Validate: validate,
    }
}

func (s *cfdiDescargadoService) GetAll(c *fiber.Ctx) ([]model.CfdiDescargado, error) {
    var cfdis []model.CfdiDescargado

    result := s.DB.WithContext(c.Context()).Order("created_at desc").Find(&cfdis)
    if result.Error != nil {
        s.Log.Errorf("Error al obtener CFDIs: %+v", result.Error)
        return nil, result.Error
    }

    return cfdis, nil
}

func (s *cfdiDescargadoService) GetByID(c *fiber.Ctx, id string) (*model.CfdiDescargado, error) {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    cfdi := new(model.CfdiDescargado)
    result := s.DB.WithContext(c.Context()).First(cfdi, "id = ?", uuidID)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, fiber.NewError(fiber.StatusNotFound, "CFDI no encontrado")
    }
    if result.Error != nil {
        s.Log.Errorf("Error al obtener CFDI: %+v", result.Error)
        return nil, result.Error
    }

    return cfdi, nil
}

func (s *cfdiDescargadoService) Create(c *fiber.Ctx, req *validation.CreateCfdiDescargado) (*model.CfdiDescargado, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateCreateCfdiDescargado(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    cfdi := req.ToModel()

    result := s.DB.WithContext(c.Context()).Create(cfdi)
    if result.Error != nil {
        s.Log.Errorf("Error al crear CFDI: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo crear el CFDI")
    }

    return cfdi, nil
}

func (s *cfdiDescargadoService) Update(c *fiber.Ctx, req *validation.UpdateCfdiDescargado, id string) (*model.CfdiDescargado, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateUpdateCfdiDescargado(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    cfdi := req.ToModel()
    cfdi.ID = uuidID.String()

    result := s.DB.WithContext(c.Context()).Model(&model.CfdiDescargado{}).Where("id = ?", uuidID).Updates(cfdi)
    if result.Error != nil {
        s.Log.Errorf("Error al actualizar CFDI: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo actualizar el CFDI")
    }

    return cfdi, nil
}

func (s *cfdiDescargadoService) Delete(c *fiber.Ctx, id string) error {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    result := s.DB.WithContext(c.Context()).Delete(&model.CfdiDescargado{}, "id = ?", uuidID)
    if result.Error != nil {
        s.Log.Errorf("Error al eliminar CFDI: %+v", result.Error)
        return fiber.NewError(fiber.StatusInternalServerError, "No se pudo eliminar el CFDI")
    }

    if result.RowsAffected == 0 {
        return fiber.NewError(fiber.StatusNotFound, "CFDI no encontrado")
    }

    return nil
}
