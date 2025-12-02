-- Remove download tracking fields from users table
DROP INDEX IF EXISTS idx_users_rfc;
ALTER TABLE users DROP COLUMN IF EXISTS ultimo_reset_descargas;
ALTER TABLE users DROP COLUMN IF EXISTS descargas_realizadas;
ALTER TABLE users DROP COLUMN IF EXISTS ciec;
ALTER TABLE users DROP COLUMN IF EXISTS rfc;

