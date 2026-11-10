-- Kembalikan status aset ke AKTIF (jika ada yang diubah)
UPDATE inventory_items SET status_aset = 'AKTIF', disposal_request_id = NULL WHERE disposal_request_id IS NOT NULL;

-- Hapus data usulan
DELETE FROM disposal_request_items;
DELETE FROM disposal_requests;
