package validation

type CreatePlanSuscripcion struct {
    Nombre                   string  `json:"nombre" validate:"required,min=3,max=100"`
    Descripcion              string  `json:"descripcion" validate:"max=255"`
    LimiteDescargasMensuales int     `json:"limite_descargas_mensuales" validate:"required,min=1"`
    Precio                   float64 `json:"precio" validate:"required,gt=0"`
    Activo                   bool    `json:"activo" validate:"required"`
}

type UpdatePlanSuscripcion struct {
    Nombre                   string  `json:"nombre" validate:"omitempty,min=3,max=100"`
    Descripcion              string  `json:"descripcion" validate:"omitempty,max=255"`
    LimiteDescargasMensuales int     `json:"limite_descargas_mensuales" validate:"omitempty,min=1"`
    Precio                   float64 `json:"precio" validate:"omitempty,gt=0"`
    Activo                   *bool   `json:"activo" validate:"omitempty"`
}

type QueryPlanSuscripcion struct {
    Page   int    `query:"page" validate:"min=1"`
    Limit  int    `query:"limit" validate:"min=1,max=100"`
    Search string `query:"search"`
}
