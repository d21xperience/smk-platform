ALTER TABLE tabel_sekolah ADD COLUMN IF NOT EXISTS lama_pendidikan INT DEFAULT 3;
-- Set default lama pendidikan untuk SMK (3 tahun)
UPDATE tabel_sekolah SET lama_pendidikan = 3;
