-- Tabel sesi absensi (per pertemuan / per jadwal)
-- Dirancang cocok dengan repository.attendance_repository.go & models.Session.
-- class_id FK ke tabel_kelas (UUID). subject_id/teacher_id sengaja TEXT tanpa FK
-- karena tabel_mapel memakai INT sedangkan model Go memakai string; FK bisa
-- ditambahkan belakangan setelah tipe dikonfirmasi.

CREATE TABLE IF NOT EXISTS attendance_sessions (
    id                  UUID PRIMARY KEY,
    school_id           UUID NOT NULL,
    academic_period_id  VARCHAR(20) NOT NULL,
    class_id            UUID NOT NULL REFERENCES tabel_kelas(rombongan_belajar_id) ON DELETE CASCADE,
    subject_id          TEXT,
    teacher_id          TEXT,
    session_date        DATE NOT NULL,
    start_time          TIME,
    end_time            TIME,
    status              VARCHAR(20) NOT NULL DEFAULT 'DRAFT'
                        CHECK (status IN ('DRAFT', 'SUBMITTED', 'CLOSED')),
    notes               TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_att_sessions_school_period
    ON attendance_sessions (school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_att_sessions_class_date
    ON attendance_sessions (class_id, session_date DESC);

-- Tabel record absensi per siswa
CREATE TABLE IF NOT EXISTS attendance_records (
    id          UUID PRIMARY KEY,
    session_id  UUID NOT NULL REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    student_id  UUID NOT NULL,
    status      VARCHAR(20) NOT NULL
                CHECK (status IN ('PRESENT', 'ABSENT', 'LATE', 'SICK', 'PERMISSION')),
    note        TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (session_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_att_records_session ON attendance_records (session_id);
CREATE INDEX IF NOT EXISTS idx_att_records_student ON attendance_records (student_id);
