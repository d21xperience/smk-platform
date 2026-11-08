CREATE TABLE IF NOT EXISTS dispositions (
    id                  UUID PRIMARY KEY,
    correspondence_id   UUID NOT NULL REFERENCES correspondences(id) ON DELETE CASCADE,
    school_id           UUID NOT NULL,
    from_user_id        UUID NOT NULL,
    to_user_id          UUID NOT NULL,
    note                TEXT,
    status              VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'READ', 'DONE')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_dispo_correspondence ON dispositions(correspondence_id);
CREATE INDEX idx_dispo_to_user ON dispositions(to_user_id);
