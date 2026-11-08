-- Drop tabel_kehadiran (legacy dapodik) yang sudah digantikan oleh
-- attendance_sessions + attendance_records. Tabel lama kosong (0 rows)
-- dan tidak direferensikan oleh kode Go/Vue.
DROP TABLE IF EXISTS tabel_kehadiran;
