-- Unlock BKU November 2026 dulu
UPDATE "buku_kas_umum"
SET is_locked = FALSE
WHERE school_id = '11111111-1111-1111-1111-111111111111'
  AND academic_period_id = '20261'
  AND EXTRACT(MONTH FROM tanggal) = 11
  AND EXTRACT(YEAR FROM tanggal) = 2026;

DELETE FROM "penutupan_buku_kas" WHERE id = 'fc000000-0000-0000-0000-000000000001';
