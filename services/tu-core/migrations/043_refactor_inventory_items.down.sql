DROP TABLE IF EXISTS inventory_items;
-- Mengembalikan ke struktur flat sebelumnya
CREATE TABLE IF NOT EXISTS inventory_items (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(100),
    category VARCHAR(100),
    condition VARCHAR(50),
    room_id UUID REFERENCES tabel_kelas(rombongan_belajar_id) ON DELETE SET NULL,
    created_at BIGINT DEFAULT EXTRACT(EPOCH FROM NOW())
);
