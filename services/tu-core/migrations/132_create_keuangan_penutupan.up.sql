-- =====================================================================
-- Migration 132: Penutupan Buku Kas
-- 1 tabel + lock mechanism pada buku_kas_umum
--
-- Constraint unique: 1 penutupan per school×period×bulan×tahun
-- Kolom e-sign (pin_hash, ttd_kepsek_*) nullable — akan diisi FASE 5
-- =====================================================================

CREATE TABLE IF NOT EXISTS "penutupan_buku_kas" (
    "id"                  UUID NOT NULL PRIMARY KEY,
    "school_id"           VARCHAR(36) NOT NULL,
    "academic_period_id"  VARCHAR(10) NOT NULL,
    "periode_bulan"       INTEGER NOT NULL,
    "periode_tahun"       INTEGER NOT NULL,

    -- Snapshot buku kas
    "saldo_awal"          NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_debet"         NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_kredit"        NUMERIC(15,2) NOT NULL DEFAULT 0,
    "saldo_akhir_buku"    NUMERIC(15,2) NOT NULL DEFAULT 0,

    -- Cash count (wajib input)
    "saldo_fisik_kas"     NUMERIC(15,2) NOT NULL,
    "selisih_kas"         NUMERIC(15,2) NOT NULL DEFAULT 0,

    -- Snapshot posisi keuangan
    "total_aset"          NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_kewajiban"     NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_ekuitas"       NUMERIC(15,2) NOT NULL DEFAULT 0,

    -- Status & audit
    "status"              VARCHAR(40) NOT NULL DEFAULT 'CLOSED',
    "catatan"             TEXT NULL,
    "closed_by"           VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "closed_at"           TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- E-sign (nullable, diisi FASE 5)
    "pin_hash"            VARCHAR(255) NULL,
    "ttd_kepsek_hash"     VARCHAR(255) NULL,
    "ttd_kepsek_by"       VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "ttd_kepsek_at"       TIMESTAMPTZ NULL,

    -- Dokumen PDF (nullable, diisi FASE 5)
    "bku_pdf_url"         TEXT NULL,
    "neraca_pdf_url"      TEXT NULL,
    "lra_pdf_url"         TEXT NULL,

    "created_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT chk_penutupan_status CHECK (status IN (
        'CLOSED',
        'MENUNGGU_TTD_KEPALA',
        'CLOSED_FULLY_SIGNED',
        'DIBUKA_KEMBALI'
    )),
    CONSTRAINT chk_penutupan_bulan CHECK (periode_bulan BETWEEN 1 AND 12),
    CONSTRAINT chk_penutupan_tahun CHECK (periode_tahun BETWEEN 2000 AND 2100),
    CONSTRAINT uq_penutupan_unique UNIQUE (
        school_id, academic_period_id, periode_bulan, periode_tahun
    )
);

CREATE INDEX IF NOT EXISTS idx_penutupan_school ON penutupan_buku_kas(school_id);
CREATE INDEX IF NOT EXISTS idx_penutupan_period
    ON penutupan_buku_kas(school_id, academic_period_id, periode_tahun, periode_bulan);
CREATE INDEX IF NOT EXISTS idx_penutupan_status ON penutupan_buku_kas(school_id, status);
