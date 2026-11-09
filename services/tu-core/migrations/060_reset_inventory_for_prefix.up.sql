-- Hapus data lama agar format ID konsisten (Development Only)
TRUNCATE TABLE inventory_item_history CASCADE;
TRUNCATE TABLE sarpras_inventory_document_links CASCADE;
TRUNCATE TABLE inventory_items CASCADE;

-- Seed ulang dengan format baru (Prefix P untuk Pembelian)
DO $$
DECLARE
    cat_eln UUID;
    cat_mbl UUID;
    cata_pc UUID;
    cata_proj UUID;
    cata_meja UUID;
    r RECORD;
    v_counter INT := 1;
    v_year INT := 2026;
    v_sch VARCHAR(20);
BEGIN
    SELECT COALESCE(kode_sekolah, 'SMKPasja') INTO v_sch FROM tabel_sekolah LIMIT 1;
    SELECT id INTO cat_eln FROM master_categories WHERE code = 'ELN';
    SELECT id INTO cat_mbl FROM master_categories WHERE code = 'MBL';
    SELECT id INTO cata_pc FROM master_catalogs WHERE name LIKE '%PC Komputer%';
    SELECT id INTO cata_proj FROM master_catalogs WHERE name LIKE '%Proyektor%';
    SELECT id INTO cata_meja FROM master_catalogs WHERE name LIKE '%Meja%';

    FOR r IN SELECT rombongan_belajar_id FROM tabel_kelas LOOP
        -- 1 PC (Pembelian)
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date, sumber_perolehan, status_aset, tanggal_pencatatan) 
        VALUES (v_sch || '/P/ELN/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_pc, 'Bagus', r.rombongan_belajar_id, make_date(v_year, 1, 1), 'PEMBELIAN', 'AKTIF', make_date(v_year, 1, 1));
        v_counter := v_counter + 1;

        -- 1 Proyektor (Hibah)
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date, sumber_perolehan, status_aset, tanggal_pencatatan) 
        VALUES (v_sch || '/H/ELN/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_proj, 'Rusak Ringan', r.rombongan_belajar_id, make_date(v_year, 1, 1), 'HIBAH', 'AKTIF', make_date(v_year, 1, 1));
        v_counter := v_counter + 1;

        -- 1 Meja (Pembelian)
        INSERT INTO inventory_items (id, catalog_id, condition, room_id, purchase_date, sumber_perolehan, status_aset, tanggal_pencatatan) 
        VALUES (v_sch || '/P/MBL/' || v_year || '/' || LPAD(v_counter::text, 3, '0'), cata_meja, 'Bagus', r.rombongan_belajar_id, make_date(v_year, 1, 1), 'PEMBELIAN', 'AKTIF', make_date(v_year, 1, 1));
        v_counter := v_counter + 1;
    END LOOP;
END $$;
