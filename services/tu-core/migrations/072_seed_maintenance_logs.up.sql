DO $$
DECLARE
    v_sch_id UUID;
    v_asset_id VARCHAR(100);
    v_asset_name VARCHAR(255);
BEGIN
    -- Ambil 1 jadwal dan 1 aset untuk dummy log
    SELECT id INTO v_sch_id FROM maintenance_schedules LIMIT 1;
    SELECT id, (SELECT name FROM master_catalogs WHERE id = catalog_id) INTO v_asset_id, v_asset_name FROM inventory_items LIMIT 1;

    IF v_sch_id IS NOT NULL AND v_asset_id IS NOT NULL THEN
        INSERT INTO maintenance_logs (schedule_id, inventory_item_id, service_date, technician, vendor, cost, notes, performed_by)
        VALUES 
        (v_sch_id, v_asset_id, CURRENT_DATE - 30, 'Pak Budi', 'CV Teknisi Mandiri', 150000, 'Pembersihan rutin dan pengecekan freon.', 'system'),
        (v_sch_id, v_asset_id, CURRENT_DATE - 90, 'Pak Budi', 'CV Teknisi Mandiri', 150000, 'Pembersihan filter udara.', 'system');
    END IF;
END $$;
