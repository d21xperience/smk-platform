-- Hapus constraint lama
ALTER TABLE inventory_items DROP CONSTRAINT IF EXISTS inventory_items_status_aset_check;

-- Tambahkan constraint baru yang mengizinkan DIUSULKAN_HAPUS
ALTER TABLE inventory_items ADD CONSTRAINT inventory_items_status_aset_check 
CHECK (status_aset IN ('AKTIF', 'TIDAK_AKTIF', 'RUSAK_BERAT', 'DIHAPUS', 'DIPINDAHKAN', 'DIUSULKAN_HAPUS'));

-- Tambahkan kolom relasi ke disposal
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS disposal_request_id UUID;
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS metode_penghapusan VARCHAR(50);
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS keterangan_penghapusan TEXT;
