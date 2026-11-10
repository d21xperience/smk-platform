INSERT INTO penilaian_periode (id, nama_periode, jenis_periode, tanggal_mulai, tanggal_selesai, deadline_pengisian, tahun_referensi, status) VALUES 
(
    '99999999-9999-9999-9999-999999999999', 
    'SKP Tahunan 2026', 
    'TAHUNAN', 
    '2026-01-01', 
    '2026-12-31', 
    '2027-01-15', 
    2026, 
    'OPEN'
)
ON CONFLICT (id) DO NOTHING;
