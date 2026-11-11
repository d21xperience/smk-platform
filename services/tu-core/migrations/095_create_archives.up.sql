CREATE TABLE IF NOT EXISTS archives (
    id                  UUID PRIMARY KEY,
    school_id           UUID NOT NULL,
    correspondence_id   UUID NOT NULL UNIQUE REFERENCES correspondences(id) ON DELETE CASCADE,
    lokasi_kode         VARCHAR(100) NOT NULL,
    kategori            VARCHAR(20) NOT NULL CHECK (kategori IN ('AKTIF', 'INAKTIF', 'STATIS')),
    tanggal_arsip       DATE NOT NULL,
    tanggal_musnah      DATE NOT NULL,
    retensi_tahun       INT NOT NULL,
    status              VARCHAR(20) NOT NULL DEFAULT 'TERARSIP' CHECK (status IN ('TERARSIP', 'DIPINJAM', 'DIMUSNAHKAN')),
    current_holder_id   UUID,
    keterangan          TEXT,
    archived_by         UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_archive_corr ON archives(correspondence_id);
CREATE INDEX IF NOT EXISTS idx_archive_status ON archives(status, tanggal_musnah);
