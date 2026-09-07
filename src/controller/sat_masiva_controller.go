package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type SatMasivaController struct {
	SatMasivaService service.SatMasivaService
}

func NewSatMasivaController(satMasivaService service.SatMasivaService) *SatMasivaController {
	return &SatMasivaController{SatMasivaService: satMasivaService}
}

// @Tags         Descarga Masiva SAT
// @Summary      FIEL status
// @Description  Valida la e.firma guardada (vigencia, FIEL vs CSD)
// @Security     BearerAuth
// @Produce      json
// @Router       /descarga-masiva/fiel [get]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) FielStatus(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	data, err := c.SatMasivaService.FielStatus(ctx, user.ID)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code: fiber.StatusOK, Status: "success", Message: "FIEL status", Data: data,
	})
}

// @Tags         Descarga Masiva SAT
// @Summary      Autentica
// @Description  Autentica contra el SAT con la FIEL (no expone el token WRAP)
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.SatServiceRequest  false  "cfdi o retenciones"
// @Router       /descarga-masiva/autenticar [post]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) Authenticate(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	req := new(validation.SatServiceRequest)
	if len(ctx.Body()) > 0 {
		if err := ctx.BodyParser(req); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}
	}
	data, err := c.SatMasivaService.Authenticate(ctx, user.ID, req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code: fiber.StatusOK, Status: "success", Message: "Authenticated", Data: data,
	})
}

// @Tags         Descarga Masiva SAT
// @Summary      SolicitaDescarga
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.SatQueryRequest  true  "Query"
// @Router       /descarga-masiva/solicitar [post]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) Query(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	req := new(validation.SatQueryRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	data, err := c.SatMasivaService.Query(ctx, user.ID, req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code: fiber.StatusOK, Status: "success", Message: "Solicitud creada", Data: data,
	})
}

// @Tags         Descarga Masiva SAT
// @Summary      VerificaSolicitudDescarga
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.SatVerifyRequest  true  "Verify"
// @Router       /descarga-masiva/verificar [post]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) Verify(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	req := new(validation.SatVerifyRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	data, err := c.SatMasivaService.Verify(ctx, user.ID, req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code: fiber.StatusOK, Status: "success", Message: "Solicitud verificada", Data: data,
	})
}

// @Tags         Descarga Masiva SAT
// @Summary      DescargaMasiva
// @Description  ZIP del paquete, o JSON si parse=cfdi|metadata
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.SatDownloadRequest  true  "Download"
// @Router       /descarga-masiva/descargar [post]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) Download(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	req := new(validation.SatDownloadRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	zip, parsed, err := c.SatMasivaService.Download(ctx, user.ID, req)
	if err != nil {
		return err
	}
	if parsed != nil {
		return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
			Code: fiber.StatusOK, Status: "success", Message: "Paquete parseado", Data: parsed,
		})
	}
	ctx.Set("Content-Type", "application/zip")
	ctx.Set("Content-Disposition", `attachment; filename="`+req.PackageID+`.zip"`)
	return ctx.Status(fiber.StatusOK).Send(zip)
}

// @Tags         Descarga Masiva SAT
// @Summary      Backfill
// @Description  Trocea, verifica, descarga y parsea un rango. Sincrónico.
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.SatBackfillRequest  true  "Backfill"
// @Router       /descarga-masiva/backfill [post]
// @Success      200  {object}  response.SuccessWithData
func (c *SatMasivaController) Backfill(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)
	req := new(validation.SatBackfillRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	data, err := c.SatMasivaService.Backfill(ctx, user.ID, req)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code: fiber.StatusOK, Status: "success", Message: "Backfill completo", Data: data,
	})
}
