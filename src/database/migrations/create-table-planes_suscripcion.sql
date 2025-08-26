CREATE TABLE planes_suscripcion(
    id  UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nombre VARCHAR NOT NULL,
    descripcion VARCHAR NOT NULL,
    limite_descargas_mensuales INT NOT NULL,
    precio DECIMAL(10, 2) NOT NULL,
    activo BOOLEAN NOT NULL DEFAULT FALSE
);