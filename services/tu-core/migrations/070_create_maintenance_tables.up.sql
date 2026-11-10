-- Tabel Tiket Kerusakan (Damage Reports)
CREATE TABLE IF NOT EXISTS damage_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    report_number VARCHAR(100) UNIQUE NOT NULL,
    inventory_item_id VARCHAR(100) REFERENCES inventory_items(id) ON DELETE CASCADE,
    asset_name VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    reporter_id VARCHAR(100) NOT NULL,
    reporter_name VARCHAR(150) NOT NULL,
    description TEXT NOT NULL,
    urgency VARCHAR(20) NOT NULL CHECK (urgency IN ('Rendah', 'Sedang', 'Tinggi')),
    status VARCHAR(20) NOT NULL DEFAULT 'Pending' CHECK (status IN ('Pending', 'Proses', 'Selesai')),
    assigned_to VARCHAR(100),
    vendor_name VARCHAR(150),
    cost NUMERIC(15,2) DEFAULT 0,
    resolution_notes TEXT,
    semester_id VARCHAR(5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel Jadwal Perawatan Berkala
CREATE TABLE IF NOT EXISTS maintenance_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_code VARCHAR(50) UNIQUE NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    inventory_item_id VARCHAR(100) REFERENCES inventory_items(id) ON DELETE SET NULL,
    scope_location VARCHAR(255),
    frequency VARCHAR(50) NOT NULL, -- Bulanan, 3 Bulanan, 6 Bulanan, Tahunan
    last_service_date DATE,
    next_service_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    semester_id VARCHAR(5),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Tabel Log Perawatan (Setiap kali servis selesai)
CREATE TABLE IF NOT EXISTS maintenance_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id UUID REFERENCES maintenance_schedules(id) ON DELETE SET NULL,
    damage_report_id UUID REFERENCES damage_reports(id) ON DELETE SET NULL,
    inventory_item_id VARCHAR(100) REFERENCES inventory_items(id) ON DELETE SET NULL,
    service_date DATE NOT NULL,
    technician VARCHAR(150),
    vendor VARCHAR(150),
    cost NUMERIC(15,2) DEFAULT 0,
    notes TEXT,
    performed_by VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_damage_status ON damage_reports(status);
CREATE INDEX IF NOT EXISTS idx_maint_next ON maintenance_schedules(next_service_date);
