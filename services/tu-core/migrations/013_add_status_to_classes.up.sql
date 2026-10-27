-- 013_add_status_to_classes.up.sql
-- Menambahkan kolom status pada tabel classes, serta kolom level jika belum ada.
-- Jika tabel belum ada, buat tabel dengan struktur lengkap.

-- Buat tabel jika belum ada (untuk environment baru)
CREATE TABLE IF NOT EXISTS classes (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    period_id VARCHAR(10) NOT NULL,
    name VARCHAR(100) NOT NULL,
    grade INTEGER NOT NULL,
    level VARCHAR(5),
    major VARCHAR(50),
    homeroom_id VARCHAR(36),
    capacity INTEGER NOT NULL DEFAULT 0,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Tambahkan kolom jika belum ada (untuk upgrade)
DO $$
BEGIN
    -- Tambahkan kolom status jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='classes' AND column_name='status') THEN
        ALTER TABLE classes ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'active';
    END IF;

    -- Tambahkan kolom level jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='classes' AND column_name='level') THEN
        ALTER TABLE classes ADD COLUMN level VARCHAR(5);
    END IF;

    -- Tambahkan kolom major jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='classes' AND column_name='major') THEN
        ALTER TABLE classes ADD COLUMN major VARCHAR(50);
    END IF;

    -- Tambahkan kolom homeroom_id jika belum ada
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns
                   WHERE table_name='classes' AND column_name='homeroom_id') THEN
        ALTER TABLE classes ADD COLUMN homeroom_id VARCHAR(36);
    END IF;
END $$;

-- Buat indeks jika belum ada
CREATE INDEX IF NOT EXISTS idx_classes_school_period ON classes(school_id, period_id);
CREATE INDEX IF NOT EXISTS idx_classes_period_status ON classes(period_id, status);
