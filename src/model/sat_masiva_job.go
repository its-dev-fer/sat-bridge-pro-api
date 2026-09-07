package model

import (
	"time"

	"github.com/google/uuid"
)

type SatMasivaJob struct {
	UUID       uuid.UUID `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Year       int       `json:"year" gorm:"not null"`
	Status     string    `json:"status" gorm:"type:varchar(20);not null"`
	Month      string    `json:"month"`
	Message    string    `json:"message"`
	Error      string    `json:"error"`
	Logs       []string  `json:"logs" gorm:"type:jsonb;serializer:json"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FinishedAt *time.Time `json:"finished_at"`
}

func (SatMasivaJob) TableName() string {
	return "sat_masiva_jobs"
}
