-- Ekstensi non-DAPODIK: kolom capacity untuk fitur manajemen rombel.
-- Semantik:
--   NULL  = tidak ada batas (unlimited)
--   > 0   = batas maksimal siswa per rombel
-- Tidak dipakai constraint NOT NULL agar kompatibel dengan data DAPODIK
-- yang diimpor tanpa kolom ini.
ALTER TABLE tabel_kelas
    ADD COLUMN IF NOT EXISTS capacity SMALLINT NULL DEFAULT NULL;

COMMENT ON COLUMN tabel_kelas.capacity IS
    'Batas maksimal siswa per rombel. NULL = unlimited (ekstensi non-DAPODIK).';
