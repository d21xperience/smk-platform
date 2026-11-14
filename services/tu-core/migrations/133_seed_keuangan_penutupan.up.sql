-- =====================================================================
-- Seed: Penutupan Buku Kas November 2026 (CLOSED)
--
-- Snapshot dari data BKU November 2026:
--   saldo_awal:        0
--   total_debet:       750.000   (3 BKM × 250k)
--   total_kredit:      1.200.000 (1 BKK)
--   saldo_akhir_buku:  -450.000
--   saldo_fisik_kas:   -450.000  (asumsi cash count match)
--   selisih_kas:       0
--
-- Snapshot posisi keuangan (sederhana):
--   total_aset:        -450.000  (kas)
--   total_kewajiban:   0
--   total_ekuitas:     -450.000
-- =====================================================================

INSERT INTO "penutupan_buku_kas" (
    id, school_id, academic_period_id, periode_bulan, periode_tahun,
    saldo_awal, total_debet, total_kredit, saldo_akhir_buku,
    saldo_fisik_kas, selisih_kas,
    total_aset, total_kewajiban, total_ekuitas,
    status, catatan, closed_by, closed_at
) VALUES (
    'fc000000-0000-0000-0000-000000000001',
    '11111111-1111-1111-1111-111111111111', '20261',
    11, 2026,
    0.00, 750000.00, 1200000.00, -450000.00,
    -450000.00, 0.00,
    -450000.00, 0.00, -450000.00,
    'CLOSED',
    'Penutupan November 2026 — cash count match dengan saldo buku',
    '11111111-1111-1111-1111-111111111111',
    '2026-11-30 17:00:00+07'
)
ON CONFLICT DO NOTHING;

-- Lock semua BKU November 2026 (sesuai design decision 3A)
UPDATE "buku_kas_umum"
SET is_locked = TRUE
WHERE school_id = '11111111-1111-1111-1111-111111111111'
  AND academic_period_id = '20261'
  AND EXTRACT(MONTH FROM tanggal) = 11
  AND EXTRACT(YEAR FROM tanggal) = 2026
  AND is_locked = FALSE;
