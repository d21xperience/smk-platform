DO $$
DECLARE
    v_asset1 VARCHAR(100);
    v_asset2 VARCHAR(100);
    v_asset3 VARCHAR(100);
BEGIN
    SELECT id INTO v_asset1 FROM inventory_items WHERE condition = 'Rusak Ringan' AND status_aset = 'AKTIF' LIMIT 1;
    SELECT id INTO v_asset2 FROM inventory_items WHERE condition = 'Bagus' AND status_aset = 'AKTIF' LIMIT 1;
    SELECT id INTO v_asset3 FROM inventory_items WHERE condition = 'Rusak Berat' AND status_aset = 'AKTIF' LIMIT 1;

    IF v_asset1 IS NOT NULL THEN
        INSERT INTO inventory_item_history (inventory_item_id, action, old_data, new_data, description, changed_by, changed_at)
        VALUES 
        (v_asset1, 'UBAH_KONDISI', '{"condition":"Bagus"}', '{"condition":"Rusak Ringan"}', 'Berubah karena sering digunakan', 'system', NOW() - INTERVAL '30 days'),
        (v_asset1, 'UBAH_KONDISI', '{"condition":"Rusak Ringan"}', '{"condition":"Rusak Berat"}', 'Kerusakan semakin parah', 'system', NOW() - INTERVAL '5 days');
    END IF;
    IF v_asset2 IS NOT NULL THEN
        INSERT INTO inventory_item_history (inventory_item_id, action, old_data, new_data, description, changed_by, changed_at)
        VALUES (v_asset2, 'UBAH_KONDISI', '{"condition":"Rusak Ringan"}', '{"condition":"Bagus"}', 'Telah diperbaiki oleh teknisi', 'system', NOW() - INTERVAL '15 days');
    END IF;
    IF v_asset3 IS NOT NULL THEN
        INSERT INTO inventory_item_history (inventory_item_id, action, old_data, new_data, description, changed_by, changed_at)
        VALUES (v_asset3, 'UBAH_KONDISI', '{"condition":"Bagus"}', '{"condition":"Rusak Berat"}', 'Rusak akibat bencana', 'system', NOW() - INTERVAL '10 days');
    END IF;
END $$;
