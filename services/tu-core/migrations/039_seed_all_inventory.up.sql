DO $$
DECLARE
    r RECORD;
    v_counter INT := 1;
BEGIN
    -- Looping untuk setiap kelas yang ada di tabel_kelas
    FOR r IN SELECT rombongan_belajar_id FROM tabel_kelas LOOP
        -- Masukkan 3 aset dummy ke setiap kelas
        INSERT INTO inventory_items (id, name, brand, category, condition, room_id) VALUES
        ('INV-' || v_counter || '-A', 'PC Komputer Rakitan Core i5', 'Asus', 'Elektronik', 'Bagus', r.rombongan_belajar_id),
        ('INV-' || v_counter || '-B', 'Proyektor BenQ Wall Mount', 'BenQ', 'Elektronik', 'Rusak Ringan', r.rombongan_belajar_id),
        ('INV-' || v_counter || '-C', 'Meja Siswa Chitose', 'Chitose', 'Mebel', 'Bagus', r.rombongan_belajar_id)
        ON CONFLICT (id) DO NOTHING;
        
        v_counter := v_counter + 1;
    END LOOP;
END $$;
