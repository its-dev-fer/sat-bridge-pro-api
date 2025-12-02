package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func FirmaElectronicaRoutes(app fiber.Router, db *gorm.DB) {
	validate := validation.New()
	firmaService := service.NewFirmaElectronicaService(db, validate)
	firmaController := controller.NewFirmaElectronicaController(firmaService)

	fiel := app.Group("/fiel")
	fiel.Use(middleware.Authenticate())

	fiel.Post("/", firmaController.CreateFirmaElectronica)
	fiel.Get("/", firmaController.GetFirmaByUser)
	fiel.Get("/all", firmaController.GetAllFirmasByUser)
	fiel.Get("/:id", firmaController.GetFirmaByID)
	fiel.Patch("/:id", firmaController.UpdateFirmaElectronica)
	fiel.Delete("/:id", firmaController.DeleteFirmaElectronica)
}

