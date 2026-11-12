-- Tabel pivot untuk menghubungkan Sekolah dengan Master Jurusan
CREATE TABLE IF NOT EXISTS "tabel_sekolah_jurusan" (
    "sekolah_id" UUID NOT NULL REFERENCES "tabel_sekolah"("sekolah_id") ON DELETE CASCADE,
    "jurusan_id" VARCHAR(25) NOT NULL REFERENCES "tabel_jurusan"("jurusan_id") ON DELETE CASCADE,
    PRIMARY KEY ("sekolah_id", "jurusan_id")
);

-- Seed data awal (opsional: hubungkan sekolah default dengan TKJ)
INSERT INTO "tabel_sekolah_jurusan" ("sekolah_id", "jurusan_id") VALUES
('11111111-1111-1111-1111-111111111111', 'TKJ')
ON CONFLICT DO NOTHING;
