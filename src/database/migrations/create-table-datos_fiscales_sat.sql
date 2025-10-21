CREATE TABLE datos_fiscales_sat (
    uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    rfc VARCHAR(225) NOT NULL,
    cer_b64_encriptado TEXT NOT NULL,
    key_b64_encriptado TEXT NOT NULL,
    password_efirma_encrip VARCHAR(255) NOT NULL,
    updated_at TIMESTAMP,
    created_at TIMESTAMP,
    updated_by UUID,
    created_by UUID,
    deleted_at TIMESTAMP,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
