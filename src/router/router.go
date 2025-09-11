package router

import (
    "app/src/config"
    "app/src/controller"
    "app/src/service"
    "app/src/validation"

    "github.com/gofiber/fiber/v2"
    "gorm.io/gorm"
)

func Routes(app *fiber.App, db *gorm.DB) {
    validate := validation.Validator()

    healthCheckService := service.NewHealthCheckService(db)
    emailService := service.NewEmailService()
    userService := service.NewUserService(db, validate)
    tokenService := service.NewTokenService(db, validate, userService)
    authService := service.NewAuthService(db, validate, userService, tokenService)

    // Plan de suscripción/
    planService := service.NewPlanSuscripcionService(db, validate)
    planController := controller.NewPlanSuscripcionController(planService)

    v1 := app.Group("/v1")

    HealthCheckRoutes(v1, healthCheckService)
    AuthRoutes(v1, authService, userService, tokenService, emailService)
    UserRoutes(v1, userService, tokenService)

    //planes de suscripción
    RegisterPlanSuscripcionRoutes(v1, planController)

    if !config.IsProd {
        DocsRoutes(v1)
    }
}
