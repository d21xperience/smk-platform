-- =============================================================================
-- SMK PLATFORM - DATABASE INITIALIZATION
-- =============================================================================
-- Script ini dijalankan otomatis saat PostgreSQL container pertama kali start.
-- Schema disesuaikan dengan repository backend (student_repository.go, dll).
-- =============================================================================

-- Create extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================================================
-- ACADEMIC TABLES
-- =============================================================================

CREATE TABLE IF NOT EXISTS academic_years (
    id VARCHAR(36) PRIMARY KEY,
    year VARCHAR(9) NOT NULL UNIQUE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS curriculums (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    academic_year_id VARCHAR(36) REFERENCES academic_years(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS subjects (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    credit_hours INTEGER NOT NULL DEFAULT 0,
    curriculum_id VARCHAR(36) REFERENCES curriculums(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    grade INTEGER NOT NULL,
    academic_year_id VARCHAR(36) REFERENCES academic_years(id),
    capacity INTEGER NOT NULL DEFAULT 30,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, academic_year_id)
);

-- =============================================================================
-- STUDENT TABLES (sesuai repository/student_repository.go)
-- =============================================================================

CREATE TABLE IF NOT EXISTS students (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    nis VARCHAR(20) NOT NULL,
    nisn VARCHAR(10) NOT NULL,
    name VARCHAR(200) NOT NULL,
    gender VARCHAR(10) NOT NULL,
    birth_date TIMESTAMP WITH TIME ZONE NOT NULL,
    birth_place VARCHAR(100),
    address TEXT,
    phone VARCHAR(20),
    email VARCHAR(100),
    parent_name VARCHAR(100),
    parent_phone VARCHAR(20),
    parent_email VARCHAR(100),
    class_name VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,

    CONSTRAINT uq_students_school_period_nis UNIQUE (school_id, academic_period_id, nis),
    CONSTRAINT uq_students_school_period_nisn UNIQUE (school_id, academic_period_id, nisn)
);

CREATE INDEX IF NOT EXISTS idx_students_school_period ON students(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_students_nis ON students(nis);
CREATE INDEX IF NOT EXISTS idx_students_nisn ON students(nisn);
CREATE INDEX IF NOT EXISTS idx_students_status ON students(status);
CREATE INDEX IF NOT EXISTS idx_students_class ON students(school_id, academic_period_id, class_name);
CREATE INDEX IF NOT EXISTS idx_students_deleted_at ON students(deleted_at);

CREATE TABLE IF NOT EXISTS enrollments (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(20) NOT NULL,
    class_id VARCHAR(36) NOT NULL,
    enrollment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_enrollments_student ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_school_period ON enrollments(school_id, period_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_class ON enrollments(class_id);

CREATE TABLE IF NOT EXISTS guardians (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    relationship VARCHAR(50) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    occupation VARCHAR(100),
    address TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- OUTBOX TABLE (sesuai outbox/worker.go)
-- =============================================================================

CREATE TABLE IF NOT EXISTS outbox_events (
    id VARCHAR(36) PRIMARY KEY,
    event_type VARCHAR(100) NOT NULL,
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id VARCHAR(36) NOT NULL,
    payload JSONB NOT NULL,
    school_id VARCHAR(36),
    academic_period_id VARCHAR(20),
    correlation_id VARCHAR(36),
    trace_id VARCHAR(36),
    version VARCHAR(10) DEFAULT 'v1',
    attempts INT NOT NULL DEFAULT 0,
    last_error TEXT,
    next_retry TIMESTAMP WITH TIME ZONE,
    published_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_outbox_unpublished ON outbox_events(created_at ASC) WHERE published_at IS NULL;
CREATE INDEX IF NOT EXISTS idx_outbox_retry ON outbox_events(next_retry) WHERE next_retry IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_outbox_type ON outbox_events(event_type);
CREATE INDEX IF NOT EXISTS idx_outbox_aggregate ON outbox_events(aggregate_type, aggregate_id);

-- =============================================================================
-- SEED DATA
-- =============================================================================

INSERT INTO academic_years(id, year, start_date, end_date, is_active)
VALUES ('ay-2025-2026', '2025/2026', '2025-07-01', '2026-06-30', TRUE)
ON CONFLICT (year) DO NOTHING;

INSERT INTO curriculums(id, name, code, description, academic_year_id, is_active)
VALUES ('cur-km-2025', 'Kurikulum Merdeka', 'KM-2025', 'Kurikulum Merdeka untuk SMK', 'ay-2025-2026', TRUE)
ON CONFLICT (code) DO NOTHING;

INSERT INTO subjects(id, name, code, description, credit_hours, curriculum_id)
VALUES
    ('sub-mtk', 'Matematika', 'MTK', 'Matematika Umum', 4, 'cur-km-2025'),
    ('sub-bind', 'Bahasa Indonesia', 'BIND', 'Bahasa Indonesia', 3, 'cur-km-2025'),
    ('sub-bing', 'Bahasa Inggris', 'BING', 'Bahasa Inggris', 3, 'cur-km-2025'),
    ('sub-pmw', 'Pemrograman Web', 'PMW', 'Pemrograman Web Dasar', 6, 'cur-km-2025')
ON CONFLICT (code) DO NOTHING;

INSERT INTO classes(id, name, grade, academic_year_id, capacity)
VALUES
    ('cls-x-ipa-1', 'X IPA 1', 10, 'ay-2025-2026', 30),
    ('cls-x-ipa-2', 'X IPA 2', 10, 'ay-2025-2026', 30),
    ('cls-xi-rpl-1', 'XI RPL 1', 11, 'ay-2025-2026', 25)
ON CONFLICT (name, academic_year_id) DO NOTHING;

-- =============================================================================
-- VERIFICATION
-- =============================================================================
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully!';
    RAISE NOTICE 'Tables created: academic_years, curriculums, subjects, classes, students, enrollments, guardians, outbox_events';
    RAISE NOTICE 'Extensions enabled: uuid-ossp, pg_trgm';
END $$;
