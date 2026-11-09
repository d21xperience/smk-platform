ALTER TABLE inventory_items DROP COLUMN IF EXISTS goods_receipt_id;
DROP TABLE IF EXISTS goods_receipt_items;
DROP TABLE IF EXISTS goods_receipts;
