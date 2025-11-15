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

type SolicitudDescargaService interface {
    GetAll(c *fiber.Ctx) ([]model.SolicitudDescarga, error)
    GetByID(c *fiber.Ctx, id string) (*model.SolicitudDescarga, error)
    Create(c *fiber.Ctx, req *validation.CreateSolicitudDescarga) (*model.SolicitudDescarga, error)
    Update(c *fiber.Ctx, req *validation.UpdateSolicitudDescarga, id string) (*model.SolicitudDescarga, error)
    Delete(c *fiber.Ctx, id string) error
}

type solicitudDescargaService struct {
    Log      *logrus.Logger
    DB       *gorm.DB
    Validate *validator.Validate
}

func NewSolicitudDescargaService(db *gorm.DB, validate *validator.Validate) SolicitudDescargaService {
    return &solicitudDescargaService{
        Log:      utils.Log,
        DB:       db,
        Validate: validate,
    }
}

func (s *solicitudDescargaService) GetAll(c *fiber.Ctx) ([]model.SolicitudDescarga, error) {
    var solicitudes []model.SolicitudDescarga

    result := s.DB.WithContext(c.Context()).Order("created_at desc").Find(&solicitudes)
    if result.Error != nil {
        s.Log.Errorf("Error al obtener solicitudes: %+v", result.Error)
        return nil, result.Error
    }

    return solicitudes, nil
}

func (s *solicitudDescargaService) GetByID(c *fiber.Ctx, id string) (*model.SolicitudDescarga, error) {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    solicitud := new(model.SolicitudDescarga)
    result := s.DB.WithContext(c.Context()).First(solicitud, "id = ?", uuidID)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, fiber.NewError(fiber.StatusNotFound, "Solicitud no encontrada")
    }
    if result.Error != nil {
        s.Log.Errorf("Error al obtener solicitud: %+v", result.Error)
        return nil, result.Error
    }

    return solicitud, nil
}

func (s *solicitudDescargaService) Create(c *fiber.Ctx, req *validation.CreateSolicitudDescarga) (*model.SolicitudDescarga, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateCreateSolicitudDescarga(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    solicitud := req.ToModel()

    result := s.DB.WithContext(c.Context()).Create(solicitud)
    if result.Error != nil {
        s.Log.Errorf("Error al crear solicitud: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo crear la solicitud")
    }

    return solicitud, nil
}

func (s *solicitudDescargaService) Update(c *fiber.Ctx, req *validation.UpdateSolicitudDescarga, id string) (*model.SolicitudDescarga, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateUpdateSolicitudDescarga(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    solicitud := req.ToModel()
    solicitud.ID = uuidID.String()

    result := s.DB.WithContext(c.Context()).Model(&model.SolicitudDescarga{}).Where("id = ?", uuidID).Updates(solicitud)
    if result.Error != nil {
        s.Log.Errorf("Error al actualizar solicitud: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo actualizar la solicitud")
    }

    return solicitud, nil
}

func (s *solicitudDescargaService) Delete(c *fiber.Ctx, id string) error {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    result := s.DB.WithContext(c.Context()).Delete(&model.SolicitudDescarga{}, "id = ?", uuidID)
    if result.Error != nil {
        s.Log.Errorf("Error al eliminar solicitud: %+v", result.Error)
        return fiber.NewError(fiber.StatusInternalServerError, "No se pudo eliminar la solicitud")
    }

    if result.RowsAffected == 0 {
        return fiber.NewError(fiber.StatusNotFound, "Solicitud no encontrada")
    }

    return nil
}
