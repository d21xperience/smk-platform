-- 1. Tabel Master Golongan & Pangkat (Standar PNS)
CREATE TABLE IF NOT EXISTS master_golongan (
    id INT PRIMARY KEY, -- Kode Golongan: 11, 12, 13, 14 (I/a - I/d), 21, 22... (II/a - II/d), dll.
    kode_golongan VARCHAR(10) NOT NULL,
    nama_pangkat VARCHAR(100) NOT NULL,
    interval_kgb_tahun INT NOT NULL DEFAULT 4,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Seed Data Standar PNS
INSERT INTO master_golongan (id, kode_golongan, nama_pangkat, interval_kgb_tahun) VALUES
(11, 'I/a', 'Juru Muda', 2), (12, 'I/b', 'Juru Muda Tingkat I', 2),
(13, 'I/c', 'Juru', 2), (14, 'I/d', 'Juru Tingkat I', 2),
(21, 'II/a', 'Pengatur Muda', 2), (22, 'II/b', 'Pengatur Muda Tingkat I', 2),
(23, 'II/c', 'Pengatur', 2), (24, 'II/d', 'Pengatur Tingkat I', 2),
(31, 'III/a', 'Penata Muda', 4), (32, 'III/b', 'Penata Muda Tingkat I', 4),
(33, 'III/c', 'Penata', 4), (34, 'III/d', 'Penata Tingkat I', 4),
(41, 'IV/a', 'Pembina', 4), (42, 'IV/b', 'Pembina Tingkat I', 4),
(43, 'IV/c', 'Pembina Utama Muda', 4), (44, 'IV/d', 'Pembina Utama Madya', 4),
(45, 'IV/e', 'Pembina Utama', 4)
ON CONFLICT (id) DO NOTHING;

-- 2. Tambahkan kolom baru ke ekstensi Buku Induk PTK
ALTER TABLE ptk_buku_induk_ext 
ADD COLUMN IF NOT EXISTS golongan_id INT REFERENCES master_golongan(id),
ADD COLUMN IF NOT EXISTS tmt_pangkat DATE,
ADD COLUMN IF NOT EXISTS tmt_gaji_berkala DATE;

-- Hapus kolom lama yang berupa string bebas jika sudah tidak dipakai (opsional, kita biarkan dulu agar aman)
-- ALTER TABLE ptk_buku_induk_ext DROP COLUMN IF EXISTS pangkat_golongan;

-- 3. Tabel Buku Kontrol
CREATE TABLE IF NOT EXISTS buku_kontrol (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE UNIQUE,
    golongan_sekarang_id INT REFERENCES master_golongan(id),
    tmt_pangkat_sekarang DATE,
    tmt_gaji_berkala_sekarang DATE,
    tanggal_kgb_berikut DATE,
    tanggal_pangkat_berikut DATE,
    tanggal_pensiun DATE,
    status_monitoring VARCHAR(30) DEFAULT 'AMAN', -- AMAN, PERLU_PERSIAPAN, SEGERA, URGENT, OVERDUE
    last_calculated_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Tabel Riwayat Kenaikan (Log)
CREATE TABLE IF NOT EXISTS riwayat_kenaikan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    tanggal_efektif DATE NOT NULL,
    jenis_kenaikan VARCHAR(50) NOT NULL, -- KGB, PANGKAT
    golongan_sebelum_id INT REFERENCES master_golongan(id),
    golongan_sesudah_id INT REFERENCES master_golongan(id),
    nomor_sk VARCHAR(100),
    keterangan TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 5. Seed Data Buku Kontrol (Untuk 3 Guru Existing sebagai contoh)
-- Mengisi TMT Palsu agar bisa diuji
UPDATE ptk_buku_induk_ext SET golongan_id = 33, tmt_pangkat = '2022-04-01', tmt_gaji_berkala = '2022-04-01' WHERE ptk_id = 'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d';
UPDATE ptk_buku_induk_ext SET golongan_id = 32, tmt_pangkat = '2023-01-01', tmt_gaji_berkala = '2023-01-01' WHERE ptk_id = 'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e';
UPDATE ptk_buku_induk_ext SET golongan_id = 31, tmt_pangkat = '2024-05-01', tmt_gaji_berkala = '2024-05-01' WHERE ptk_id = 'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f';

-- Trigger dummy untuk inisialisasi (akan di-handle oleh Go Service nanti)
INSERT INTO buku_kontrol (ptk_id, golongan_sekarang_id, tmt_pangkat_sekarang, tmt_gaji_berkala_sekarang, tanggal_kgb_berikut, tanggal_pangkat_berikut, tanggal_pensiun, status_monitoring)
VALUES 
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 33, '2022-04-01', '2022-04-01', '2026-04-01', '2026-04-01', '2045-01-01', 'PERLU_PERSIAPAN'),
('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 32, '2023-01-01', '2023-01-01', '2027-01-01', '2027-01-01', '2050-02-02', 'AMAN'),
('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 31, '2024-05-01', '2024-05-01', '2028-05-01', '2028-05-01', '2055-03-03', 'AMAN')
ON CONFLICT (ptk_id) DO NOTHING;
