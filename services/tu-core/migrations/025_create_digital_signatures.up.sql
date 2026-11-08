CREATE TABLE IF NOT EXISTS digital_signatures (
    id                  UUID PRIMARY KEY,
    correspondence_id   UUID NOT NULL REFERENCES correspondences(id) ON DELETE CASCADE,
    signer_user_id      UUID NOT NULL,
    signer_role         VARCHAR(50),
    image_url           TEXT NOT NULL,
    verification_code   VARCHAR(32) UNIQUE NOT NULL,
    signed_at           TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sig_correspondence ON digital_signatures(correspondence_id);
CREATE INDEX idx_sig_verification ON digital_signatures(verification_code);
