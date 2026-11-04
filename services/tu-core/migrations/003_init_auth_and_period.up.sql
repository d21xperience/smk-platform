-- Menggunakan skema dasar untuk autentikasi dan periode akademik
CREATE TABLE IF NOT EXISTS "semesters" (
    "id" VARCHAR(10) PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "academic_year" VARCHAR(10) NOT NULL,
    "semester" INTEGER NOT NULL,
    "is_active" BOOLEAN DEFAULT FALSE
);

CREATE TABLE IF NOT EXISTS "user_login" (
    "user_id" VARCHAR(36) PRIMARY KEY,
    "username" VARCHAR(50) UNIQUE NOT NULL,
    "password_hash" VARCHAR(255) NOT NULL,
    "role" VARCHAR(20) NOT NULL,
    "is_active" BOOLEAN DEFAULT TRUE,
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS "user_period_preference" (
    "user_id" VARCHAR(36) NOT NULL REFERENCES "user_login"("user_id") ON DELETE CASCADE,
    "period_id" VARCHAR(10) NOT NULL REFERENCES "semesters"("id") ON DELETE CASCADE,
    PRIMARY KEY ("user_id", "period_id")
);
