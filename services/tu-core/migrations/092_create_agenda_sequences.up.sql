CREATE TABLE IF NOT EXISTS agenda_sequences (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    school_id           UUID NOT NULL,
    jenis               VARCHAR(10) NOT NULL CHECK (jenis IN ('IN', 'OUT')),
    tahun               INT NOT NULL,
    last_number         INT NOT NULL DEFAULT 0,
    UNIQUE (school_id, jenis, tahun)
);
