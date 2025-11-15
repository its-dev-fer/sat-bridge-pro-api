package model

import (
    "time"

    "gorm.io/gorm"
)

type CfdiDescargado struct {
    ID            string         `json:"id" gorm:"type:uuid;primaryKey"`
    SolicitudID   string         `json:"solicitud_id" gorm:"type:uuid;not null"`
    TipoCFDI      string         `json:"tipo_cfdi" gorm:"type:text;check:tipo_cfdi IN ('ingresos','egresos','nomina','pagos');not null"`
    CfdiUUID      string         `json:"cfdi_uuid" gorm:"type:uuid;not null"`
    RFCEmisor     string         `json:"rfc_emisor" gorm:"type:varchar(13);not null"`
    RFCReceptor   string         `json:"rfc_receptor" gorm:"type:varchar(13);not null"`
    StatusSAT     string         `json:"status_sat" gorm:"type:text;check:status_sat IN ('activo','cancelado');not null"`
    FechaEmision  time.Time      `json:"fecha_emision" gorm:"type:date;not null"`
    MontoTotal    float64        `json:"monto_total" gorm:"type:decimal(18,2);not null"`
    ArchivoXML    string         `json:"archivo_xml" gorm:"type:text;not null"`
    CreatedAt     time.Time      `json:"created_at"`
    UpdatedAt     time.Time      `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

    Solicitud *SolicitudDescarga `json:"solicitud,omitempty" gorm:"foreignKey:SolicitudID"`
}
