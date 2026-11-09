ALTER TABLE procurement_requests ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
CREATE INDEX IF NOT EXISTS idx_procurement_semester ON procurement_requests(semester_id);
