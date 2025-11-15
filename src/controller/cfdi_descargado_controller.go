package controller

import (
    "app/src/service"
    "app/src/validation"
    "github.com/gofiber/fiber/v2"
)

// CfdiDescargadoController gestiona los endpoints de CFDIs descargados
type CfdiDescargadoController struct {
    Service service.CfdiDescargadoService
}

func NewCfdiDescargadoController(s service.CfdiDescargadoService) *CfdiDescargadoController {
    return &CfdiDescargadoController{Service: s}
}

// GetCfdis godoc
// @Summary Obtener lista de CFDIs descargados
// @Tags CFDIs Descargados
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} utils.CommonErrorResponse
// @Router /cfdis-descargados [get]
func (ctrl *CfdiDescargadoController) GetCfdis(c *fiber.Ctx) error {
    cfdis, err := ctrl.Service.GetAll(c)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "data":  cfdis,
        "total": len(cfdis),
    })
}

// GetCfdiByID godoc
// @Summary Obtener un CFDI descargado por ID
// @Tags CFDIs Descargados
// @Accept json
// @Produce json
// @Param id path string true "ID del CFDI"
// @Success 200 {object} model.CfdiDescargado
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /cfdis-descargados/{id} [get]
func (ctrl *CfdiDescargadoController) GetCfdiByID(c *fiber.Ctx) error {
    id := c.Params("id")
    cfdi, err := ctrl.Service.GetByID(c, id)
    if err != nil {
        return err
    }
    return c.JSON(cfdi)
}

// CreateCfdi godoc
// @Summary Crear un nuevo CFDI descargado
// @Tags CFDIs Descargados
// @Accept json
// @Produce json
// @Param cfdi body validation.CreateCfdiDescargado true "Datos del CFDI"
// @Success 201 {object} model.CfdiDescargado
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 422 {object} utils.DetailedErrorResponse
// @Router /cfdis-descargados [post]
func (ctrl *CfdiDescargadoController) CreateCfdi(c *fiber.Ctx) error {
    req := new(validation.CreateCfdiDescargado)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    cfdi, err := ctrl.Service.Create(c, req)
    if err != nil {
        return err
    }

    return c.Status(fiber.StatusCreated).JSON(cfdi)
}

// UpdateCfdi godoc
// @Summary Actualizar un CFDI descargado por ID
// @Tags CFDIs Descargados
// @Accept json
// @Produce json
// @Param id path string true "ID del CFDI"
// @Param cfdi body validation.UpdateCfdiDescargado true "Datos actualizados"
// @Success 200 {object} model.CfdiDescargado
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /cfdis-descargados/{id} [put]
func (ctrl *CfdiDescargadoController) UpdateCfdi(c *fiber.Ctx) error {
    id := c.Params("id")
    req := new(validation.UpdateCfdiDescargado)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Cuerpo de solicitud inválido")
    }

    cfdi, err := ctrl.Service.Update(c, req, id)
    if err != nil {
        return err
    }

    return c.JSON(cfdi)
}

// DeleteCfdi godoc
// @Summary Eliminar un CFDI descargado por ID
// @Tags CFDIs Descargados
// @Accept json
// @Produce json
// @Param id path string true "ID del CFDI"
// @Success 204 "CFDI eliminado exitosamente"
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /cfdis-descargados/{id} [delete]
func (ctrl *CfdiDescargadoController) DeleteCfdi(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := ctrl.Service.Delete(c, id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
