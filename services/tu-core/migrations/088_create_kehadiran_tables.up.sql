-- 1. Upgrade Tabel ptk_attendance (dari Sub-Modul 3)
ALTER TABLE ptk_attendance 
ADD COLUMN IF NOT EXISTS jam_masuk TIME,
ADD COLUMN IF NOT EXISTS jam_pulang TIME,
ADD COLUMN IF NOT EXISTS durasi_kerja INT, -- dalam menit
ADD COLUMN IF NOT EXISTS menit_terlambat INT DEFAULT 0,
ADD COLUMN IF NOT EXISTS source_data VARCHAR(50) DEFAULT 'MANUAL', -- FINGERPRINT, MOBILE, MANUAL
ADD COLUMN IF NOT EXISTS is_corrected BOOLEAN DEFAULT FALSE,
ADD COLUMN IF NOT EXISTS correction_id UUID;

-- 2. Tabel Master Jam Kerja
CREATE TABLE IF NOT EXISTS master_jam_kerja (
    id SERIAL PRIMARY KEY,
    jenis_pegawai VARCHAR(20) NOT NULL, -- GURU, TENDIK
    hari INT NOT NULL, -- 1=Senin, 2=Selasa, ..., 6=Sabtu
    jam_masuk TIME NOT NULL,
    jam_pulang TIME NOT NULL,
    toleransi_menit INT DEFAULT 15
);

-- Seed Master Jam Kerja (Contoh: Guru masuk 07:00 pulang 14:00, Tendik masuk 07:00 pulang 15:00)
INSERT INTO master_jam_kerja (jenis_pegawai, hari, jam_masuk, jam_pulang, toleransi_menit) VALUES
('GURU', 1, '07:00', '14:00', 15), ('GURU', 2, '07:00', '14:00', 15),
('GURU', 3, '07:00', '14:00', 15), ('GURU', 4, '07:00', '14:00', 15),
('GURU', 5, '07:00', '11:30', 15), ('GURU', 6, '08:00', '12:00', 15),
('TENDIK', 1, '07:00', '15:00', 15), ('TENDIK', 2, '07:00', '15:00', 15),
('TENDIK', 3, '07:00', '15:00', 15), ('TENDIK', 4, '07:00', '15:00', 15),
('TENDIK', 5, '07:00', '11:30', 15), ('TENDIK', 6, '08:00', '13:00', 15)
ON CONFLICT DO NOTHING;

-- 3. Tabel Koreksi Kehadiran
CREATE TABLE IF NOT EXISTS ptk_attendance_corrections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    tanggal DATE NOT NULL,
    status_lama VARCHAR(30),
    status_koreksi VARCHAR(30) NOT NULL,
    jam_masuk_lama TIME,
    jam_masuk_koreksi TIME,
    jam_pulang_lama TIME,
    jam_pulang_koreksi TIME,
    alasan TEXT NOT NULL,
    bukti_url TEXT,
    approval_level VARCHAR(20) NOT NULL, -- KTU, KEPSEK
    status VARCHAR(30) NOT NULL DEFAULT 'PENDING_KTU',
    requested_by UUID,
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    approved_by UUID,
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT
);

-- 4. Tabel Rekap Bulanan
CREATE TABLE IF NOT EXISTS ptk_rekap_bulanan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bulan INT NOT NULL,
    tahun INT NOT NULL,
    total_hari_kerja INT,
    total_pegawai INT,
    rata_persentase NUMERIC(5,2),
    status VARCHAR(30) DEFAULT 'GENERATING',
    excel_url TEXT,
    pdf_url TEXT,
    generated_by UUID,
    generated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(bulan, tahun)
);

CREATE TABLE IF NOT EXISTS ptk_rekap_bulanan_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    rekap_id UUID NOT NULL REFERENCES ptk_rekap_bulanan(id) ON DELETE CASCADE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    total_hadir INT DEFAULT 0,
    total_terlambat INT DEFAULT 0,
    total_sakit INT DEFAULT 0,
    total_izin INT DEFAULT 0,
    total_alpha INT DEFAULT 0,
    total_dinas_luar INT DEFAULT 0,
    persentase_kehadiran NUMERIC(5,2),
    total_menit_terlambat INT DEFAULT 0
);
