INSERT INTO ptk_performances (id, ptk_id, tahun, semester, nilai_skp, predikat, catatan, evaluator) VALUES 
(
    '44444444-4444-4444-4444-444444444444', 
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', -- Budi Santoso
    '2026', 'Ganjil', 92.50, 'Sangat Baik', 'Guru berprestasi, inovatif dalam metode mengajar.', 'Kepala Sekolah'
),
(
    '55555555-5555-5555-5555-555555555555', 
    'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', -- Siti Aminah
    '2026', 'Ganjil', 85.00, 'Baik', 'Aktif dalam kegiatan sekolah, perlu peningkatan kedisiplinan waktu.', 'Kepala Sekolah'
),
(
    '66666666-6666-6666-6666-666666666666', 
    'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', -- Ahmad Fauzi
    '2026', 'Ganjil', 78.00, 'Cukup', 'Perlu bimbingan dalam penyusunan RPP dan administrasi kelas.', 'Kepala Sekolah'
) ON CONFLICT (id) DO NOTHING;
