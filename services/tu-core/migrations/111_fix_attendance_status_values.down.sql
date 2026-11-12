UPDATE attendance_records SET status = 'HADIR' WHERE status = 'PRESENT';
UPDATE attendance_records SET status = 'SAKIT' WHERE status = 'SICK';
UPDATE attendance_records SET status = 'IZIN' WHERE status = 'PERMISSION';
UPDATE attendance_records SET status = 'ALPA' WHERE status = 'ABSENT';
UPDATE attendance_records SET status = 'TERLAMBAT' WHERE status = 'LATE';
