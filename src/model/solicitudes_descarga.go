package model

import (
    "time"

    "gorm.io/gorm"
)

type SolicitudDescarga struct {
    ID                string         `json:"id" gorm:"type:uuid;primaryKey"`
    UsuarioID         string         `json:"usuario_id" gorm:"type:uuid;not null"`
    TipoCFDI          string         `json:"tipo_cfdi" gorm:"type:text;check:tipo_cfdi IN ('ingresos','egresos','nomina','pagos','todos');not null"`
    RFCSolicitante    string         `json:"rfc_solicitante" gorm:"type:varchar(13);not null"`
    FechaSolicitud    time.Time      `json:"fecha_solicitud" gorm:"type:date;not null"`
    ResultadosSATJSON string         `json:"resultados_sat_json" gorm:"type:text"`
    Status            string         `json:"status" gorm:"type:text;check:status IN ('activo','cancelado','pendiente');default:'pendiente';not null"`
    CreatedAt         time.Time      `json:"created_at"`
    UpdatedAt         time.Time      `json:"updated_at"`
    DeletedAt         gorm.DeletedAt `json:"-" gorm:"index"`

    Usuario *User `json:"usuario,omitempty" gorm:"foreignKey:UsuarioID"`
}
