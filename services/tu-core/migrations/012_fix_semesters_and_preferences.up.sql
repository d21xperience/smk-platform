-- 012_fix_semesters_and_preferences.up.sql
-- 1. Hapus tabel preferensi lama yang tidak punya school_id
DROP TABLE IF EXISTS "user_period_preference";

-- 2. Buat ulang tabel preferensi dengan school_id dan referensi ke semesters(id)
CREATE TABLE IF NOT EXISTS "user_period_preference" (
    "user_id" VARCHAR(36) NOT NULL REFERENCES "user_login"("user_id") ON DELETE CASCADE,
    "school_id" UUID NOT NULL REFERENCES "tabel_sekolah"("sekolah_id") ON DELETE CASCADE,
    "semester_id" VARCHAR(10) NOT NULL REFERENCES "semesters"("id") ON DELETE CASCADE,
    "updated_at" TIMESTAMPTZ DEFAULT NOW(),
    PRIMARY KEY ("user_id", "school_id")
);

-- 3. Seed data ke tabel semesters (jika belum ada)
INSERT INTO "semesters" ("id", "name", "academic_year", "semester", "is_active") VALUES
('20251', 'Semester 1', '2025/2026', 1, FALSE),
('20252', 'Semester 2', '2025/2026', 2, FALSE),
('20261', 'Semester 1', '2026/2027', 1, TRUE),
('20262', 'Semester 2', '2026/2027', 2, FALSE)
ON CONFLICT ("id") DO NOTHING;

-- 4. Set default preference admin ke semester aktif
INSERT INTO "user_period_preference" ("user_id", "school_id", "semester_id") VALUES
('22222222-2222-2222-2222-222222222222', '11111111-1111-1111-1111-111111111111', '20261')
ON CONFLICT ("user_id", "school_id") DO UPDATE SET "semester_id" = EXCLUDED."semester_id";
