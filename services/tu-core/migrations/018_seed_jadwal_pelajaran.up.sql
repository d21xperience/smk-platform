-- Seed jadwal pelajaran demo untuk rombel "X TKJ 1" (aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa)
-- Mereferensikan pembelajaran_id yang sudah ada dari migrasi 016.
-- Header yang dibutuhkan frontend: x-school-id=1111..., x-academic-period-id=20261.
-- Idempotent via ON CONFLICT (id) DO UPDATE.

INSERT INTO tabel_jadwal_pelajaran (id, pembelajaran_id, hari, jam_mulai, jam_selesai) VALUES
    ('cccccccc-cccc-cccc-cccc-000000000001', '99999999-9999-9999-9999-999999999991', 'Senin',  '07:00', '08:30'),
    ('cccccccc-cccc-cccc-cccc-000000000002', '99999999-9999-9999-9999-999999999992', 'Senin',  '08:30', '10:00'),
    ('cccccccc-cccc-cccc-cccc-000000000003', '99999999-9999-9999-9999-999999999993', 'Senin',  '10:15', '11:45'),
    ('cccccccc-cccc-cccc-cccc-000000000004', '99999999-9999-9999-9999-999999999991', 'Selasa', '07:00', '08:30'),
    ('cccccccc-cccc-cccc-cccc-000000000005', '99999999-9999-9999-9999-999999999992', 'Selasa', '08:30', '10:00'),
    ('cccccccc-cccc-cccc-cccc-000000000006', '99999999-9999-9999-9999-999999999993', 'Rabu',   '07:00', '08:30'),
    ('cccccccc-cccc-cccc-cccc-000000000007', '99999999-9999-9999-9999-999999999991', 'Kamis',  '07:00', '08:30'),
    ('cccccccc-cccc-cccc-cccc-000000000008', '99999999-9999-9999-9999-999999999992', 'Jumat',  '07:00', '08:30')
ON CONFLICT (id) DO UPDATE SET
    pembelajaran_id = EXCLUDED.pembelajaran_id,
    hari            = EXCLUDED.hari,
    jam_mulai       = EXCLUDED.jam_mulai,
    jam_selesai     = EXCLUDED.jam_selesai;

-- Validasi: pastikan pembelajaran yang dirujuk memang ada, agar seed tidak diam-diam gagal FK.
DO $$
DECLARE
    missing INT;
BEGIN
    SELECT count(*) INTO missing
    FROM (VALUES
        ('99999999-9999-9999-9999-999999999991'::uuid),
        ('99999999-9999-9999-9999-999999999992'::uuid),
        ('99999999-9999-9999-9999-999999999993'::uuid)
    ) AS v(pid)
    WHERE NOT EXISTS (
        SELECT 1 FROM tabel_pembelajaran tp WHERE tp.pembelajaran_id = v.pid
    );

    IF missing > 0 THEN
        RAISE EXCEPTION 'Seed 018 gagal: % pembelajaran_id tidak ditemukan. Jalankan migrasi 016 dulu.', missing;
    END IF;
END $$;
