CREATE TABLE IF NOT EXISTS "student_mutations" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "student_id" UUID NULL DEFAULT NULL,
    "mutation_type" VARCHAR(20) NOT NULL,
    "nis" VARCHAR(20) NULL DEFAULT NULL,
    "nisn" VARCHAR(20) NULL DEFAULT NULL,
    "nama" VARCHAR(200) NOT NULL,
    "jenis_kelamin" VARCHAR(1) NULL DEFAULT NULL,
    "tempat_lahir" VARCHAR(100) NULL DEFAULT NULL,
    "tanggal_lahir" DATE NULL DEFAULT NULL,
    "asal_sekolah" VARCHAR(200) NULL DEFAULT NULL,
    "tujuan_sekolah" VARCHAR(200) NULL DEFAULT NULL,
    "alasan" TEXT NULL DEFAULT NULL,
    "dokumen_url" TEXT NULL DEFAULT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "catatan" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_mutations_school_period ON student_mutations(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_mutations_student ON student_mutations(student_id);
CREATE INDEX IF NOT EXISTS idx_mutations_status ON student_mutations(status);
