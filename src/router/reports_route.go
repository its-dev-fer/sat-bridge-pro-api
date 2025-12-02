package router

import (
	"app/src/controller"
	"app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func ReportsRoutes(app fiber.Router, db *gorm.DB) {
	reportsService := service.NewReportsService(db)
	reportsController := controller.NewReportsController(reportsService)

	reports := app.Group("/reports")
	reports.Use(middleware.Authenticate())

	reports.Get("/monthly", reportsController.GetMonthlyReport)
	reports.Get("/yearly", reportsController.GetYearlyReport)
	reports.Get("/date-range", reportsController.GetCFDIsByDateRange)
	reports.Get("/expenses", reportsController.GetExpensesSummary)
	reports.Get("/income", reportsController.GetIncomeSummary)
}

