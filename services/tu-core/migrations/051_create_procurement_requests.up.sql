CREATE TABLE IF NOT EXISTS procurement_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    request_number VARCHAR(50) UNIQUE NOT NULL,
    requester_id VARCHAR(100) NOT NULL, -- Menyimpan ID user/guru dari context
    requester_name VARCHAR(150) NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    specification TEXT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    estimated_price NUMERIC(15, 2) DEFAULT 0,
    total_estimated NUMERIC(15, 2) DEFAULT 0,
    priority VARCHAR(20) NOT NULL CHECK (priority IN ('LOW', 'NORMAL', 'HIGH', 'URGENT')),
    justification TEXT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'PENDING' CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED', 'PROCURED')),
    approved_by VARCHAR(100),
    approved_at TIMESTAMPTZ,
    rejection_reason TEXT,
    notes TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_procurement_status ON procurement_requests(status);
CREATE INDEX IF NOT EXISTS idx_procurement_requester ON procurement_requests(requester_id);
CREATE INDEX IF NOT EXISTS idx_procurement_created_at ON procurement_requests(created_at);

-- Trigger auto update updated_at
CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_procurement_modtime
BEFORE UPDATE ON procurement_requests
FOR EACH ROW EXECUTE PROCEDURE update_modified_column();
