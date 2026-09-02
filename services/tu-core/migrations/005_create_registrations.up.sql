-- +migrate Up
CREATE TABLE IF NOT EXISTS registrations (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,

    -- Data Calon Siswa
    calon_nisn VARCHAR(20),
    calon_nama VARCHAR(200) NOT NULL,
    calon_jenis_kelamin VARCHAR(10) NOT NULL,
    calon_tempat_lahir VARCHAR(100),
    calon_tanggal_lahir TIMESTAMP WITH TIME ZONE NOT NULL,
    calon_alamat TEXT,
    calon_telepon VARCHAR(20),
    calon_email VARCHAR(100),

    -- Data Orang Tua/Wali
    nama_ayah VARCHAR(200),
    nama_ibu VARCHAR(200),
    telepon_ortu VARCHAR(20),
    pekerjaan_ayah VARCHAR(100),
    pekerjaan_ibu VARCHAR(100),

    -- Data Pendaftaran
    asal_sekolah VARCHAR(200),
    jurusan_dipilih VARCHAR(100),
    alasan_memilih TEXT,

    -- Status & Verifikasi
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    catatan_verifikasi TEXT,
    verified_by VARCHAR(36),
    verified_at TIMESTAMP WITH TIME ZONE,

    -- Metadata
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_registrations_nisn_period UNIQUE (calon_nisn, academic_period_id)
);

CREATE INDEX IF NOT EXISTS idx_registrations_school_period ON registrations(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);
CREATE INDEX IF NOT EXISTS idx_registrations_nisn ON registrations(calon_nisn);
CREATE INDEX IF NOT EXISTS idx_registrations_created_at ON registrations(created_at DESC);
