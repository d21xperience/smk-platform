ALTER TABLE inventory_items DROP COLUMN IF EXISTS current_borrowing_id;
ALTER TABLE inventory_items DROP CONSTRAINT IF EXISTS inventory_items_status_aset_check;
ALTER TABLE inventory_items ADD CONSTRAINT inventory_items_status_aset_check 
CHECK (status_aset IN ('AKTIF', 'TIDAK_AKTIF', 'RUSAK_BERAT', 'DIHAPUS', 'DIPINDAHKAN', 'DIUSULKAN_HAPUS'));
