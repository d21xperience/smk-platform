CREATE TABLE IF NOT EXISTS disposal_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_number VARCHAR(100) UNIQUE NOT NULL,
    proposed_by VARCHAR(100) NOT NULL,
    proposal_date DATE NOT NULL,
    reason_category VARCHAR(50) NOT NULL, -- RUSAK_BERAT, USIA, TIDAK_EKONOMIS, HILANG, USANG, BENCANA, LAINNYA
    proposed_method VARCHAR(50) NOT NULL, -- DIMUSNAHKAN, DIJUAL, DISUMBANGKAN, DITUKAR, DIMANFAATKAN
    total_assets INT NOT NULL,
    total_value NUMERIC(15,2) NOT NULL,
    approval_level INT NOT NULL, -- 1 (<10jt), 2 (10-50jt), 3 (>50jt)
    status VARCHAR(50) NOT NULL DEFAULT 'DRAFT', -- DRAFT, VERIFIED_KTU, REJECTED, APPROVED, EXECUTED
    current_approval_step INT DEFAULT 0,
    rejection_reason TEXT,
    execution_date DATE,
    execution_method VARCHAR(50),
    execution_location VARCHAR(255),
    executors JSONB,
    witnesses JSONB,
    execution_documentation JSONB,
    execution_notes TEXT,
    sale_proceeds NUMERIC(15,2) DEFAULT 0,
    sale_evidence_url TEXT,
    bast_number VARCHAR(100),
    bast_pdf_url TEXT,
    verification_hash VARCHAR(255),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS disposal_request_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disposal_request_id UUID NOT NULL REFERENCES disposal_requests(id) ON DELETE CASCADE,
    inventory_item_id VARCHAR(100) NOT NULL REFERENCES inventory_items(id) ON DELETE CASCADE,
    specific_reason TEXT NOT NULL,
    photo_urls JSONB,
    verification_status VARCHAR(50) DEFAULT 'PENDING',
    verified_condition VARCHAR(50),
    verification_notes TEXT
);

CREATE TABLE IF NOT EXISTS disposal_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disposal_request_id UUID NOT NULL REFERENCES disposal_requests(id) ON DELETE CASCADE,
    verification_date DATE NOT NULL,
    verification_team JSONB NOT NULL,
    verification_document_url TEXT,
    team_recommendation VARCHAR(50) NOT NULL, -- SETUJU, TOLAK, REVISI
    team_notes TEXT,
    verified_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS disposal_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    disposal_request_id UUID NOT NULL REFERENCES disposal_requests(id) ON DELETE CASCADE,
    approval_step INT NOT NULL,
    approver_id VARCHAR(100) NOT NULL,
    decision VARCHAR(50) NOT NULL, -- APPROVED, REJECTED
    approval_document_url TEXT,
    notes TEXT,
    approval_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_disposal_status ON disposal_requests(status);
CREATE INDEX IF NOT EXISTS idx_disposal_items_req ON disposal_request_items(disposal_request_id);
