-- =====================================================================
-- Seed: Transaksi Inti Keuangan
-- Adaptif — hanya insert kalau tabel terkait kosong
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. Seed tagihan_siswa (3 tagihan SPP untuk 3 siswa pertama)
-- ---------------------------------------------------------------------
INSERT INTO "tagihan_siswa" (
    id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
    jenis_iuran_id, periode_bulan, periode_tahun, nominal, status,
    tanggal_jatuh_tempo, created_by
)
SELECT
    gen_random_uuid(),
    '11111111-1111-1111-1111-111111111111',
    '20261',
    s.peserta_didik_id,
    s.nis,
    COALESCE(
        (SELECT k.nm_kelas
         FROM tabel_anggotakelas ak
         JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
         WHERE ak.peserta_didik_id = s.peserta_didik_id
         LIMIT 1),
        ''
    ),
    'f1000000-0000-0000-0000-000000000001',
    11, 2026, 250000.00, 'BELUM_LUNAS', '2026-11-30',
    '11111111-1111-1111-1111-111111111111'
FROM tabel_siswa s
WHERE NOT EXISTS (SELECT 1 FROM tagihan_siswa)
LIMIT 3;

-- ---------------------------------------------------------------------
-- 2. Seed penerimaan_iuran (3 penerimaan LUNAS dari tagihan di atas)
-- ---------------------------------------------------------------------
INSERT INTO "penerimaan_iuran" (
    id, school_id, academic_period_id, no_bkm, tanggal,
    tagihan_siswa_id, peserta_didik_id, nis, jenis_iuran_id,
    periode_bulan, periode_tahun, jumlah, metode_bayar,
    status, created_by
)
SELECT
    gen_random_uuid(),
    t.school_id,
    t.academic_period_id,
    'BKM/00' || ROW_NUMBER() OVER (ORDER BY t.created_at)::text || '/SMK/XI/2026',
    '2026-11-10',
    t.id,
    t.peserta_didik_id,
    t.nis,
    t.jenis_iuran_id,
    t.periode_bulan,
    t.periode_tahun,
    t.nominal,
    'TRANSFER',
    'LUNAS',
    '11111111-1111-1111-1111-111111111111'
FROM tagihan_siswa t
WHERE NOT EXISTS (SELECT 1 FROM penerimaan_iuran)
LIMIT 3;

-- Update tagihan yang sudah dibayar
UPDATE "tagihan_siswa" t
SET status = 'LUNAS',
    total_dibayar = COALESCE((SELECT SUM(pi.jumlah) FROM penerimaan_iuran pi WHERE pi.tagihan_siswa_id = t.id), 0),
    tanggal_bayar = CURRENT_DATE,
    updated_at = NOW()
WHERE EXISTS (SELECT 1 FROM penerimaan_iuran pi WHERE pi.tagihan_siswa_id = t.id);

-- ---------------------------------------------------------------------
-- 3. Seed buku_kas_umum — DEBET dari penerimaan_iuran
-- ---------------------------------------------------------------------
WITH ordered AS (
    SELECT
        pi.id, pi.school_id, pi.academic_period_id, pi.tanggal, pi.no_bkm,
        pi.jumlah, pi.nis, ji.nama AS jenis_nama,
        ROW_NUMBER() OVER (ORDER BY pi.tanggal, pi.created_at) AS rn,
        SUM(pi.jumlah) OVER (ORDER BY pi.tanggal, pi.created_at
                             ROWS BETWEEN UNBOUNDED PRECEDING AND CURRENT ROW) AS running_saldo
    FROM penerimaan_iuran pi
    JOIN jenis_iuran ji ON ji.id = pi.jenis_iuran_id
)
INSERT INTO "buku_kas_umum" (
    id, school_id, academic_period_id, tanggal, no_bukti, uraian,
    kode_rekening_id, rekening_bank_id, debet, kredit, saldo_berjalan,
    referensi_tabel, referensi_id, created_by
)
SELECT
    gen_random_uuid(),
    o.school_id, o.academic_period_id, o.tanggal, o.no_bkm,
    'Penerimaan ' || o.jenis_nama || ' - NIS ' || o.nis,
    'f3000000-0000-0000-0000-000000000032', -- 4.1.1 Pendapatan SPP
    'f4000000-0000-0000-0000-000000000001', -- BCA
    o.jumlah, 0, o.running_saldo,
    'penerimaan_iuran', o.id,
    '11111111-1111-1111-1111-111111111111'
FROM ordered o
WHERE NOT EXISTS (SELECT 1 FROM buku_kas_umum WHERE referensi_tabel = 'penerimaan_iuran');

-- ---------------------------------------------------------------------
-- 4. Seed pengajuan_pengeluaran (3 status berbeda)
-- ---------------------------------------------------------------------
INSERT INTO "pengajuan_pengeluaran" (
    id, school_id, academic_period_id, nomor_pengajuan, tanggal,
    pemohon_user_id, keperluan, jumlah, prioritas, status,
    catatan_pemohon, verified_by, verified_at, approved_by, approved_at, no_bkk
)
VALUES
(
    'f6000000-0000-0000-0000-000000000001',
    '11111111-1111-1111-1111-111111111111', '20261',
    'PP/001/SMK/XI/2026', '2026-11-01',
    '11111111-1111-1111-1111-111111111111',
    'Pembelian ATK untuk kegiatan ujian semester', 500000.00, 'SEDANG',
    'MENUNGGU_VERIFIKASI_BENDAHARA',
    'Mohon segera diproses', NULL, NULL, NULL, NULL, NULL
),
(
    'f6000000-0000-0000-0000-000000000002',
    '11111111-1111-1111-1111-111111111111', '20261',
    'PP/002/SMK/XI/2026', '2026-11-05',
    '11111111-1111-1111-1111-111111111111',
    'Konsumsi rapat komite sekolah', 750000.00, 'TINGGI',
    'MENUNGGU_PERSETUJUAN_KEPALA',
    'Rapat tanggal 10 November', '11111111-1111-1111-1111-111111111111', '2026-11-06 09:00:00+07', NULL, NULL, NULL
),
(
    'f6000000-0000-0000-0000-000000000003',
    '11111111-1111-1111-1111-111111111111', '20261',
    'PP/003/SMK/XI/2026', '2026-11-08',
    '11111111-1111-1111-1111-111111111111',
    'Perbaikan ringan AC ruang guru', 1200000.00, 'DARURAT',
    'DISETUJUI',
    'AC mati 2 unit', '11111111-1111-1111-1111-111111111111', '2026-11-09 08:00:00+07',
    '11111111-1111-1111-1111-111111111111', '2026-11-09 14:00:00+07', NULL
)
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------
-- 5. Seed pengeluaran (1 baris untuk pengajuan DISETUJUI)
-- ---------------------------------------------------------------------
INSERT INTO "pengeluaran" (
    id, school_id, academic_period_id, no_bkk, tanggal,
    pengajuan_id, jumlah, keperluan, status, created_by
)
SELECT
    gen_random_uuid(),
    p.school_id,
    p.academic_period_id,
    'BKK/001/SMK/XI/2026',
    '2026-11-10',
    p.id,
    p.jumlah,
    p.keperluan,
    'SIAP_CAIR',
    '11111111-1111-1111-1111-111111111111'
FROM pengajuan_pengeluaran p
WHERE p.status = 'DISETUJUI'
AND NOT EXISTS (SELECT 1 FROM pengeluaran);

-- Update no_bkk di pengajuan
UPDATE pengajuan_pengeluaran p
SET no_bkk = (SELECT no_bkk FROM pengeluaran pe WHERE pe.pengajuan_id = p.id),
    updated_at = NOW()
WHERE EXISTS (SELECT 1 FROM pengeluaran pe WHERE pe.pengajuan_id = p.id);

-- ---------------------------------------------------------------------
-- 6. Seed buku_kas_umum — KREDIT dari pengeluaran
-- ---------------------------------------------------------------------
WITH last_saldo AS (
    SELECT COALESCE(MAX(saldo_berjalan), 0) AS s
    FROM buku_kas_umum
)
INSERT INTO "buku_kas_umum" (
    id, school_id, academic_period_id, tanggal, no_bukti, uraian,
    kode_rekening_id, rekening_bank_id, debet, kredit, saldo_berjalan,
    referensi_tabel, referensi_id, created_by
)
SELECT
    gen_random_uuid(),
    pe.school_id, pe.academic_period_id, pe.tanggal, pe.no_bkk,
    'Pengeluaran: ' || pe.keperluan,
    'f3000000-0000-0000-0000-000000000041', -- 5.1 Belanja Operasional
    'f4000000-0000-0000-0000-000000000001',
    0, pe.jumlah,
    ls.s - pe.jumlah,
    'pengeluaran', pe.id,
    '11111111-1111-1111-1111-111111111111'
FROM pengeluaran pe, last_saldo ls
WHERE NOT EXISTS (SELECT 1 FROM buku_kas_umum WHERE referensi_tabel = 'pengeluaran');
