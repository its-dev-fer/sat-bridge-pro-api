CREATE TABLE cfdis (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    folio_fiscal VARCHAR(36) NOT NULL,
    rfc_emisor VARCHAR(13) NOT NULL,
    nombre_emisor VARCHAR(255) NOT NULL,
    rfc_receptor VARCHAR(13) NOT NULL,
    nombre_receptor VARCHAR(255) NOT NULL,
    fecha_emision TIMESTAMP,
    fecha_certificacion TIMESTAMP,
    pac_certifico VARCHAR(255),
    total NUMERIC(18, 6),
    efecto_comprobante VARCHAR(20),
    estatus_cancelacion VARCHAR(50),
    estado_comprobante VARCHAR(50),
    entra BOOLEAN,
    comentarios TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_cfdis_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uq_cfdis_user_folio UNIQUE (user_id, folio_fiscal)
);

CREATE INDEX idx_cfdis_user_id ON cfdis (user_id);
