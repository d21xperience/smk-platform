ALTER TABLE inventory_items DROP COLUMN IF EXISTS disposal_request_id;
ALTER TABLE inventory_items DROP COLUMN IF EXISTS metode_penghapusan;
ALTER TABLE inventory_items DROP COLUMN IF EXISTS keterangan_penghapusan;

ALTER TABLE inventory_items DROP CONSTRAINT IF EXISTS inventory_items_status_aset_check;
ALTER TABLE inventory_items ADD CONSTRAINT inventory_items_status_aset_check 
CHECK (status_aset IN ('AKTIF', 'TIDAK_AKTIF', 'RUSAK_BERAT', 'DIHAPUS', 'DIPINDAHKAN'));
