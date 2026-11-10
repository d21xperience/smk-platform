DO $$
DECLARE
    r RECORD;
    v_count INT;
    v_year INT;
    v_sch VARCHAR(20);
    v_month_roman VARCHAR(5);
BEGIN
    SELECT COALESCE(kode_sekolah, 'SMKPasja') INTO v_sch FROM tabel_sekolah LIMIT 1;

    FOR r IN SELECT id, purchase_date FROM inventory_items WHERE nomor_registrasi IS NULL OR nomor_registrasi = '' ORDER BY purchase_date ASC, id ASC LOOP
        v_year := EXTRACT(YEAR FROM r.purchase_date);
        
        SELECT COUNT(*) INTO v_count FROM inventory_items WHERE nomor_registrasi IS NOT NULL AND nomor_registrasi != '' AND EXTRACT(YEAR FROM purchase_date) = v_year;
        
        CASE EXTRACT(MONTH FROM r.purchase_date)
            WHEN 1 THEN v_month_roman := 'I';
            WHEN 2 THEN v_month_roman := 'II';
            WHEN 3 THEN v_month_roman := 'III';
            WHEN 4 THEN v_month_roman := 'IV';
            WHEN 5 THEN v_month_roman := 'V';
            WHEN 6 THEN v_month_roman := 'VI';
            WHEN 7 THEN v_month_roman := 'VII';
            WHEN 8 THEN v_month_roman := 'VIII';
            WHEN 9 THEN v_month_roman := 'IX';
            WHEN 10 THEN v_month_roman := 'X';
            WHEN 11 THEN v_month_roman := 'XI';
            WHEN 12 THEN v_month_roman := 'XII';
        END CASE;

        UPDATE inventory_items 
        SET nomor_registrasi = LPAD((v_count + 1)::text, 3, '0') || '/INV/' || v_sch || '/' || v_month_roman || '/' || v_year
        WHERE id = r.id;
    END LOOP;
END $$;
