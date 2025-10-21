package model

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DatosFiscalesSAT struct {
	UUID                  uuid.UUID      `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID                uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	User                  User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	RFC                   string         `json:"rfc" gorm:"type:varchar(13);not null"`
	CerB64Encriptado      string         `json:"-" gorm:"type:text;not null"` // No exponer en JSON
	KeyB64Encriptado      string         `json:"-" gorm:"type:text;not null"` // No exponer en JSON  
	PasswordEfirmaEncrip  string         `json:"-" gorm:"type:varchar(255);not null"` // No exponer en JSON
	CreatedAt             time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt             time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	CreatedBy             *uuid.UUID     `json:"created_by,omitempty" gorm:"type:uuid"`
	UpdatedBy             *uuid.UUID     `json:"updated_by,omitempty" gorm:"type:uuid"`
	DeletedAt             gorm.DeletedAt `json:"-" gorm:"index"`
}

func (DatosFiscalesSAT) TableName() string {
	return "datos_fiscales_sat"
}