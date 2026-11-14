-- =====================================================================
-- Migration 134: Index untuk aging piutang
--
-- Piutang = tagihan_siswa dengan status BELUM_LUNAS/SEBAGIAN.
-- Partial index berikut membuat query aging report cepat karena hanya
-- scan row yang belum lunas (biasanya <20% dari total tagihan).
-- =====================================================================

-- Partial index untuk filter status piutang + sort tanggal
CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_piutang
    ON tagihan_siswa (school_id, academic_period_id, tanggal_jatuh_tempo)
    WHERE status IN ('BELUM_LUNAS', 'SEBAGIAN');

-- Partial index untuk lookup per NIS (untuk endpoint detail per siswa)
CREATE INDEX IF NOT EXISTS idx_tagihan_siswa_piutang_nis
    ON tagihan_siswa (school_id, academic_period_id, nis)
    WHERE status IN ('BELUM_LUNAS', 'SEBAGIAN');
