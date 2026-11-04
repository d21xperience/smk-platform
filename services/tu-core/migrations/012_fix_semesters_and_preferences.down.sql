DROP TABLE IF EXISTS "user_period_preference";

CREATE TABLE "user_period_preference" (
    "user_id" VARCHAR(36) NOT NULL REFERENCES "user_login"("user_id") ON DELETE CASCADE,
    "period_id" VARCHAR(10) NOT NULL,
    PRIMARY KEY ("user_id", "period_id")
);
