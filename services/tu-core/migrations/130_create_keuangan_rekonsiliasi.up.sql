-- =====================================================================
-- Migration 130: Rekonsiliasi Bank (header + detail + unmatched)
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. rekonsiliasi_bank — 1 rekonsiliasi per school×period×rekening×bulan×tahun
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rekonsiliasi_bank" (
    "id"                        UUID NOT NULL PRIMARY KEY,
    "school_id"                 VARCHAR(36) NOT NULL,
    "academic_period_id"        VARCHAR(10) NOT NULL,
    "rekening_bank_id"          UUID NOT NULL REFERENCES rekening_bank(id) ON DELETE RESTRICT,
    "periode_bulan"             INTEGER NOT NULL,
    "periode_tahun"             INTEGER NOT NULL,
    "saldo_buku"                NUMERIC(15,2) NOT NULL DEFAULT 0,
    "saldo_bank"                NUMERIC(15,2) NOT NULL DEFAULT 0,
    "selisih"                   NUMERIC(15,2) NOT NULL DEFAULT 0,
    "jumlah_matched"            INTEGER NOT NULL DEFAULT 0,
    "jumlah_unmatched_bank"     INTEGER NOT NULL DEFAULT 0,
    "jumlah_unmatched_buku"     INTEGER NOT NULL DEFAULT 0,
    "saldo_rekonsiliasi_bank"   NUMERIC(15,2) NOT NULL DEFAULT 0,
    "saldo_rekonsiliasi_buku"   NUMERIC(15,2) NOT NULL DEFAULT 0,
    "status"                    VARCHAR(20) NOT NULL DEFAULT 'UNBALANCE',
    "file_mutasi_url"           TEXT NULL,
    "catatan"                   TEXT NULL,
    "created_by"                VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "created_at"                TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rekonsiliasi_status CHECK (status IN ('BALANCE', 'UNBALANCE')),
    CONSTRAINT chk_rekonsiliasi_bulan CHECK (periode_bulan BETWEEN 1 AND 12),
    CONSTRAINT uq_rekonsiliasi_unique UNIQUE
        (school_id, academic_period_id, rekening_bank_id, periode_bulan, periode_tahun)
);

CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_school ON rekonsiliasi_bank(school_id);
CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_period
    ON rekonsiliasi_bank(school_id, academic_period_id, periode_tahun, periode_bulan);
CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_rekening ON rekonsiliasi_bank(rekening_bank_id);

-- ---------------------------------------------------------------------
-- 2. rekonsiliasi_detail — matched pairs (mutasi bank ↔ buku kas)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rekonsiliasi_detail" (
    "id"                 UUID NOT NULL PRIMARY KEY,
    "rekonsiliasi_id"    UUID NOT NULL REFERENCES rekonsiliasi_bank(id) ON DELETE CASCADE,
    "mutasi_tanggal"     DATE NOT NULL,
    "mutasi_keterangan"  TEXT NULL,
    "mutasi_debet"       NUMERIC(15,2) NOT NULL DEFAULT 0,
    "mutasi_kredit"      NUMERIC(15,2) NOT NULL DEFAULT 0,
    "mutasi_saldo"       NUMERIC(15,2) NULL,
    "mutasi_referensi"   VARCHAR(100) NULL,
    "buku_kas_umum_id"   UUID NOT NULL REFERENCES buku_kas_umum(id) ON DELETE RESTRICT,
    "selisih_hari"       INTEGER NOT NULL DEFAULT 0,
    "created_at"         TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_detail_header ON rekonsiliasi_detail(rekonsiliasi_id);
CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_detail_bku ON rekonsiliasi_detail(buku_kas_umum_id);

-- ---------------------------------------------------------------------
-- 3. rekonsiliasi_unmatched — item yang belum match (bank atau buku)
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rekonsiliasi_unmatched" (
    "id"                    UUID NOT NULL PRIMARY KEY,
    "rekonsiliasi_id"       UUID NOT NULL REFERENCES rekonsiliasi_bank(id) ON DELETE CASCADE,
    "tipe"                  VARCHAR(10) NOT NULL,
    "data_json"             JSONB NOT NULL,
    "status_tindak_lanjut"  VARCHAR(20) NOT NULL DEFAULT 'PENDING',
    "created_at"            TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rekonsiliasi_unmatched_tipe CHECK (tipe IN ('BANK', 'BUKU')),
    CONSTRAINT chk_rekonsiliasi_unmatched_status CHECK (
        status_tindak_lanjut IN ('PENDING', 'DIABAIKAN', 'DIKOREKSI')
    )
);

CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_unmatched_header ON rekonsiliasi_unmatched(rekonsiliasi_id);
CREATE INDEX IF NOT EXISTS idx_rekonsiliasi_unmatched_tipe ON rekonsiliasi_unmatched(rekonsiliasi_id, tipe);
