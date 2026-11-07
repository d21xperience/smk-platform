-- Tambah kolom Data Pribadi Dapodik ke tabel_ptk
ALTER TABLE tabel_ptk 
    ADD COLUMN IF NOT EXISTS nik VARCHAR(16),
    ADD COLUMN IF NOT EXISTS agama_id INT,
    ADD COLUMN IF NOT EXISTS status_kepegawaian_id INT,
    ADD COLUMN IF NOT EXISTS foto_url TEXT,
    ADD COLUMN IF NOT EXISTS no_hp VARCHAR(15),
    ADD COLUMN IF NOT EXISTS email VARCHAR(50),
    ADD COLUMN IF NOT EXISTS status_perkawinan INT,
    ADD COLUMN IF NOT EXISTS nama_ibu_kandung VARCHAR(100);

-- Hubungkan tabel users dengan ptk_id (Mencegah error jika tabel users tidak ada/beda nama)
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        ALTER TABLE users ADD COLUMN IF NOT EXISTS ptk_id UUID REFERENCES tabel_ptk(ptk_id);
    END IF;
END $$;
