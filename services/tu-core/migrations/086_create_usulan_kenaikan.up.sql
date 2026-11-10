-- 1. Tabel Utama Usulan Kenaikan
CREATE TABLE IF NOT EXISTS usulan_kenaikan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nomor_usulan VARCHAR(100) NOT NULL UNIQUE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis VARCHAR(50) NOT NULL, -- PANGKAT_REGULER, PANGKAT_PILIHAN, KGB, SATYALANCANA
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT', 
    -- Status: DRAFT, LENGKAP, SUBMITTED_KEPSEK, APPROVED_KEPSEK, SUBMITTED_YAYASAN, APPROVED_YAYASAN, SUBMITTED_DINAS, PROSES_BKN, SK_TERBIT, DITOLAK, REVISI
    
    golongan_sekarang_id INT REFERENCES master_golongan(id),
    golongan_baru_id INT REFERENCES master_golongan(id),
    tmt_baru DATE,
    
    nomor_sk_terbit VARCHAR(100),
    tanggal_sk_terbit DATE,
    scan_sk_url TEXT,
    
    catatan_ktu TEXT,
    catatan_kepsek TEXT,
    alasan_penolakan TEXT,
    
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_usulan_ptk_status ON usulan_kenaikan(ptk_id, status);

-- 2. Tabel Batch KGB (Untuk Usulan Massal)
CREATE TABLE IF NOT EXISTS usulan_kgb_batch (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nomor_batch VARCHAR(100) NOT NULL UNIQUE,
    total_pegawai INT NOT NULL DEFAULT 0,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT',
    catatan_ktu TEXT,
    surat_pengantar_url TEXT,
    daftar_usul_url TEXT,
    rekap_url TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Relasi Batch ke Usulan Individual
ALTER TABLE usulan_kenaikan ADD COLUMN IF NOT EXISTS batch_id UUID REFERENCES usulan_kgb_batch(id);

-- 3. Tabel Dokumen Pendukung Usulan
CREATE TABLE IF NOT EXISTS usulan_dokumen (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usulan_id UUID NOT NULL REFERENCES usulan_kenaikan(id) ON DELETE CASCADE,
    jenis_dokumen VARCHAR(50) NOT NULL, -- SK_PANGKAT_TERAKHIR, SKP_2_TAHUN, SURAT_PERNYATAAN, PAK, IJAZAH
    file_url TEXT NOT NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Tabel Timeline (History Aksi)
CREATE TABLE IF NOT EXISTS usulan_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usulan_id UUID NOT NULL REFERENCES usulan_kenaikan(id) ON DELETE CASCADE,
    status_dari VARCHAR(30),
    status_ke VARCHAR(30) NOT NULL,
    action_by UUID,
    action_at TIMESTAMPTZ DEFAULT NOW(),
    catatan TEXT
);

-- 5. Seed Data Dummy untuk Pengujian (Menggunakan Data Guru Existing)
-- Budi Santoso (III/c -> III/d)
INSERT INTO usulan_kenaikan (
    id, nomor_usulan, ptk_id, jenis, status, 
    golongan_sekarang_id, golongan_baru_id, tmt_baru, catatan_ktu, created_by
) VALUES (
    '11111111-1111-1111-1111-111111111111', 
    'USUL/PANGKAT/2026/XI/001', 
    'a1b2c3d4-e5f6-7a8b-9c0d-1e2f3a4b5c6d', 
    'PANGKAT_REGULER', 'DRAFT', 
    33, 34, '2026-04-01', 
    'Mohon diproses untuk kenaikan pangkat reguler.', 
    '123e4567-e89b-12d3-a456-426614174000'
) ON CONFLICT (id) DO NOTHING;

-- Timeline Awal
INSERT INTO usulan_timeline (
    usulan_id, status_dari, status_ke, action_by, catatan
) VALUES (
    '11111111-1111-1111-1111-111111111111', 
    NULL, 'DRAFT', 
    '123e4567-e89b-12d3-a456-426614174000', 
    'Usulan dibuat oleh KTU.'
) ON CONFLICT DO NOTHING;
