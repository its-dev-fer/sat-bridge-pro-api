package validation

type SatServiceRequest struct {
	Service string `json:"service" validate:"omitempty,oneof=cfdi retenciones"`
}

type SatQueryRequest struct {
	Service            string   `json:"service" validate:"omitempty,oneof=cfdi retenciones"`
	From               string   `json:"from" validate:"required"`
	To                 string   `json:"to" validate:"required"`
	Download           string   `json:"download" validate:"required,oneof=received issued recibidas emitidas"`
	Request            string   `json:"request" validate:"required,oneof=cfdi metadata"`
	RfcEmisor          string   `json:"rfc_emisor" validate:"omitempty,min=12,max=13"`
	RfcReceptores      []string `json:"rfc_receptores" validate:"omitempty,max=5,dive,min=12,max=13"`
	TipoComprobante    string   `json:"tipo_comprobante" validate:"omitempty,oneof=I E T N P"`
	EstadoComprobante  string   `json:"estado_comprobante" validate:"omitempty,oneof=0 1"`
	RfcACuentaTerceros string   `json:"rfc_a_cuenta_terceros" validate:"omitempty,min=12,max=13"`
	Complemento        string   `json:"complemento"`
}

type SatVerifyRequest struct {
	Service   string `json:"service" validate:"omitempty,oneof=cfdi retenciones"`
	RequestID string `json:"request_id" validate:"required"`
}

type SatDownloadRequest struct {
	Service   string `json:"service" validate:"omitempty,oneof=cfdi retenciones"`
	PackageID string `json:"package_id" validate:"required"`
	Parse     string `json:"parse" validate:"omitempty,oneof=cfdi metadata"`
}

type SatBackfillRequest struct {
	Service           string `json:"service" validate:"omitempty,oneof=cfdi retenciones"`
	From              string `json:"from" validate:"required"`
	To                string `json:"to" validate:"required"`
	Download          string `json:"download" validate:"required,oneof=received issued recibidas emitidas"`
	Request           string `json:"request" validate:"required,oneof=cfdi metadata"`
	RfcEmisor         string `json:"rfc_emisor" validate:"omitempty,min=12,max=13"`
	TipoComprobante   string `json:"tipo_comprobante" validate:"omitempty,oneof=I E T N P"`
	EstadoComprobante string `json:"estado_comprobante" validate:"omitempty,oneof=0 1"`
	Complemento       string `json:"complemento"`
	ChunkDays         int    `json:"chunk_days" validate:"omitempty,min=1,max=365"`
	PollSeconds       int    `json:"poll_seconds" validate:"omitempty,min=5,max=300"`
}
