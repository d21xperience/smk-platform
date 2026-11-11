CREATE TABLE IF NOT EXISTS "nis_sequences" (
    "academic_period_id" VARCHAR(10) NOT NULL,
    "prefix" VARCHAR(10) NOT NULL,
    "last_sequence" INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY ("academic_period_id")
);
