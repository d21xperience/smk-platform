CREATE TABLE IF NOT EXISTS ptk_archives (
    id UUID PRIMARY KEY,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis_dokumen VARCHAR(100) NOT NULL,
    nomor_dokumen VARCHAR(100),
    tanggal_terbit DATE,
    tanggal_kadaluarsa DATE,
    catatan TEXT,
    file_url TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_ptk_archives_ptk_id ON ptk_archives(ptk_id);
