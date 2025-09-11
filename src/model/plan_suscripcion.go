package model

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type PlanSuscripcion struct {
    ID                      uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
    Nombre                  string         `gorm:"type:varchar(100);not null" json:"nombre"`
    Descripcion             string         `gorm:"type:varchar(255)" json:"descripcion"`
    LimiteDescargasMensuales int           `gorm:"not null" json:"limite_descargas_mensuales"`
    Precio                  float64        `gorm:"type:decimal(10,2);not null" json:"precio"`
    Activo                  bool           `gorm:"default:true" json:"activo"`
    CreatedAt               time.Time      `json:"created_at"`
    UpdatedAt               time.Time      `json:"updated_at"`
    DeletedAt               gorm.DeletedAt `gorm:"index" json:"-"`
}
