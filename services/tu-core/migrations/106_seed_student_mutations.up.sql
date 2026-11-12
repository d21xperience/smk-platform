INSERT INTO "student_mutations" (
    "id", "school_id", "academic_period_id", "student_id", "mutation_type",
    "nis", "nisn", "nama", "jenis_kelamin", "asal_sekolah", "tujuan_sekolah",
    "alasan", "status", "created_at", "updated_at"
) VALUES
-- Skenario Mutasi MASUK (Siswa Pindahan)
(
    'b1111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    '20261',
    NULL,
    'MASUK',
    '262700099',
    '0099887766',
    'Rian Hidayat',
    'L',
    'SMKN 2 Bandung',
    NULL,
    'Ikut Orang Tua',
    'PENDING',
    NOW(), NOW()
),
-- Skenario Mutasi KELUAR (Siswa Pindah)
(
    'b2222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    '20261',
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba', -- ID Baldhan Faqihuddin
    'KELUAR',
    '25260002',
    '0098661629',
    'Baldhan Faqihuddin',
    'L',
    NULL,
    'SMKN 1 Jakarta',
    'Pindah Domisili',
    'PENDING',
    NOW(), NOW()
)
ON CONFLICT ("id") DO NOTHING;
