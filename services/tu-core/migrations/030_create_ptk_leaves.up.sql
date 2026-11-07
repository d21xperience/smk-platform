CREATE TABLE IF NOT EXISTS ptk_leaves (
    id UUID PRIMARY KEY,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis_cuti VARCHAR(50) NOT NULL,
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    alasan TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    catatan_approval TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_ptk_leaves_ptk_id ON ptk_leaves(ptk_id);
