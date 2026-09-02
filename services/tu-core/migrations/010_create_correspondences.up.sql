-- +migrate Up
CREATE TABLE IF NOT EXISTS correspondences (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    type VARCHAR(20) NOT NULL, -- INCOMING, OUTGOING
    number VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    subject VARCHAR(200) NOT NULL,
    from_party VARCHAR(200) NOT NULL,
    to_party VARCHAR(200) NOT NULL,
    description TEXT,
    attachment_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT', -- DRAFT, PROCESSED, ARCHIVED
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_correspondences_school_period ON correspondences(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_correspondences_type ON correspondences(type);
CREATE INDEX IF NOT EXISTS idx_correspondences_number ON correspondences(number);
CREATE INDEX IF NOT EXISTS idx_correspondences_date ON correspondences(date);
CREATE INDEX IF NOT EXISTS idx_correspondences_status ON correspondences(status);
