-- =====================================================================
-- Seed: 2 rekonsiliasi bank
--   Rekon 1 (Nov 2026) — BALANCE: 3 matched + 1 unmatched buku
--   Rekon 2 (Sep 2026) — UNBALANCE: 2 unmatched bank
--
-- Total: 2 header + 3 detail + 3 unmatched
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. Seed rekonsiliasi_bank
-- ---------------------------------------------------------------------
INSERT INTO "rekonsiliasi_bank" (
    id, school_id, academic_period_id, rekening_bank_id,
    periode_bulan, periode_tahun,
    saldo_buku, saldo_bank, selisih,
    jumlah_matched, jumlah_unmatched_bank, jumlah_unmatched_buku,
    saldo_rekonsiliasi_bank, saldo_rekonsiliasi_buku,
    status, catatan, created_by
) VALUES
(
    'fb000000-0000-0000-0000-000000000001',
    '11111111-1111-1111-1111-111111111111', '20261',
    'f4000000-0000-0000-0000-000000000001',
    11, 2026,
    -450000.00, 750000.00, 1200000.00,
    3, 0, 1,
    -450000.00, -450000.00,
    'BALANCE', 'Rekonsiliasi November 2026 - semua cocok setelah adjustment',
    '11111111-1111-1111-1111-111111111111'
),
(
    'fb000000-0000-0000-0000-000000000002',
    '11111111-1111-1111-1111-111111111111', '20261',
    'f4000000-0000-0000-0000-000000000001',
    9, 2026,
    0.00, -20000.00, -20000.00,
    0, 2, 0,
    -20000.00, 0.00,
    'UNBALANCE', 'Rekonsiliasi September 2026 - ada biaya admin & transfer belum dicatat',
    '11111111-1111-1111-1111-111111111111'
)
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------
-- 2. Seed rekonsiliasi_detail (3 matched pairs — Nov 2026)
-- 3 BKM 250k ↔ 3 bank mutasi 250k, tanggal sama (selisih 0 hari)
-- ---------------------------------------------------------------------
INSERT INTO "rekonsiliasi_detail" (
    id, rekonsiliasi_id,
    mutasi_tanggal, mutasi_keterangan, mutasi_debet, mutasi_kredit,
    mutasi_saldo, mutasi_referensi,
    buku_kas_umum_id, selisih_hari
)
SELECT
    gen_random_uuid(),
    'fb000000-0000-0000-0000-000000000001',
    '2026-11-10',
    'Transfer masuk - ' || b.uraian,
    b.debet,
    0,
    b.debet,
    b.no_bukti,
    b.id,
    0
FROM buku_kas_umum b
WHERE b.school_id = '11111111-1111-1111-1111-111111111111'
  AND b.no_bukti IN ('BKM/001/SMK/XI/2026', 'BKM/002/SMK/XI/2026', 'BKM/003/SMK/XI/2026')
  AND NOT EXISTS (SELECT 1 FROM rekonsiliasi_detail WHERE rekonsiliasi_id = 'fb000000-0000-0000-0000-000000000001');

-- ---------------------------------------------------------------------
-- 3. Seed rekonsiliasi_unmatched
--   Rekon 1: 1 unmatched BUKU (BKK 1.2jt)
--   Rekon 2: 2 unmatched BANK (biaya admin & transfer masuk)
-- ---------------------------------------------------------------------
INSERT INTO "rekonsiliasi_unmatched" (
    id, rekonsiliasi_id, tipe, data_json, status_tindak_lanjut
)
SELECT
    gen_random_uuid(),
    'fb000000-0000-0000-0000-000000000001',
    'BUKU',
    jsonb_build_object(
        'no_bukti', b.no_bukti,
        'tanggal', b.tanggal,
        'uraian', b.uraian,
        'debet', b.debet,
        'kredit', b.kredit,
        'saldo_berjalan', b.saldo_berjalan,
        'kode_rekening_id', b.kode_rekening_id
    ),
    'PENDING'
FROM buku_kas_umum b
WHERE b.school_id = '11111111-1111-1111-1111-111111111111'
  AND b.no_bukti = 'BKK/001/SMK/XI/2026'
  AND NOT EXISTS (SELECT 1 FROM rekonsiliasi_unmatched WHERE rekonsiliasi_id = 'fb000000-0000-0000-0000-000000000001');

INSERT INTO "rekonsiliasi_unmatched" (
    id, rekonsiliasi_id, tipe, data_json, status_tindak_lanjut
) VALUES
(
    gen_random_uuid(),
    'fb000000-0000-0000-0000-000000000002',
    'BANK',
    '{"tanggal":"2026-09-15","keterangan":"Biaya admin bulanan","debet":0,"kredit":20000,"saldo":-20000,"referensi":null}'::jsonb,
    'PENDING'
),
(
    gen_random_uuid(),
    'fb000000-0000-0000-0000-000000000002',
    'BANK',
    '{"tanggal":"2026-09-20","keterangan":"Transfer masuk belum tercatat","debet":0,"kredit":500000,"saldo":480000,"referensi":null}'::jsonb,
    'PENDING'
);
