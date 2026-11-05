INSERT INTO tabel_ptk_pelengkap (
    ptk_pelengkap_id, ptk_id, gelar_depan, gelar_belakang, nip_niy
) VALUES 
(
    '11111111-1111-1111-1111-111111111111', 'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 
    NULL, 'S.Kom.', '198001012005011001'
),
(
    '22222222-2222-2222-2222-222222222222', 'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 
    NULL, 'S.Pd.', '198502022010012002'
),
(
    '33333333-3333-3333-3333-333333333333', 'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 
    NULL, 'M.Pd.', '199003032015011003'
) ON CONFLICT (ptk_pelengkap_id) DO NOTHING;
