package validation

import "time"

// CreateFirmaElectronica validation for creating a new FIEL
type CreateFirmaElectronica struct {
	RFC                 string    `json:"rfc" validate:"required,min=12,max=13"`
	CertificadoCER      string    `json:"certificado_cer" validate:"required,base64"` // Base64 encoded .cer file
	ClavePrivadaKEY     string    `json:"clave_privada_key" validate:"required,base64"` // Base64 encoded .key file
	PasswordKEY         string    `json:"password_key" validate:"required,min=8"`
	NombreCertificado   string    `json:"nombre_certificado" validate:"omitempty,max=255"`
	FechaVigenciaInicio time.Time `json:"fecha_vigencia_inicio" validate:"required"`
	FechaVigenciaFin    time.Time `json:"fecha_vigencia_fin" validate:"required,gtfield=FechaVigenciaInicio"`
}

// UpdateFirmaElectronica validation for updating FIEL
type UpdateFirmaElectronica struct {
	CertificadoCER      string    `json:"certificado_cer" validate:"omitempty,base64"`
	ClavePrivadaKEY     string    `json:"clave_privada_key" validate:"omitempty,base64"`
	PasswordKEY         string    `json:"password_key" validate:"omitempty,min=8"`
	NombreCertificado   string    `json:"nombre_certificado" validate:"omitempty,max=255"`
	FechaVigenciaInicio time.Time `json:"fecha_vigencia_inicio" validate:"omitempty"`
	FechaVigenciaFin    time.Time `json:"fecha_vigencia_fin" validate:"omitempty"`
	Activo              *bool     `json:"activo" validate:"omitempty"`
}

// CreateCIEC validation for CIEC authentication
type CreateCIEC struct {
	RFC  string `json:"rfc" validate:"required,min=12,max=13"`
	CIEC string `json:"ciec" validate:"required,min=8"`
}

// UpdateCIEC validation for updating CIEC
type UpdateCIEC struct {
	CIEC string `json:"ciec" validate:"required,min=8"`
}

// QueryFirmaElectronica validation for querying FIEL records
type QueryFirmaElectronica struct {
	Page   int    `json:"page" validate:"required,min=1"`
	Limit  int    `json:"limit" validate:"required,min=1,max=100"`
	RFC    string `json:"rfc" validate:"omitempty,min=12,max=13"`
	Activo *bool  `json:"activo" validate:"omitempty"`
}

