INSERT INTO "pkl_plotting" (
    "id", "school_id", "academic_period_id", "nama", "kelas", 
    "jurusan", "perusahaan", "kota", "pembimbing", "created_at", "updated_at"
) VALUES
(
    'ee111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Ahmad Fauzi', 'XII RPL 1', 'RPL', 'PT. Telekomunikasi Indonesia Tbk', 'Bandung', 'Budi Rahardjo, S.T.', NOW(), NOW()
),
(
    'ee222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Budi Santoso', 'XII TKRO 3', 'TKRO', 'PT. Astra Honda Motor', 'Jakarta', 'Eko Sulistyo, M.Pd.', NOW(), NOW()
),
(
    'ee333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Chandra Wijaya', 'XI RPL 2', 'RPL', 'PT. Solusi Teknologi Nusantara', 'Cimahi', 'Budi Rahardjo, S.T.', NOW(), NOW()
),
(
    'ee444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Dinda Lestari', 'XI AKL 1', 'AKL', 'Bank Mandiri (Persero) Tbk', 'Bandung', 'Hj. Nenden, S.E.', NOW(), NOW()
)
ON CONFLICT DO NOTHING;
