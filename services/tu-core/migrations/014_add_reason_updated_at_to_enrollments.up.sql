-- 014_add_reason_updated_at_to_enrollments.up.sql
-- Menambahkan kolom reason dan updated_at pada tabel enrollments.
-- Jika tabel belum ada, buat tabel dengan struktur lengkap.

-- Buat tabel jika belum ada (untuk environment baru)
CREATE TABLE IF NOT EXISTS enrollments (
    id VARCHAR(36) PRIMARY KEY,
    student_id VARCHAR(36) NOT NULL,
    school_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(10) NOT NULL,
    class_id VARCHAR(36) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Tambahkan kolom jika belum ada (untuk upgrade)
DO $$
BEGIN
    -- Tambahkan kolom reason jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='enrollments' AND column_name='reason') THEN
        ALTER TABLE enrollments ADD COLUMN reason TEXT;
    END IF;

    -- Tambahkan kolom updated_at jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='enrollments' AND column_name='updated_at') THEN
        ALTER TABLE enrollments ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
    END IF;

    -- Tambahkan kolom deleted_at jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='enrollments' AND column_name='deleted_at') THEN
        ALTER TABLE enrollments ADD COLUMN deleted_at TIMESTAMP WITH TIME ZONE;
    END IF;

    -- Tambahkan kolom school_id jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='enrollments' AND column_name='school_id') THEN
        ALTER TABLE enrollments ADD COLUMN school_id VARCHAR(36) NOT NULL DEFAULT '';
    END IF;
END $$;

-- Buat indeks jika belum ada
CREATE INDEX IF NOT EXISTS idx_enrollments_student_period ON enrollments(student_id, period_id);
CREATE INDEX IF NOT EXISTS idx_enrollments_class_status ON enrollments(class_id, status);
CREATE INDEX IF NOT EXISTS idx_enrollments_school_period ON enrollments(school_id, period_id);
