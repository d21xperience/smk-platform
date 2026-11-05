INSERT INTO tabel_pembelajaran (
    pembelajaran_id, rombongan_belajar_id, mata_pelajaran_id, semester_id, ptk_terdaftar_id, nama_mata_pelajaran
) VALUES 
('99999999-9999-9999-9999-999999999991', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 1, '20261', 'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'Matematika'),
('99999999-9999-9999-9999-999999999992', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 2, '20261', 'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'Bahasa Indonesia'),
('99999999-9999-9999-9999-999999999993', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 3, '20261', 'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 'Bahasa Inggris')
ON CONFLICT (pembelajaran_id) DO NOTHING;
