-- +migrate Up
CREATE TABLE IF NOT EXISTS invoices (
    id VARCHAR(36) PRIMARY KEY,
    school_id VARCHAR(36) NOT NULL,
    academic_period_id VARCHAR(20) NOT NULL,
    student_id VARCHAR(36) NOT NULL,
    component_type VARCHAR(20) NOT NULL, -- SPP, UNIFORM, BOOK, EXAM
    month INTEGER, -- NULL untuk UNIFORM, BOOK, EXAM
    year INTEGER, -- NULL untuk UNIFORM, BOOK, EXAM
    amount DECIMAL(12,2) NOT NULL,
    due_date DATE NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'UNPAID', -- UNPAID, PARTIAL, PAID, CANCELLED
    paid_amount DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    paid_at TIMESTAMP WITH TIME ZONE,
    notes TEXT,
    created_by VARCHAR(36),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),

    CONSTRAINT uq_invoices_student_component_period UNIQUE (school_id, academic_period_id, student_id, component_type, month, year)
);

CREATE INDEX IF NOT EXISTS idx_invoices_school_period ON invoices(school_id, academic_period_id);
CREATE INDEX IF NOT EXISTS idx_invoices_student ON invoices(student_id);
CREATE INDEX IF NOT EXISTS idx_invoices_component ON invoices(component_type);
CREATE INDEX IF NOT EXISTS idx_invoices_status ON invoices(status);
CREATE INDEX IF NOT EXISTS idx_invoices_due_date ON invoices(due_date);
CREATE INDEX IF NOT EXISTS idx_invoices_month_year ON invoices(month, year);
