CREATE TABLE IF NOT EXISTS "tabel_ptk_pelengkap" (
    "ptk_pelengkap_id" UUID NOT NULL PRIMARY KEY,
    "ptk_id" UUID NOT NULL REFERENCES "tabel_ptk"("ptk_id") ON DELETE CASCADE,
    "gelar_depan" VARCHAR(20) NULL DEFAULT NULL,
    "gelar_belakang" VARCHAR(20) NULL DEFAULT NULL,
    "nip_niy" VARCHAR(18) NULL DEFAULT NULL
);
