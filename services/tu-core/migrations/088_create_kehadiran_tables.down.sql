DROP TABLE IF EXISTS ptk_rekap_bulanan_detail;
DROP TABLE IF EXISTS ptk_rekap_bulanan;
DROP TABLE IF EXISTS ptk_attendance_corrections;
DROP TABLE IF EXISTS master_jam_kerja;
ALTER TABLE ptk_attendance 
DROP COLUMN IF EXISTS correction_id,
DROP COLUMN IF EXISTS is_corrected,
DROP COLUMN IF EXISTS source_data,
DROP COLUMN IF EXISTS menit_terlambat,
DROP COLUMN IF EXISTS durasi_kerja,
DROP COLUMN IF EXISTS jam_pulang,
DROP COLUMN IF EXISTS jam_masuk;
