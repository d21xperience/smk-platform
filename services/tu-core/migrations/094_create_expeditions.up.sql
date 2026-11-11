CREATE TABLE IF NOT EXISTS expeditions (
    id                  UUID PRIMARY KEY,
    school_id           UUID NOT NULL,
    correspondence_id   UUID NOT NULL REFERENCES correspondences(id) ON DELETE CASCADE,
    penerima_nama       VARCHAR(255) NOT NULL,
    penerima_instansi   VARCHAR(255),
    penerima_alamat     TEXT,
    metode_kirim        VARCHAR(50) NOT NULL,
    nomor_resi          VARCHAR(100),
    tanggal_kirim       DATE NOT NULL,
    jumlah_lembar       INT NOT NULL DEFAULT 1,
    status              VARCHAR(20) NOT NULL DEFAULT 'TERKIRIM',
    bukti_pdf_url       TEXT,
    verify_token        VARCHAR(50) UNIQUE NOT NULL,
    created_by          UUID NOT NULL,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_exp_corr ON expeditions(correspondence_id);
CREATE INDEX IF NOT EXISTS idx_exp_token ON expeditions(verify_token);
