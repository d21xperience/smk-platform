-- Seed 117: student_academic_scores (REVISED)
-- 8 baris seed: 4 siswa x 3 assessment_type (DAILY/MIDTERM/FINAL)
-- student_id merujuk tabel_siswa yang sudah ada (seed 115)
-- Format academic_period_id = '20261' (Ganjil 2026/2027), semester = 1

INSERT INTO "student_academic_scores" (
    "id", "school_id", "academic_period_id", "semester", "student_id",
    "pembelajaran_id", "mata_pelajaran_id", "nama_mata_pelajaran",
    "assessment_type", "score", "assessment_date", "catatan",
    "created_at", "updated_at"
) VALUES
(
    'bb000001-1111-1111-1111-000000000001',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba',
    NULL, 1, 'Pendidikan Agama dan Budi Pekerti',
    'DAILY', 88.00, '2026-08-15', 'Ulangan harian bab 1 - pemahaman baik',
    NOW(), NOW()
),
(
    'bb000002-1111-1111-1111-000000000002',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba',
    NULL, 1, 'Pendidikan Agama dan Budi Pekerti',
    'MIDTERM', 85.00, '2026-09-20', NULL,
    NOW(), NOW()
),
(
    'bb000003-1111-1111-1111-000000000003',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba',
    NULL, 1, 'Pendidikan Agama dan Budi Pekerti',
    'FINAL', 90.00, '2026-12-10', 'UAS semester 1 - sangat baik',
    NOW(), NOW()
),
(
    'bb000004-1111-1111-1111-000000000004',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    '9830ce3f-cd87-4b10-aa0e-a3fb89d16317',
    NULL, 3, 'Bahasa Indonesia',
    'DAILY', 82.00, '2026-08-20', 'Ulangan harian teks eksposisi',
    NOW(), NOW()
),
(
    'bb000005-1111-1111-1111-000000000005',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    '9830ce3f-cd87-4b10-aa0e-a3fb89d16317',
    NULL, 3, 'Bahasa Indonesia',
    'MIDTERM', 80.00, '2026-09-25', 'Perlu peningkatan pada struktur teks',
    NOW(), NOW()
),
(
    'bb000006-1111-1111-1111-000000000006',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    '7308aff7-3b22-4d2f-8648-dece2107e5fd',
    NULL, 4, 'Matematika',
    'DAILY', 76.00, '2026-08-22', 'Perlu pendampingan materi statistika',
    NOW(), NOW()
),
(
    'bb000007-1111-1111-1111-000000000007',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    '15c822ae-c9b5-4fbf-bba0-6a85088555eb',
    NULL, 5, 'Bahasa Inggris',
    'DAILY', 84.00, '2026-08-25', 'Listening & speaking baik',
    NOW(), NOW()
),
(
    'bb000008-1111-1111-1111-000000000008',
    '11111111-1111-1111-1111-111111111111', '20261', 1,
    '15c822ae-c9b5-4fbf-bba0-6a85088555eb',
    NULL, 5, 'Bahasa Inggris',
    'FINAL', 87.00, '2026-12-12', 'UAS semester 1 - konsisten',
    NOW(), NOW()
)
ON CONFLICT DO NOTHING;
