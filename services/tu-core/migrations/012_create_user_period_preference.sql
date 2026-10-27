CREATE TABLE IF NOT EXISTS "user_period_preference" (
    "user_id" VARCHAR(50) NOT NULL,
    "school_id" VARCHAR(50) NOT NULL,
    "semester_id" VARCHAR(5) NOT NULL,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("user_id", "school_id")
);
