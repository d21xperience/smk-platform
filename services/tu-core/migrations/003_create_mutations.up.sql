-- +migrate Up
CREATE TABLE IF NOT EXISTS mutations (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    student_id VARCHAR(36), -- NULL untuk mutasi MASUK
    mutation_type VARCHAR(10) NOT NULL, -- 'MASUK' atau 'KELUAR'

    -- Data khusus Mutasi MASUK
    calon_nama VARCHAR(200),
    calon_nisn VARCHAR(20),
    asal_sekolah VARCHAR(200),

    -- Data khusus Mutasi KELUAR
    tujuan_sekolah VARCHAR(200),

    alasan TEXT NOT NULL,
    dokumen_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING', -- 'PENDING', 'DISETUJUI', 'DITOLAK'
    catatan_verifikasi TEXT,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mutations_school_period ON mutations(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_mutations_student ON mutations(student_id);
CREATE INDEX IF NOT EXISTS idx_mutations_status ON mutations(status);
CREATE INDEX IF NOT EXISTS idx_mutations_type ON mutations(mutation_type);
