CREATE TABLE IF NOT EXISTS retention_schedules (
    id                  UUID PRIMARY KEY,
    kategori            VARCHAR(100) NOT NULL UNIQUE,
    retensi_tahun       INT NOT NULL,
    keterangan          TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Seed data standar ANRI
INSERT INTO retention_schedules (id, kategori, retensi_tahun, keterangan) VALUES
(gen_random_uuid(), 'Surat Keputusan (SK)', 5, 'Retensi 5 tahun'),
(gen_random_uuid(), 'Surat Tugas (ST)', 3, 'Retensi 3 tahun'),
(gen_random_uuid(), 'SPPD', 5, 'Retensi 5 tahun'),
(gen_random_uuid(), 'MoU / Kerjasama', 10, 'Retensi 10 tahun'),
(gen_random_uuid(), 'Surat Masuk Umum', 2, 'Retensi 2 tahun')
ON CONFLICT (kategori) DO NOTHING;
