DROP INDEX IF EXISTS idx_inventory_items_semester;
ALTER TABLE inventory_items DROP COLUMN IF EXISTS semester_id;
