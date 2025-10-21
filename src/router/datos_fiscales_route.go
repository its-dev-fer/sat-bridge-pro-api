package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func DatosFiscalesRoutes(v1 fiber.Router, d service.DatosFiscalesService, u service.UserService) {
	datosFiscalesController := controller.NewDatosFiscalesController(d)

	datosFiscales := v1.Group("/datos-fiscales")

	datosFiscales.Use(m.Auth(u))

	datosFiscales.Post("/", datosFiscalesController.CreateDatosFiscales)
	datosFiscales.Get("/", datosFiscalesController.GetDatosFiscales)
	datosFiscales.Patch("/", datosFiscalesController.UpdateDatosFiscales)
	datosFiscales.Delete("/", datosFiscalesController.DeleteDatosFiscales)
}