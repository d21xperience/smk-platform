DO $$
DECLARE
    v_asset_1 VARCHAR(100);
    v_asset_2 VARCHAR(100);
BEGIN
    SELECT id INTO v_asset_1 FROM inventory_items WHERE status_aset = 'AKTIF' LIMIT 1;
    SELECT id INTO v_asset_2 FROM inventory_items WHERE status_aset = 'AKTIF' OFFSET 1 LIMIT 1;

    -- Seed Tiket Kerusakan
    IF v_asset_1 IS NOT NULL THEN
        INSERT INTO damage_reports (report_number, inventory_item_id, asset_name, location, reporter_id, reporter_name, description, urgency, status, semester_id)
        VALUES 
        ('DMG/2026/001', v_asset_1, 'AC Split Sharp 1 PK', 'Ruang Kelas 11-B', 'user-1', 'Bu Sri (Wali Kelas)', 'AC mengeluarkan air bocor deras ke lantai kelas.', 'Tinggi', 'Pending', '20261'),
        ('DMG/2026/002', v_asset_2, 'Proyektor BenQ', 'Lab Fisika', 'user-2', 'Pak Bambang', 'Lampu proyektor berkedip merah lalu mati.', 'Tinggi', 'Proses', '20261');
    END IF;

    -- Seed Jadwal Perawatan
    INSERT INTO maintenance_schedules (schedule_code, asset_name, scope_location, frequency, last_service_date, next_service_date, semester_id)
    VALUES
    ('PM-01', 'AC Split Sharp 1 PK (Semua Ruang)', 'Gedung A & B', '3 Bulanan', CURRENT_DATE - 90, CURRENT_DATE + 5, '20261'),
    ('PM-02', 'Deep Cleaning PC Lab', 'Lab Komputer', 'Bulanan', CURRENT_DATE - 30, CURRENT_DATE + 1, '20261'),
    ('PM-03', 'Servis Rutin Mobil Sekolah', 'Garasi Samping', '6 Bulanan', CURRENT_DATE - 180, CURRENT_DATE + 45, '20261');
END $$;
