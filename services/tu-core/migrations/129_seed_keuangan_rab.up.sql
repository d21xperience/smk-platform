-- =====================================================================
-- Seed: RAB (tahun_anggaran 2026) + COA children + migrate BKU
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. Tambah COA children untuk belanja (5.1.1, 5.1.2, 5.1.3)
-- ---------------------------------------------------------------------
INSERT INTO "kode_rekening" (id, school_id, kode, nama, tipe, parent_id, level, is_aktif) VALUES
('f8000000-0000-0002-0000-000000000001', '11111111-1111-1111-1111-111111111111', '5.1.1', 'Belanja ATK & Operasional', 'BELANJA', 'f3000000-0000-0000-0000-000000000041', 3, TRUE),
('f8000000-0000-0002-0000-000000000002', '11111111-1111-1111-1111-111111111111', '5.1.2', 'Belanja Pemeliharaan',      'BELANJA', 'f3000000-0000-0000-0000-000000000041', 3, TRUE),
('f8000000-0000-0002-0000-000000000003', '11111111-1111-1111-1111-111111111111', '5.1.3', 'Belanja Modal',             'BELANJA', 'f3000000-0000-0000-0000-000000000041', 3, TRUE)
ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------
-- 2. Migrate BKU existing (dari seed 127) ke child 5.1.1
-- agar realisasi bisa dilacak per kategori
-- ---------------------------------------------------------------------
UPDATE "buku_kas_umum"
SET kode_rekening_id = 'f8000000-0000-0002-0000-000000000001'
WHERE kode_rekening_id = 'f3000000-0000-0000-0000-000000000041'
  AND referensi_tabel = 'pengeluaran';

-- ---------------------------------------------------------------------
-- 3. Seed rab_header (1 RAB DISAHKAN untuk tahun 2026)
-- Pendapatan = 200jt, Belanja = 200jt, Surplus = 0
-- ---------------------------------------------------------------------
INSERT INTO "rab_header" (
    id, school_id, academic_period_id, tahun_anggaran,
    total_pendapatan, total_belanja, surplus_defisit, status,
    disusun_oleh, tanggal_penyusunan, disahkan_oleh, tanggal_pengesahan,
    catatan
) VALUES (
    'f9000000-0000-0000-0000-000000000001',
    '11111111-1111-1111-1111-111111111111',
    '20261', 2026,
    200000000.00, 200000000.00, 0.00, 'DISAHKAN',
    '11111111-1111-1111-1111-111111111111', CURRENT_DATE,
    '11111111-1111-1111-1111-111111111111', CURRENT_DATE,
    'RAB awal tahun anggaran 2026'
) ON CONFLICT DO NOTHING;

-- ---------------------------------------------------------------------
-- 4. Seed rab_detail (5 item: 2 pendapatan + 3 belanja)
-- Subtotal = volume × harga_satuan (volume=1 untuk semua → subtotal=harga)
-- ---------------------------------------------------------------------
INSERT INTO "rab_detail" (
    id, rab_header_id, kode_rekening_id, uraian, volume, satuan,
    harga_satuan, subtotal, sumber_dana
) VALUES
('fa000000-0000-0001-0000-000000000001', 'f9000000-0000-0000-0000-000000000001',
 'f3000000-0000-0000-0000-000000000032', 'Penerimaan SPP tahun 2026',
 1, 'tahun', 100000000.00, 100000000.00, 'MANDIRI'),
('fa000000-0000-0001-0000-000000000002', 'f9000000-0000-0000-0000-000000000001',
 'f3000000-0000-0000-0000-000000000033', 'Dana BOS tahun 2026',
 1, 'tahun', 100000000.00, 100000000.00, 'BOS'),
('fa000000-0000-0001-0000-000000000011', 'f9000000-0000-0000-0000-000000000001',
 'f8000000-0000-0002-0000-000000000001', 'Belanja ATK & operasional harian',
 1, 'tahun', 50000000.00, 50000000.00, 'BOS'),
('fa000000-0000-0001-0000-000000000012', 'f9000000-0000-0000-0000-000000000001',
 'f8000000-0000-0002-0000-000000000002', 'Pemeliharaan sarana sekolah',
 1, 'tahun', 80000000.00, 80000000.00, 'BOS'),
('fa000000-0000-0001-0000-000000000013', 'f9000000-0000-0000-0000-000000000001',
 'f8000000-0000-0002-0000-000000000003', 'Belanja modal (peralatan)',
 1, 'tahun', 70000000.00, 70000000.00, 'BOS')
ON CONFLICT DO NOTHING;
