CREATE TABLE IF NOT EXISTS ptk_performances (
    id UUID PRIMARY KEY,
    ptk_id UUID NOT NULL REFERENCES tabel_ptk(ptk_id) ON DELETE CASCADE,
    tahun VARCHAR(4) NOT NULL,
    semester VARCHAR(10) NOT NULL,
    nilai_skp NUMERIC(5,2) NOT NULL,
    predikat VARCHAR(20) NOT NULL,
    catatan TEXT,
    evaluator VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE INDEX idx_ptk_performances_ptk_id ON ptk_performances(ptk_id);
