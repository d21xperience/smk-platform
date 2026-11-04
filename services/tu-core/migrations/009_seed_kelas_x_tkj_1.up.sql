-- Seed 1 rombel contoh: "X TKJ 1", semester 20261, capacity=36.
-- ON CONFLICT DO UPDATE agar idempotent (kalau 007 sudah pernah di-apply,
-- row-nya akan di-update dengan capacity baru).
INSERT INTO tabel_kelas (
    rombongan_belajar_id,
    sekolah_id,
    semester_id,
    nm_kelas,
    tingkat_pendidikan_id,
    ptk_id,
    nama_jurusan_sp,
    jenis_rombel,
    kurikulum_id,
    capacity
) VALUES (
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
    '11111111-1111-1111-1111-111111111111',
    '20261',
    'X TKJ 1',
    10,
    NULL,
    'Teknik Komputer dan Jaringan',
    1,
    1,
    36
)
ON CONFLICT (rombongan_belajar_id) DO UPDATE SET
    nm_kelas              = EXCLUDED.nm_kelas,
    tingkat_pendidikan_id = EXCLUDED.tingkat_pendidikan_id,
    nama_jurusan_sp       = EXCLUDED.nama_jurusan_sp,
    capacity              = EXCLUDED.capacity;
