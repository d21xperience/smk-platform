-- Seed data untuk pengujian UI Cuti Kerja
INSERT INTO ptk_leaves (id, ptk_id, jenis_cuti, tanggal_mulai, tanggal_selesai, alasan, status, catatan_approval) VALUES 
(
    '11111111-1111-1111-1111-111111111111', 
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', -- Budi Santoso
    'Cuti Tahunan', 
    '2026-11-10', 
    '2026-11-12', 
    'Acara keluarga di luar kota', 
    'pending', 
    NULL
),
(
    '22222222-2222-2222-2222-222222222222', 
    'b2c3d4e5-f6a7-8b9c-0d1e-2f3a4b5c6d7e', -- Siti Aminah
    'Cuti Sakit', 
    '2026-11-15', 
    '2026-11-16', 
    'Demam dan perlu istirahat total', 
    'disetujui', 
    'Semoga cepat sembuh, kelas sudah dicarikan guru pengganti.'
),
(
    '33333333-3333-3333-3333-333333333333', 
    'c3d4e5f6-a7b8-9c0d-1e2f-3a4b5c6d7e8f', -- Ahmad Fauzi
    'SPT (Surat Perintah Tugas)', 
    '2026-11-20', 
    '2026-11-20', 
    'Menghadiri undangan Dinas Pendidikan', 
    'ditolak', 
    'Bentrok dengan jadwal Ujian Akhir Semester, mohon dijadwalkan ulang.'
) ON CONFLICT (id) DO NOTHING;
