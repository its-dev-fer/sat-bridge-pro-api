package validation

import (
    "app/src/model"
    "errors"
    "time"

    "github.com/google/uuid"
)

type CreateCfdiDescargado struct {
    SolicitudID  string    `json:"solicitud_id" validate:"required,uuid4"`
    TipoCFDI     string    `json:"tipo_cfdi" validate:"required,oneof=ingresos egresos nomina pagos"`
    CfdiUUID     string    `json:"cfdi_uuid" validate:"required,uuid4"`
    RFCEmisor    string    `json:"rfc_emisor" validate:"required,min=12,max=13"`
    RFCReceptor  string    `json:"rfc_receptor" validate:"required,min=12,max=13"`
    StatusSAT    string    `json:"status_sat" validate:"required,oneof=activo cancelado"`
    FechaEmision time.Time `json:"fecha_emision" validate:"required"`
    MontoTotal   float64   `json:"monto_total" validate:"required,gt=0"`
    ArchivoXML   string    `json:"archivo_xml" validate:"required"`
}

func (v *CreateCfdiDescargado) ToModel() *model.CfdiDescargado {
    return &model.CfdiDescargado{
        ID:           uuid.New().String(),
        SolicitudID:  v.SolicitudID,
        TipoCFDI:     v.TipoCFDI,
        CfdiUUID:     v.CfdiUUID,
        RFCEmisor:    v.RFCEmisor,
        RFCReceptor:  v.RFCReceptor,
        StatusSAT:    v.StatusSAT,
        FechaEmision: v.FechaEmision,
        MontoTotal:   v.MontoTotal,
        ArchivoXML:   v.ArchivoXML,
    }
}

type UpdateCfdiDescargado struct {
    SolicitudID  string    `json:"solicitud_id" validate:"omitempty,uuid4"`
    TipoCFDI     string    `json:"tipo_cfdi" validate:"omitempty,oneof=ingresos egresos nomina pagos"`
    CfdiUUID     string    `json:"cfdi_uuid" validate:"omitempty,uuid4"`
    RFCEmisor    string    `json:"rfc_emisor" validate:"omitempty,min=12,max=13"`
    RFCReceptor  string    `json:"rfc_receptor" validate:"omitempty,min=12,max=13"`
    StatusSAT    string    `json:"status_sat" validate:"omitempty,oneof=activo cancelado"`
    FechaEmision time.Time `json:"fecha_emision" validate:"omitempty"`
    MontoTotal   float64   `json:"monto_total" validate:"omitempty,gt=0"`
    ArchivoXML   string    `json:"archivo_xml" validate:"omitempty"`
}

func (v *UpdateCfdiDescargado) ToModel() *model.CfdiDescargado {
    return &model.CfdiDescargado{
        SolicitudID:  v.SolicitudID,
        TipoCFDI:     v.TipoCFDI,
        CfdiUUID:     v.CfdiUUID,
        RFCEmisor:    v.RFCEmisor,
        RFCReceptor:  v.RFCReceptor,
        StatusSAT:    v.StatusSAT,
        FechaEmision: v.FechaEmision,
        MontoTotal:   v.MontoTotal,
        ArchivoXML:   v.ArchivoXML,
    }
}

type QueryCfdiDescargado struct {
    Page   int    `query:"page" validate:"min=1"`
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Search string `query:"search"`
}

func ValidateCreateCfdiDescargado(v *CreateCfdiDescargado) error {
    if v.FechaEmision.After(time.Now()) {
        return errors.New("fecha_emision no puede ser futura")
    }
    return nil
}

func ValidateUpdateCfdiDescargado(v *UpdateCfdiDescargado) error {
    if !v.FechaEmision.IsZero() && v.FechaEmision.After(time.Now()) {
        return errors.New("fecha_emision no puede ser futura")
    }
    return nil
}
