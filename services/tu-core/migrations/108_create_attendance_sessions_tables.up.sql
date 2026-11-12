-- Tabel Sesi Absensi (per hari per mapel)
CREATE TABLE IF NOT EXISTS "attendance_sessions" (
    "id" UUID NOT NULL PRIMARY KEY,
    "school_id" VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "class_id" UUID NOT NULL,
    "subject_id" INTEGER NULL DEFAULT NULL,
    "teacher_id" UUID NULL DEFAULT NULL,
    "session_date" DATE NOT NULL,
    "start_time" TIME NULL DEFAULT NULL,
    "end_time" TIME NULL DEFAULT NULL,
    "status" VARCHAR(20) NOT NULL DEFAULT 'OPEN',
    "notes" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_att_sessions_school_period ON attendance_sessions(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_att_sessions_class_date ON attendance_sessions(class_id, session_date);
CREATE INDEX IF NOT EXISTS idx_att_sessions_date ON attendance_sessions(session_date);

-- Tabel Detail Absensi (per siswa per sesi)
CREATE TABLE IF NOT EXISTS "attendance_records" (
    "id" UUID NOT NULL PRIMARY KEY,
    "session_id" UUID NOT NULL REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    "student_id" UUID NOT NULL REFERENCES tabel_siswa(peserta_didik_id) ON DELETE CASCADE,
    "status" VARCHAR(20) NOT NULL,
    "note" TEXT NULL DEFAULT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_att_records_session_student UNIQUE ("session_id", "student_id")
);

CREATE INDEX IF NOT EXISTS idx_att_records_session ON attendance_records(session_id);
CREATE INDEX IF NOT EXISTS idx_att_records_student ON attendance_records(student_id);
CREATE INDEX IF NOT EXISTS idx_att_records_status ON attendance_records(status);
