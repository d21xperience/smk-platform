DROP INDEX IF EXISTS idx_logs_semester;
ALTER TABLE maintenance_logs DROP COLUMN IF EXISTS semester_id;
