-- =====================================================================
-- Migration 126: Transaksi Inti Keuangan
-- Tabel: tagihan_siswa, penerimaan_iuran, buku_kas_umum,
--        pengajuan_pengeluaran, pengeluaran
-- =====================================================================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ---------------------------------------------------------------------
-- 1. tagihan_siswa — tagihan per siswa × jenis iuran × periode
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "tagihan_siswa" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "peserta_didik_id"   UUID NOT NULL REFERENCES tabel_siswa(peserta_didik_id) ON DELETE CASCADE,
    "nis"                VARCHAR(20) NOT NULL,
    "kelas_nama"         VARCHAR(50) NULL,
    "jenis_iuran_id"     UUID NOT NULL REFERENCES jenis_iuran(id) ON DELETE RESTRICT,
    "periode_bulan"      INTEGER NOT NULL,
    "periode_tahun"      INTEGER NOT NULL,
    "nominal"            NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_dibayar"      NUMERIC(15,2) NOT NULL DEFAULT 0,
    "status"             VARCHAR(20) NOT NULL DEFAULT 'BELUM_LUNAS',
    "tanggal_jatuh_tempo" DATE NULL,
    "tanggal_bayar"      DATE NULL,
    "created_by"         VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_tagihan_siswa_status CHECK (
        status IN ('BELUM_LUNAS', 'LUNAS', 'SEBAGIAN', 'DIBATALKAN')
    ),
    CONSTRAINT chk_tagihan_siswa_periode_bulan CHECK (periode_bulan BETWEEN 1 AND 12),
    CONSTRAINT uq_tagihan_siswa_unique UNIQUE
        (school_id, academic_period_id, nis, jenis_iuran_id, periode_bulan, periode_tahun)
);

CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_school_period ON tagihan_siswa(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_nis ON tagihan_siswa(nis);
CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_status ON tagihan_siswa(status);
CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_jenis ON tagihan_siswa(jenis_iuran_id);

-- ---------------------------------------------------------------------
-- 2. penerimaan_iuran — transaksi pembayaran iuran siswa (auto BKM)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "penerimaan_iuran" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "no_bkm"             VARCHAR(50) NOT NULL UNIQUE,
    "tanggal"            DATE NOT NULL,
    "tagihan_siswa_id"   UUID NOT NULL REFERENCES tagihan_siswa(id) ON DELETE RESTRICT,
    "peserta_didik_id"   UUID NOT NULL REFERENCES tabel_siswa(peserta_didik_id) ON DELETE RESTRICT,
    "nis"                VARCHAR(20) NOT NULL,
    "jenis_iuran_id"     UUID NOT NULL REFERENCES jenis_iuran(id) ON DELETE RESTRICT,
    "periode_bulan"      INTEGER NOT NULL,
    "periode_tahun"      INTEGER NOT NULL,
    "jumlah"             NUMERIC(15,2) NOT NULL,
    "metode_bayar"       VARCHAR(20) NOT NULL DEFAULT 'TUNAI',
    "no_referensi"       VARCHAR(100) NULL,
    "keterangan"         TEXT NULL,
    "status"             VARCHAR(20) NOT NULL DEFAULT 'LUNAS',
    "created_by"         VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_penerimaan_metode CHECK (metode_bayar IN ('TUNAI', 'TRANSFER', 'QRIS')),
    CONSTRAINT chk_penerimaan_status CHECK (status IN ('LUNAS', 'DIBATALKAN')),
    CONSTRAINT chk_penerimaan_bulan CHECK (periode_bulan BETWEEN 1 AND 12),
    CONSTRAINT chk_penerimaan_jumlah CHECK (jumlah > 0)
);

CREATE INDEX IF NOT EXISTS idx_penerimaan_school_period ON penerimaan_iuran(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_penerimaan_tanggal ON penerimaan_iuran(tanggal);
CREATE INDEX IF NOT EXISTS idx_penerimaan_nis ON penerimaan_iuran(nis);
CREATE INDEX IF NOT EXISTS idx_penerimaan_tagihan ON penerimaan_iuran(tagihan_siswa_id);

-- ---------------------------------------------------------------------
-- 3. buku_kas_umum (BKU) — setiap transaksi tercatat di sini
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "buku_kas_umum" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "tanggal"            DATE NOT NULL,
    "no_bukti"           VARCHAR(50) NOT NULL,
    "uraian"             TEXT NOT NULL,
    "kode_rekening_id"   UUID NOT NULL REFERENCES kode_rekening(id) ON DELETE RESTRICT,
    "rekening_bank_id"   UUID NULL REFERENCES rekening_bank(id) ON DELETE SET NULL,
    "debet"              NUMERIC(15,2) NOT NULL DEFAULT 0,
    "kredit"             NUMERIC(15,2) NOT NULL DEFAULT 0,
    "saldo_berjalan"     NUMERIC(15,2) NOT NULL DEFAULT 0,
    "referensi_tabel"    VARCHAR(50) NULL,
    "referensi_id"       UUID NULL,
    "is_locked"          BOOLEAN NOT NULL DEFAULT FALSE,
    "is_rekonsiliasi"    BOOLEAN NOT NULL DEFAULT FALSE,
    "created_by"         VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_bku_debet_kredit CHECK (
        (debet > 0 AND kredit = 0) OR (kredit > 0 AND debet = 0)
    )
);

CREATE INDEX IF NOT EXISTS idx_bku_school_period_tgl ON buku_kas_umum(school_id, academic_period_id, tanggal);
CREATE INDEX IF NOT EXISTS idx_bku_kode_rekening ON buku_kas_umum(kode_rekening_id);
CREATE INDEX IF NOT EXISTS idx_bku_referensi ON buku_kas_umum(referensi_tabel, referensi_id);
CREATE INDEX IF NOT EXISTS idx_bku_locked ON buku_kas_umum(is_locked) WHERE is_locked = TRUE;

-- ---------------------------------------------------------------------
-- 4. pengajuan_pengeluaran — workflow 3 tingkat (Opsi A)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "pengajuan_pengeluaran" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "nomor_pengajuan"    VARCHAR(50) NOT NULL UNIQUE,
    "tanggal"            DATE NOT NULL,
    "pemohon_user_id"    VARCHAR(36) NOT NULL REFERENCES user_login(user_id) ON DELETE RESTRICT,
    "keperluan"          TEXT NOT NULL,
    "jumlah"             NUMERIC(15,2) NOT NULL,
    "prioritas"          VARCHAR(20) NOT NULL DEFAULT 'SEDANG',
    "status"             VARCHAR(40) NOT NULL DEFAULT 'MENUNGGU_VERIFIKASI_BENDAHARA',
    "catatan_pemohon"    TEXT NULL,
    "catatan_bendahara"  TEXT NULL,
    "catatan_kepala"     TEXT NULL,
    "verified_by"        VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "verified_at"        TIMESTAMPTZ NULL,
    "approved_by"        VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "approved_at"        TIMESTAMPTZ NULL,
    "no_bkk"             VARCHAR(50) NULL,
    "pos_anggaran_id"    UUID NULL,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_pengajuan_jumlah CHECK (jumlah > 0),
    CONSTRAINT chk_pengajuan_prioritas CHECK (prioritas IN ('RENDAH', 'SEDANG', 'TINGGI', 'DARURAT')),
    CONSTRAINT chk_pengajuan_status CHECK (status IN (
        'MENUNGGU_VERIFIKASI_BENDAHARA',
        'DITOLAK_BENDAHARA',
        'PERLU_REVISI',
        'MENUNGGU_PERSETUJUAN_KEPALA',
        'DITOLAK_KEPALA',
        'DISETUJUI',
        'CAIR',
        'DIBATALKAN'
    ))
);

CREATE INDEX IF NOT EXISTS idx_pengajuan_school_period ON pengajuan_pengeluaran(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_pengajuan_status ON pengajuan_pengeluaran(status);
CREATE INDEX IF NOT EXISTS idx_pengajuan_pemohon ON pengajuan_pengeluaran(pemohon_user_id);

-- ---------------------------------------------------------------------
-- 5. pengeluaran — transaksi kas keluar (auto BKK saat approved)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "pengeluaran" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "academic_period_id" VARCHAR(10) NOT NULL,
    "no_bkk"             VARCHAR(50) NOT NULL UNIQUE,
    "tanggal"            DATE NOT NULL,
    "pengajuan_id"       UUID NOT NULL UNIQUE REFERENCES pengajuan_pengeluaran(id) ON DELETE RESTRICT,
    "jumlah"             NUMERIC(15,2) NOT NULL,
    "keperluan"          TEXT NOT NULL,
    "status"             VARCHAR(20) NOT NULL DEFAULT 'SIAP_CAIR',
    "created_by"         VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_pengeluaran_jumlah CHECK (jumlah > 0),
    CONSTRAINT chk_pengeluaran_status CHECK (status IN ('SIAP_CAIR', 'CAIR', 'DIBATALKAN'))
);

CREATE INDEX IF NOT EXISTS idx_pengeluaran_school_period ON pengeluaran(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_pengeluaran_tanggal ON pengeluaran(tanggal);
CREATE INDEX IF NOT EXISTS idx_pengeluaran_status ON pengeluaran(status);
