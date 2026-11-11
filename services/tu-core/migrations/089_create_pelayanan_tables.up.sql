-- 1. Master Katalog Layanan
CREATE TABLE IF NOT EXISTS master_layanan (
    kode VARCHAR(30) PRIMARY KEY, -- SK-AKTIF, SK-GAJI, LG-IJAZAH, dll
    nama VARCHAR(100) NOT NULL,
    kategori VARCHAR(50) NOT NULL, -- SURAT_KETERANGAN, SURAT_TUGAS, LEGALISIR, REKOMENDASI
    deskripsi TEXT,
    sla_hari INT NOT NULL DEFAULT 1,
    template_surat VARCHAR(50) NOT NULL, -- Nama template di Go (SK_AKTIF, SK_GAJI, dll)
    is_active BOOLEAN DEFAULT TRUE
);

-- Seed Katalog Layanan (Sesuai Spesifikasi)
INSERT INTO master_layanan (kode, nama, kategori, deskripsi, sla_hari, template_surat) VALUES
('SK-AKTIF', 'Surat Keterangan Aktif', 'SURAT_KETERANGAN', 'Keterangan pegawai masih aktif bekerja', 1, 'SK_AKTIF'),
('SK-GAJI', 'Surat Keterangan Gaji', 'SURAT_KETERANGAN', 'Keterangan besaran gaji untuk bank/kredit', 1, 'SK_GAJI'),
('SK-TUGAS', 'Surat Keterangan Tugas', 'SURAT_KETERANGAN', 'Keterangan tugas mengajar/tambahan', 1, 'SK_TUGAS'),
('SK-PENGHASILAN', 'Surat Keterangan Penghasilan', 'SURAT_KETERANGAN', 'Untuk kredit/pajak', 2, 'SK_PENGHASILAN'),
('ST-DINAS', 'Surat Tugas Dinas Luar', 'SURAT_TUGAS', 'Untuk perjalanan dinas', 1, 'ST_DINAS'),
('ST-DIKLAT', 'Surat Tugas Diklat', 'SURAT_TUGAS', 'Untuk mengikuti pelatihan', 2, 'ST_DIKLAT'),
('LG-SK', 'Legalisir SK', 'LEGALISIR', 'Legalisir SK pangkat/penugasan', 1, 'LEGALISIR'),
('LG-IJAZAH', 'Legalisir Ijazah', 'LEGALISIR', 'Legalisir ijazah terakhir', 1, 'LEGALISIR'),
('RC-DIKLAT', 'Rekomendasi Diklat', 'REKOMENDASI', 'Rekomendasi mengikuti diklat', 3, 'RC_DIKLAT'),
('RC-BEASISWA', 'Rekomendasi Beasiswa', 'REKOMENDASI', 'Rekomendasi beasiswa', 5, 'RC_BEASISWA')
ON CONFLICT (kode) DO NOTHING;

-- 2. Header Pengajuan Layanan
CREATE TABLE IF NOT EXISTS pelayanan_pengajuan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nomor_pengajuan VARCHAR(50) NOT NULL UNIQUE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    layanan_kode VARCHAR(30) NOT NULL REFERENCES master_layanan(kode) ON DELETE RESTRICT,
    status VARCHAR(30) NOT NULL DEFAULT 'DRAFT', 
    -- Status: DRAFT, SUBMITTED, VERIFIED, ISSUED, REJECTED, REVISION
    
    -- Detail Form
    tujuan_surat TEXT,
    kepada TEXT,
    keterangan TEXT,
    dokumen_pendukung_url TEXT,
    
    -- Tracking Waktu
    estimasi_selesai DATE,
    requested_at TIMESTAMPTZ DEFAULT NOW(),
    verified_at TIMESTAMPTZ,
    approved_at TIMESTAMPTZ,
    issued_at TIMESTAMPTZ,
    
    -- Aksi
    verified_by UUID REFERENCES tabel_ptk(ptk_id),
    approved_by UUID REFERENCES tabel_ptk(ptk_id),
    catatan_verifikasi TEXT,
    alasan_penolakan TEXT
);
CREATE INDEX idx_pelayanan_ptk ON pelayanan_pengajuan(ptk_id);
CREATE INDEX idx_pelayanan_status ON pelayanan_pengajuan(status);

-- 3. Timeline Log
CREATE TABLE IF NOT EXISTS pelayanan_timeline (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pengajuan_id UUID NOT NULL REFERENCES pelayanan_pengajuan(id) ON DELETE CASCADE,
    status_dari VARCHAR(30),
    status_ke VARCHAR(30) NOT NULL,
    action_by UUID REFERENCES tabel_ptk(ptk_id),
    action_at TIMESTAMPTZ DEFAULT NOW(),
    catatan TEXT
);

-- 4. Surat Terbit (Dokumen Final)
CREATE TABLE IF NOT EXISTS pelayanan_surat (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pengajuan_id UUID NOT NULL UNIQUE REFERENCES pelayanan_pengajuan(id) ON DELETE CASCADE,
    nomor_surat VARCHAR(100) NOT NULL UNIQUE,
    file_url TEXT NOT NULL,
    qr_url TEXT,
    issued_at TIMESTAMPTZ DEFAULT NOW()
);
