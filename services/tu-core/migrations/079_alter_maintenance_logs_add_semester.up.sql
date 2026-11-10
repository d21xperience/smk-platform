ALTER TABLE maintenance_logs ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE maintenance_logs SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_logs_semester ON maintenance_logs(semester_id);
