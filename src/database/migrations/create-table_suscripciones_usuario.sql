CREATE TABLE suscripciones_usuario(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    usuario_id UUID,
    plan_id UUID,
    limite_descargas_mensuales INT NOT NULL,
    fecha_inicio DATE NOT NULL,
    fecha_fin DATE NOT NULL,
    estatus ENUM('activo', 'cancelado', 'pendiente') NOT NULL,
    CONSTRAINT fk_usuario_id FOREIGN KEY (usuario_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_plan_id FOREIGN KEY (plan_id) REFERENCES planes_suscripcion(id) ON DELETE CASCADE
)