package validation

type DatosFiscalesRequest struct {
	RFC      string `json:"rfc" validate:"required,min=12,max=13" example:"XAXX010101000"`
	Password string `json:"password" validate:"required,min=8,max=50" example:"efirma_password"`
}
