package router

import (
	"app/src/config"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
	validate := validation.Validator()

	// Servicios existentes
	healthCheckService := service.NewHealthCheckService(db)
	emailService := service.NewEmailService()
	userService := service.NewUserService(db, validate)
	tokenService := service.NewTokenService(db, validate, userService)
	authService := service.NewAuthService(db, validate, userService, tokenService)
	
	// NUEVO: Servicio de datos fiscales usando config.EncryptionKey
	datosFiscalesService := service.NewDatosFiscalesService(db, validate, config.EncryptionKey)

	v1 := app.Group("/v1")

	// Rutas existentes
	HealthCheckRoutes(v1, healthCheckService)
	AuthRoutes(v1, authService, userService, tokenService, emailService)
	UserRoutes(v1, userService, tokenService)
	
	// NUEVA: Ruta de datos fiscales
	DatosFiscalesRoutes(v1, datosFiscalesService, userService)

	if !config.IsProd {
		DocsRoutes(v1)
	}
}