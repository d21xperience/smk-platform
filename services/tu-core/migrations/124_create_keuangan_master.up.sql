CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- =====================================================================
-- Migration 124: Create Master Tables for Keuangan Module
-- Tables: jenis_iuran, tarif_iuran, kode_rekening, rekening_bank
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. jenis_iuran (Master Jenis Iuran)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "jenis_iuran" (
    "id"         UUID NOT NULL PRIMARY KEY,
    "school_id"  VARCHAR(36) NOT NULL,
    "kode"       VARCHAR(30) NOT NULL,
    "nama"       VARCHAR(150) NOT NULL,
    "kategori"   VARCHAR(30) NOT NULL,
    "is_wajib"   BOOLEAN NOT NULL DEFAULT TRUE,
    "is_aktif"   BOOLEAN NOT NULL DEFAULT TRUE,
    "keterangan" TEXT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_jenis_iuran_kode UNIQUE ("school_id", "kode"),
    CONSTRAINT chk_jenis_iuran_kategori CHECK (
        "kategori" IN ('SPP', 'GEDUNG', 'EKSTRA', 'SERAGAM', 'KEGIATAN', 'LAINNYA')
    )
);

CREATE INDEX IF NOT EXISTS idx_jenis_iuran_school ON jenis_iuran(school_id);
CREATE INDEX IF NOT EXISTS idx_jenis_iuran_kategori ON jenis_iuran(school_id, kategori);

-- ---------------------------------------------------------------------
-- 2. tarif_iuran (Tarif per kelas & tahun ajaran)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "tarif_iuran" (
    "id"                  UUID NOT NULL PRIMARY KEY,
    "school_id"           VARCHAR(36) NOT NULL,
    "academic_period_id"  VARCHAR(10) NOT NULL,
    "jenis_iuran_id"      UUID NOT NULL REFERENCES jenis_iuran(id) ON DELETE CASCADE,
    "kelas"               VARCHAR(50) NULL,
    "nominal"             NUMERIC(15,2) NOT NULL DEFAULT 0,
    "tahun_berlaku"       INTEGER NOT NULL,
    "created_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_tarif_iuran_school_period ON tarif_iuran(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_tarif_iuran_jenis ON tarif_iuran(jenis_iuran_id);
CREATE UNIQUE INDEX IF NOT EXISTS uq_tarif_iuran_unique
    ON tarif_iuran(school_id, academic_period_id, jenis_iuran_id, COALESCE(kelas, ''));

-- ---------------------------------------------------------------------
-- 3. kode_rekening (Chart of Accounts - COA)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "kode_rekening" (
    "id"         UUID NOT NULL PRIMARY KEY,
    "school_id"  VARCHAR(36) NOT NULL,
    "kode"       VARCHAR(30) NOT NULL,
    "nama"       VARCHAR(200) NOT NULL,
    "tipe"       VARCHAR(20) NOT NULL,
    "parent_id"  UUID NULL REFERENCES kode_rekening(id) ON DELETE SET NULL,
    "level"      INTEGER NOT NULL DEFAULT 1,
    "is_aktif"   BOOLEAN NOT NULL DEFAULT TRUE,
    "keterangan" TEXT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_kode_rekening_kode UNIQUE ("school_id", "kode"),
    CONSTRAINT chk_kode_rekening_tipe CHECK (
        "tipe" IN ('ASET', 'KEWAJIBAN', 'EKUITAS', 'PENDAPATAN', 'BELANJA')
    )
);

CREATE INDEX IF NOT EXISTS idx_kode_rekening_school ON kode_rekening(school_id);
CREATE INDEX IF NOT EXISTS idx_kode_rekening_tipe ON kode_rekening(school_id, tipe);
CREATE INDEX IF NOT EXISTS idx_kode_rekening_parent ON kode_rekening(parent_id);

-- ---------------------------------------------------------------------
-- 4. rekening_bank (Master Rekening Bank Sekolah)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rekening_bank" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "school_id"          VARCHAR(36) NOT NULL,
    "nama_bank"          VARCHAR(100) NOT NULL,
    "nomor_rekening"     VARCHAR(50) NOT NULL,
    "atas_nama"          VARCHAR(150) NOT NULL,
    "cabang"             VARCHAR(150) NULL,
    "saldo_awal"         NUMERIC(15,2) NOT NULL DEFAULT 0,
    "saldo_terkini"      NUMERIC(15,2) NOT NULL DEFAULT 0,
    "kode_rekening_id"   UUID NULL REFERENCES kode_rekening(id) ON DELETE SET NULL,
    "format_file_mutasi" VARCHAR(20) NULL,
    "is_aktif"           BOOLEAN NOT NULL DEFAULT TRUE,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT uq_rekening_bank_norek UNIQUE ("school_id", "nomor_rekening")
);

CREATE INDEX IF NOT EXISTS idx_rekening_bank_school ON rekening_bank(school_id);
