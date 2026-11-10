DROP INDEX IF EXISTS idx_disposal_semester;
ALTER TABLE disposal_requests DROP COLUMN IF EXISTS semester_id;
