CREATE TABLE IF NOT EXISTS "registrations" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "calon_nisn" VARCHAR(20) NOT NULL,
    "nis" VARCHAR(20) NULL DEFAULT NULL,
    "calon_nama" VARCHAR(200) NOT NULL,
    "calon_jenis_kelamin" VARCHAR(1) NOT NULL,
    "calon_tempat_lahir" VARCHAR(100) NULL DEFAULT NULL,
    "calon_tanggal_lahir" DATE NULL DEFAULT NULL,
    "asal_sekolah" VARCHAR(200) NULL DEFAULT NULL,
    "jurusan_dipilih" VARCHAR(50) NOT NULL,
    "biaya_standar" NUMERIC(12,0) NOT NULL DEFAULT 0,
    "biaya_pendaftaran" NUMERIC(12,0) NOT NULL DEFAULT 0,
    "berkas_ijazah" BOOLEAN NOT NULL DEFAULT FALSE,
    "berkas_kk" BOOLEAN NOT NULL DEFAULT FALSE,
    "berkas_akta" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" VARCHAR(20) NOT NULL DEFAULT 'REGISTERED',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_registrations_nisn_period UNIQUE ("calon_nisn", "academic_period_id")
);

CREATE INDEX IF NOT EXISTS idx_registrations_school_period ON registrations(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_registrations_status ON registrations(status);
