-- +migrate Up
CREATE TABLE IF NOT EXISTS payments (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    invoice_id VARCHAR(36) NOT NULL REFERENCES invoices(id) ON DELETE RESTRICT,
    amount DECIMAL(12,2) NOT NULL,
    payment_method VARCHAR(20) NOT NULL, -- CASH, TRANSFER, VA, QRIS
    payment_date TIMESTAMP WITH TIME ZONE NOT NULL,
    reference_number VARCHAR(100),
    notes TEXT,
    created_by VARCHAR(36),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_payments_school_period ON payments(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_payments_invoice ON payments(invoice_id);
CREATE INDEX IF NOT EXISTS idx_payments_method ON payments(payment_method);
CREATE INDEX IF NOT EXISTS idx_payments_date ON payments(payment_date);
