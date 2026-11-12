CREATE TABLE IF NOT EXISTS "student_letters" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "student_id" UUID NOT NULL REFERENCES "tabel_siswa"("peserta_didik_id") ON DELETE CASCADE,
    "letter_number" VARCHAR(100) NOT NULL,
    "parent_name" VARCHAR(100) NULL DEFAULT NULL,
    "purpose" TEXT NULL DEFAULT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "file_url" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_letters_number UNIQUE ("letter_number")
);

CREATE INDEX IF NOT EXISTS idx_letters_school_period ON student_letters(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_letters_student ON student_letters(student_id);
CREATE INDEX IF NOT EXISTS idx_letters_status ON student_letters(status);
