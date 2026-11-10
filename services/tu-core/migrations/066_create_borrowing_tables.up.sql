CREATE TABLE IF NOT EXISTS borrowings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrowing_number VARCHAR(100) UNIQUE NOT NULL,
    borrower_id VARCHAR(100) NOT NULL,      -- User ID / PTK ID peminjam
    borrower_name VARCHAR(150) NOT NULL,    -- Nama peminjam (snapshot)
    inventory_item_id VARCHAR(100) NOT NULL REFERENCES inventory_items(id) ON DELETE RESTRICT,
    purpose_type VARCHAR(50) NOT NULL,      -- PEMBELAJARAN, RAPAT, KEGIATAN, LOMBA, DINAS, LAINNYA
    purpose_description TEXT,
    borrow_date DATE NOT NULL,
    expected_return_date DATE NOT NULL,
    usage_location VARCHAR(255),
    duration_days INT NOT NULL,
    approval_level INT NOT NULL DEFAULT 1,
    status VARCHAR(50) NOT NULL DEFAULT 'PENDING_APPROVAL', -- PENDING_APPROVAL, APPROVED, REJECTED, PICKED_UP, RETURNED
    current_approval_step INT DEFAULT 0,
    rejection_reason TEXT,
    
    -- Pickup Data
    actual_pickup_date DATE,
    pickup_by VARCHAR(100),
    pickup_condition VARCHAR(50),
    pickup_photos JSONB,
    pickup_notes TEXT,
    bast_number VARCHAR(100),
    bast_pdf_url TEXT,
    
    -- Return Data
    actual_return_date DATE,
    return_by VARCHAR(100),
    return_condition VARCHAR(50),
    return_photos JSONB,
    return_notes TEXT,
    days_late INT DEFAULT 0,
    fine_amount NUMERIC(15,2) DEFAULT 0,
    fine_paid BOOLEAN DEFAULT TRUE,
    bap_number VARCHAR(100),
    bap_pdf_url TEXT,
    verification_hash VARCHAR(255),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS borrowing_approvals (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrowing_id UUID NOT NULL REFERENCES borrowings(id) ON DELETE CASCADE,
    approval_step INT NOT NULL,
    approver_id VARCHAR(100) NOT NULL,
    decision VARCHAR(50) NOT NULL, -- APPROVED, REJECTED
    notes TEXT,
    approval_date TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS borrowing_fines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    borrowing_id UUID NOT NULL REFERENCES borrowings(id) ON DELETE CASCADE,
    borrower_id VARCHAR(100) NOT NULL,
    fine_amount NUMERIC(15,2) NOT NULL,
    days_late INT NOT NULL,
    fine_per_day NUMERIC(15,2) NOT NULL,
    fine_date DATE NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PAID
    paid_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_borrowing_status ON borrowings(status);
CREATE INDEX IF NOT EXISTS idx_borrowing_item ON borrowings(inventory_item_id);
CREATE INDEX IF NOT EXISTS idx_borrowing_borrower ON borrowings(borrower_id);
