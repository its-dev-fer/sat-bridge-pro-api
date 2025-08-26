package controller

import (
	"app/src/model"
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type DatosFiscalesController struct {
	DatosFiscalesService service.DatosFiscalesService
}

func NewDatosFiscalesController(datosFiscalesService service.DatosFiscalesService) *DatosFiscalesController {
	return &DatosFiscalesController{
		DatosFiscalesService: datosFiscalesService,
	}
}

// @Tags         Datos Fiscales
// @Summary      Register fiscal data
// @Description  Registra los datos fiscales del usuario autenticado
// @Security     BearerAuth
// @Accept       multipart/form-data
// @Produce      json
// @Param        rfc formData string true "RFC"
// @Param        password formData string true "E-firma password"
// @Param        cer_file formData file true "Certificate file (.cer)"
// @Param        key_file formData file true "Key file (.key)"
// @Router       /datos-fiscales [post]
// @Success      201  {object}  response.Common
// @Failure      400  {object}  response.Common  "Bad request"
// @Failure      401  {object}  response.Common  "Unauthorized"
// @Failure      409  {object}  response.Common  "Fiscal data already exists"
// @Failure      500  {object}  response.Common  "Internal server error"
func (c *DatosFiscalesController) CreateDatosFiscales(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)

	// Parse form data
	req := new(validation.DatosFiscalesRequest)
	req.RFC = ctx.FormValue("rfc")
	req.Password = ctx.FormValue("password")

	if req.RFC == "" || req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "RFC and password are required")
	}

	// Obtener archivos
	cerFile, err := ctx.FormFile("cer_file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Certificate file is required")
	}

	keyFile, err := ctx.FormFile("key_file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Key file is required")
	}

	// Llamar al servicio
	if err := c.DatosFiscalesService.CreateDatosFiscales(ctx, user.ID, req, cerFile, keyFile); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(response.Common{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Fiscal data registered successfully",
	})
}

// @Tags         Datos Fiscales
// @Summary      Get user fiscal data
// @Description  Get fiscal data for the authenticated user (without sensitive data)
// @Security     BearerAuth
// @Produce      json
// @Router       /datos-fiscales [get]
// @Success      200  {object}  response.Common
// @Failure      401  {object}  response.Common  "Unauthorized"
// @Failure      404  {object}  response.Common  "Not found"
func (c *DatosFiscalesController) GetDatosFiscales(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)

	datosFiscales, err := c.DatosFiscalesService.GetDatosFiscalesByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.SuccessWithData{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Fiscal data retrieved successfully",
		Data:    datosFiscales,
	})
}

// @Tags         Datos Fiscales
// @Summary      Update fiscal data
// @Description  Update RFC or password for fiscal data
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.DatosFiscalesRequest  true  "Request body"
// @Router       /datos-fiscales [patch]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.Common  "Bad request"
// @Failure      401  {object}  response.Common  "Unauthorized"
// @Failure      404  {object}  response.Common  "Not found"
func (c *DatosFiscalesController) UpdateDatosFiscales(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)

	req := new(validation.DatosFiscalesRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := c.DatosFiscalesService.UpdateDatosFiscales(ctx, user.ID, req); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Fiscal data updated successfully",
	})
}

// @Tags         Datos Fiscales
// @Summary      Delete fiscal data
// @Description  Delete all fiscal data for the authenticated user
// @Security     BearerAuth
// @Produce      json
// @Router       /datos-fiscales [delete]
// @Success      200  {object}  response.Common
// @Failure      401  {object}  response.Common  "Unauthorized"
// @Failure      404  {object}  response.Common  "Not found"
func (c *DatosFiscalesController) DeleteDatosFiscales(ctx *fiber.Ctx) error {
	user, _ := ctx.Locals("user").(*model.User)

	if err := c.DatosFiscalesService.DeleteDatosFiscales(ctx, user.ID); err != nil {
		return err
	}

	return ctx.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Fiscal data deleted successfully",
	})
}
