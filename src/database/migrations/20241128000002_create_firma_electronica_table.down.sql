-- Drop firma_electronica table and indexes
DROP INDEX IF EXISTS idx_firma_electronicas_deleted_at;
DROP INDEX IF EXISTS idx_firma_electronicas_activo;
DROP INDEX IF EXISTS idx_firma_electronicas_rfc;
DROP INDEX IF EXISTS idx_firma_electronicas_usuario_id;
DROP TABLE IF EXISTS firma_electronicas;

