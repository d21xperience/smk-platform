ALTER TABLE inventory_items 
ADD COLUMN IF NOT EXISTS nomor_registrasi VARCHAR(100),
ADD COLUMN IF NOT EXISTS sumber_perolehan VARCHAR(50),
ADD COLUMN IF NOT EXISTS keterangan TEXT,
ADD COLUMN IF NOT EXISTS masa_manfaat_tahun INT DEFAULT 5,
ADD COLUMN IF NOT EXISTS nilai_buku NUMERIC(15, 2) DEFAULT 0,
ADD COLUMN IF NOT EXISTS status_aset VARCHAR(20) DEFAULT 'AKTIF' CHECK (status_aset IN ('AKTIF', 'TIDAK_AKTIF', 'RUSAK_BERAT', 'DIHAPUS', 'DIPINDAHKAN')),
ADD COLUMN IF NOT EXISTS tanggal_pencatatan DATE,
ADD COLUMN IF NOT EXISTS tanggal_penghapusan DATE,
ADD COLUMN IF NOT EXISTS alasan_penghapusan TEXT;

CREATE TABLE IF NOT EXISTS inventory_item_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    inventory_item_id VARCHAR(100) NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    action VARCHAR(50) NOT NULL,
    old_data JSONB,
    new_data JSONB,
    description TEXT,
    changed_by VARCHAR(100) NOT NULL,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inv_history_item ON inventory_item_history(inventory_item_id);
