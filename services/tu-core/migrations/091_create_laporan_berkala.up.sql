-- 1. Master Penerima Laporan
CREATE TABLE IF NOT EXISTS laporan_recipients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_instansi VARCHAR(100) NOT NULL,
    tipe VARCHAR(20) NOT NULL, -- 'DINAS', 'YAYASAN'
    email VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO laporan_recipients (nama_instansi, tipe, email) VALUES 
('Dinas Pendidikan', 'DINAS', 'dinas@example.go.id'),
('Yayasan Pendidikan', 'YAYASAN', 'info@yayasan-example.org');

-- 2. Laporan Berkala (Header)
CREATE TABLE IF NOT EXISTS laporan_berkala (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id UUID NOT NULL,
    tahun INT NOT NULL,
    triwulan INT NOT NULL, -- 1, 2, 3, 4
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT', -- DRAFT, REVIEWED, APPROVED, SENT, PARTIAL_SENT, ARCHIVED, REJECTED
    data_json JSONB, -- Snapshot data agregat
    excel_dinas_url TEXT,
    excel_yayasan_url TEXT,
    pdf_url TEXT,
    ringkasan_url TEXT,
    surat_url TEXT,
    pdf_final_url TEXT,
    qr_url TEXT,
    catatan_ktu TEXT,
    catatan_approve TEXT,
    reviewed_by UUID,
    approved_by UUID,
    sent_by UUID,
    sent_at TIMESTAMPTZ,
    confirmed_at TIMESTAMPTZ,
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT unique_laporan_periode UNIQUE (school_id, tahun, triwulan)
);

-- 3. Timeline Audit
CREATE TABLE IF NOT EXISTS laporan_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    laporan_id UUID NOT NULL REFERENCES laporan_berkala(id) ON DELETE CASCADE,
    status_dari VARCHAR(30),
    status_ke VARCHAR(30) NOT NULL,
    action_by UUID NOT NULL,
    action_at TIMESTAMPTZ DEFAULT NOW(),
    catatan TEXT
);

CREATE INDEX idx_laporan_timeline_id ON laporan_timeline(laporan_id);
