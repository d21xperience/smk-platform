ALTER TABLE attendance_sessions ALTER COLUMN subject_id TYPE VARCHAR(36) USING subject_id::VARCHAR;
