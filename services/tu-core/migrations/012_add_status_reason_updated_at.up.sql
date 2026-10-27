-- Tambah status pada kelas (active/inactive)
ALTER TABLE classes ADD COLUMN status VARCHAR(20) DEFAULT 'active';

-- Tambah reason dan updated_at pada enrollments
ALTER TABLE enrollments ADD COLUMN reason TEXT;
ALTER TABLE enrollments ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT now();

-- Index untuk performa query histori
CREATE INDEX idx_enrollments_student_period ON enrollments (student_id, period_id);
CREATE INDEX idx_enrollments_class_status ON enrollments (class_id, status);

-- Update existing records
UPDATE classes SET status = 'active' WHERE status IS NULL;
UPDATE enrollments SET updated_at = created_at WHERE updated_at IS NULL;
ALTER TABLE classes ALTER COLUMN status SET NOT NULL;
