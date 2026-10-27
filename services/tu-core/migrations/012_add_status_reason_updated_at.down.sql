DROP INDEX idx_enrollments_class_status;
DROP INDEX idx_enrollments_student_period;

ALTER TABLE enrollments DROP COLUMN updated_at;
ALTER TABLE enrollments DROP COLUMN reason;

ALTER TABLE classes DROP COLUMN status;
