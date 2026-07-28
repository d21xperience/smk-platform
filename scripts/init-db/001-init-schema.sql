-- =============================================================================
-- SMK PLATFORM - DATABASE INITIALIZATION
-- =============================================================================
-- Script ini dijalankan otomatis saat PostgreSQL container pertama kali start.
-- File akan dijalankan sesuai urutan nama file (001, 002, dst).
--
-- CATATAN:
-- - Script ini HANYA dijalankan sekali (saat volume pertama kali dibuat)
-- - Jika ingin re-run, hapus volume: docker-compose down -v
-- - Untuk migration selanjutnya, gunakan golang-migrate
-- =============================================================================

-- Create schemas untuk setiap Bounded Context
CREATE SCHEMA IF NOT EXISTS academic;
CREATE SCHEMA IF NOT EXISTS student;
CREATE SCHEMA IF NOT EXISTS attendance;
CREATE SCHEMA IF NOT EXISTS assessment;
CREATE SCHEMA IF NOT EXISTS teaching;
CREATE SCHEMA IF NOT EXISTS reporting;

-- Grant permissions
GRANT ALL PRIVILEGES ON SCHEMA academic TO smk_dev;
GRANT ALL PRIVILEGES ON SCHEMA student TO smk_dev;
GRANT ALL PRIVILEGES ON SCHEMA attendance TO smk_dev;
GRANT ALL PRIVILEGES ON SCHEMA assessment TO smk_dev;
GRANT ALL PRIVILEGES ON SCHEMA teaching TO smk_dev;
GRANT ALL PRIVILEGES ON SCHEMA reporting TO smk_dev;

-- Create extension untuk UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create extension untuk full-text search
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- =============================================================================
-- ACADEMIC SCHEMA
-- =============================================================================

-- Academic Years
CREATE TABLE IF NOT EXISTS academic.academic_years (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    year VARCHAR(9) NOT NULL UNIQUE,  -- Format: 2025/2026
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Curriculums
CREATE TABLE IF NOT EXISTS academic.curriculums (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    academic_year_id UUID REFERENCES academic.academic_years(id),
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Subjects
CREATE TABLE IF NOT EXISTS academic.subjects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    description TEXT,
    credit_hours INTEGER NOT NULL DEFAULT 0,
    curriculum_id UUID REFERENCES academic.curriculums(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Classes
CREATE TABLE IF NOT EXISTS academic.classes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) NOT NULL,  -- Contoh: X IPA 1
    grade INTEGER NOT NULL,  -- 10, 11, 12
    academic_year_id UUID REFERENCES academic.academic_years(id),
    capacity INTEGER NOT NULL DEFAULT 30,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, academic_year_id)
);

-- =============================================================================
-- STUDENT SCHEMA
-- =============================================================================

-- Students
CREATE TABLE IF NOT EXISTS student.students (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    nis VARCHAR(20) NOT NULL UNIQUE,  -- Nomor Induk Siswa
    nisn VARCHAR(20) UNIQUE,  -- NISN Nasional
    full_name VARCHAR(200) NOT NULL,
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    birth_date DATE,
    gender VARCHAR(10) CHECK (gender IN ('M', 'F')),
    address TEXT,
    photo_url VARCHAR(500),
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'GRADUATED', 'DROPPED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Student Enrollments
CREATE TABLE IF NOT EXISTS student.enrollments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES student.students(id) ON DELETE CASCADE,
    class_id UUID NOT NULL REFERENCES academic.classes(id),
    academic_year_id UUID NOT NULL REFERENCES academic.academic_years(id),
    enrollment_date DATE NOT NULL DEFAULT CURRENT_DATE,
    status VARCHAR(20) DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'INACTIVE', 'TRANSFERRED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(student_id, academic_year_id)
);

-- Guardians
CREATE TABLE IF NOT EXISTS student.guardians (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    student_id UUID NOT NULL REFERENCES student.students(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    relationship VARCHAR(50) NOT NULL,  -- Father, Mother, Guardian
    phone VARCHAR(20),
    email VARCHAR(100),
    occupation VARCHAR(100),
    address TEXT,
    is_primary BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================================================
-- ATTENDANCE SCHEMA
-- =============================================================================

-- Attendance Sessions
CREATE TABLE IF NOT EXISTS attendance.sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL,
    class_id UUID NOT NULL REFERENCES academic.classes(id),
    subject_id UUID REFERENCES academic.subjects(id),
    teacher_id UUID REFERENCES student.students(id),  -- Simplified, should reference teachers table
    status VARCHAR(20) DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'SUBMITTED', 'CLOSED')),
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, class_id, subject_id)
);

-- Attendance Records
CREATE TABLE IF NOT EXISTS attendance.records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID NOT NULL REFERENCES attendance.sessions(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES student.students(id),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PRESENT', 'ABSENT', 'SICK', 'PERMISSION')),
    note TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(session_id, student_id)
);

-- =============================================================================
-- ASSESSMENT SCHEMA
-- =============================================================================

-- Assessments
CREATE TABLE IF NOT EXISTS assessment.assessments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('EXAM', 'QUIZ', 'ASSIGNMENT', 'PRACTICE')),
    subject_id UUID NOT NULL REFERENCES academic.subjects(id),
    class_id UUID NOT NULL REFERENCES academic.classes(id),
    max_score DECIMAL(5,2) DEFAULT 100.00,
    weight DECIMAL(5,2) DEFAULT 1.00,
    assessment_date DATE NOT NULL,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Grades
CREATE TABLE IF NOT EXISTS assessment.grades (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    assessment_id UUID NOT NULL REFERENCES assessment.assessments(id) ON DELETE CASCADE,
    student_id UUID NOT NULL REFERENCES student.students(id),
    score DECIMAL(5,2) NOT NULL,
    predicate VARCHAR(5),  -- A, B, C, D, E
    note TEXT,
    graded_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(assessment_id, student_id)
);

-- =============================================================================
-- INDEXES - Performance Optimization
-- =============================================================================

-- Student indexes
CREATE INDEX IF NOT EXISTS idx_students_nis ON student.students(nis);
CREATE INDEX IF NOT EXISTS idx_students_nisn ON student.students(nisn);
CREATE INDEX IF NOT EXISTS idx_students_status ON student.students(status);

-- Enrollment indexes
CREATE INDEX IF NOT EXISTS idx_enrollments_student ON student.enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_class ON student.enrollments(class_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_year ON student.enrollments(academic_year_id);

-- Attendance indexes
CREATE INDEX IF NOT EXISTS idx_attendance_sessions_date ON attendance.sessions(date);
CREATE INDEX IF NOT EXISTS idx_attendance_sessions_class ON attendance.sessions(class_id);
CREATE INDEX IF NOT EXISTS idx_attendance_records_session ON attendance.records(session_id);
CREATE INDEX IF NOT EXISTS idx_attendance_records_student ON attendance.records(student_id);

-- Assessment indexes
CREATE INDEX IF NOT EXISTS idx_assessments_subject ON assessment.assessments(subject_id);
CREATE INDEX IF NOT EXISTS idx_assessments_class ON assessment.assessments(class_id);
CREATE INDEX IF NOT EXISTS idx_grades_assessment ON assessment.grades(assessment_id);
CREATE INDEX IF NOT EXISTS idx_grades_student ON assessment.grades(student_id);

-- =============================================================================
-- SEED DATA - Initial Data for Development
-- =============================================================================

-- Insert sample academic year
INSERT INTO academic.academic_years (year, start_date, end_date, is_active)
VALUES ('2025/2026', '2025-07-01', '2026-06-30', TRUE)
ON CONFLICT (year) DO NOTHING;

-- Insert sample curriculum
INSERT INTO academic.curriculums (name, code, description, academic_year_id, is_active)
VALUES (
    'Kurikulum Merdeka',
    'KM-2025',
    'Kurikulum Merdeka untuk SMK',
    (SELECT id FROM academic.academic_years WHERE year = '2025/2026'),
    TRUE
)
ON CONFLICT (code) DO NOTHING;

-- Insert sample subjects
INSERT INTO academic.subjects (name, code, description, credit_hours, curriculum_id)
VALUES
    ('Matematika', 'MTK', 'Matematika Umum', 4, (SELECT id FROM academic.curriculums WHERE code = 'KM-2025')),
    ('Bahasa Indonesia', 'BIND', 'Bahasa Indonesia', 3, (SELECT id FROM academic.curriculums WHERE code = 'KM-2025')),
    ('Bahasa Inggris', 'BING', 'Bahasa Inggris', 3, (SELECT id FROM academic.curriculums WHERE code = 'KM-2025')),
    ('Pemrograman Web', 'PMW', 'Pemrograman Web Dasar', 6, (SELECT id FROM academic.curriculums WHERE code = 'KM-2025'))
ON CONFLICT (code) DO NOTHING;

-- Insert sample classes
INSERT INTO academic.classes (name, grade, academic_year_id, capacity)
VALUES
    ('X IPA 1', 10, (SELECT id FROM academic.academic_years WHERE year = '2025/2026'), 30),
    ('X IPA 2', 10, (SELECT id FROM academic.academic_years WHERE year = '2025/2026'), 30),
    ('XI RPL 1', 11, (SELECT id FROM academic.academic_years WHERE year = '2025/2026'), 25)
ON CONFLICT (name, academic_year_id) DO NOTHING;

-- =============================================================================
-- VERIFICATION
-- =============================================================================

-- Verify schemas created
DO $$
BEGIN
    RAISE NOTICE 'Database initialization completed successfully!';
    RAISE NOTICE 'Schemas created: academic, student, attendance, assessment, teaching, reporting';
    RAISE NOTICE 'Extensions enabled: uuid-ossp, pg_trgm';
END $$;