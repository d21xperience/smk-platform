ALTER TABLE sarpras_inventory_documents ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE sarpras_inventory_documents SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_documents_semester ON sarpras_inventory_documents(semester_id);
