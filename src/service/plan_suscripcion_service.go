package service

import (
	"app/src/model"
	"app/src/utils"
	"app/src/validation"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PlanSuscripcionService interface {
	GetPlanes(c *fiber.Ctx, params *validation.QueryPlanSuscripcion) ([]model.PlanSuscripcion, int64, error)
	GetPlanByID(c *fiber.Ctx, id string) (*model.PlanSuscripcion, error)
	CreatePlan(c *fiber.Ctx, req *validation.CreatePlanSuscripcion) (*model.PlanSuscripcion, error)
	UpdatePlan(c *fiber.Ctx, req *validation.UpdatePlanSuscripcion, id string) (*model.PlanSuscripcion, error)
	DeletePlan(c *fiber.Ctx, id string) error
}

type planSuscripcionService struct {
	Log      *logrus.Logger
	DB       *gorm.DB
	Validate *validator.Validate
}

func NewPlanSuscripcionService(db *gorm.DB, validate *validator.Validate) PlanSuscripcionService {
	return &planSuscripcionService{
		Log:      utils.Log,
		DB:       db,
		Validate: validate,
	}
}

func (s *planSuscripcionService) GetPlanes(c *fiber.Ctx, params *validation.QueryPlanSuscripcion) ([]model.PlanSuscripcion, int64, error) {
	var planes []model.PlanSuscripcion
	var totalResults int64

	if err := s.Validate.Struct(params); err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.Limit
	query := s.DB.WithContext(c.Context()).Order("created_at asc")

	if search := params.Search; search != "" {
		query = query.Where("nombre LIKE ? OR descripcion LIKE ?", 
			"%"+search+"%", "%"+search+"%")
	}

	// Contar total de resultados
	result := query.Model(&model.PlanSuscripcion{}).Count(&totalResults)
	if result.Error != nil {
		s.Log.Errorf("Failed to count planes: %+v", result.Error)
		return nil, 0, result.Error
	}

	// Obtener resultados paginados
	result = query.Limit(params.Limit).Offset(offset).Find(&planes)
	if result.Error != nil {
		s.Log.Errorf("Failed to get planes: %+v", result.Error)
		return nil, 0, result.Error
	}

	return planes, totalResults, nil
}

func (s *planSuscripcionService) GetPlanByID(c *fiber.Ctx, id string) (*model.PlanSuscripcion, error) {
	plan := new(model.PlanSuscripcion)

	// Convertir string ID a UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	result := s.DB.WithContext(c.Context()).First(plan, "id = ?", uuidID)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, fiber.NewError(fiber.StatusNotFound, "Plan not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed get plan by id: %+v", result.Error)
		return nil, result.Error
	}

	return plan, nil
}

func (s *planSuscripcionService) CreatePlan(c *fiber.Ctx, req *validation.CreatePlanSuscripcion) (*model.PlanSuscripcion, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	plan := &model.PlanSuscripcion{
		ID:                      uuid.New(),
		Nombre:                  req.Nombre,
		Descripcion:             req.Descripcion,
		LimiteDescargasMensuales: req.LimiteDescargasMensuales,
		Precio:                  req.Precio,
		Activo:                  req.Activo,
	}

	result := s.DB.WithContext(c.Context()).Create(plan)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Plan name already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to create plan: %+v", result.Error)
		return nil, result.Error
	}

	return plan, nil
}

func (s *planSuscripcionService) UpdatePlan(c *fiber.Ctx, req *validation.UpdatePlanSuscripcion, id string) (*model.PlanSuscripcion, error) {
	if err := s.Validate.Struct(req); err != nil {
		return nil, err
	}

	// Verificar que al menos un campo sea proporcionado
	if req.Nombre == "" && req.Descripcion == "" && req.LimiteDescargasMensuales == 0 && req.Precio == 0 && req.Activo == nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "No fields to update")
	}

	// Convertir string ID a UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	// Preparar updates
	updates := make(map[string]interface{})
	if req.Nombre != "" {
		updates["nombre"] = req.Nombre
	}
	if req.Descripcion != "" {
		updates["descripcion"] = req.Descripcion
	}
	if req.LimiteDescargasMensuales > 0 {
		updates["limite_descargas_mensuales"] = req.LimiteDescargasMensuales
	}
	if req.Precio > 0 {
		updates["precio"] = req.Precio
	}
	if req.Activo != nil {
		updates["activo"] = *req.Activo
	}

	result := s.DB.WithContext(c.Context()).
		Model(&model.PlanSuscripcion{}).
		Where("id = ?", uuidID).
		Updates(updates)

	if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
		return nil, fiber.NewError(fiber.StatusConflict, "Plan name already exists")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to update plan: %+v", result.Error)
		return nil, result.Error
	}

	// Obtener el plan actualizado
	updatedPlan, err := s.GetPlanByID(c, id)
	if err != nil {
		return nil, err
	}

	return updatedPlan, nil
}

func (s *planSuscripcionService) DeletePlan(c *fiber.Ctx, id string) error {
	// Convertir string ID a UUID
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID format")
	}

	result := s.DB.WithContext(c.Context()).Delete(&model.PlanSuscripcion{}, "id = ?", uuidID)

	if result.RowsAffected == 0 {
		return fiber.NewError(fiber.StatusNotFound, "Plan not found")
	}

	if result.Error != nil {
		s.Log.Errorf("Failed to delete plan: %+v", result.Error)
		return result.Error
	}

	return nil
}