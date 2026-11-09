CREATE TABLE IF NOT EXISTS goods_receipts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bast_number VARCHAR(100) UNIQUE NOT NULL,
    receipt_date DATE NOT NULL,
    source_type VARCHAR(50) NOT NULL CHECK (source_type IN ('PROCUREMENT', 'HIBAH', 'DONASI', 'MUTASI', 'OTHER')),
    supplier_name VARCHAR(255) NOT NULL,
    invoice_number VARCHAR(100),
    funding_source VARCHAR(100),
    procurement_request_id UUID REFERENCES procurement_requests(id) ON DELETE SET NULL,
    received_by VARCHAR(100) NOT NULL,
    verified_by VARCHAR(100),
    status VARCHAR(20) NOT NULL DEFAULT 'DRAFT' CHECK (status IN ('DRAFT', 'VERIFIED', 'REJECTED', 'DISTRIBUTED')),
    semester_id VARCHAR(5) NOT NULL,
    rejection_reason TEXT,
    document_urls JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS goods_receipt_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    goods_receipt_id UUID NOT NULL REFERENCES goods_receipts(id) ON DELETE CASCADE,
    catalog_id UUID REFERENCES master_catalogs(id),
    item_name VARCHAR(255) NOT NULL,
    specification TEXT,
    quantity INT NOT NULL CHECK (quantity > 0),
    unit VARCHAR(50) NOT NULL,
    unit_price NUMERIC(15, 2) DEFAULT 0,
    total_price NUMERIC(15, 2) DEFAULT 0,
    condition VARCHAR(50) NOT NULL CHECK (condition IN ('BAIK', 'RUSAK_RINGAN', 'RUSAK_BERAT')),
    is_inventory BOOLEAN DEFAULT FALSE
);

-- Tambah kolom tracking di inventory_items
ALTER TABLE inventory_items ADD COLUMN IF NOT EXISTS goods_receipt_id UUID REFERENCES goods_receipts(id) ON DELETE SET NULL;
