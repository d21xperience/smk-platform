-- Tabel Riwayat Pendidikan Formal
CREATE TABLE IF NOT EXISTS tabel_rwy_pend_formal (
    rwy_pend_formal_id UUID PRIMARY KEY,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenjang_pendidikan VARCHAR(20) NOT NULL,
    jurusan VARCHAR(100),
    nama_sekolah VARCHAR(150),
    tahun_lulus INT,
    gelar VARCHAR(50)
);

-- Tabel Riwayat Sertifikasi
CREATE TABLE IF NOT EXISTS tabel_rwy_sertifikasi (
    rwy_sertifikasi_id UUID PRIMARY KEY,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis_sertifikasi VARCHAR(100) NOT NULL,
    nomor_sertifikat VARCHAR(100),
    tahun_sertifikasi INT
);
