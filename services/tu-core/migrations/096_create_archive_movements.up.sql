CREATE TABLE IF NOT EXISTS archive_movements (
    id                      UUID PRIMARY KEY,
    archive_id              UUID NOT NULL REFERENCES archives(id) ON DELETE CASCADE,
    type                    VARCHAR(20) NOT NULL CHECK (type IN ('ARCHIVE', 'BORROW', 'RETURN', 'MOVE', 'DESTROY')),
    petugas_id              UUID NOT NULL,
    peminjam_id             UUID,
    tanggal                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    tanggal_kembali_estimasi DATE,
    tanggal_kembali_aktual  DATE,
    tujuan                  TEXT,
    keterangan              TEXT,
    status                  VARCHAR(20) NOT NULL DEFAULT 'AKTIF'
);
CREATE INDEX IF NOT EXISTS idx_movement_archive ON archive_movements(archive_id);
