-- Menambahkan kolom audit standar (non-DAPODIK) untuk melacak waktu.
ALTER TABLE tabel_kelas
    ADD COLUMN IF NOT EXISTS created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    ADD COLUMN IF NOT EXISTS updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW();
