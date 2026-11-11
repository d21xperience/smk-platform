-- 1. Tabel Auto-Numbering SK
CREATE TABLE IF NOT EXISTS sk_numbers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    jenis VARCHAR(50) NOT NULL,
    tahun INT NOT NULL,
    bulan INT NOT NULL,
    sequence INT NOT NULL,
    nomor_sk VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_sk_numbers_lookup ON sk_numbers(jenis, tahun, bulan);

-- 2. Tabel Header SK
CREATE TABLE IF NOT EXISTS sk_pembagian_tugas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nomor_sk VARCHAR(100) NOT NULL UNIQUE,
    tahun_ajaran VARCHAR(10) NOT NULL,
    semester INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT', -- DRAFT, SUBMITTED, APPROVED, ISSUED, REJECTED, REVISION
    draft_pdf_url TEXT,
    final_pdf_url TEXT,
    lampiran_pdf_url TEXT,
    qr_tte_url TEXT,
    created_by UUID,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    rejection_reason TEXT,
    revision_notes TEXT
);

-- 3. Tabel Detail Pembagian Tugas
CREATE TABLE IF NOT EXISTS sk_pembagian_tugas_detail (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sk_id UUID NOT NULL REFERENCES sk_pembagian_tugas(id) ON DELETE CASCADE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    mata_pelajaran_id INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- 4. Tabel Kelas & Jam per Detail
CREATE TABLE IF NOT EXISTS sk_pembagian_tugas_kelas (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    detail_id UUID NOT NULL REFERENCES sk_pembagian_tugas_detail(id) ON DELETE CASCADE,
    rombongan_belajar_id UUID NOT NULL REFERENCES tabel_kelas(rombongan_belajar_id) ON DELETE CASCADE,
    jam_per_minggu INT NOT NULL
);

-- 5. Tabel Tugas Tambahan
CREATE TABLE IF NOT EXISTS sk_pembagian_tugas_tambahan (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    sk_id UUID NOT NULL REFERENCES sk_pembagian_tugas(id) ON DELETE CASCADE,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    jenis_tugas VARCHAR(50) NOT NULL,
    referensi_id UUID,
    keterangan TEXT
);
