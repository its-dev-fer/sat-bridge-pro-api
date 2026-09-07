package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CFDI struct {
	UUID                uuid.UUID      `json:"uuid" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID              uuid.UUID      `json:"user_id" gorm:"type:uuid;not null;uniqueIndex:uq_cfdis_user_folio"`
	User                User           `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	FolioFiscal         string         `json:"folio_fiscal" gorm:"type:varchar(36);not null;uniqueIndex:uq_cfdis_user_folio"`
	RFCEmisor           string         `json:"rfc_emisor" gorm:"column:rfc_emisor;type:varchar(13);not null"`
	NombreEmisor        string         `json:"nombre_emisor" gorm:"type:varchar(255);not null"`
	RFCReceptor         string         `json:"rfc_receptor" gorm:"column:rfc_receptor;type:varchar(13);not null"`
	NombreReceptor      string         `json:"nombre_receptor" gorm:"type:varchar(255);not null"`
	FechaEmision        *time.Time     `json:"fecha_emision"`
	FechaCertificacion  *time.Time     `json:"fecha_certificacion"`
	PACCertifico        string         `json:"pac_certifico" gorm:"column:pac_certifico;type:varchar(255)"`
	Total               *float64       `json:"total" gorm:"type:numeric(18,6)"`
	EfectoComprobante   string         `json:"efecto_comprobante" gorm:"type:varchar(20)"`
	EstatusCancelacion  string         `json:"estatus_cancelacion" gorm:"type:varchar(50)"`
	EstadoComprobante   string         `json:"estado_comprobante" gorm:"type:varchar(50)"`
	Entra               *bool          `json:"entra"`
	Comentarios         string         `json:"comentarios" gorm:"type:text"`
	CreatedAt           time.Time      `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time      `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
	DeletedAt           gorm.DeletedAt `json:"-" gorm:"index"`
}

func (CFDI) TableName() string {
	return "cfdis"
}
