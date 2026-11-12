CREATE TABLE IF NOT EXISTS "student_scholarships" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "student_id" UUID NOT NULL REFERENCES "tabel_siswa"("peserta_didik_id") ON DELETE CASCADE,
    "category" VARCHAR(50) NOT NULL,
    "card_number" VARCHAR(100) NULL DEFAULT NULL,
    "reason" TEXT NULL DEFAULT NULL,
    "is_document_valid" BOOLEAN NOT NULL DEFAULT FALSE,
    "status" VARCHAR(30) NOT NULL DEFAULT 'SK NOMINASI',
    "bank_name" VARCHAR(50) NULL DEFAULT NULL,
    "bank_account" VARCHAR(50) NULL DEFAULT NULL,
    "disbursed_amount" NUMERIC(12,0) NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_scholarships_school_period ON student_scholarships(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_scholarships_student ON student_scholarships(student_id);
CREATE INDEX IF NOT EXISTS idx_scholarships_status ON student_scholarships(status);
