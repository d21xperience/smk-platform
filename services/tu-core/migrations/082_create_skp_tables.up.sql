-- Hapus tabel lama yang tidak sesuai standar Permenpan 6/2022
DROP TABLE IF EXISTS ptk_performances CASCADE;

-- 1. Tabel Template Penilaian
CREATE TABLE IF NOT EXISTS penilaian_template (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jenis_pegawai VARCHAR(50) NOT NULL, -- GURU_ASN, TENDIK_ASN, NON_ASN
    nama_template VARCHAR(100) NOT NULL,
    bobot_perilaku INT NOT NULL DEFAULT 50,
    bobot_hasil_utama INT NOT NULL DEFAULT 50,
    bobot_hasil_tambahan INT NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);

-- 2. Tabel Periode Penilaian
CREATE TABLE IF NOT EXISTS penilaian_periode (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nama_periode VARCHAR(100) NOT NULL,
    jenis_periode VARCHAR(20) NOT NULL, -- BULANAN, SEMESTERAN, TAHUNAN
    tanggal_mulai DATE NOT NULL,
    tanggal_selesai DATE NOT NULL,
    deadline_pengisian DATE NOT NULL,
    tahun_referensi INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN', -- OPEN, CLOSED, OVERDUE
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 3. Tabel Penilaian Pegawai (Header)
CREATE TABLE IF NOT EXISTS penilaian_pegawai (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    periode_id UUID NOT NULL REFERENCES penilaian_periode(id) ON DELETE CASCADE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    template_id UUID NOT NULL REFERENCES penilaian_template(id) ON DELETE RESTRICT,
    penilai_utama_id UUID REFERENCES tabel_ptk(ptk_id), 
    penilai_berikat_id UUID REFERENCES tabel_ptk(ptk_id),
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING', -- PENDING, DRAFT, SUBMITTED, FINAL
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(periode_id, ptk_id)
);

-- 4. Tabel Hasil Akhir Penilaian
CREATE TABLE IF NOT EXISTS penilaian_hasil (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    penilaian_id UUID NOT NULL UNIQUE REFERENCES penilaian_pegawai(id) ON DELETE CASCADE,
    nilai_perilaku NUMERIC(5,2) NOT NULL,
    predikat_perilaku VARCHAR(20) NOT NULL,
    nilai_hasil_utama NUMERIC(5,2) NOT NULL,
    nilai_hasil_tambahan NUMERIC(5,2) DEFAULT 0,
    nilai_akhir NUMERIC(5,2) NOT NULL,
    predikat_akhir VARCHAR(20) NOT NULL,
    catatan_penilai TEXT,
    catatan_berikat TEXT,
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Tabel Detail Penilaian per Aspek
CREATE TABLE IF NOT EXISTS penilaian_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    penilaian_id UUID NOT NULL REFERENCES penilaian_pegawai(id) ON DELETE CASCADE,
    kelompok VARCHAR(30) NOT NULL, -- PERILAKU, HASIL_UTAMA, HASIL_TAMBAHAN
    aspek VARCHAR(100) NOT NULL,
    nilai INT NOT NULL
);

-- INDEXES
CREATE INDEX idx_penilaian_pegawai_periode ON penilaian_pegawai(periode_id);
CREATE INDEX idx_penilaian_pegawai_ptk ON penilaian_pegawai(ptk_id);
CREATE INDEX idx_penilaian_pegawai_penilai ON penilaian_pegawai(penilai_utama_id);

-- SEED DATA TEMPLATE
INSERT INTO penilaian_template (id, jenis_pegawai, nama_template, bobot_perilaku, bobot_hasil_utama, bobot_hasil_tambahan) VALUES
('11111111-1111-1111-1111-111111111111', 'GURU_ASN', 'Template Guru ASN (Permenpan 6/2022)', 50, 30, 20),
('22222222-2222-2222-2222-222222222222', 'TENDIK_ASN', 'Template Tendik ASN (Pelaksana)', 50, 50, 0),
('33333333-3333-3333-3333-333333333333', 'NON_ASN', 'Template Guru/Tendik Non-ASN', 40, 60, 0)
ON CONFLICT (id) DO NOTHING;
