DROP TABLE IF EXISTS inventory_item_history;
ALTER TABLE inventory_items 
DROP COLUMN IF EXISTS nomor_registrasi,
DROP COLUMN IF EXISTS sumber_perolehan,
DROP COLUMN IF EXISTS keterangan,
DROP COLUMN IF EXISTS masa_manfaat_tahun,
DROP COLUMN IF EXISTS nilai_buku,
DROP COLUMN IF EXISTS status_aset,
DROP COLUMN IF EXISTS tanggal_pencatatan,
DROP COLUMN IF EXISTS tanggal_penghapusan,
DROP COLUMN IF EXISTS alasan_penghapusan;
