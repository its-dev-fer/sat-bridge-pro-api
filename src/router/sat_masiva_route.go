package router

import (
	"app/src/controller"
	m "app/src/middleware"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

func SatMasivaRoutes(v1 fiber.Router, sat service.SatMasivaService, u service.UserService) {
	ctrl := controller.NewSatMasivaController(sat)
	g := v1.Group("/descarga-masiva")
	g.Use(m.Auth(u))
	g.Get("/fiel", ctrl.FielStatus)
	g.Post("/autenticar", ctrl.Authenticate)
	g.Post("/solicitar", ctrl.Query)
	g.Post("/verificar", ctrl.Verify)
	g.Post("/descargar", ctrl.Download)
	g.Post("/backfill", ctrl.Backfill)
	g.Get("/sync", ctrl.SyncStatus)
	g.Post("/sync", ctrl.StartSync)
	g.Post("/sync/abort", ctrl.AbortSync)
}
