CREATE TABLE IF NOT EXISTS letter_types (
    id                  UUID PRIMARY KEY,
    kategori            VARCHAR(100) NOT NULL,
    peruntukan          VARCHAR(20) NOT NULL CHECK (peruntukan IN ('SISWA', 'GURU', 'UMUM')),
    value               VARCHAR(50) UNIQUE NOT NULL,
    label               VARCHAR(150) NOT NULL,
    kode_klasifikasi    VARCHAR(50) NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
