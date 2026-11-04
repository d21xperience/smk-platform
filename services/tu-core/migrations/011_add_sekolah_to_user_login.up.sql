-- 011_add_sekolah_to_user_login.up.sql
-- Menambahkan tabel_sekolah (DAPODIK) dan relasi single-tenant user_login.sekolah_id.

-- 1. Tabel sekolah (skema DAPODIK). IF NOT EXISTS karena mungkin sudah dibuat migrasi lain.
CREATE TABLE IF NOT EXISTS "tabel_sekolah" (
    "sekolah_id" UUID NOT NULL,
    "nama" VARCHAR(100) NOT NULL,
    "npsn" VARCHAR(8) NULL DEFAULT NULL,
    "nss" VARCHAR(12) NULL DEFAULT NULL,
    "alamat" TEXT NULL DEFAULT NULL,
    "kd_pos" VARCHAR(6) NULL DEFAULT NULL,
    "telepon" VARCHAR(20) NULL DEFAULT NULL,
    "fax" VARCHAR(20) NULL DEFAULT NULL,
    "kelurahan" VARCHAR(60) NULL DEFAULT NULL,
    "kecamatan" VARCHAR(60) NULL DEFAULT NULL,
    "kab_kota" VARCHAR(60) NULL DEFAULT NULL,
    "propinsi" VARCHAR(60) NULL DEFAULT NULL,
    "website" VARCHAR(100) NULL DEFAULT NULL,
    "email" VARCHAR(50) NULL DEFAULT NULL,
    "nm_kepsek" VARCHAR(100) NULL DEFAULT NULL,
    "nip_kepsek" VARCHAR(25) NULL DEFAULT NULL,
    "niy_kepsek" VARCHAR(30) NULL DEFAULT NULL,
    "status_kepemilikan_id" NUMERIC(1,0) NOT NULL,
    "kode_aktivasi" VARCHAR(30) NULL DEFAULT NULL,
    "jenjang" VARCHAR(20) NULL DEFAULT NULL,
    "bentuk_pendidikan_id" SMALLINT NULL DEFAULT NULL,
    PRIMARY KEY ("sekolah_id")
);

-- 1b. Jika tabel_sekolah sudah ada sebelumnya (dari migrasi lain) tanpa PK, tambahkan PK.
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE table_name = 'tabel_sekolah' AND constraint_type = 'PRIMARY KEY'
    ) THEN
        ALTER TABLE "tabel_sekolah" ADD CONSTRAINT "tabel_sekolah_pkey" PRIMARY KEY ("sekolah_id");
    END IF;
END $$;

-- 2. Seed sekolah default (ID harus sama dengan yang dipakai frontend saat dev).
INSERT INTO "tabel_sekolah" (
    "sekolah_id", "nama", "npsn", "status_kepemilikan_id", "jenjang", "bentuk_pendidikan_id"
) VALUES (
    '11111111-1111-1111-1111-111111111111',
    'SMK Negeri 1 Contoh',
    '12345678',
    1,
    'SMK',
    24
)
ON CONFLICT ("sekolah_id") DO NOTHING;

-- 3. Tambah kolom sekolah_id ke user_login (single-tenant: 1 user = 1 sekolah).
ALTER TABLE "user_login" ADD COLUMN IF NOT EXISTS "sekolah_id" UUID NULL;

-- 4. Backfill user lama ke sekolah default.
UPDATE "user_login"
SET "sekolah_id" = '11111111-1111-1111-1111-111111111111'
WHERE "sekolah_id" IS NULL;

-- 5. Foreign key constraint (hanya jika belum ada).
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'user_login_sekolah_id_fkey'
          AND table_name = 'user_login'
    ) THEN
        ALTER TABLE "user_login" ADD CONSTRAINT "user_login_sekolah_id_fkey"
            FOREIGN KEY ("sekolah_id") REFERENCES "tabel_sekolah"("sekolah_id");
    END IF;
END $$;

-- 6. Seed user admin untuk testing login (username: admin, password: admin123).
--    NOTE: password_hash disimpan plaintext. auth_service.go punya fallback
--    yang membandingkan plaintext jika bcrypt gagal. Untuk PRODUKSI,
--    ganti nilai ini dengan hash bcrypt asli (mis. via `htpasswd -bnBC 10 "" admin123`).
INSERT INTO "user_login" (
    "user_id", "username", "password_hash", "role", "is_active", "sekolah_id"
) VALUES (
    '22222222-2222-2222-2222-222222222222',
    'admin',
    'admin123',
    'admin',
    TRUE,
    '11111111-1111-1111-1111-111111111111'
)
ON CONFLICT ("username") DO NOTHING;
