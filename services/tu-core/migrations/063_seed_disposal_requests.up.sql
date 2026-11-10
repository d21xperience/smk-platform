DO $$
DECLARE
    v_req_id UUID;
    v_asset_1 VARCHAR(100);
    v_asset_2 VARCHAR(100);
    v_price_1 NUMERIC;
    v_price_2 NUMERIC;
    v_total NUMERIC;
    v_level INT;
BEGIN
    -- Ambil 2 aset yang berstatus AKTIF untuk dijadikan dummy
    SELECT id, COALESCE(nilai_buku, 0) INTO v_asset_1, v_price_1 FROM inventory_items WHERE status_aset = 'AKTIF' LIMIT 1;
    SELECT id, COALESCE(nilai_buku, 0) INTO v_asset_2, v_price_2 FROM inventory_items WHERE status_aset = 'AKTIF' OFFSET 1 LIMIT 1;

    IF v_asset_1 IS NOT NULL THEN
        v_total := COALESCE(v_price_1, 0) + COALESCE(v_price_2, 0);
        IF v_total < 10000000 THEN v_level := 1;
        ELSIF v_total < 50000000 THEN v_level := 2;
        ELSE v_level := 3; END IF;

        -- 1. Buat Header Usulan
        INSERT INTO disposal_requests (id, request_number, proposed_by, proposal_date, reason_category, proposed_method, total_assets, total_value, approval_level, status, current_approval_step)
        VALUES (gen_random_uuid(), 'USUL-HPS/SMKPASJA/2026/001', 'system', CURRENT_DATE, 'RUSAK_BERAT', 'DIMUSNAHKAN', 2, v_total, v_level, 'DRAFT', 0)
        RETURNING id INTO v_req_id;

        -- 2. Buat Detail Item 1
        INSERT INTO disposal_request_items (disposal_request_id, inventory_item_id, specific_reason) 
        VALUES (v_req_id, v_asset_1, 'Layar rusak dan mati total, tidak dapat diperbaiki.');
        UPDATE inventory_items SET status_aset = 'DIUSULKAN_HAPUS', disposal_request_id = v_req_id WHERE id = v_asset_1;

        -- 3. Buat Detail Item 2 (Jika ada)
        IF v_asset_2 IS NOT NULL THEN
            INSERT INTO disposal_request_items (disposal_request_id, inventory_item_id, specific_reason) 
            VALUES (v_req_id, v_asset_2, 'Komponen motherboard terbakar akibat korsleting.');
            UPDATE inventory_items SET status_aset = 'DIUSULKAN_HAPUS', disposal_request_id = v_req_id WHERE id = v_asset_2;
        END IF;
    END IF;
END $$;
