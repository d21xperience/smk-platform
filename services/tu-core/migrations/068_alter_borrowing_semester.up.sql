ALTER TABLE borrowings ADD COLUMN IF NOT EXISTS semester_id VARCHAR(5);
UPDATE borrowings SET semester_id = '20261' WHERE semester_id IS NULL;
CREATE INDEX IF NOT EXISTS idx_borrowing_semester ON borrowings(semester_id);
