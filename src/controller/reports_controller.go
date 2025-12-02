package controller

import (
	"app/src/response"
	"app/src/service"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ReportsController struct {
	ReportsService service.ReportsService
}

func NewReportsController(reportsService service.ReportsService) *ReportsController {
	return &ReportsController{
		ReportsService: reportsService,
	}
}

// @Tags         Reports
// @Summary      Get monthly report
// @Description  Get accounting report for a specific month
// @Security     BearerAuth
// @Produce      json
// @Param        year   query  int  true  "Year"
// @Param        month  query  int  true  "Month (1-12)"
// @Router       /reports/monthly [get]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *ReportsController) GetMonthlyReport(c *fiber.Ctx) error {
	year := c.QueryInt("year", time.Now().Year())
	month := c.QueryInt("month", int(time.Now().Month()))

	if month < 1 || month > 12 {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid month, must be between 1 and 12")
	}

	userID := c.Locals("userId").(string)

	report, err := ctrl.ReportsService.GetMonthlyReport(c, userID, year, month)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Monthly report retrieved successfully",
		Data:    report,
	})
}

// @Tags         Reports
// @Summary      Get yearly report
// @Description  Get accounting report for a specific year
// @Security     BearerAuth
// @Produce      json
// @Param        year  query  int  true  "Year"
// @Router       /reports/yearly [get]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *ReportsController) GetYearlyReport(c *fiber.Ctx) error {
	year := c.QueryInt("year", time.Now().Year())

	userID := c.Locals("userId").(string)

	report, err := ctrl.ReportsService.GetYearlyReport(c, userID, year)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Yearly report retrieved successfully",
		Data:    report,
	})
}

// @Tags         Reports
// @Summary      Get CFDIs by date range
// @Description  Get all CFDIs within a specific date range
// @Security     BearerAuth
// @Produce      json
// @Param        start_date  query  string  true  "Start date (YYYY-MM-DD)"
// @Param        end_date    query  string  true  "End date (YYYY-MM-DD)"
// @Router       /reports/date-range [get]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *ReportsController) GetCFDIsByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
	}

	userID := c.Locals("userId").(string)

	report, err := ctrl.ReportsService.GetCFDIsByDateRange(c, userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "CFDIs retrieved successfully",
		Data:    report,
	})
}

// @Tags         Reports
// @Summary      Get expenses summary
// @Description  Get detailed expenses summary by supplier and month
// @Security     BearerAuth
// @Produce      json
// @Param        start_date  query  string  true  "Start date (YYYY-MM-DD)"
// @Param        end_date    query  string  true  "End date (YYYY-MM-DD)"
// @Router       /reports/expenses [get]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *ReportsController) GetExpensesSummary(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
	}

	userID := c.Locals("userId").(string)

	summary, err := ctrl.ReportsService.GetExpensesSummary(c, userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Expenses summary retrieved successfully",
		Data:    summary,
	})
}

// @Tags         Reports
// @Summary      Get income summary
// @Description  Get detailed income summary by client and month
// @Security     BearerAuth
// @Produce      json
// @Param        start_date  query  string  true  "Start date (YYYY-MM-DD)"
// @Param        end_date    query  string  true  "End date (YYYY-MM-DD)"
// @Router       /reports/income [get]
// @Success      200  {object}  response.Success
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
func (ctrl *ReportsController) GetIncomeSummary(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return fiber.NewError(fiber.StatusBadRequest, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid start_date format, use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid end_date format, use YYYY-MM-DD")
	}

	userID := c.Locals("userId").(string)

	summary, err := ctrl.ReportsService.GetIncomeSummary(c, userID, startDate, endDate)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(response.Success{
		Code:    fiber.StatusOK,
		Status:  "success",
		Message: "Income summary retrieved successfully",
		Data:    summary,
	})
}

