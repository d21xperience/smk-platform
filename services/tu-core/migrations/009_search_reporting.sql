-- ============================================
-- MIGRATION 009: Search Index & Reporting Tables
-- ============================================

-- 1. Tabel Outbox Events (Transactional Outbox Pattern)
CREATE TABLE IF NOT EXISTS outbox_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    aggregate_type VARCHAR(50) NOT NULL,
    aggregate_id UUID NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    payload JSONB NOT NULL,
    correlation_id VARCHAR(100),
    trace_id VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_events_unpublished ON outbox_events(published_at)
WHERE published_at IS NULL;

-- 2. Tabel Search Index (Full Text Search dengan PostgreSQL)
CREATE TABLE IF NOT EXISTS students_search_index (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL,
    academic_period_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    class_name VARCHAR(100),
    search_vector tsvector,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Generate search vector dari name (weight A) dan class_name (weight B)
CREATE INDEX idx_students_search_vector ON students_search_index USING GIN(search_vector);
CREATE INDEX idx_students_search_school ON students_search_index(school_id, academic_period_id);

-- Trigger untuk auto-update search_vector
CREATE OR REPLACE FUNCTION update_students_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('indonesian', COALESCE(NEW.name, '')), 'A') ||
        setweight(to_tsvector('indonesian', COALESCE(NEW.class_name, '')), 'B');
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_students_search_vector
BEFORE INSERT OR UPDATE ON students_search_index
FOR EACH ROW EXECUTE FUNCTION update_students_search_vector();

-- 3. Tabel Reporting Attendance Daily (Read Model)
CREATE TABLE IF NOT EXISTS reporting_attendance_daily (
    school_id UUID NOT NULL,
    academic_period_id UUID NOT NULL,
    class_id UUID NOT NULL,
    date DATE NOT NULL,
    total_students INT NOT NULL DEFAULT 0,
    present_count INT NOT NULL DEFAULT 0,
    absent_count INT NOT NULL DEFAULT 0,
    late_count INT NOT NULL DEFAULT 0,
    attendance_rate DECIMAL(5,2) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (school_id, academic_period_id, class_id, date)
);

CREATE INDEX idx_reporting_attendance_period ON reporting_attendance_daily(school_id, academic_period_id, date);

-- 4. Tabel Reporting Finance Summary (Read Model)
CREATE TABLE IF NOT EXISTS reporting_finance_monthly (
    school_id UUID NOT NULL,
    academic_period_id UUID NOT NULL,
    month INT NOT NULL, -- 1-12
    year INT NOT NULL,
    component_type VARCHAR(50) NOT NULL, -- 'SPP', 'UNIFORM', dll
    total_billed DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_paid DECIMAL(15,2) NOT NULL DEFAULT 0,
    total_outstanding DECIMAL(15,2) NOT NULL DEFAULT 0,
    payment_count INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (school_id, academic_period_id, month, year, component_type)
);

-- 5. Tabel Notification Log (untuk audit dan debugging)
CREATE TABLE IF NOT EXISTS notification_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    correlation_id VARCHAR(100),
    channel VARCHAR(50) NOT NULL, -- 'email', 'whatsapp', 'firebase'
    recipient VARCHAR(255) NOT NULL,
    template_code VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL, -- 'SENT', 'FAILED', 'PENDING'
    error_message TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_notification_logs_correlation ON notification_logs(correlation_id);
CREATE INDEX idx_notification_logs_school ON notification_logs(school_id, created_at);
