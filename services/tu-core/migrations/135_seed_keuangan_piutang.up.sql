-- =====================================================================
-- Seed: 3 tagihan BELUM_LUNAS untuk demo aging
--
-- Skenario (mengacu CURRENT_DATE saat seed dijalankan):
--   - 25260001: jatuh_tempo H-15  → bucket <30
--   - 25260297: jatuh_tempo H-45  → bucket 30-60
--   - 25260003: jatuh_tempo H-100 → bucket >90
--
-- Pakai jenis_iuran GEDUNG (bukan SPP yang sudah LUNAS) agar unique
-- constraint (school, period, nis, jenis_iuran, bulan, tahun) tidak conflict.
-- Pakai CURRENT_DATE relatif supaya seed tetap relevan sepanjang waktu.
-- =====================================================================

-- ---------------------------------------------------------------------
-- 25260001 — GEDUNG, jatuh_tempo H-15 → <30
-- ---------------------------------------------------------------------
INSERT INTO "tagihan_siswa" (
    id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
    jenis_iuran_id, periode_bulan, periode_tahun, nominal, total_dibayar,
    status, tanggal_jatuh_tempo, created_by
)
SELECT
    gen_random_uuid(),
    '11111111-1111-1111-1111-111111111111', '20261',
    s.peserta_didik_id, s.nis,
    COALESCE((SELECT k.nm_kelas FROM tabel_anggotakelas ak
              JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
              WHERE ak.peserta_didik_id = s.peserta_didik_id LIMIT 1), ''),
    'f1000000-0000-0000-0000-000000000002', -- GEDUNG
    11, 2026, 1500000.00, 0,
    'BELUM_LUNAS', (CURRENT_DATE - INTERVAL '15 days')::date,
    '11111111-1111-1111-1111-111111111111'
FROM tabel_siswa s
WHERE s.nis = '25260001'
AND NOT EXISTS (
    SELECT 1 FROM tagihan_siswa
    WHERE nis = '25260001'
      AND jenis_iuran_id = 'f1000000-0000-0000-0000-000000000002'
      AND periode_bulan = 11 AND periode_tahun = 2026
);

-- ---------------------------------------------------------------------
-- 25260297 — GEDUNG, jatuh_tempo H-45 → 30-60
-- ---------------------------------------------------------------------
INSERT INTO "tagihan_siswa" (
    id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
    jenis_iuran_id, periode_bulan, periode_tahun, nominal, total_dibayar,
    status, tanggal_jatuh_tempo, created_by
)
SELECT
    gen_random_uuid(),
    '11111111-1111-1111-1111-111111111111', '20261',
    s.peserta_didik_id, s.nis,
    COALESCE((SELECT k.nm_kelas FROM tabel_anggotakelas ak
              JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
              WHERE ak.peserta_didik_id = s.peserta_didik_id LIMIT 1), ''),
    'f1000000-0000-0000-0000-000000000002',
    11, 2026, 1500000.00, 0,
    'BELUM_LUNAS', (CURRENT_DATE - INTERVAL '45 days')::date,
    '11111111-1111-1111-1111-111111111111'
FROM tabel_siswa s
WHERE s.nis = '25260297'
AND NOT EXISTS (
    SELECT 1 FROM tagihan_siswa
    WHERE nis = '25260297'
      AND jenis_iuran_id = 'f1000000-0000-0000-0000-000000000002'
      AND periode_bulan = 11 AND periode_tahun = 2026
);

-- ---------------------------------------------------------------------
-- 25260003 — GEDUNG, jatuh_tempo H-100 → >90
-- ---------------------------------------------------------------------
INSERT INTO "tagihan_siswa" (
    id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
    jenis_iuran_id, periode_bulan, periode_tahun, nominal, total_dibayar,
    status, tanggal_jatuh_tempo, created_by
)
SELECT
    gen_random_uuid(),
    '11111111-1111-1111-1111-111111111111', '20261',
    s.peserta_didik_id, s.nis,
    COALESCE((SELECT k.nm_kelas FROM tabel_anggotakelas ak
              JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
              WHERE ak.peserta_didik_id = s.peserta_didik_id LIMIT 1), ''),
    'f1000000-0000-0000-0000-000000000002',
    11, 2026, 1500000.00, 0,
    'BELUM_LUNAS', (CURRENT_DATE - INTERVAL '100 days')::date,
    '11111111-1111-1111-1111-111111111111'
FROM tabel_siswa s
WHERE s.nis = '25260003'
AND NOT EXISTS (
    SELECT 1 FROM tagihan_siswa
    WHERE nis = '25260003'
      AND jenis_iuran_id = 'f1000000-0000-0000-0000-000000000002'
      AND periode_bulan = 11 AND periode_tahun = 2026
);

-- ---------------------------------------------------------------------
-- 1 tagihan SEBAGIAN (untuk test status SEBAGIAN di aging)
-- 25260001 — EKSTRA dengan DP 50%
-- ---------------------------------------------------------------------
INSERT INTO "tagihan_siswa" (
    id, school_id, academic_period_id, peserta_didik_id, nis, kelas_nama,
    jenis_iuran_id, periode_bulan, periode_tahun, nominal, total_dibayar,
    status, tanggal_jatuh_tempo, created_by
)
SELECT
    gen_random_uuid(),
    '11111111-1111-1111-1111-111111111111', '20261',
    s.peserta_didik_id, s.nis,
    COALESCE((SELECT k.nm_kelas FROM tabel_anggotakelas ak
              JOIN tabel_kelas k ON k.rombongan_belajar_id = ak.rombongan_belajar_id
              WHERE ak.peserta_didik_id = s.peserta_didik_id LIMIT 1), ''),
    'f1000000-0000-0000-0000-000000000003', -- EKSTRA
    11, 2026, 75000.00, 37500.00,
    'SEBAGIAN', (CURRENT_DATE - INTERVAL '20 days')::date,
    '11111111-1111-1111-1111-111111111111'
FROM tabel_siswa s
WHERE s.nis = '25260001'
AND NOT EXISTS (
    SELECT 1 FROM tagihan_siswa
    WHERE nis = '25260001'
      AND jenis_iuran_id = 'f1000000-0000-0000-0000-000000000003'
      AND periode_bulan = 11 AND periode_tahun = 2026
);
