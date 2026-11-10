DROP INDEX IF EXISTS idx_documents_semester;
ALTER TABLE sarpras_inventory_documents DROP COLUMN IF EXISTS semester_id;
