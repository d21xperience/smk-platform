CREATE TABLE IF NOT EXISTS statistic_filter_presets (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    filter_data JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quarterly_reports (
    id UUID PRIMARY KEY,
    school_id UUID NOT NULL,
    year INT NOT NULL,
    quarter INT NOT NULL,
    report_type VARCHAR(20) NOT NULL,
    file_url TEXT,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
