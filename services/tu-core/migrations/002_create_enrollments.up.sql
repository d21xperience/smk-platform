-- + migrate Up
CREATE TABLE IF NOT EXISTS enrollments (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL REFERENCES students(id) ON DELETE CASCADE,
    school_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(10) NOT NULL,
    class_id VARCHAR(36) NOT NULL,
    enrollment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_enrollments_student ON enrollments(student_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_school_period ON enrollments(school_id, period_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_class ON enrollments(class_id);


