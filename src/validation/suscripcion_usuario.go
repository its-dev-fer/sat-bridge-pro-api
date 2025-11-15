package validation

import (
    "app/src/model"
    "time"
	"errors"
	"github.com/google/uuid"

)

type CreateSuscripcionUsuario struct {
    UsuarioID   string    `json:"usuario_id" validate:"required,uuid4"`
    PlanID      string    `json:"plan_id" validate:"required,uuid4"`
    FechaInicio time.Time `json:"fecha_inicio" validate:"required"`
    FechaFin    time.Time `json:"fecha_fin" validate:"required,gtfield=FechaInicio"`
    Status      string    `json:"status" validate:"required,oneof=activo cancelado pendiente"`
}

func (v *CreateSuscripcionUsuario) ToModel() *model.SuscripcionUsuario {
    return &model.SuscripcionUsuario{
        ID:          uuid.New().String(),
        UsuarioID:   v.UsuarioID,
        PlanID:      v.PlanID,
        FechaInicio: v.FechaInicio,
        FechaFin:    v.FechaFin,
        Status:      v.Status,
    }
}

type UpdateSuscripcionUsuario struct {
    UsuarioID   string    `json:"usuario_id" validate:"omitempty,uuid4"`
    PlanID      string    `json:"plan_id" validate:"omitempty,uuid4"`
    FechaInicio time.Time `json:"fecha_inicio" validate:"omitempty"`
    FechaFin    time.Time `json:"fecha_fin" validate:"omitempty,gtfield=FechaInicio"`
    Status      string    `json:"status" validate:"omitempty,oneof=activo cancelado pendiente"`
}

func (v *UpdateSuscripcionUsuario) ToModel() *model.SuscripcionUsuario {
    return &model.SuscripcionUsuario{
        UsuarioID:   v.UsuarioID,
        PlanID:      v.PlanID,
        FechaInicio: v.FechaInicio,
        FechaFin:    v.FechaFin,
        Status:      v.Status,
    }
}

type QuerySuscripcionUsuario struct {
    Page   int    `query:"page" validate:"min=1"`
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Search string `query:"search"`
}

func ValidateCreateSuscripcionUsuario(v *CreateSuscripcionUsuario) error {
    if v.FechaFin.Before(v.FechaInicio) {
        return errors.New("fecha_fin debe ser posterior a fecha_inicio")
    }
    return nil
}

func ValidateUpdateSuscripcionUsuario(v *UpdateSuscripcionUsuario) error {
    if !v.FechaFin.IsZero() && !v.FechaInicio.IsZero() && v.FechaFin.Before(v.FechaInicio) {
        return errors.New("fecha_fin debe ser posterior a fecha_inicio")
    }
    return nil
}
