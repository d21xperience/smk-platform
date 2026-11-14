-- =====================================================================
-- Seed: PIN TTD default untuk admin (PIN: 123456)
--
-- Menggunakan PostgreSQL `crypt` dengan `gen_salt('bf', 10)` — menghasilkan
-- bcrypt hash format $2a$ yang kompatibel dengan Go bcrypt
-- (`bcrypt.CompareHashAndPassword` support $2a$/$2b$/$2y$).
-- =====================================================================

INSERT INTO "user_pin_ttd" (
    user_id, pin_hash, max_attempt, current_attempt, is_locked
)
SELECT
    u.user_id,
    crypt('123456', gen_salt('bf', 10)),
    3, 0, FALSE
FROM user_login u
WHERE u.user_id = '11111111-1111-1111-1111-111111111111'
ON CONFLICT (user_id) DO NOTHING;
