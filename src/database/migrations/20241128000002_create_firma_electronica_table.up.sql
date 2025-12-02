-- Create firma_electronica table for FIEL storage
CREATE TABLE IF NOT EXISTS firma_electronicas (
    id UUID PRIMARY KEY,
    usuario_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rfc VARCHAR(13) NOT NULL,
    certificado_cer TEXT NOT NULL,
    clave_privada_key TEXT NOT NULL,
    password_key TEXT NOT NULL,
    nombre_certificado VARCHAR(255),
    fecha_vigencia_inicio DATE,
    fecha_vigencia_fin DATE,
    activo BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_firma_electronicas_usuario_id ON firma_electronicas(usuario_id);
CREATE INDEX IF NOT EXISTS idx_firma_electronicas_rfc ON firma_electronicas(rfc);
CREATE INDEX IF NOT EXISTS idx_firma_electronicas_activo ON firma_electronicas(activo);
CREATE INDEX IF NOT EXISTS idx_firma_electronicas_deleted_at ON firma_electronicas(deleted_at);

-- Add comment
COMMENT ON TABLE firma_electronicas IS 'Stores encrypted FIEL (Firma Electrónica) certificates and private keys for SAT authentication';

