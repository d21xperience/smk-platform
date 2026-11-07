-- 1. Bersihkan duplikat ptk_id di tabel_ptk_pelengkap (simpan baris terakhir)
DELETE FROM tabel_ptk_pelengkap a
USING tabel_ptk_pelengkap b
WHERE a.ptk_id = b.ptk_id 
  AND a.ctid < b.ctid;

-- 2. Tambahkan UNIQUE constraint agar ON CONFLICT berfungsi
ALTER TABLE tabel_ptk_pelengkap 
ADD CONSTRAINT tabel_ptk_pelengkap_ptk_id_key UNIQUE (ptk_id);
