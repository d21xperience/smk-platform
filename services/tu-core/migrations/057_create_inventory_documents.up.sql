CREATE TABLE IF NOT EXISTS sarpras_inventory_documents (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_type VARCHAR(50) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    document_number VARCHAR(100),
    document_date DATE,
    bucket_name VARCHAR(100) NOT NULL,
    object_key TEXT NOT NULL,
    original_filename VARCHAR(255) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    file_size BIGINT NOT NULL,
    checksum VARCHAR(64),
    visibility VARCHAR(20) NOT NULL DEFAULT 'INTERNAL' CHECK (visibility IN ('PUBLIC', 'INTERNAL', 'RESTRICTED')),
    status VARCHAR(20) NOT NULL DEFAULT 'UPLOADED' CHECK (status IN ('UPLOADED', 'VERIFIED', 'REJECTED', 'ARCHIVED')),
    rejection_reason TEXT,
    uploaded_by VARCHAR(100) NOT NULL,
    verified_by VARCHAR(100),
    verified_at TIMESTAMPTZ,
    archived_by VARCHAR(100),
    archived_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sarpras_inventory_document_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    document_id UUID NOT NULL REFERENCES sarpras_inventory_documents(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL CHECK (entity_type IN ('INVENTORY_ITEM', 'GOODS_RECEIPT', 'MASTER_CATALOG', 'MAINTENANCE_LOG', 'DISPOSAL_REQUEST')),
    entity_id VARCHAR(100) NOT NULL,
    relation_type VARCHAR(50) DEFAULT 'PRIMARY',
    created_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_doc_links_entity ON sarpras_inventory_document_links(entity_type, entity_id);
