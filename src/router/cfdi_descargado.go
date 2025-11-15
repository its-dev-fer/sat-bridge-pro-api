package router

import (
    "app/src/controller"
    "github.com/gofiber/fiber/v2"
)

func RegisterCfdiDescargadoRoutes(r fiber.Router, ctrl *controller.CfdiDescargadoController) {
    cfdis := r.Group("/cfdis-descargados")
    cfdis.Get("/", ctrl.GetCfdis)
    cfdis.Get("/:id", ctrl.GetCfdiByID)
    cfdis.Post("/", ctrl.CreateCfdi)
    cfdis.Put("/:id", ctrl.UpdateCfdi)
    cfdis.Delete("/:id", ctrl.DeleteCfdi)
}
