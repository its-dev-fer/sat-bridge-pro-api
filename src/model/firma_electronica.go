package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// FirmaElectronica representa la FIEL de un usuario para autenticación con el SAT
type FirmaElectronica struct {
	ID              uuid.UUID      `gorm:"type:uuid;primaryKey" json:"id"`
	UsuarioID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"usuario_id"`
	RFC             string         `gorm:"type:varchar(13);not null" json:"rfc"`
	CertificadoCER  string         `gorm:"type:text;not null" json:"-"` // Base64 encoded .cer file
	ClavePrivadaKEY string         `gorm:"type:text;not null" json:"-"` // Base64 encoded .key file
	PasswordKEY     string         `gorm:"type:text;not null" json:"-"` // Encrypted password for .key file
	NombreCertificado string       `gorm:"type:varchar(255)" json:"nombre_certificado"`
	FechaVigenciaInicio time.Time  `gorm:"type:date" json:"fecha_vigencia_inicio"`
	FechaVigenciaFin    time.Time  `gorm:"type:date" json:"fecha_vigencia_fin"`
	Activo          bool           `gorm:"default:true" json:"activo"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	Usuario *User `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}

func (f *FirmaElectronica) BeforeCreate(_ *gorm.DB) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	return nil
}

// IsValid verifica si la FIEL está vigente
func (f *FirmaElectronica) IsValid() bool {
	now := time.Now()
	return f.Activo && 
		now.After(f.FechaVigenciaInicio) && 
		now.Before(f.FechaVigenciaFin)
}

