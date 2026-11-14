-- =====================================================================
-- Migration 136: Tabel user_pin_ttd untuk PIN Tanda Tangan Digital
-- =====================================================================
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS "user_pin_ttd" (
    "id"               UUID NOT NULL PRIMARY KEY DEFAULT gen_random_uuid(),
    "user_id"          VARCHAR(36) NOT NULL UNIQUE REFERENCES user_login(user_id) ON DELETE CASCADE,
    "pin_hash"         VARCHAR(255) NOT NULL,
    "max_attempt"      INTEGER NOT NULL DEFAULT 3,
    "current_attempt"  INTEGER NOT NULL DEFAULT 0,
    "last_attempt_at"  TIMESTAMPTZ NULL,
    "is_locked"        BOOLEAN NOT NULL DEFAULT FALSE,
    "locked_until"     TIMESTAMPTZ NULL,
    "created_at"       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    "updated_at"       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT chk_pin_max_attempt CHECK (max_attempt BETWEEN 1 AND 10),
    CONSTRAINT chk_pin_current_attempt CHECK (current_attempt >= 0)
);

CREATE INDEX IF NOT EXISTS idx_user_pin_ttd_user ON user_pin_ttd(user_id);
CREATE INDEX IF NOT EXISTS idx_user_pin_ttd_locked ON user_pin_ttd(is_locked) WHERE is_locked = TRUE;
