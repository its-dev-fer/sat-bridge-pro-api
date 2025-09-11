CREATE TABLE plan_suscripcion (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nombre VARCHAR(100) NOT NULL,
    descripcion VARCHAR(255),
    limite_descargas_mensuales INTEGER NOT NULL,
    precio DECIMAL(10,2) NOT NULL,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_plan_suscripcion_activo ON plan_suscripcion(activo);
CREATE INDEX idx_plan_suscripcion_precio ON plan_suscripcion(precio);
CREATE INDEX idx_plan_suscripcion_deleted_at ON plan_suscripcion(deleted_at);
CREATE UNIQUE INDEX idx_plan_suscripcion_nombre ON plan_suscripcion(nombre) WHERE deleted_at IS NULL;