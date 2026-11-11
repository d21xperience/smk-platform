-- Seed Data untuk Pengujian Modul PSB
-- Menggunakan school_id statis Single-Tenant sesuai Hard Rule #19
-- Menggunakan academic_period_id 20261 (asumsi aktif)

INSERT INTO "registrations" (
    "id", "school_id", "academic_period_id", "calon_nisn", "nis", "calon_nama", 
    "calon_jenis_kelamin", "calon_tempat_lahir", "calon_tanggal_lahir", "asal_sekolah", 
    "jurusan_dipilih", "biaya_standar", "biaya_pendaftaran", 
    "berkas_ijazah", "berkas_kk", "berkas_akta", "status", "created_at", "updated_at"
) VALUES
-- Skenario 1: Lunas & Berkas Lengkap (Siap Aktivasi)
(
    'a1111111-1111-1111-1111-111111111111', 
    '11111111-1111-1111-1111-111111111111', 
    '20261', 
    '0091122331', 
    '', 
    'Bagus Setiawan', 
    'L', 
    'Bandung', 
    '2007-05-12', 
    'SMPN 1 Soreang', 
    'TKRO', 
    450000, 
    450000, 
    TRUE, TRUE, TRUE, 
    'REGISTERED', 
    NOW(), NOW()
),
-- Skenario 2: Dicicil (Belum Lunas, Berkas Lengkap)
(
    'a2222222-2222-2222-2222-222222222222', 
    '11111111-1111-1111-1111-111111111111', 
    '20261', 
    '0091122332', 
    '', 
    'Aditya Pratama', 
    'L', 
    'Bandung', 
    '2007-08-20', 
    'SMPN 2 Bandung', 
    'TKRO', 
    450000, 
    200000, 
    TRUE, TRUE, TRUE, 
    'REGISTERED', 
    NOW(), NOW()
),
-- Skenario 3: Belum Bayar & Berkas Kurang
(
    'a3333333-3333-3333-3333-333333333333', 
    '11111111-1111-1111-1111-111111111111', 
    '20261', 
    '0091122333', 
    '', 
    'Chandra Kirana', 
    'P', 
    'Cimahi', 
    '2007-11-05', 
    'MTs As-Syifa', 
    'RPL', 
    450000, 
    0, 
    FALSE, FALSE, FALSE, 
    'REGISTERED', 
    NOW(), NOW()
),
-- Skenario 4: Lunas, Berkas Kurang (Tidak Bisa Aktivasi)
(
    'a4444444-4444-4444-4444-444444444444', 
    '11111111-1111-1111-1111-111111111111', 
    '20261', 
    '0091122334', 
    '', 
    'Deni Sukma', 
    'L', 
    'Garut', 
    '2007-03-15', 
    'SMP Pasundan 1', 
    'AKL', 
    450000, 
    450000, 
    FALSE, TRUE, FALSE, 
    'REGISTERED', 
    NOW(), NOW()
),
-- Skenario 5: Sudah Punya NIS (Simulasi setelah Generate NIS Massal)
(
    'a5555555-5555-5555-5555-555555555555', 
    '11111111-1111-1111-1111-111111111111', 
    '20261', 
    '0091122335', 
    '262700001', 
    'Eka Putri', 
    'P', 
    'Sumedang', 
    '2007-07-07', 
    'SMPN 3 Sumedang', 
    'RPL', 
    450000, 
    450000, 
    TRUE, TRUE, TRUE, 
    'REGISTERED', 
    NOW(), NOW()
)
ON CONFLICT ("calon_nisn", "academic_period_id") DO NOTHING;

-- Sinkronisasi sequence agar NIS Massal berikutnya melanjutkan dari nomor terakhir
INSERT INTO "nis_sequences" ("academic_period_id", "prefix", "last_sequence") 
VALUES ('20261', '2627', 1)
ON CONFLICT ("academic_period_id") DO UPDATE SET "last_sequence" = 1;
