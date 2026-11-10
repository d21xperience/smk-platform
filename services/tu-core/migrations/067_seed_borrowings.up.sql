DO $$
DECLARE
    v_asset_1 VARCHAR(100);
    v_asset_2 VARCHAR(100);
    v_asset_3 VARCHAR(100);
    v_asset_4 VARCHAR(100);
    v_bw_1 UUID;
    v_bw_2 UUID;
    v_bw_3 UUID;
    v_bw_4 UUID;
BEGIN
    -- Ambil 4 aset untuk 4 skenario berbeda
    SELECT id INTO v_asset_1 FROM inventory_items WHERE status_aset = 'AKTIF' LIMIT 1;
    SELECT id INTO v_asset_2 FROM inventory_items WHERE status_aset = 'AKTIF' OFFSET 1 LIMIT 1;
    SELECT id INTO v_asset_3 FROM inventory_items WHERE status_aset = 'AKTIF' OFFSET 2 LIMIT 1;
    SELECT id INTO v_asset_4 FROM inventory_items WHERE status_aset = 'AKTIF' OFFSET 3 LIMIT 1;

    IF v_asset_1 IS NOT NULL THEN
        -- 1. PENDING_APPROVAL (Baru diajukan, belum disetujui)
        INSERT INTO borrowings (id, borrowing_number, borrower_id, borrower_name, inventory_item_id, purpose_type, purpose_description, borrow_date, expected_return_date, usage_location, duration_days, approval_level, status)
        VALUES (gen_random_uuid(), 'PJM/SMKPASJA/2026/001', 'user-guru-1', 'Pak Hermawan, S.Pd', v_asset_1, 'PEMBELAJARAN', 'Untuk mengajar materi Jaringan Komputer', CURRENT_DATE, CURRENT_DATE + 2, 'Lab TKJ 1', 3, 2, 'PENDING_APPROVAL');
    END IF;

    IF v_asset_2 IS NOT NULL THEN
        -- 2. APPROVED (Sudah disetujui, belum diambil)
        INSERT INTO borrowings (id, borrowing_number, borrower_id, borrower_name, inventory_item_id, purpose_type, purpose_description, borrow_date, expected_return_date, usage_location, duration_days, approval_level, status, current_approval_step)
        VALUES (gen_random_uuid(), 'PJM/SMKPASJA/2026/002', 'user-guru-2', 'Bu Linda, M.Si', v_asset_2, 'RAPAT', 'Untuk presentasi Rapat Dinas', CURRENT_DATE, CURRENT_DATE + 1, 'Ruang Rapat Utama', 2, 2, 'APPROVED', 2) RETURNING id INTO v_bw_2;

        -- Simulasi approval sudah selesai (step 1 & 2)
        INSERT INTO borrowing_approvals (borrowing_id, approval_step, approver_id, decision, notes) VALUES (v_bw_2, 1, 'user-pjr-1', 'APPROVED', 'Silakan digunakan dengan baik.');
        INSERT INTO borrowing_approvals (borrowing_id, approval_step, approver_id, decision, notes) VALUES (v_bw_2, 2, 'user-ktu-1', 'APPROVED', 'Disetujui untuk rapat.');
    END IF;

    IF v_asset_3 IS NOT NULL THEN
        -- 3. PICKED_UP (Sedang dipinjam)
        INSERT INTO borrowings (id, borrowing_number, borrower_id, borrower_name, inventory_item_id, purpose_type, purpose_description, borrow_date, expected_return_date, usage_location, duration_days, approval_level, status, current_approval_step, actual_pickup_date, pickup_by, pickup_condition, bast_number)
        VALUES (gen_random_uuid(), 'PJM/SMKPASJA/2026/003', 'user-guru-3', 'Pak Danang, S.T', v_asset_3, 'LOMBA', 'Untuk kegiatan LKS tingkat Kabupaten', CURRENT_DATE - 1, CURRENT_DATE + 3, 'SMK Negeri 1 Bandung', 5, 3, 'PICKED_UP', 3, CURRENT_DATE - 1, 'user-sarpras-1', 'BAIK', 'BAST-PJM/SMKPASJA/2026/001');
        
        -- Update status aset jadi DIPINJAM
        UPDATE inventory_items SET status_aset = 'DIPINJAM' WHERE id = v_asset_3;
    END IF;

    IF v_asset_4 IS NOT NULL THEN
        -- 4. RETURNED (Sudah dikembalikan, ada denda keterlambatan)
        INSERT INTO borrowings (id, borrowing_number, borrower_id, borrower_name, inventory_item_id, purpose_type, purpose_description, borrow_date, expected_return_date, usage_location, duration_days, approval_level, status, current_approval_step, actual_pickup_date, pickup_by, pickup_condition, bast_number, actual_return_date, return_by, return_condition, days_late, fine_amount, fine_paid, bap_number)
        VALUES (gen_random_uuid(), 'PJM/SMKPASJA/2026/004', 'user-guru-4', 'Ibu Nenden, S.Pd', v_asset_4, 'KEGIATAN', 'Untuk kegiatan Pramuka', CURRENT_DATE - 10, CURRENT_DATE - 5, 'Lapangan Sekolah', 6, 3, 'RETURNED', 3, CURRENT_DATE - 10, 'user-sarpras-1', 'BAIK', 'BAST-PJM/SMKPASJA/2026/002', CURRENT_DATE - 3, 'user-sarpras-1', 'BAIK', 2, 20000, FALSE, 'BAP-PJM/SMKPASJA/2026/001')
        RETURNING id INTO v_bw_4;

        -- Insert denda
        INSERT INTO borrowing_fines (borrowing_id, borrower_id, fine_amount, days_late, fine_per_day, fine_date, status)
        VALUES (v_bw_4, 'user-guru-4', 20000, 2, 10000, CURRENT_DATE - 3, 'UNPAID');
    END IF;
END $$;
