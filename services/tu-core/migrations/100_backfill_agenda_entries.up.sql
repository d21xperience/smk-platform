DO $$
DECLARE
    r RECORD;
    v_year INT;
    v_month INT;
    v_nomor_urut INT;
    v_nomor_agenda VARCHAR;
    v_roman VARCHAR;
    v_type VARCHAR(10);
    v_school_id UUID;
BEGIN
    -- Ambil school_id dari salah satu record atau hardcode sesuai env
    SELECT school_id INTO v_school_id FROM correspondences LIMIT 1;
    IF v_school_id IS NULL THEN
        v_school_id := '11111111-1111-1111-1111-111111111111';
    END IF;

    FOR r IN SELECT id, type, date FROM correspondences WHERE id NOT IN (SELECT correspondence_id FROM agenda_entries) ORDER BY date ASC LOOP
        v_year := EXTRACT(YEAR FROM r.date);
        v_month := EXTRACT(MONTH FROM r.date);
        v_type := CASE WHEN r.type = 'INCOMING' THEN 'IN' ELSE 'OUT' END;

        -- Generate Sequence
        INSERT INTO agenda_sequences (school_id, jenis, tahun, last_number)
        VALUES (v_school_id, v_type, v_year, 1)
        ON CONFLICT (school_id, jenis, tahun)
        DO UPDATE SET last_number = agenda_sequences.last_number + 1
        RETURNING last_number INTO v_nomor_urut;

        v_roman := CASE v_month
            WHEN 1 THEN 'I' WHEN 2 THEN 'II' WHEN 3 THEN 'III' WHEN 4 THEN 'IV'
            WHEN 5 THEN 'V' WHEN 6 THEN 'VI' WHEN 7 THEN 'VII' WHEN 8 THEN 'VIII'
            WHEN 9 THEN 'IX' WHEN 10 THEN 'X' WHEN 11 THEN 'XI' WHEN 12 THEN 'XII'
        END;

        v_nomor_agenda := LPAD(v_nomor_urut::text, 3, '0') || '/AGD-' || v_type || '/' || v_roman || '/' || v_year;

        INSERT INTO agenda_entries (
            id, school_id, correspondence_id, jenis, nomor_agenda, nomor_urut, tahun, tanggal_catat, petugas_id
        ) VALUES (
            gen_random_uuid(), v_school_id, r.id, v_type, v_nomor_agenda, v_nomor_urut, v_year, NOW(), '11111111-1111-1111-1111-111111111111'
        );
    END LOOP;
END $$;
