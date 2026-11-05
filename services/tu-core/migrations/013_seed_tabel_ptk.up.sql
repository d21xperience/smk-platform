INSERT INTO tabel_ptk (
    ptk_id, nama, nip, jenis_ptk_id, jenis_kelamin, tempat_lahir, tanggal_lahir, 
    nuptk, alamat_jalan, status_keaktifan_id, soft_delete
) VALUES 
(
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 'Budi Santoso, S.Kom.', '198001012005011001', 
    1, 'L', 'Jakarta', '1980-01-01', '1234567890123456', 'Jl. Merdeka No. 1', 1, 0
),
(
    'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 'Siti Aminah, S.Pd.', '198502022010012002', 
    1, 'P', 'Bandung', '1985-02-02', '2345678901234567', 'Jl. Sudirman No. 2', 1, 0
),
(
    'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 'Ahmad Fauzi, M.Pd.', '199003032015011003', 
    1, 'L', 'Surabaya', '1990-03-03', '3456789012345678', 'Jl. Diponegoro No. 3', 1, 0
) ON CONFLICT (ptk_id) DO NOTHING;
