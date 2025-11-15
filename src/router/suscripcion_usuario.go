package router

import (
    "app/src/controller"
    "github.com/gofiber/fiber/v2"
)

func RegisterSuscripcionUsuarioRoutes(r fiber.Router, ctrl *controller.SuscripcionUsuarioController) {
    suscripciones := r.Group("/suscripciones")
    suscripciones.Get("/", ctrl.GetSuscripciones)
    suscripciones.Get("/:id", ctrl.GetSuscripcionByID)
    suscripciones.Post("/", ctrl.CreateSuscripcion)
    suscripciones.Put("/:id", ctrl.UpdateSuscripcion)
    suscripciones.Delete("/:id", ctrl.DeleteSuscripcion)
}
