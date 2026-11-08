CREATE TABLE IF NOT EXISTS "tabel_berkas_siswa" (
    "id" UUID NOT NULL PRIMARY KEY,
    "peserta_didik_id" UUID NOT NULL,
    "jenis_berkas" VARCHAR(50) NOT NULL,
    "nama_file" VARCHAR(255) NOT NULL,
    "url_file" TEXT NOT NULL,
    "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT "tabel_berkas_siswa_peserta_didik_id_fkey" FOREIGN KEY ("peserta_didik_id") REFERENCES "tabel_siswa" ("peserta_didik_id") ON UPDATE NO ACTION ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_berkas_siswa_peserta ON "tabel_berkas_siswa"("peserta_didik_id");
