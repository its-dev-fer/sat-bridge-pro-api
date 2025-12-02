package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SatDownloadRoutes(app fiber.Router, db *gorm.DB) {
	validate := validation.New()
	
	// Services
	limitService := service.NewDownloadLimitService(db)
	firmaService := service.NewFirmaElectronicaService(db, validate)
	ciecService := service.NewCIECService(db, validate, firmaService)
	satService := service.NewSatDownloadService(db, limitService, firmaService, ciecService)
	
	// Controller
	satController := controller.NewSatDownloadController(satService, limitService)

	sat := app.Group("/sat")
	sat.Use(middleware.Authenticate())

	sat.Post("/download", satController.DownloadCFDIs)
	sat.Post("/query-metadata", satController.QueryMetadata)
	sat.Post("/download-uuid", satController.DownloadByUUID)
	sat.Get("/stats", satController.GetDownloadStats)
}

