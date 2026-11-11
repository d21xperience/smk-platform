CREATE TABLE IF NOT EXISTS agenda_entries (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID NOT NULL,
    correspondence_id   UUID NOT NULL REFERENCES correspondences(id) ON DELETE CASCADE,
    jenis               VARCHAR(10) NOT NULL CHECK (jenis IN ('IN', 'OUT')),
    nomor_agenda        VARCHAR(50) NOT NULL,
    nomor_urut          INT NOT NULL,
    tahun               INT NOT NULL,
    tanggal_catat       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    petugas_id          UUID NOT NULL,
    UNIQUE (school_id, jenis, tahun, nomor_urut)
);
CREATE INDEX IF NOT EXISTS idx_agenda_corr ON agenda_entries(correspondence_id);
