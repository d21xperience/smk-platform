-- Menghubungkan 3 Guru existing ke Periode SKP Tahunan 2026
-- Menggunakan Template Guru ASN (11111111-1111-1111-1111-111111111111)
INSERT INTO penilaian_pegawai (id, periode_id, ptk_id, template_id, penilai_utama_id, penilai_berikat_id, status) VALUES 
(
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 
    '99999999-9999-9999-9999-999999999999', -- Periode Seed
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', -- Budi Santoso
    '11111111-1111-1111-1111-111111111111', -- Template Guru ASN
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', -- Penilai Utama (Diri Sendiri sebagai contoh)
    'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', -- Penilai Berikat (Ahmad Fauzi / Kepsek)
    'PENDING'
),
(
    'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 
    '99999999-9999-9999-9999-999999999999', 
    'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', -- Siti Aminah
    '11111111-1111-1111-1111-111111111111', 
    'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', 
    'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', 
    'PENDING'
)
ON CONFLICT (id) DO NOTHING;
