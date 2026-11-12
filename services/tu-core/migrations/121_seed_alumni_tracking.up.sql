INSERT INTO "alumni_tracking" (
    "id", "school_id", "academic_period_id", "nama", "tahun_lulus", 
    "jurusan", "status", "instansi", "posisi_usaha", "created_at", "updated_at"
) VALUES
(
    'dd111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Rian Hidayat', '2025', 'RPL', 'Bekerja', 'PT. Telekomunikasi Indonesia Tbk', 'Frontend Developer', NOW(), NOW()
),
(
    'dd222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Siti Sarah', '2025', 'AKL', 'Melanjutkan Kuliah', 'Universitas Indonesia', 'S1 Akuntansi Perpajakan', NOW(), NOW()
),
(
    'dd333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Denny Wahyudi', '2025', 'TKRO', 'Wirausaha', 'Bengkel Wahyudi Motor', 'Pemilik Bengkel Mandiri', NOW(), NOW()
),
(
    'dd444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Gita Permata', '2025', 'RPL', 'Belum Bekerja', '', '', NOW(), NOW()
)
ON CONFLICT DO NOTHING;
