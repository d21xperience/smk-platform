-- Migration 116: Create student_academic_scores (REVISED)
-- Tabel nilai asesmen granular per mapel (Buku Induk Siswa)
-- Filter wajib: school_id + academic_period_id (konsisten dengan 114/115)
-- Kolom "semester" eksplisit (1/2) sesuai format Dapodik "20261"/"20262"

CREATE TABLE IF NOT EXISTS "student_academic_scores" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "semester" SMALLINT NOT NULL,
    "student_id" UUID NOT NULL REFERENCES "tabel_siswa"("peserta_didik_id") ON DELETE CASCADE,
    "pembelajaran_id" UUID NULL DEFAULT NULL REFERENCES "tabel_pembelajaran"("pembelajaran_id") ON DELETE SET NULL,
    "mata_pelajaran_id" INTEGER NOT NULL,
    "nama_mata_pelajaran" VARCHAR(100) NOT NULL,
    "assessment_type" VARCHAR(20) NOT NULL,
    "score" NUMERIC(5,2) NOT NULL,
    "assessment_date" DATE NULL DEFAULT NULL,
    "catatan" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT ck_academic_scores_semester
        CHECK ("semester" IN (1, 2)),
    CONSTRAINT ck_academic_scores_type
        CHECK ("assessment_type" IN ('DAILY', 'MIDTERM', 'FINAL')),
    CONSTRAINT ck_academic_scores_score
        CHECK ("score" >= 0 AND "score" <= 100)
);

CREATE INDEX IF NOT EXISTS idx_academic_scores_school_period
    ON student_academic_scores(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_academic_scores_student
    ON student_academic_scores(student_id);
CREATE INDEX IF NOT EXISTS idx_academic_scores_pembelajaran
    ON student_academic_scores(pembelajaran_id);
CREATE INDEX IF NOT EXISTS idx_academic_scores_type
    ON student_academic_scores(assessment_type);

-- Partial unique: 1 nilai MIDTERM & 1 nilai FINAL per mapel per semester per siswa.
-- DAILY boleh banyak (dibatasi di UI/service).
CREATE UNIQUE INDEX IF NOT EXISTS uq_academic_scores_exam
    ON student_academic_scores(student_id, academic_period_id, semester, mata_pelajaran_id, assessment_type)
    WHERE assessment_type IN ('MIDTERM', 'FINAL');
