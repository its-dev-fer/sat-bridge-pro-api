package controller

import (
	"app/src/response"
	"app/src/service"
	"app/src/validation"

	"github.com/gofiber/fiber/v2"
)

type CIECController struct {
	CIECService service.CIECService
}

func NewCIECController(ciecService service.CIECService) *CIECController {
	return &CIECController{
		CIECService: ciecService,
	}
}

// @Tags         CIEC
// @Summary      Save CIEC credentials
// @Description  Save CIEC (Clave de Identificación Electrónica Confidencial) for the authenticated user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.CreateCIEC  true  "CIEC data"
// @Router       /ciec [post]
// @Success      201  {object}  response.Common
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *CIECController) SaveCIEC(c *fiber.Ctx) error {
	req := new(validation.CreateCIEC)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	if err := ctrl.CIECService.SaveCIEC(c, req, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(response.Common{
		Code:    fiber.StatusCreated,
		Status:  "success",
		Message: "CIEC saved successfully",
	})
}

// @Tags         CIEC
// @Summary      Update CIEC credentials
// @Description  Update CIEC for the authenticated user
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  validation.UpdateCIEC  true  "New CIEC"
// @Router       /ciec [patch]
// @Success      200  {object}  response.Common
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *CIECController) UpdateCIEC(c *fiber.Ctx) error {
	req := new(validation.UpdateCIEC)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	if err := ctrl.CIECService.UpdateCIEC(c, req, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CIEC updated successfully",
	})
}

// @Tags         CIEC
// @Summary      Delete CIEC credentials
// @Description  Delete CIEC for the authenticated user
// @Security     BearerAuth
// @Produce      json
// @Router       /ciec [delete]
// @Success      200  {object}  response.Common
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *CIECController) DeleteCIEC(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	if err := ctrl.CIECService.DeleteCIEC(c, userID); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Common{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CIEC deleted successfully",
	})
}

// @Tags         CIEC
// @Summary      Check if user has CIEC
// @Description  Check if the authenticated user has CIEC configured
// @Security     BearerAuth
// @Produce      json
// @Router       /ciec/status [get]
// @Success      200  {object}  response.Success
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *CIECController) HasCIEC(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	hasCIEC, err := ctrl.CIECService.HasCIEC(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CIEC status retrieved successfully",
		Data: map[string]bool{
			"has_ciec": hasCIEC,
		},
	})
}

