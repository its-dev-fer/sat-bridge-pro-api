package controller

import (
    "app/src/service"
    "app/src/validation"
    "github.com/gofiber/fiber/v2"
)

// SuscripcionUsuarioController gestiona los endpoints de suscripciones de usuario
type SuscripcionUsuarioController struct {
    Service service.SuscripcionUsuarioService
}

func NewSuscripcionUsuarioController(s service.SuscripcionUsuarioService) *SuscripcionUsuarioController {
    return &SuscripcionUsuarioController{Service: s}
}

// GetSuscripciones godoc
// @Summary Obtener lista de suscripciones de usuario
// @Tags Suscripciones de Usuario
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.CommonErrorResponse
// @Router /suscripciones [get]
func (ctrl *SuscripcionUsuarioController) GetSuscripciones(c *fiber.Ctx) error {
    subs, err := ctrl.Service.GetAll(c)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "data":  subs,
        "total": len(subs),
    })
}

// GetSuscripcionByID godoc
// @Summary Obtener una suscripción por ID
// @Tags Suscripciones de Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Success 200 {object} model.SuscripcionUsuario
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /suscripciones/{id} [get]
func (ctrl *SuscripcionUsuarioController) GetSuscripcionByID(c *fiber.Ctx) error {
    id := c.Params("id")
    sub, err := ctrl.Service.GetByID(c, id)
    if err != nil {
        return err
    }
    return c.JSON(sub)
}

// CreateSuscripcion godoc
// @Summary Crear una nueva suscripción de usuario
// @Tags Suscripciones de Usuario
// @Accept json
// @Produce json
// @Param suscripcion body validation.CreateSuscripcionUsuario true "Datos de la suscripción"
// @Success 201 {object} model.SuscripcionUsuario
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 422 {object} utils.DetailedErrorResponse
// @Router /suscripciones [post]
func (ctrl *SuscripcionUsuarioController) CreateSuscripcion(c *fiber.Ctx) error {
    req := new(validation.CreateSuscripcionUsuario)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    sub, err := ctrl.Service.Create(c, req)
    if err != nil {
        return err
    }

    return c.Status(fiber.StatusCreated).JSON(sub)
}

// UpdateSuscripcion godoc
// @Summary Actualizar una suscripción por ID
// @Tags Suscripciones de Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Param suscripcion body validation.UpdateSuscripcionUsuario true "Datos actualizados"
// @Success 200 {object} model.SuscripcionUsuario
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /suscripciones/{id} [put]
func (ctrl *SuscripcionUsuarioController) UpdateSuscripcion(c *fiber.Ctx) error {
    id := c.Params("id")
    req := new(validation.UpdateSuscripcionUsuario)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    sub, err := ctrl.Service.Update(c, req, id)
    if err != nil {
        return err
    }

    return c.JSON(sub)
}

// DeleteSuscripcion godoc
// @Summary Eliminar una suscripción por ID
// @Tags Suscripciones de Usuario
// @Accept json
// @Produce json
// @Param id path string true "ID de la suscripción"
// @Success 204 "Suscripción eliminada exitosamente"
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /suscripciones/{id} [delete]
func (ctrl *SuscripcionUsuarioController) DeleteSuscripcion(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := ctrl.Service.Delete(c, id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
