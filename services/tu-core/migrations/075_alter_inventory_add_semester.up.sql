ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE inventory_items SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_inventory_items_semester ON inventory_items(semester_id);
