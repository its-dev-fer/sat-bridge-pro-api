package router

import (
    "app/src/controller"
    "github.com/gofiber/fiber/v2"
)

func RegisterPlanSuscripcionRoutes(r fiber.Router, ctrl *controller.PlanSuscripcionController) {
    planes := r.Group("/planes")
    planes.Get("/", ctrl.GetPlanes)
    planes.Get("/:id", ctrl.GetPlanByID)
    planes.Post("/", ctrl.CreatePlan)
    planes.Put("/:id", ctrl.UpdatePlan)
    planes.Delete("/:id", ctrl.DeletePlan)
}
