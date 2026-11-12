-- Seed: 2 Surat Aktif untuk siswa yang sudah ada
-- Baldhan (fa319cb5) dan Basudewa (9830ce3f)

INSERT INTO "student_letters" (
    "id", "school_id", "academic_period_id", "student_id",
    "letter_number", "parent_name", "purpose", "status", "created_at", "updated_at"
) VALUES
(
    'e1111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111',
    '20261',
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba',
    '421.5/084-SMK/2026',
    'Bapak Faqih',
    'Persyaratan BPJS Kesehatan',
    'PENDING',
    NOW(), NOW()
),
(
    'e2222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111',
    '20261',
    '9830ce3f-cd87-4b10-aa0e-a3fb89d16317',
    '421.5/085-SMK/2026',
    'Bapak Arif',
    'Pengajuan Beasiswa Luar Sekolah',
    'COMPLETED',
    NOW(), NOW()
)
ON CONFLICT ("id") DO NOTHING;
