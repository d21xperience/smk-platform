DO $$
DECLARE
    v_room_id UUID;
BEGIN
    -- Ambil 1 rombel/kelas secara acak untuk dijadikan tempat seed aset
    SELECT rombongan_belajar_id INTO v_room_id FROM tabel_kelas LIMIT 1;

    IF v_room_id IS NOT NULL THEN
        INSERT INTO inventory_items (id, name, brand, category, condition, room_id) VALUES
        ('YYS/SAR-COMP/2025/001', 'PC Komputer Rakitan Core i5', 'Asus', 'Elektronik', 'Bagus', v_room_id),
        ('YYS/SAR-COMP/2025/002', 'PC Komputer Rakitan Core i5', 'Asus', 'Elektronik', 'Bagus', v_room_id),
        ('YYS/SAR-PROJ/2024/012', 'Proyektor BenQ Wall Mount', 'BenQ', 'Elektronik', 'Rusak Ringan', v_room_id),
        ('YYS/SAR-MEBL/2023/105', 'Meja Komputer Kayu Panjang', 'Custom', 'Mebel', 'Bagus', v_room_id)
        ON CONFLICT (id) DO NOTHING;
    END IF;
END $$;
