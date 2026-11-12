UPDATE attendance_records SET status = 'PRESENT' WHERE status = 'HADIR';
UPDATE attendance_records SET status = 'SICK' WHERE status = 'SAKIT';
UPDATE attendance_records SET status = 'PERMISSION' WHERE status = 'IZIN';
UPDATE attendance_records SET status = 'ABSENT' WHERE status = 'ALPA';
UPDATE attendance_records SET status = 'LATE' WHERE status = 'TERLAMBAT';
