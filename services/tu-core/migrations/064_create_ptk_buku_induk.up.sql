CREATE TABLE IF NOT EXISTS ptk_buku_induk_ext (
    ptk_id UUID PRIMARY KEY REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis_pegawai VARCHAR(50) NOT NULL DEFAULT 'Guru',
    status_pegawai VARCHAR(50) NOT NULL DEFAULT 'Honorer',
    tmt_masuk DATE,
    unit_kerja VARCHAR(100),
    jabatan VARCHAR(100),
    pangkat_golongan VARCHAR(50),
    no_sk_pengangkatan VARCHAR(100),
    tanggal_sk DATE
);

-- Seed data untuk 3 guru existing
INSERT INTO ptk_buku_induk_ext (ptk_id, jenis_pegawai, status_pegawai, tmt_masuk, unit_kerja, jabatan) VALUES
('a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'Guru', 'PNS', '2010-01-01', 'Kurikulum', 'Guru Mata Pelajaran'),
('b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'Guru', 'PPPK', '2015-07-01', 'Kesiswaan', 'Wali Kelas'),
('c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 'Guru', 'GTY', '2018-01-01', 'Sarpras', 'Staf Sarpras')
ON CONFLICT (ptk_id) DO NOTHING;
