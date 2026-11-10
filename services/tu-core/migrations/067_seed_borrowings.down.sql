-- Kembalikan status aset yang dipinjam
UPDATE inventory_items SET status_aset = 'AKTIF', current_borrowing_id = NULL WHERE status_aset = 'DIPINJAM';
DELETE FROM borrowing_fines;
DELETE FROM borrowing_approvals;
DELETE FROM borrowings;
