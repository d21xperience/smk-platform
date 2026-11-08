-- Tabel utama persuratan (menggabungkan surat masuk, keluar, dan internal)
CREATE TABLE IF NOT EXISTS correspondences (
    id                  UUID PRIMARY KEY,
    school_id           UUID NOT NULL,
    academic_period_id  VARCHAR(20) NOT NULL,
    type                VARCHAR(20) NOT NULL CHECK (type IN ('INCOMING', 'OUTGOING')),
    jenis_surat         VARCHAR(50) NOT NULL, -- SK, ST, SPPD, SP, SE, SU, SPENGUMUMAN, BA, dll.
    number              VARCHAR(100) NOT NULL,
    date                DATE NOT NULL,
    subject             TEXT NOT NULL,
    from_party          TEXT NOT NULL,
    to_party            TEXT NOT NULL,
    description         TEXT,
    attachment_url      TEXT,
    file_url            TEXT,       -- URL ke PDF yang di-generate
    ttd_url             TEXT,       -- URL ke gambar TTD digital
    konten              JSONB,      -- Field dinamis sesuai jenis surat
    status              VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'PROCESSED', 'ARCHIVED')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_corr_school_period ON correspondences (school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_corr_type_status ON correspondences (type, status);
CREATE INDEX IF NOT EXISTS idx_corr_jenis_surat ON correspondences (jenis_surat);

-- Tabel counter untuk auto-numbering
CREATE TABLE IF NOT EXISTS correspondence_counters (
    id                  UUID PRIMARY KEY,
    school_id           UUID NOT NULL,
    academic_period_id  VARCHAR(20) NOT NULL,
    jenis_surat         VARCHAR(50) NOT NULL,
    last_number         INT NOT NULL DEFAULT 0,
    UNIQUE (school_id, academic_period_id, jenis_surat)
);
