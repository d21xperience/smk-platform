ALTER TABLE disposal_requests ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE disposal_requests SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_disposal_semester ON disposal_requests(semester_id);
