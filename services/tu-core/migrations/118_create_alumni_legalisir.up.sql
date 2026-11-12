CREATE TABLE IF NOT EXISTS "alumni_legalisir" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "nama" VARCHAR(150) NOT NULL,
    "tahun_lulus" VARCHAR(4) NOT NULL,
    "jurusan" VARCHAR(100) NOT NULL,
    "no_ijazah" VARCHAR(100) NOT NULL,
    "jumlah_lembar" INTEGER NOT NULL DEFAULT 1,
    "status" VARCHAR(50) NOT NULL DEFAULT 'PROSES TTD',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alumni_legalisir_school ON alumni_legalisir(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_alumni_legalisir_status ON alumni_legalisir(status);
CREATE INDEX IF NOT EXISTS idx_alumni_legalisir_nama ON alumni_legalisir(nama);
