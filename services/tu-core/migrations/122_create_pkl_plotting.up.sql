CREATE TABLE IF NOT EXISTS "pkl_plotting" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "nama" VARCHAR(150) NOT NULL,
    "kelas" VARCHAR(50) NOT NULL,
    "jurusan" VARCHAR(100) NOT NULL,
    "perusahaan" VARCHAR(150) NOT NULL,
    "kota" VARCHAR(100) NULL DEFAULT NULL,
    "pembimbing" VARCHAR(150) NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_pkl_plotting_school ON pkl_plotting(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_pkl_plotting_jurusan ON pkl_plotting(jurusan);
