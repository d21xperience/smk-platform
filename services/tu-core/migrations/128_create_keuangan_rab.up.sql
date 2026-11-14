-- =====================================================================
-- Migration 128: RAB (Rencana Anggaran Belanja) Header + Detail
-- + FK constraint pengajuan_pengeluaran.pos_anggaran_id → rab_detail.id
-- =====================================================================

-- ---------------------------------------------------------------------
-- 1. rab_header — 1 RAB per school × period × tahun_anggaran
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rab_header" (
    "id"                  UUID NOT NULL PRIMARY KEY,
    "school_id"           VARCHAR(36) NOT NULL,
    "academic_period_id"  VARCHAR(10) NOT NULL,
    "tahun_anggaran"      INTEGER NOT NULL,
    "total_pendapatan"    NUMERIC(15,2) NOT NULL DEFAULT 0,
    "total_belanja"       NUMERIC(15,2) NOT NULL DEFAULT 0,
    "surplus_defisit"     NUMERIC(15,2) NOT NULL DEFAULT 0,
    "status"              VARCHAR(20) NOT NULL DEFAULT 'DRAFT',
    "disusun_oleh"        VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "tanggal_penyusunan"  DATE NOT NULL DEFAULT CURRENT_DATE,
    "disahkan_oleh"       VARCHAR(36) NULL REFERENCES user_login(user_id) ON DELETE SET NULL,
    "tanggal_pengesahan"  DATE NULL,
    "catatan"             TEXT NULL,
    "created_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rab_header_status CHECK (
        status IN ('DRAFT', 'DIAJUKAN', 'DISAHKAN', 'DIREVISI')
    ),
    CONSTRAINT chk_rab_header_tahun CHECK (tahun_anggaran BETWEEN 2000 AND 2100),
    CONSTRAINT uq_rab_header_periode UNIQUE (school_id, academic_period_id, tahun_anggaran)
);

CREATE INDEX IF NOT EXISTS idx_rab_header_school ON rab_header(school_id);
CREATE INDEX IF NOT EXISTS idx_rab_header_tahun ON rab_header(school_id, tahun_anggaran);
CREATE INDEX IF NOT EXISTS idx_rab_header_status ON rab_header(school_id, status);

-- ---------------------------------------------------------------------
-- 2. rab_detail — item RAB per kode_rekening
-- ---------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS "rab_detail" (
    "id"                UUID NOT NULL PRIMARY KEY,
    "rab_header_id"     UUID NOT NULL REFERENCES rab_header(id) ON DELETE CASCADE,
    "kode_rekening_id"  UUID NOT NULL REFERENCES kode_rekening(id) ON DELETE RESTRICT,
    "uraian"            TEXT NOT NULL,
    "volume"            NUMERIC(15,2) NOT NULL DEFAULT 1,
    "satuan"            VARCHAR(50) NOT NULL DEFAULT 'unit',
    "harga_satuan"      NUMERIC(15,2) NOT NULL,
    "subtotal"          NUMERIC(15,2) NOT NULL,
    "sumber_dana"       VARCHAR(20) NULL,
    "created_at"        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_rab_detail_volume CHECK (volume > 0),
    CONSTRAINT chk_rab_detail_harga CHECK (harga_satuan > 0),
    CONSTRAINT chk_rab_detail_subtotal CHECK (subtotal >= 0),
    CONSTRAINT chk_rab_detail_sumber_dana CHECK (
        sumber_dana IS NULL OR sumber_dana IN ('BOS', 'BOP', 'HIBAH', 'MANDIRI')
    )
);

CREATE INDEX IF NOT EXISTS idx_rab_detail_header ON rab_detail(rab_header_id);
CREATE INDEX IF NOT EXISTS idx_rab_detail_kode ON rab_detail(kode_rekening_id);

-- ---------------------------------------------------------------------
-- 3. FK: pengajuan_pengeluaran.pos_anggaran_id → rab_detail.id
-- (nullable; hanya pengajuan yang punya budget link yang akan dicek)
-- ---------------------------------------------------------------------
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'pengajuan_pos_anggaran_fk'
          AND table_name = 'pengajuan_pengeluaran'
    ) THEN
        ALTER TABLE pengajuan_pengeluaran
        ADD CONSTRAINT pengajuan_pos_anggaran_fk
        FOREIGN KEY (pos_anggaran_id) REFERENCES rab_detail(id) ON DELETE SET NULL;
    END IF;
END $$;
