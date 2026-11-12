INSERT INTO "alumni_legalisir" (
    "id", "school_id", "academic_period_id", "nama", "tahun_lulus", 
    "jurusan", "no_ijazah", "jumlah_lembar", "status", "created_at", "updated_at"
) VALUES
(
    'cc111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Rian Hidayat', '2025', 'RPL', 'M-SMK/23/0012391', 5, 'SIAP DIAMBIL', NOW(), NOW()
),
(
    'cc222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Siti Sarah', '2025', 'AKL', 'M-SMK/23/0052312', 8, 'PROSES TTD', NOW(), NOW()
),
(
    'cc333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111', '20261',
    'Denny Wahyudi', '2024', 'TKRO', 'M-SMK/22/0088123', 5, 'SIAP DIAMBIL', NOW(), NOW()
)
ON CONFLICT DO NOTHING;
