-- Drop tabel lama (karena data sebelumnya hanya dummy)
DROP TABLE IF EXISTS inventory_items;

CREATE TABLE IF NOT EXISTS inventory_items (
    id VARCHAR(100) PRIMARY KEY,
    catalog_id UUID REFERENCES master_catalogs(id) ON DELETE RESTRICT,
    condition VARCHAR(50),
    room_id UUID REFERENCES tabel_kelas(rombongan_belajar_id) ON DELETE SET NULL,
    purchase_date DATE NOT NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);
CREATE INDEX IF NOT EXISTS idx_inventory_items_room ON inventory_items(room_id);
