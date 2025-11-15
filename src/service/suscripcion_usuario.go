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

type SuscripcionUsuarioService interface {
    GetAll(c *fiber.Ctx) ([]model.SuscripcionUsuario, error)
    GetByID(c *fiber.Ctx, id string) (*model.SuscripcionUsuario, error)
    Create(c *fiber.Ctx, req *validation.CreateSuscripcionUsuario) (*model.SuscripcionUsuario, error)
    Update(c *fiber.Ctx, req *validation.UpdateSuscripcionUsuario, id string) (*model.SuscripcionUsuario, error)
    Delete(c *fiber.Ctx, id string) error
}

type suscripcionUsuarioService struct {
    Log      *logrus.Logger
    DB       *gorm.DB
    Validate *validator.Validate
}

func NewSuscripcionUsuarioService(db *gorm.DB, validate *validator.Validate) SuscripcionUsuarioService {
    return &suscripcionUsuarioService{
        Log:      utils.Log,
        DB:       db,
        Validate: validate,
    }
}

func (s *suscripcionUsuarioService) GetAll(c *fiber.Ctx) ([]model.SuscripcionUsuario, error) {
    var subs []model.SuscripcionUsuario

    result := s.DB.WithContext(c.Context()).Order("created_at desc").Find(&subs)
    if result.Error != nil {
        s.Log.Errorf("Error al obtener suscripciones: %+v", result.Error)
        return nil, result.Error
    }

    return subs, nil
}

func (s *suscripcionUsuarioService) GetByID(c *fiber.Ctx, id string) (*model.SuscripcionUsuario, error) {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    sub := new(model.SuscripcionUsuario)
    result := s.DB.WithContext(c.Context()).First(sub, "id = ?", uuidID)

    if errors.Is(result.Error, gorm.ErrRecordNotFound) {
        return nil, fiber.NewError(fiber.StatusNotFound, "Suscripción no encontrada")
    }
    if result.Error != nil {
        s.Log.Errorf("Error al obtener suscripción: %+v", result.Error)
        return nil, result.Error
    }

    return sub, nil
}

func (s *suscripcionUsuarioService) Create(c *fiber.Ctx, req *validation.CreateSuscripcionUsuario) (*model.SuscripcionUsuario, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateCreateSuscripcionUsuario(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    sub := req.ToModel()

    result := s.DB.WithContext(c.Context()).Create(sub)
    if result.Error != nil {
        s.Log.Errorf("Error al crear suscripción: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo crear la suscripción")
    }

    return sub, nil
}

func (s *suscripcionUsuarioService) Update(c *fiber.Ctx, req *validation.UpdateSuscripcionUsuario, id string) (*model.SuscripcionUsuario, error) {
    if err := s.Validate.Struct(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    if err := validation.ValidateUpdateSuscripcionUsuario(req); err != nil {
        return nil, fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
    }

    uuidID, err := uuid.Parse(id)
    if err != nil {
        return nil, fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    sub := req.ToModel()
    sub.ID = uuidID.String()

    result := s.DB.WithContext(c.Context()).Model(&model.SuscripcionUsuario{}).Where("id = ?", uuidID).Updates(sub)
    if result.Error != nil {
        s.Log.Errorf("Error al actualizar suscripción: %+v", result.Error)
        return nil, fiber.NewError(fiber.StatusInternalServerError, "No se pudo actualizar la suscripción")
    }

    return sub, nil
}

func (s *suscripcionUsuarioService) Delete(c *fiber.Ctx, id string) error {
    uuidID, err := uuid.Parse(id)
    if err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Formato de ID inválido")
    }

    result := s.DB.WithContext(c.Context()).Delete(&model.SuscripcionUsuario{}, "id = ?", uuidID)
    if result.Error != nil {
        s.Log.Errorf("Error al eliminar suscripción: %+v", result.Error)
        return fiber.NewError(fiber.StatusInternalServerError, "No se pudo eliminar la suscripción")
    }

    if result.RowsAffected == 0 {
        return fiber.NewError(fiber.StatusNotFound, "Suscripción no encontrada")
    }

    return nil
}
