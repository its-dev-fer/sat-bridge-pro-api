package service

import (
	"app/src/model"
	"app/src/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ReportsService interface {
	GetMonthlyReport(c *fiber.Ctx, userID string, year int, month int) (*MonthlyReport, error)
	GetYearlyReport(c *fiber.Ctx, userID string, year int) (*YearlyReport, error)
	GetCFDIsByDateRange(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*DateRangeReport, error)
	GetExpensesSummary(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*ExpensesSummary, error)
	GetIncomeSummary(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*IncomeSummary, error)
}

type reportsService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

type MonthlyReport struct {
	UserID           string           `json:"user_id"`
	Year             int              `json:"year"`
	Month            int              `json:"month"`
	TotalIngresos    float64          `json:"total_ingresos"`
	TotalEgresos     float64          `json:"total_egresos"`
	TotalCFDIs       int64            `json:"total_cfdis"`
	CFDIsByType      map[string]int64 `json:"cfdis_by_type"`
	CFDIsByStatus    map[string]int64 `json:"cfdis_by_status"`
	TopProveedores   []ProveedorSummary `json:"top_proveedores"`
	TopClientes      []ClienteSummary   `json:"top_clientes"`
}

type YearlyReport struct {
	UserID           string                 `json:"user_id"`
	Year             int                    `json:"year"`
	TotalIngresos    float64                `json:"total_ingresos"`
	TotalEgresos     float64                `json:"total_egresos"`
	TotalCFDIs       int64                  `json:"total_cfdis"`
	MonthlyBreakdown []MonthlyBreakdown     `json:"monthly_breakdown"`
	CFDIsByType      map[string]int64       `json:"cfdis_by_type"`
}

type DateRangeReport struct {
	UserID        string              `json:"user_id"`
	StartDate     time.Time           `json:"start_date"`
	EndDate       time.Time           `json:"end_date"`
	TotalIngresos float64             `json:"total_ingresos"`
	TotalEgresos  float64             `json:"total_egresos"`
	TotalCFDIs    int64               `json:"total_cfdis"`
	CFDIs         []model.CfdiDescargado `json:"cfdis"`
}

type ExpensesSummary struct {
	UserID         string               `json:"user_id"`
	StartDate      time.Time            `json:"start_date"`
	EndDate        time.Time            `json:"end_date"`
	TotalExpenses  float64              `json:"total_expenses"`
	TotalCFDIs     int64                `json:"total_cfdis"`
	ByProveedor    []ProveedorSummary   `json:"by_proveedor"`
	ByMonth        []MonthlyExpense     `json:"by_month"`
}

type IncomeSummary struct {
	UserID       string             `json:"user_id"`
	StartDate    time.Time          `json:"start_date"`
	EndDate      time.Time          `json:"end_date"`
	TotalIncome  float64            `json:"total_income"`
	TotalCFDIs   int64              `json:"total_cfdis"`
	ByCliente    []ClienteSummary   `json:"by_cliente"`
	ByMonth      []MonthlyIncome    `json:"by_month"`
}

type MonthlyBreakdown struct {
	Month    int     `json:"month"`
	Ingresos float64 `json:"ingresos"`
	Egresos  float64 `json:"egresos"`
	CFDIs    int64   `json:"cfdis"`
}

type ProveedorSummary struct {
	RFC        string  `json:"rfc"`
	Total      float64 `json:"total"`
	NumCFDIs   int64   `json:"num_cfdis"`
}

type ClienteSummary struct {
	RFC        string  `json:"rfc"`
	Total      float64 `json:"total"`
	NumCFDIs   int64   `json:"num_cfdis"`
}

type MonthlyExpense struct {
	Year   int     `json:"year"`
	Month  int     `json:"month"`
	Total  float64 `json:"total"`
	CFDIs  int64   `json:"cfdis"`
}

type MonthlyIncome struct {
	Year   int     `json:"year"`
	Month  int     `json:"month"`
	Total  float64 `json:"total"`
	CFDIs  int64   `json:"cfdis"`
}

func NewReportsService(db *gorm.DB) ReportsService {
	return &reportsService{
		Log: utils.Log,
		DB:  db,
	}
}

func (s *reportsService) GetMonthlyReport(c *fiber.Ctx, userID string, year int, month int) (*MonthlyReport, error) {
	// Get user's RFC
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		s.Log.Errorf("Failed to get user: %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.RFC == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "User has no RFC configured")
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Second)

	report := &MonthlyReport{
		UserID:        userID,
		Year:          year,
		Month:         month,
		CFDIsByType:   make(map[string]int64),
		CFDIsByStatus: make(map[string]int64),
	}

	// Get all CFDIs for the month
	var cfdis []model.CfdiDescargado
	err := s.DB.WithContext(c.Context()).
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND cfdi_descargados.fecha_emision BETWEEN ? AND ?", 
			userID, startDate, endDate).
		Find(&cfdis).Error

	if err != nil {
		s.Log.Errorf("Failed to get CFDIs: %+v", err)
		return nil, err
	}

	report.TotalCFDIs = int64(len(cfdis))

	// Calculate totals and group by type/status
	for _, cfdi := range cfdis {
		if cfdi.RFCEmisor == user.RFC {
			// This is an income (ingreso)
			report.TotalIngresos += cfdi.MontoTotal
			report.CFDIsByType["ingresos"]++
		} else {
			// This is an expense (egreso)
			report.TotalEgresos += cfdi.MontoTotal
			report.CFDIsByType["egresos"]++
		}

		report.CFDIsByStatus[cfdi.StatusSAT]++
	}

	// Get top proveedores (suppliers - where user is receptor)
	type ProveedorResult struct {
		RFC      string
		Total    float64
		NumCFDIs int64
	}
	var proveedores []ProveedorResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("rfc_emisor as rfc, SUM(monto_total) as total, COUNT(*) as num_cfdis").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_receptor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("rfc_emisor").
		Order("total DESC").
		Limit(10).
		Scan(&proveedores)

	for _, p := range proveedores {
		report.TopProveedores = append(report.TopProveedores, ProveedorSummary{
			RFC:      p.RFC,
			Total:    p.Total,
			NumCFDIs: p.NumCFDIs,
		})
	}

	// Get top clientes (customers - where user is emisor)
	type ClienteResult struct {
		RFC      string
		Total    float64
		NumCFDIs int64
	}
	var clientes []ClienteResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("rfc_receptor as rfc, SUM(monto_total) as total, COUNT(*) as num_cfdis").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_emisor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("rfc_receptor").
		Order("total DESC").
		Limit(10).
		Scan(&clientes)

	for _, cl := range clientes {
		report.TopClientes = append(report.TopClientes, ClienteSummary{
			RFC:      cl.RFC,
			Total:    cl.Total,
			NumCFDIs: cl.NumCFDIs,
		})
	}

	return report, nil
}

func (s *reportsService) GetYearlyReport(c *fiber.Ctx, userID string, year int) (*YearlyReport, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.RFC == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "User has no RFC configured")
	}

	report := &YearlyReport{
		UserID:       userID,
		Year:         year,
		CFDIsByType:  make(map[string]int64),
	}

	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	// Get all CFDIs for the year
	var cfdis []model.CfdiDescargado
	err := s.DB.WithContext(c.Context()).
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND cfdi_descargados.fecha_emision BETWEEN ? AND ?", 
			userID, startDate, endDate).
		Find(&cfdis).Error

	if err != nil {
		s.Log.Errorf("Failed to get CFDIs: %+v", err)
		return nil, err
	}

	report.TotalCFDIs = int64(len(cfdis))

	// Monthly breakdown
	monthlyData := make(map[int]*MonthlyBreakdown)
	for i := 1; i <= 12; i++ {
		monthlyData[i] = &MonthlyBreakdown{Month: i}
	}

	for _, cfdi := range cfdis {
		month := int(cfdi.FechaEmision.Month())
		
		if cfdi.RFCEmisor == user.RFC {
			report.TotalIngresos += cfdi.MontoTotal
			monthlyData[month].Ingresos += cfdi.MontoTotal
			report.CFDIsByType["ingresos"]++
		} else {
			report.TotalEgresos += cfdi.MontoTotal
			monthlyData[month].Egresos += cfdi.MontoTotal
			report.CFDIsByType["egresos"]++
		}
		monthlyData[month].CFDIs++
	}

	for i := 1; i <= 12; i++ {
		report.MonthlyBreakdown = append(report.MonthlyBreakdown, *monthlyData[i])
	}

	return report, nil
}

func (s *reportsService) GetCFDIsByDateRange(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*DateRangeReport, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	report := &DateRangeReport{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	var cfdis []model.CfdiDescargado
	err := s.DB.WithContext(c.Context()).
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND cfdi_descargados.fecha_emision BETWEEN ? AND ?", 
			userID, startDate, endDate).
		Find(&cfdis).Error

	if err != nil {
		s.Log.Errorf("Failed to get CFDIs: %+v", err)
		return nil, err
	}

	report.TotalCFDIs = int64(len(cfdis))
	report.CFDIs = cfdis

	for _, cfdi := range cfdis {
		if cfdi.RFCEmisor == user.RFC {
			report.TotalIngresos += cfdi.MontoTotal
		} else {
			report.TotalEgresos += cfdi.MontoTotal
		}
	}

	return report, nil
}

func (s *reportsService) GetExpensesSummary(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*ExpensesSummary, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.RFC == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "User has no RFC configured")
	}

	summary := &ExpensesSummary{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get total expenses
	var result struct {
		Total float64
		Count int64
	}
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("SUM(monto_total) as total, COUNT(*) as count").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_receptor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Scan(&result)

	summary.TotalExpenses = result.Total
	summary.TotalCFDIs = result.Count

	// Group by proveedor
	type ProveedorResult struct {
		RFC      string
		Total    float64
		NumCFDIs int64
	}
	var proveedores []ProveedorResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("rfc_emisor as rfc, SUM(monto_total) as total, COUNT(*) as num_cfdis").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_receptor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("rfc_emisor").
		Order("total DESC").
		Scan(&proveedores)

	for _, p := range proveedores {
		summary.ByProveedor = append(summary.ByProveedor, ProveedorSummary{
			RFC:      p.RFC,
			Total:    p.Total,
			NumCFDIs: p.NumCFDIs,
		})
	}

	// Group by month
	type MonthResult struct {
		Year  int
		Month int
		Total float64
		Count int64
	}
	var months []MonthResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("EXTRACT(YEAR FROM fecha_emision)::int as year, EXTRACT(MONTH FROM fecha_emision)::int as month, SUM(monto_total) as total, COUNT(*) as count").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_receptor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("year, month").
		Order("year, month").
		Scan(&months)

	for _, m := range months {
		summary.ByMonth = append(summary.ByMonth, MonthlyExpense{
			Year:  m.Year,
			Month: m.Month,
			Total: m.Total,
			CFDIs: m.Count,
		})
	}

	return summary, nil
}

func (s *reportsService) GetIncomeSummary(c *fiber.Ctx, userID string, startDate, endDate time.Time) (*IncomeSummary, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		return nil, fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	if user.RFC == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "User has no RFC configured")
	}

	summary := &IncomeSummary{
		UserID:    userID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Get total income
	var result struct {
		Total float64
		Count int64
	}
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("SUM(monto_total) as total, COUNT(*) as count").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_emisor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Scan(&result)

	summary.TotalIncome = result.Total
	summary.TotalCFDIs = result.Count

	// Group by cliente
	type ClienteResult struct {
		RFC      string
		Total    float64
		NumCFDIs int64
	}
	var clientes []ClienteResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("rfc_receptor as rfc, SUM(monto_total) as total, COUNT(*) as num_cfdis").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_emisor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("rfc_receptor").
		Order("total DESC").
		Scan(&clientes)

	for _, cl := range clientes {
		summary.ByCliente = append(summary.ByCliente, ClienteSummary{
			RFC:      cl.RFC,
			Total:    cl.Total,
			NumCFDIs: cl.NumCFDIs,
		})
	}

	// Group by month
	type MonthResult struct {
		Year  int
		Month int
		Total float64
		Count int64
	}
	var months []MonthResult
	s.DB.WithContext(c.Context()).
		Model(&model.CfdiDescargado{}).
		Select("EXTRACT(YEAR FROM fecha_emision)::int as year, EXTRACT(MONTH FROM fecha_emision)::int as month, SUM(monto_total) as total, COUNT(*) as count").
		Joins("JOIN solicitudes_descarga ON cfdi_descargados.solicitud_id = solicitudes_descarga.id").
		Where("solicitudes_descarga.usuario_id = ? AND rfc_emisor = ? AND fecha_emision BETWEEN ? AND ?", 
			userID, user.RFC, startDate, endDate).
		Group("year, month").
		Order("year, month").
		Scan(&months)

	for _, m := range months {
		summary.ByMonth = append(summary.ByMonth, MonthlyIncome{
			Year:  m.Year,
			Month: m.Month,
			Total: m.Total,
			CFDIs: m.Count,
		})
	}

	return summary, nil
}

