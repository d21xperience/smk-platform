CREATE TABLE IF NOT EXISTS semester (
    semester_id      VARCHAR(5) PRIMARY KEY,
    academic_year_id UUID NOT NULL REFERENCES academic_years(id) ON DELETE RESTRICT,
    nama_semester    VARCHAR(50) NOT NULL,
    semester         INTEGER NOT NULL CHECK (semester IN (1, 2)),
    is_active        BOOLEAN NOT NULL DEFAULT false,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT semester_academic_year_semester_uniq UNIQUE (academic_year_id, semester)
);

CREATE INDEX IF NOT EXISTS idx_semester_active ON semester(is_active) WHERE is_active = true;

INSERT INTO semester (semester_id, academic_year_id, nama_semester, semester, is_active)
VALUES
  ('20251', 'abf81815-b02c-40f7-b638-79ee0f4d935f', 'Semester 1', 1, true),
  ('20252', 'abf81815-b02c-40f7-b638-79ee0f4d935f', 'Semester 2', 2, false)
ON CONFLICT (semester_id) DO NOTHING;
