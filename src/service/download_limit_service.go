package service

import (
	"app/src/model"
	"app/src/utils"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DownloadLimitService interface {
	CanDownload(c *fiber.Ctx, userID string) (bool, string, error)
	IncrementDownloadCount(c *fiber.Ctx, userID string) error
	GetDownloadStats(c *fiber.Ctx, userID string) (*DownloadStats, error)
	ResetMonthlyDownloads(c *fiber.Ctx, userID string) error
}

type downloadLimitService struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

type DownloadStats struct {
	UsuarioID              uuid.UUID `json:"usuario_id"`
	NombrePlan             string    `json:"nombre_plan"`
	LimiteDescargasMensuales int     `json:"limite_descargas_mensuales"`
	DescargasRealizadas    int       `json:"descargas_realizadas"`
	DescargasRestantes     int       `json:"descargas_restantes"`
	UltimoResetDescargas   time.Time `json:"ultimo_reset_descargas"`
	PlanActivo             bool      `json:"plan_activo"`
	SuscripcionStatus      string    `json:"suscripcion_status"`
}

func NewDownloadLimitService(db *gorm.DB) DownloadLimitService {
	return &downloadLimitService{
		Log: utils.Log,
		DB:  db,
	}
}

// CanDownload verifica si el usuario puede realizar una descarga según su plan
func (s *downloadLimitService) CanDownload(c *fiber.Ctx, userID string) (bool, string, error) {
	// Get user with download count
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		s.Log.Errorf("Failed to get user: %+v", err)
		return false, "User not found", err
	}

	// Check if we need to reset monthly downloads
	now := time.Now()
	lastReset := user.UltimoResetDescargas
	
	// Reset if it's a new month or first time
	if lastReset.IsZero() || (now.Year() > lastReset.Year()) || 
	   (now.Year() == lastReset.Year() && now.Month() > lastReset.Month()) {
		if err := s.ResetMonthlyDownloads(c, userID); err != nil {
			s.Log.Errorf("Failed to reset monthly downloads: %+v", err)
		} else {
			user.DescargasRealizadas = 0
		}
	}

	// Get active subscription
	var suscripcion model.SuscripcionUsuario
	err := s.DB.WithContext(c.Context()).
		Preload("Plan").
		Where("usuario_id = ? AND status = ? AND fecha_inicio <= ? AND fecha_fin >= ?", 
			userID, "activo", now, now).
		Order("created_at DESC").
		First(&suscripcion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, "No active subscription found", nil
		}
		s.Log.Errorf("Failed to get user subscription: %+v", err)
		return false, "Error checking subscription", err
	}

	if suscripcion.Plan == nil {
		return false, "Invalid subscription plan", nil
	}

	// Check plan limits
	plan := suscripcion.Plan

	// Empresarial plan with unlimited downloads (if more than 5 users, for example)
	if plan.Nombre == "Empresarial" && plan.LimiteDescargasMensuales == -1 {
		return true, "Unlimited downloads", nil
	}

	// Check if user has reached the limit
	if user.DescargasRealizadas >= plan.LimiteDescargasMensuales {
		return false, "Monthly download limit reached", nil
	}

	return true, "Download allowed", nil
}

// IncrementDownloadCount incrementa el contador de descargas del usuario
func (s *downloadLimitService) IncrementDownloadCount(c *fiber.Ctx, userID string) error {
	result := s.DB.WithContext(c.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("descargas_realizadas", gorm.Expr("descargas_realizadas + ?", 1))

	if result.Error != nil {
		s.Log.Errorf("Failed to increment download count: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}

// GetDownloadStats obtiene las estadísticas de descarga del usuario
func (s *downloadLimitService) GetDownloadStats(c *fiber.Ctx, userID string) (*DownloadStats, error) {
	var user model.User
	if err := s.DB.WithContext(c.Context()).First(&user, "id = ?", userID).Error; err != nil {
		s.Log.Errorf("Failed to get user: %+v", err)
		return nil, err
	}

	// Get active subscription
	var suscripcion model.SuscripcionUsuario
	now := time.Now()
	err := s.DB.WithContext(c.Context()).
		Preload("Plan").
		Where("usuario_id = ? AND status = ? AND fecha_inicio <= ? AND fecha_fin >= ?", 
			userID, "activo", now, now).
		Order("created_at DESC").
		First(&suscripcion).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &DownloadStats{
				UsuarioID:              user.ID,
				NombrePlan:             "Sin plan",
				LimiteDescargasMensuales: 0,
				DescargasRealizadas:    user.DescargasRealizadas,
				DescargasRestantes:     0,
				UltimoResetDescargas:   user.UltimoResetDescargas,
				PlanActivo:             false,
				SuscripcionStatus:      "inactive",
			}, nil
		}
		s.Log.Errorf("Failed to get user subscription: %+v", err)
		return nil, err
	}

	plan := suscripcion.Plan
	descargasRestantes := plan.LimiteDescargasMensuales - user.DescargasRealizadas
	if descargasRestantes < 0 {
		descargasRestantes = 0
	}

	// Handle unlimited downloads
	if plan.LimiteDescargasMensuales == -1 {
		descargasRestantes = -1 // Indicates unlimited
	}

	stats := &DownloadStats{
		UsuarioID:              user.ID,
		NombrePlan:             plan.Nombre,
		LimiteDescargasMensuales: plan.LimiteDescargasMensuales,
		DescargasRealizadas:    user.DescargasRealizadas,
		DescargasRestantes:     descargasRestantes,
		UltimoResetDescargas:   user.UltimoResetDescargas,
		PlanActivo:             plan.Activo,
		SuscripcionStatus:      suscripcion.Status,
	}

	return stats, nil
}

// ResetMonthlyDownloads reinicia el contador mensual de descargas
func (s *downloadLimitService) ResetMonthlyDownloads(c *fiber.Ctx, userID string) error {
	result := s.DB.WithContext(c.Context()).
		Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"descargas_realizadas":   0,
			"ultimo_reset_descargas": time.Now(),
		})

	if result.Error != nil {
		s.Log.Errorf("Failed to reset monthly downloads: %+v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "User not found")
	}

	return nil
}

