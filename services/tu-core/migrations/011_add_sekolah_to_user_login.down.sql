-- 011_add_sekolah_to_user_login.down.sql
-- Rollback migrasi 011. Tabel tabel_sekolah SENGAJA tidak di-drop
-- (mungkin sudah ada sebelumnya dari migrasi lain).
-- Drop manual jika memang dibuat oleh migrasi ini.

DELETE FROM "user_login" WHERE "user_id" = '22222222-2222-2222-2222-222222222222';

ALTER TABLE "user_login" DROP CONSTRAINT IF EXISTS "user_login_sekolah_id_fkey";
ALTER TABLE "user_login" DROP COLUMN IF EXISTS "sekolah_id";

DELETE FROM "tabel_sekolah" WHERE "sekolah_id" = '11111111-1111-1111-1111-111111111111';
