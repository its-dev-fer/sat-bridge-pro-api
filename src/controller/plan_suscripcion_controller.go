package controller

import (
    "app/src/service"
    "app/src/validation"
    "github.com/gofiber/fiber/v2"
)

// PlanSuscripcionController gestiona los endpoints de planes de suscripción
type PlanSuscripcionController struct {
    Service service.PlanSuscripcionService
}

func NewPlanSuscripcionController(s service.PlanSuscripcionService) *PlanSuscripcionController {
    return &PlanSuscripcionController{Service: s}
}

// GetPlanes godoc
// @Summary Obtener lista de planes de suscripción
// @Tags Planes de Suscripción
// @Accept json
// @Produce json
// @Param page query int true "Número de página"
// @Param limit query int true "Cantidad por página"
// @Param search query string false "Buscar por nombre o descripción"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 500 {object} utils.CommonErrorResponse
// @Router /planes [get]
func (ctrl *PlanSuscripcionController) GetPlanes(c *fiber.Ctx) error {
    params := new(validation.QueryPlanSuscripcion)
    if err := c.QueryParser(params); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid query parameters")
    }

    planes, total, err := ctrl.Service.GetPlanes(c, params)
    if err != nil {
        return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }

    return c.JSON(fiber.Map{
        "data":  planes,
        "total": total,
    })
}

// GetPlanByID godoc
// @Summary Obtener un plan por ID
// @Tags Planes de Suscripción
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Success 200 {object} model.PlanSuscripcion
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /planes/{id} [get]
func (ctrl *PlanSuscripcionController) GetPlanByID(c *fiber.Ctx) error {
    id := c.Params("id")
    plan, err := ctrl.Service.GetPlanByID(c, id)
    if err != nil {
        return err
    }
    return c.JSON(plan)
}

// CreatePlan godoc
// @Summary Crear un nuevo plan de suscripción
// @Tags Planes de Suscripción
// @Accept json
// @Produce json
// @Param plan body validation.CreatePlanSuscripcion true "Datos del plan"
// @Success 201 {object} model.PlanSuscripcion
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 422 {object} utils.DetailedErrorResponse
// @Router /planes [post]
func (ctrl *PlanSuscripcionController) CreatePlan(c *fiber.Ctx) error {
    req := new(validation.CreatePlanSuscripcion)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }

    plan, err := ctrl.Service.CreatePlan(c, req)
    if err != nil {
        return err
    }
    return c.Status(fiber.StatusCreated).JSON(plan)
}

// UpdatePlan godoc
// @Summary Actualizar un plan por ID
// @Tags Planes de Suscripción
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Param plan body validation.UpdatePlanSuscripcion true "Datos actualizados"
// @Success 200 {object} model.PlanSuscripcion
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /planes/{id} [put]
func (ctrl *PlanSuscripcionController) UpdatePlan(c *fiber.Ctx) error {
    id := c.Params("id")
    req := new(validation.UpdatePlanSuscripcion)
    if err := c.BodyParser(req); err != nil {
        return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
    }

    plan, err := ctrl.Service.UpdatePlan(c, req, id)
    if err != nil {
        return err
    }
    return c.JSON(plan)
}

// DeletePlan godoc
// @Summary Eliminar un plan por ID
// @Tags Planes de Suscripción
// @Accept json
// @Produce json
// @Param id path string true "ID del plan"
// @Success 204 "Plan eliminado exitosamente"
// @Failure 400 {object} utils.CommonErrorResponse
// @Failure 404 {object} utils.CommonErrorResponse
// @Router /planes/{id} [delete]
func (ctrl *PlanSuscripcionController) DeletePlan(c *fiber.Ctx) error {
    id := c.Params("id")
    if err := ctrl.Service.DeletePlan(c, id); err != nil {
        return err
    }
    return c.SendStatus(fiber.StatusNoContent)
}
