-- Seed contoh surat untuk Kelompok Hubungan Industri (PKL/MoU) dan Orang Tua/Kesiswaan
-- UUID statis agar idempotent dan mudah dilacak.
INSERT INTO correspondences (
    id, school_id, academic_period_id, type, jenis_surat, number, date, subject,
    from_party, to_party, description, status, created_at, updated_at
) VALUES
-- 1. Surat Permohonan PKL (Keluar)
('dddddddd-dddd-dddd-dddd-000000000001', '11111111-1111-1111-1111-111111111111', '20261', 'OUTGOING', 'PERMOHONAN_PKL', '001/PERMOHONAN_PKL/SMK.1/X/2026', '2026-10-15', 'Permohonan Tempat Praktik Kerja Lapangan', 'SMK Pasundan Jatinangor', 'PT Teknologi Maju', 'Mohon kesediaan menerima siswa PKL jurusan TKJ', 'PROCESSED', NOW(), NOW()),

-- 2. Surat Perjanjian MoU (Masuk)
('dddddddd-dddd-dddd-dddd-000000000002', '11111111-1111-1111-1111-111111111111', '20261', 'INCOMING', 'MOU', '002/MOU/PT-TM/X/2026', '2026-10-20', 'Penandatanganan MoU Kerjasama Pendidikan', 'PT Teknologi Maju', 'SMK Pasundan Jatinangor', 'Kerjasama sinkronisasi kurikulum dan PKL', 'ARCHIVED', NOW(), NOW()),

-- 3. Surat Undangan Orang Tua (Keluar)
('dddddddd-dddd-dddd-dddd-000000000003', '11111111-1111-1111-1111-111111111111', '20261', 'OUTGOING', 'UNDANGAN_ORTU', '003/UNDANGAN_ORTU/SMK.1/X/2026', '2026-10-25', 'Undangan Rapat Pleno Komite Sekolah', 'SMK Pasundan Jatinangor', 'Orang Tua/Wali Kelas X', 'Menghadiri rapat pleno pembagian rapor tengah semester', 'DRAFT', NOW(), NOW()),

-- 4. Surat Keterangan Lulus (SKL) (Keluar)
('dddddddd-dddd-dddd-dddd-000000000004', '11111111-1111-1111-1111-111111111111', '20261', 'OUTGOING', 'SKL', '004/SKL/SMK.1/XI/2026', '2026-11-01', 'Surat Keterangan Lulus (Sementara)', 'SMK Pasundan Jatinangor', 'Siswa Ybs', 'Menerangkan bahwa siswa telah lulus dan menunggu ijazah asli', 'PROCESSED', NOW(), NOW()),

-- 5. Surat Panggilan Orang Tua (Keluar)
('dddddddd-dddd-dddd-dddd-000000000005', '11111111-1111-1111-1111-111111111111', '20261', 'OUTGOING', 'PANGGILAN_ORTU', '005/PANGGILAN_ORTU/SMK.1/XI/2026', '2026-11-05', 'Panggilan Orang Tua/Wali Siswa', 'SMK Pasundan Jatinangor', 'Orang Tua dari Budi', 'Adanya pelanggaran tata tertib yang perlu diselesaikan', 'DRAFT', NOW(), NOW()),

-- 6. Surat Rekomendasi Lomba (Keluar)
('dddddddd-dddd-dddd-dddd-000000000006', '11111111-1111-1111-1111-111111111111', '20261', 'OUTGOING', 'REKOMENDASI_LOMBA', '006/REKOMENDASI_LOMBA/SMK.1/XI/2026', '2026-11-06', 'Rekomendasi Mengikuti LKS Tingkat Provinsi', 'SMK Pasundan Jatinangor', 'Panitia LKS', 'Merekomendasikan siswa terbaik untuk berpartisipasi dalam lomba', 'PROCESSED', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Tambahkan counter agar auto-number berikutnya tidak bentrok
INSERT INTO correspondence_counters (id, school_id, academic_period_id, jenis_surat, last_number)
VALUES 
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', '20261', 'PERMOHONAN_PKL', 1),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', '20261', 'UNDANGAN_ORTU', 1),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', '20261', 'SKL', 1),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', '20261', 'PANGGILAN_ORTU', 1),
(gen_random_uuid(), '11111111-1111-1111-1111-111111111111', '20261', 'REKOMENDASI_LOMBA', 1)
ON CONFLICT (school_id, academic_period_id, jenis_surat) DO NOTHING;
