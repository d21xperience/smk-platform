-- +migrate Up
CREATE TABLE IF NOT EXISTS attendance_sessions (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    class_id VARCHAR(36) NOT NULL,
    subject_id VARCHAR(36),
    teacher_id VARCHAR(36),
    session_date DATE NOT NULL,
    start_time TIME,
    end_time TIME,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_sessions_date_class UNIQUE (school_id, academic_period_id, session_date, class_id)
);

CREATE TABLE IF NOT EXISTS attendance_records (
    id VARCHAR(36) PRIMARY KEY,
    session_id VARCHAR(36) NOT NULL REFERENCES attendance_sessions(id) ON DELETE CASCADE,
    student_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL,
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_records_session_student UNIQUE (session_id, student_id)
);

CREATE INDEX IF NOT EXISTS idx_sessions_school_period ON attendance_sessions(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_sessions_date ON attendance_sessions(session_date);
CREATE INDEX IF NOT EXISTS idx_sessions_class ON attendance_sessions(class_id);
CREATE INDEX IF NOT EXISTS idx_records_session ON attendance_records(session_id);
CREATE INDEX IF NOT EXISTS idx_records_student ON attendance_records(student_id);
CREATE INDEX IF NOT EXISTS idx_records_status ON attendance_records(status);
