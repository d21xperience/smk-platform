ALTER TABLE sarpras_logistics ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE sarpras_logistics SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_logistics_semester ON sarpras_logistics(semester_id);
