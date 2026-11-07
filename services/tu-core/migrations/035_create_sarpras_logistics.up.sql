CREATE TABLE IF NOT EXISTS sarpras_logistics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    kode VARCHAR(50) UNIQUE NOT NULL,
    nama VARCHAR(255) NOT NULL,
    kategori VARCHAR(100) NOT NULL,
    satuan VARCHAR(50) NOT NULL,
    stok INT NOT NULL DEFAULT 0,
    stok_minimum INT NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sarpras_logistic_mutations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    logistic_id UUID NOT NULL REFERENCES sarpras_logistics(id) ON DELETE CASCADE,
    tipe VARCHAR(10) NOT NULL CHECK (tipe IN ('IN', 'OUT')),
    jumlah INT NOT NULL CHECK (jumlah > 0),
    keterangan TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_sarpras_logistics_kode ON sarpras_logistics(kode);
CREATE INDEX IF NOT EXISTS idx_sarpras_mutations_logistic_id ON sarpras_logistic_mutations(logistic_id);
