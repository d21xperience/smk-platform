CREATE TABLE IF NOT EXISTS master_catalogs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES master_categories(id) ON DELETE RESTRICT,
    name VARCHAR(255) NOT NULL,
    brand VARCHAR(100)
);
