-- Tabel Absensi Harian Guru (PTK)
CREATE TABLE IF NOT EXISTS ptk_attendance (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'HADIR', -- HADIR, SAKIT, IZIN, ALPHA, CUTI, DINAS_LUAR
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(ptk_id, date)
);
CREATE INDEX idx_ptk_attendance_date ON ptk_attendance(date);

-- Tabel Laporan Triwulanan (Placeholder untuk Fase Lanjutan)
CREATE TABLE IF NOT EXISTS quarterly_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    year INT NOT NULL,
    quarter INT NOT NULL,
    report_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    file_url TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Tabel Preset Filter
CREATE TABLE IF NOT EXISTS statistic_filter_presets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(100) NOT NULL,
    filter_json JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed Data Kehadiran (Bulan Ini) agar Chart langsung berisi data
INSERT INTO ptk_attendance (ptk_id, date, status) VALUES 
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', CURRENT_DATE - INTERVAL '1 day', 'HADIR'),
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', CURRENT_DATE - INTERVAL '2 day', 'HADIR'),
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', CURRENT_DATE - INTERVAL '3 day', 'HADIR'),

('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', CURRENT_DATE - INTERVAL '1 day', 'HADIR'),
('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', CURRENT_DATE - INTERVAL '2 day', 'IZIN'),
('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', CURRENT_DATE - INTERVAL '3 day', 'HADIR'),

('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', CURRENT_DATE - INTERVAL '1 day', 'HADIR'),
('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', CURRENT_DATE - INTERVAL '2 day', 'HADIR'),
('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', CURRENT_DATE - INTERVAL '3 day', 'ALPHA')
ON CONFLICT DO NOTHING;
