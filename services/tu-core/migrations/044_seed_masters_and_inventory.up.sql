DO $$
DECLARE
    cat_eln UUID;
    cat_mbl UUID;
    cat_lab UUID;
    cata_pc UUID;
    cata_proj UUID;
    cata_meja UUID;
    r RECORD;
    v_counter INT := 1;
    v_year INT := EXTRACT(YEAR FROM NOW());
    v_kode_sekolah VARCHAR(20);
BEGIN
    -- Set Kode Sekolah Default
    UPDATE tabel_sekolah SET kode_sekolah = 'SMKPasja' WHERE kode_sekolah IS NULL;
    SELECT COALESCE(kode_sekolah, 'SMKPasja') INTO v_kode_sekolah FROM tabel_sekolah LIMIT 1;

    -- Insert Master Categories
    INSERT INTO master_categories (code, name) VALUES ('ELN', 'Elektronik & Multimedia') RETURNING id INTO cat_eln;
    INSERT INTO master_categories (code, name) VALUES ('MBL', 'Mebel & Furnitur') RETURNING id INTO cat_mbl;
    INSERT INTO master_categories (code, name) VALUES ('LAB', 'Sarana Pembelajaran & Laboratorium') RETURNING id INTO cat_lab;

    -- Insert Master Catalogs
    INSERT INTO master_catalogs (category_id, name, brand) VALUES (cat_eln, 'PC Komputer Rakitan Core i5', 'Asus') RETURNING id INTO cata_pc;
    INSERT INTO master_catalogs (category_id, name, brand) VALUES (cat_eln, 'Proyektor Epson EB-X400', 'Epson') RETURNING id INTO cata_proj;
    INSERT INTO master_catalogs (category_id, name, brand) VALUES (cat_mbl, 'Meja Siswa Chitose', 'Chitose') RETURNING id INTO cata_meja;

    -- Loop setiap kelas untuk masukkan aset
    FOR r IN SELECT rombongan_belajar_id FROM tabel_kelas LOOP
        -- 1 PC
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date) 
        VALUES (v_kode_sekolah || '/ELN/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_pc, 'Bagus', r.rombongan_belajar_id, make_date(v_year, 1, 1));
        v_counter := v_counter + 1;

        -- 1 Proyektor
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date) 
        VALUES (v_kode_sekolah || '/ELN/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_proj, 'Rusak Ringan', r.rombongan_belajar_id, make_date(v_year, 1, 1));
        v_counter := v_counter + 1;

        -- 1 Meja
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date) 
        VALUES (v_kode_sekolah || '/MBL/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_meja, 'Bagus', r.rombongan_belajar_id, make_date(v_year, 1, 1));
        v_counter := v_counter + 1;
    END LOOP;
END $$;
