-- Add RFC, CIEC, and download tracking fields to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS rfc VARCHAR(13) UNIQUE;
ALTER TABLE users ADD COLUMN IF NOT EXISTS ciec TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS descargas_realizadas INTEGER DEFAULT 0 NOT NULL;
ALTER TABLE users ADD COLUMN IF NOT EXISTS ultimo_reset_descargas DATE;

-- Create index for RFC lookups
CREATE INDEX IF NOT EXISTS idx_users_rfc ON users(rfc);

