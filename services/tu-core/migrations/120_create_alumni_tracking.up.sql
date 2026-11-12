CREATE TABLE IF NOT EXISTS "alumni_tracking" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "nama" VARCHAR(150) NOT NULL,
    "tahun_lulus" VARCHAR(4) NOT NULL,
    "jurusan" VARCHAR(100) NOT NULL,
    "status" VARCHAR(50) NOT NULL DEFAULT 'Belum Bekerja',
    "instansi" VARCHAR(150) NULL DEFAULT NULL,
    "posisi_usaha" VARCHAR(150) NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_alumni_tracking_school ON alumni_tracking(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_alumni_tracking_status ON alumni_tracking(status);
CREATE INDEX IF NOT EXISTS idx_alumni_tracking_tahun ON alumni_tracking(tahun_lulus);
