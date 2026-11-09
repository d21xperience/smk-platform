-- Update data seed yang sudah ada agar memiliki semester_id (sesuai context 20261)
UPDATE procurement_requests SET semester_id = '20261' WHERE semester_id IS NULL;
