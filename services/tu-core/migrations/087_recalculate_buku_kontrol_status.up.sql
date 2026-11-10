-- Recalculate status_monitoring untuk semua pegawai berdasarkan tanggal hari ini
-- Aturan:
-- < 0 hari: OVERDUE
-- <= 7 hari: URGENT
-- <= 30 hari: SEGERA
-- <= 90 hari: PERLU_PERSIAPAN
-- > 90 hari: AMAN

UPDATE buku_kontrol 
SET status_monitoring = CASE
    WHEN LEAST(
        COALESCE(tanggal_kgb_berikut - CURRENT_DATE, 99999),
        COALESCE(tanggal_pangkat_berikut - CURRENT_DATE, 99999)
    ) < 0 THEN 'OVERDUE'
    WHEN LEAST(
        COALESCE(tanggal_kgb_berikut - CURRENT_DATE, 99999),
        COALESCE(tanggal_pangkat_berikut - CURRENT_DATE, 99999)
    ) <= 7 THEN 'URGENT'
    WHEN LEAST(
        COALESCE(tanggal_kgb_berikut - CURRENT_DATE, 99999),
        COALESCE(tanggal_pangkat_berikut - CURRENT_DATE, 99999)
    ) <= 30 THEN 'SEGERA'
    WHEN LEAST(
        COALESCE(tanggal_kgb_berikut - CURRENT_DATE, 99999),
        COALESCE(tanggal_pangkat_berikut - CURRENT_DATE, 99999)
    ) <= 90 THEN 'PERLU_PERSIAPAN'
    ELSE 'AMAN'
END,
last_calculated_at = NOW();
