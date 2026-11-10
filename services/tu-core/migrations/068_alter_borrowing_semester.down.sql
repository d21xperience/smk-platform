DROP INDEX IF EXISTS idx_borrowing_semester;
ALTER TABLE borrowings DROP COLUMN IF EXISTS semester_id;
