package model

import (
    "time"

    "gorm.io/gorm"
)

type SuscripcionUsuario struct {
    ID          string         `json:"id" gorm:"type:uuid;primaryKey"`
    UsuarioID   string         `json:"usuario_id" gorm:"type:uuid;not null"`
    PlanID      string         `json:"plan_id" gorm:"type:uuid;not null"`
    FechaInicio time.Time      `json:"fecha_inicio" gorm:"type:date;not null"`
    FechaFin    time.Time      `json:"fecha_fin" gorm:"type:date;not null"`
    Status      string         `json:"status" gorm:"type:text;check:status IN ('activo','cancelado','pendiente');default:'pendiente';not null"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

    Usuario *User             `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
    Plan    *PlanSuscripcion  `json:"plan,omitempty" gorm:"foreignKey:PlanID"`
}
