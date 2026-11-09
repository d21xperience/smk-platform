CREATE TABLE IF NOT EXISTS "user_period_preference" (
    "user_id" VARCHAR(36) NOT NULL,
    "school_id" VARCHAR(36) NOT NULL,
    "semester_id" VARCHAR(10) NOT NULL,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    PRIMARY KEY ("user_id", "school_id")
);
