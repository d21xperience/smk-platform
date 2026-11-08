-- Recreate tabel_kehadiran (rollback)
CREATE TABLE IF NOT EXISTS tabel_kehadiran (
    kd_kehadiran UUID PRIMARY KEY,
    semester_id VARCHAR(5) NOT NULL,
    anggota_rombel_id UUID NOT NULL REFERENCES tabel_anggotakelas(anggota_rombel_id) ON DELETE CASCADE,
    sakit NUMERIC(3,0),
    izin NUMERIC(3,0),
    tanpa_keterangan NUMERIC(3,0),
    peserta_didik_id UUID,
    semester NUMERIC(1,0)
);
