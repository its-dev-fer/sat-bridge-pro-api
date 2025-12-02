package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type FirmaElectronicaController struct {
	FirmaService service.FirmaElectronicaService
}

func NewFirmaElectronicaController(firmaService service.FirmaElectronicaService) *FirmaElectronicaController {
	return &FirmaElectronicaController{
		FirmaService: firmaService,
	}
}

// @Tags         FIEL
// @Summary      Create digital signature (FIEL)
// @Description  Create a new FIEL for the authenticated user with .key and .cer files in base64
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.CreateFirmaElectronica  true  "FIEL data"
// @Router       /fiel [post]
// @Success      201  {object}  response.Common
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) CreateFirmaElectronica(c *fiber.Ctx) error {
	req := new(validation.CreateFirmaElectronica)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	firma, err := ctrl.FirmaService.CreateFirmaElectronica(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.Success{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "Firma electronica created successfully",
		Data:    firma.ID,
	})
}

// @Tags         FIEL
// @Summary      Get user's active FIEL
// @Description  Get the active digital signature for the authenticated user
// @Security     BearerAuth
// @Produce      json
// @Router       /fiel [get]
// @Success      200  {object}  response.Success
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) GetFirmaByUser(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	firma, err := ctrl.FirmaService.GetFirmaByUserID(c, userID)
	if err != nil {
		return err
	}

	// Remove sensitive data before sending
	firma.ClavePrivadaKEY = ""
	firma.PasswordKEY = ""

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Firma electronica retrieved successfully",
		Data:    firma,
	})
}

// @Tags         FIEL
// @Summary      Get all FIELs for user
// @Description  Get all digital signatures for the authenticated user
// @Security     BearerAuth
// @Produce      json
// @Router       /fiel/all [get]
// @Success      200  {object}  response.Success
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) GetAllFirmasByUser(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	firmas, err := ctrl.FirmaService.GetAllFirmasByUser(c, userID)
	if err != nil {
		return err
	}

	// Remove sensitive data
	for i := range firmas {
		firmas[i].ClavePrivadaKEY = ""
		firmas[i].PasswordKEY = ""
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Firmas electronicas retrieved successfully",
		Data:    firmas,
	})
}

// @Tags         FIEL
// @Summary      Get FIEL by ID
// @Description  Get a specific digital signature by ID (only owner can access)
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "FIEL ID"
// @Router       /fiel/{id} [get]
// @Success      200  {object}  response.Success
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) GetFirmaByID(c *fiber.Ctx) error {
	firmaID := c.Params("id")

	if _, err := uuid.Parse(firmaID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid firma ID")
	}

	userID := c.Locals("userId").(string)

	firma, err := ctrl.FirmaService.GetFirmaByID(c, firmaID)
	if err != nil {
		return err
	}

	// Verify ownership
	if firma.UsuarioID.String() != userID {
		return fiber.NewError(fiber.StatusForbidden, "Access denied")
	}

	// Remove sensitive data
	firma.ClavePrivadaKEY = ""
	firma.PasswordKEY = ""

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Firma electronica retrieved successfully",
		Data:    firma,
	})
}

// @Tags         FIEL
// @Summary      Update FIEL
// @Description  Update digital signature information
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "FIEL ID"
// @Param        request  body  validation.UpdateFirmaElectronica  true  "Update data"
// @Router       /fiel/{id} [patch]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) UpdateFirmaElectronica(c *fiber.Ctx) error {
	firmaID := c.Params("id")

	if _, err := uuid.Parse(firmaID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid firma ID")
	}

	req := new(validation.UpdateFirmaElectronica)
	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	firma, err := ctrl.FirmaService.UpdateFirmaElectronica(c, req, firmaID, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Firma electronica updated successfully",
		Data:    firma.ID,
	})
}

// @Tags         FIEL
// @Summary      Delete FIEL
// @Description  Delete a digital signature
// @Security     BearerAuth
// @Produce      json
// @Param        id  path  string  true  "FIEL ID"
// @Router       /fiel/{id} [delete]
// @Success      200  {object}  response.Common
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
func (ctrl *FirmaElectronicaController) DeleteFirmaElectronica(c *fiber.Ctx) error {
	firmaID := c.Params("id")

	if _, err := uuid.Parse(firmaID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid firma ID")
	}

	userID := c.Locals("userId").(string)

	if err := ctrl.FirmaService.DeleteFirmaElectronica(c, firmaID, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Firma electronica deleted successfully",
	})
}

