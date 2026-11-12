INSERT INTO "student_scholarships" (
    "id", "school_id", "academic_period_id", "student_id",
    "category", "card_number", "reason", "is_document_valid", "status",
    "bank_name", "bank_account", "disbursed_amount", "created_at", "updated_at"
) VALUES
(
    'f1111111-1111-1111-1111-111111111111',
    '11111111-1111-1111-1111-111111111111', '20261',
    'fa319cb5-eb56-44ee-90e0-eb4e743e54ba',
    'Pemegang KIP', 'KIP-2026-8871A', 'Keluarga kurang mampu',
    TRUE, 'SK PEMBERIAN',
    'BNI', '1172836210', 1800000, NOW(), NOW()
),
(
    'f2222222-2222-2222-2222-222222222222',
    '11111111-1111-1111-1111-111111111111', '20261',
    '9830ce3f-cd87-4b10-aa0e-a3fb89d16317',
    'Keluarga PKH / KKS', 'KKS-3204-1192B', 'Penerima PKH',
    TRUE, 'SK NOMINASI',
    'BNI', '1172836220', NULL, NOW(), NOW()
),
(
    'f3333333-3333-3333-3333-333333333333',
    '11111111-1111-1111-1111-111111111111', '20261',
    '7308aff7-3b22-4d2f-8648-dece2107e5fd',
    'Pertimbangan Khusus (SKTM)', '', 'Yatim piatu, orang tua meninggal',
    FALSE, 'SK NOMINASI',
    NULL, NULL, NULL, NOW(), NOW()
),
(
    'f4444444-4444-4444-4444-444444444444',
    '11111111-1111-1111-1111-111111111111', '20261',
    '15c822ae-c9b5-4fbf-bba0-6a85088555eb',
    'Pemegang KIP', 'KIP-2026-9921C', 'Penerima KIP tahun lalu',
    TRUE, 'SK PEMBERIAN',
    'BNI', '1172836230', 1800000, NOW(), NOW()
)
ON CONFLICT ("id") DO NOTHING;
