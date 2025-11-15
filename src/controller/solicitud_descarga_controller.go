package controller

import (
    "app/src/service"
    "app/src/validation"
    "github.com/gofiber/fiber/v2"
)

// SolicitudDescargaController gestiona los endpoints de solicitudes de descarga
type SolicitudDescargaController struct {
    Service service.SolicitudDescargaService
}

func NewSolicitudDescargaController(s service.SolicitudDescargaService) *SolicitudDescargaController {
    return &SolicitudDescargaController{Service: s}
}

// GetSolicitudes godoc
// @Summary Obtener lista de solicitudes de descarga
// @Tags Solicitudes de Descarga
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.CommonErrorResponse
// @Router /solicitudes-descarga [get]
func (ctrl *SolicitudDescargaController) GetSolicitudes(c *fiber.Ctx) error {
    solicitudes, err := ctrl.Service.GetAll(c)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "data":  solicitudes,
        "total": len(solicitudes),
    })
}

// GetSolicitudByID godoc
// @Summary Obtener una solicitud de descarga por ID
// @Tags Solicitudes de Descarga
// @Accept json
// @Produce json
// @Param id path string true "ID de la solicitud"
// @Success 200 {object} model.SolicitudDescarga
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /solicitudes-descarga/{id} [get]
func (ctrl *SolicitudDescargaController) GetSolicitudByID(c *fiber.Ctx) error {
    id := c.Params("id")
    solicitud, err := ctrl.Service.GetByID(c, id)
    if err != nil {
        return err
    }
    return c.JSON(solicitud)
}

// CreateSolicitud godoc
// @Summary Crear una nueva solicitud de descarga
// @Tags Solicitudes de Descarga
// @Accept json
// @Produce json
// @Param solicitud body validation.CreateSolicitudDescarga true "Datos de la solicitud"
// @Success 201 {object} model.SolicitudDescarga
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 422 {object} utils.DetailedErrorResponse
// @Router /solicitudes-descarga [post]
func (ctrl *SolicitudDescargaController) CreateSolicitud(c *fiber.Ctx) error {
    req := new(validation.CreateSolicitudDescarga)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    solicitud, err := ctrl.Service.Create(c, req)
    if err != nil {
        return err
    }

    return c.Status(fiber.StatusCreated).JSON(solicitud)
}

// UpdateSolicitud godoc
// @Summary Actualizar una solicitud de descarga por ID
// @Tags Solicitudes de Descarga
// @Accept json
// @Produce json
// @Param id path string true "ID de la solicitud"
// @Param solicitud body validation.UpdateSolicitudDescarga true "Datos actualizados"
// @Success 200 {object} model.SolicitudDescarga
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /solicitudes-descarga/{id} [put]
func (ctrl *SolicitudDescargaController) UpdateSolicitud(c *fiber.Ctx) error {
    id := c.Params("id")
    req := new(validation.UpdateSolicitudDescarga)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    solicitud, err := ctrl.Service.Update(c, req, id)
    if err != nil {
        return err
    }

    return c.JSON(solicitud)
}

// DeleteSolicitud godoc
// @Summary Eliminar una solicitud de descarga por ID
// @Tags Solicitudes de Descarga
// @Accept json
// @Produce json
// @Param id path string true "ID de la solicitud"
// @Success 204 "Solicitud eliminada exitosamente"
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /solicitudes-descarga/{id} [delete]
func (ctrl *SolicitudDescargaController) DeleteSolicitud(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := ctrl.Service.Delete(c, id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
