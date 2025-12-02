package controller

import (
	"app/src/response"
	"app/src/service"

	"github.com/gofiber/fiber/v2"
)

type SatDownloadController struct {
	SatService   service.SatDownloadService
	LimitService service.DownloadLimitService
}

func NewSatDownloadController(satService service.SatDownloadService, limitService service.DownloadLimitService) *SatDownloadController {
	return &SatDownloadController{
		SatService:   satService,
		LimitService: limitService,
	}
}

// @Tags         SAT Downloads
// @Summary      Download CFDIs from SAT
// @Description  Download CFDIs from SAT portal using CIEC or FIEL authentication
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  service.DownloadCFDIRequest  true  "Download parameters"
// @Router       /sat/download [post]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
func (ctrl *SatDownloadController) DownloadCFDIs(c *fiber.Ctx) error {
	req := new(service.DownloadCFDIRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	result, err := ctrl.SatService.DownloadCFDIs(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CFDIs downloaded successfully",
		Data:    result,
	})
}

// @Tags         SAT Downloads
// @Summary      Query CFDI metadata
// @Description  Query CFDI metadata without downloading the XMLs
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  service.DownloadCFDIRequest  true  "Query parameters"
// @Router       /sat/query-metadata [post]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *SatDownloadController) QueryMetadata(c *fiber.Ctx) error {
	req := new(service.DownloadCFDIRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	result, err := ctrl.SatService.QueryMetadata(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Metadata queried successfully",
		Data:    result,
	})
}

// @Tags         SAT Downloads
// @Summary      Download CFDI by UUID
// @Description  Download a specific CFDI by its UUID
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        request  body  service.DownloadByUUIDRequest  true  "UUID and date range"
// @Router       /sat/download-uuid [post]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      403  {object}  response.ErrorResponse
func (ctrl *SatDownloadController) DownloadByUUID(c *fiber.Ctx) error {
	req := new(service.DownloadByUUIDRequest)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	userID := c.Locals("userId").(string)

	result, err := ctrl.SatService.DownloadByUUID(c, req, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CFDI downloaded successfully",
		Data:    result,
	})
}

// @Tags         SAT Downloads
// @Summary      Get download statistics
// @Description  Get download limits and usage statistics for the authenticated user
// @Security     BearerAuth
// @Produce      json
// @Router       /sat/stats [get]
// @Success      200  {object}  response.Success
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *SatDownloadController) GetDownloadStats(c *fiber.Ctx) error {
	userID := c.Locals("userId").(string)

	stats, err := ctrl.LimitService.GetDownloadStats(c, userID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Download statistics retrieved successfully",
		Data:    stats,
	})
}

