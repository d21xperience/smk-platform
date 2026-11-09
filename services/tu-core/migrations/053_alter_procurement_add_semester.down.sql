DROP INDEX IF EXISTS idx_procurement_semester;
ALTER TABLE procurement_requests DROP COLUMN IF EXISTS semester_id;
