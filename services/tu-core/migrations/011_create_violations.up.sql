CREATE TABLE IF NOT EXISTS violations (
    id SERIAL PRIMARY KEY,
    school_id VARCHAR(50) NOT NULL,
    academic_period_id VARCHAR(50) NOT NULL,
    student_id VARCHAR(50) NOT NULL,
    student_name VARCHAR(255) NOT NULL,
    class_name VARCHAR(100) NOT NULL,
    title TEXT NOT NULL,
    category VARCHAR(50) NOT NULL,
    points INT NOT NULL DEFAULT 0,
    violation_date TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_violations_school_period ON violations(school_id, academic_period_id);
CREATE INDEX idx_violations_student ON violations(student_id);
CREATE INDEX idx_violations_date ON violations(violation_date);
