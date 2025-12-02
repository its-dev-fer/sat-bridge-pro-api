package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CIECRoutes(app fiber.Router, db *gorm.DB) {
	validate := validation.New()
	firmaService := service.NewFirmaElectronicaService(db, validate)
	ciecService := service.NewCIECService(db, validate, firmaService)
	ciecController := controller.NewCIECController(ciecService)

	ciec := app.Group("/ciec")
	ciec.Use(middleware.Authenticate())

	ciec.Post("/", ciecController.SaveCIEC)
	ciec.Patch("/", ciecController.UpdateCIEC)
	ciec.Delete("/", ciecController.DeleteCIEC)
	ciec.Get("/status", ciecController.HasCIEC)
}

