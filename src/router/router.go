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

    // Plan de suscripción
    planService := service.NewPlanSuscripcionService(db, validate)
    planController := controller.NewPlanSuscripcionController(planService)

    // Suscripción de usuario
    suscripcionUsuarioService := service.NewSuscripcionUsuarioService(db, validate)
    suscripcionUsuarioController := controller.NewSuscripcionUsuarioController(suscripcionUsuarioService)

    //solicitudes de descarga
    solicitudDescargaService := service.NewSolicitudDescargaService(db, validate)
    solicitudDescargaController := controller.NewSolicitudDescargaController(solicitudDescargaService)

     // CFDIs descargados
    cfdiDescargadoService := service.NewCfdiDescargadoService(db, validate)
    cfdiDescargadoController := controller.NewCfdiDescargadoController(cfdiDescargadoService)



    v1 := app.Group("/v1")

    HealthCheckRoutes(v1, healthCheckService)
    AuthRoutes(v1, authService, userService, tokenService, emailService)
    UserRoutes(v1, userService, tokenService)



    // Planes de suscripción
    RegisterPlanSuscripcionRoutes(v1, planController)
    
    // Suscripciones de usuario
    RegisterSuscripcionUsuarioRoutes(v1, suscripcionUsuarioController)

     // Solicitudes de descarga
    RegisterSolicitudDescargaRoutes(v1, solicitudDescargaController)

     // CFDIs descargados
    RegisterCfdiDescargadoRoutes(v1, cfdiDescargadoController)

    if !config.IsProd {
        DocsRoutes(v1)
    }
}