package validation

import (
    "app/src/model"
    "errors"
    "time"

    "github.com/google/uuid"
)

type CreateSolicitudDescarga struct {
    UsuarioID      string    `json:"usuario_id" validate:"required,uuid4"`
    TipoCFDI       string    `json:"tipo_cfdi" validate:"required,oneof=ingresos egresos nomina pagos todos"`
    RFCSolicitante string    `json:"rfc_solicitante" validate:"required,min=12,max=13"`
    FechaSolicitud time.Time `json:"fecha_solicitud" validate:"required"`
    ResultadosSAT  string    `json:"resultados_sat_json" validate:"omitempty"`
    Status         string    `json:"status" validate:"required,oneof=activo cancelado pendiente"`
}

func (v *CreateSolicitudDescarga) ToModel() *model.SolicitudDescarga {
    return &model.SolicitudDescarga{
        ID:                uuid.New().String(),
        UsuarioID:         v.UsuarioID,
        TipoCFDI:          v.TipoCFDI,
        RFCSolicitante:    v.RFCSolicitante,
        FechaSolicitud:    v.FechaSolicitud,
        ResultadosSATJSON: v.ResultadosSAT,
        Status:            v.Status,
    }
}

type UpdateSolicitudDescarga struct {
    UsuarioID      string    `json:"usuario_id" validate:"omitempty,uuid4"`
    TipoCFDI       string    `json:"tipo_cfdi" validate:"omitempty,oneof=ingresos egresos nomina pagos todos"`
    RFCSolicitante string    `json:"rfc_solicitante" validate:"omitempty,min=12,max=13"`
    FechaSolicitud time.Time `json:"fecha_solicitud" validate:"omitempty"`
    ResultadosSAT  string    `json:"resultados_sat_json" validate:"omitempty"`
    Status         string    `json:"status" validate:"omitempty,oneof=activo cancelado pendiente"`
}

func (v *UpdateSolicitudDescarga) ToModel() *model.SolicitudDescarga {
    return &model.SolicitudDescarga{
        UsuarioID:         v.UsuarioID,
        TipoCFDI:          v.TipoCFDI,
        RFCSolicitante:    v.RFCSolicitante,
        FechaSolicitud:    v.FechaSolicitud,
        ResultadosSATJSON: v.ResultadosSAT,
        Status:            v.Status,
    }
}

type QuerySolicitudDescarga struct {
    Page   int    `query:"page" validate:"min=1"`
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Search string `query:"search"`
}

func ValidateCreateSolicitudDescarga(v *CreateSolicitudDescarga) error {
    if v.FechaSolicitud.After(time.Now()) {
        return errors.New("fecha_solicitud no puede ser futura")
    }
    return nil
}

func ValidateUpdateSolicitudDescarga(v *UpdateSolicitudDescarga) error {
    if !v.FechaSolicitud.IsZero() && v.FechaSolicitud.After(time.Now()) {
        return errors.New("fecha_solicitud no puede ser futura")
    }
    return nil
}
