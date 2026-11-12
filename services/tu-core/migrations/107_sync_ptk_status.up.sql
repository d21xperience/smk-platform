-- Sinkronisasi status_kepegawaian_id dari ptk_buku_induk_ext.status_pegawai
UPDATE tabel_ptk tp
SET status_kepegawaian_id = 
    CASE ext.status_pegawai
        WHEN 'PNS' THEN 1
        WHEN 'PPPK' THEN 2
        WHEN 'GTY' THEN 3
        WHEN 'Honorer' THEN 4
        ELSE 0
    END
FROM ptk_buku_induk_ext ext
WHERE tp.ptk_id = ext.ptk_id;

-- Sinkronisasi jenis_ptk_id (Hanya jika kosong/0)
UPDATE tabel_ptk tp
SET jenis_ptk_id = 
    CASE ext.jenis_pegawai
        WHEN 'Guru' THEN 1
        WHEN 'Tendik' THEN 4
        ELSE 1
    END
FROM ptk_buku_induk_ext ext
WHERE tp.ptk_id = ext.ptk_id AND (tp.jenis_ptk_id IS NULL OR tp.jenis_ptk_id = 0);
