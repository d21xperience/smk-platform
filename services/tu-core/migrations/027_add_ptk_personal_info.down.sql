ALTER TABLE tabel_ptk 
    DROP COLUMN IF EXISTS nik,
    DROP COLUMN IF EXISTS agama_id,
    DROP COLUMN IF EXISTS status_kepegawaian_id,
    DROP COLUMN IF EXISTS foto_url,
    DROP COLUMN IF EXISTS no_hp,
    DROP COLUMN IF EXISTS email,
    DROP COLUMN IF EXISTS status_perkawinan,
    DROP COLUMN IF EXISTS nama_ibu_kandung;

DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'users') THEN
        ALTER TABLE users DROP COLUMN IF EXISTS ptk_id;
    END IF;
END $$;
