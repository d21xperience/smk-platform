DELETE FROM correspondences WHERE id IN (
    'dddddddd-dddd-dddd-dddd-000000000001',
    'dddddddd-dddd-dddd-dddd-000000000002',
    'dddddddd-dddd-dddd-dddd-000000000003',
    'dddddddd-dddd-dddd-dddd-000000000004',
    'dddddddd-dddd-dddd-dddd-000000000005',
    'dddddddd-dddd-dddd-dddd-000000000006'
);
DELETE FROM correspondence_counters WHERE jenis_surat IN (
    'PERMOHONAN_PKL', 'UNDANGAN_ORTU', 'SKL', 'PANGGILAN_ORTU', 'REKOMENDASI_LOMBA'
);
