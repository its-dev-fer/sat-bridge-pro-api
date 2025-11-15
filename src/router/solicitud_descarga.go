package router

import (
    "app/src/controller"
    "github.com/gofiber/fiber/v2"
)

func RegisterSolicitudDescargaRoutes(r fiber.Router, ctrl *controller.SolicitudDescargaController) {
    solicitudes := r.Group("/solicitudes-descarga")
    solicitudes.Get("/", ctrl.GetSolicitudes)
    solicitudes.Get("/:id", ctrl.GetSolicitudByID)
    solicitudes.Post("/", ctrl.CreateSolicitud)
    solicitudes.Put("/:id", ctrl.UpdateSolicitud)
    solicitudes.Delete("/:id", ctrl.DeleteSolicitud)
}
