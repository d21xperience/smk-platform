package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DamageReport struct {
	ID              string    `json:"id"`
	ReportNumber    string    `json:"report_number"`
	InventoryItemID string    `json:"inventory_item_id"`
	AssetName       string    `json:"asset_name"`
	Location        string    `json:"location"`
	ReporterID      string    `json:"reporter_id"`
	ReporterName    string    `json:"reporter_name"`
	Description     string    `json:"description"`
	Urgency         string    `json:"urgency"`
	Status          string    `json:"status"`
	AssignedTo      string    `json:"assigned_to"`
	VendorName      string    `json:"vendor_name"`
	Cost            float64   `json:"cost"`
	ResolutionNotes string    `json:"resolution_notes"`
	SemesterID      string    `json:"semester_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type MaintenanceSchedule struct {
	ID              string     `json:"id"`
	ScheduleCode    string     `json:"schedule_code"`
	AssetName       string     `json:"asset_name"`
	InventoryItemID string     `json:"inventory_item_id"`
	ScopeLocation   string     `json:"scope_location"`
	Frequency       string     `json:"frequency"`
	LastServiceDate *time.Time `json:"last_service_date"`
	NextServiceDate time.Time  `json:"next_service_date"`
	IsActive        bool       `json:"is_active"`
	SemesterID      string     `json:"semester_id"`
}

type MaintenanceLog struct {
	ID              string    `json:"id"`
	ScheduleID      string    `json:"schedule_id"`
	DamageReportID  string    `json:"damage_report_id"`
	InventoryItemID string    `json:"inventory_item_id"`
	ServiceDate     time.Time `json:"service_date"`
	Technician      string    `json:"technician"`
	Vendor          string    `json:"vendor"`
	Cost            float64   `json:"cost"`
	Notes           string    `json:"notes"`
	PerformedBy     string    `json:"performed_by"`
	CreatedAt       time.Time `json:"created_at"`
	
	// Joined Fields
	AssetCode string `json:"asset_code"`
	AssetName string `json:"asset_name"`
}

type MaintenanceRepository interface {
	GetAllDamageReports(ctx context.Context, semesterID string) ([]*DamageReport, error)
	CreateDamageReport(ctx context.Context, d *DamageReport) error
	UpdateDamageReport(ctx context.Context, id string, status, vendor string, cost float64, notes string) error
	GenerateDamageNumber(ctx context.Context) (string, error)

	GetAllSchedules(ctx context.Context, semesterID string) ([]*MaintenanceSchedule, error)
	CreateSchedule(ctx context.Context, s *MaintenanceSchedule) error
	CompleteSchedule(ctx context.Context, id string, serviceDate time.Time, cost float64, notes string, performedBy string) error
	GenerateScheduleCode(ctx context.Context) (string, error)

	GetMaintenanceLogs(ctx context.Context, startDate, endDate string, semesterID string) ([]*MaintenanceLog, error)
}

type pgxMaintenanceRepository struct {
	db *pgxpool.Pool
}

func NewMaintenanceRepository(db *pgxpool.Pool) MaintenanceRepository {
	return &pgxMaintenanceRepository{db: db}
}

func (r *pgxMaintenanceRepository) GenerateDamageNumber(ctx context.Context) (string, error) {
	var count int
	year := time.Now().Year()
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM damage_reports WHERE EXTRACT(YEAR FROM created_at) = $1`, year).Scan(&count)
	return fmt.Sprintf("DMG/%d/%03d", year, count+1), nil
}

func (r *pgxMaintenanceRepository) GenerateScheduleCode(ctx context.Context) (string, error) {
	var count int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM maintenance_schedules`).Scan(&count)
	return fmt.Sprintf("PM-%02d", count+1), nil
}

func (r *pgxMaintenanceRepository) GetAllDamageReports(ctx context.Context, semesterID string) ([]*DamageReport, error) {
	query := `
		SELECT id::text, report_number, COALESCE(inventory_item_id, ''), asset_name, COALESCE(location, ''), 
		       reporter_id, reporter_name, description, urgency, status, COALESCE(assigned_to, ''), 
		       COALESCE(vendor_name, ''), COALESCE(cost, 0), COALESCE(resolution_notes, ''), COALESCE(semester_id, ''), created_at, updated_at
		FROM damage_reports
		WHERE ($1 = '' OR semester_id = $1)
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*DamageReport
	for rows.Next() {
		var d DamageReport
		err := rows.Scan(&d.ID, &d.ReportNumber, &d.InventoryItemID, &d.AssetName, &d.Location,
			&d.ReporterID, &d.ReporterName, &d.Description, &d.Urgency, &d.Status, &d.AssignedTo,
			&d.VendorName, &d.Cost, &d.ResolutionNotes, &d.SemesterID, &d.CreatedAt, &d.UpdatedAt)
		if err != nil { return nil, err }
		list = append(list, &d)
	}
	return list, nil
}

func (r *pgxMaintenanceRepository) CreateDamageReport(ctx context.Context, d *DamageReport) error {
	query := `
		INSERT INTO damage_reports (report_number, inventory_item_id, asset_name, location, reporter_id, reporter_name, description, urgency, status, semester_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'Pending', $9)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, d.ReportNumber, d.InventoryItemID, d.AssetName, d.Location, d.ReporterID, d.ReporterName, d.Description, d.Urgency, d.SemesterID).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
}

func (r *pgxMaintenanceRepository) UpdateDamageReport(ctx context.Context, id string, status, vendor string, cost float64, notes string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	var itemID, assetName string
	err = tx.QueryRow(ctx, `SELECT COALESCE(inventory_item_id, ''), asset_name FROM damage_reports WHERE id = $1`, id).Scan(&itemID, &assetName)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `UPDATE damage_reports SET status = $1, vendor_name = $2, cost = $3, resolution_notes = $4, updated_at = NOW() WHERE id = $5`, status, vendor, cost, notes, id)
	if err != nil { return err }

	if itemID != "" {
		switch status {
		case "Proses":
			tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'DALAM_PERBAIKAN', updated_at = NOW() WHERE id = $1`, itemID)
			tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'MULAI_PERBAIKAN', $2, 'system')`, itemID, "Servis dimulai: "+assetName)
		case "Selesai":
			tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'AKTIF', updated_at = NOW() WHERE id = $1`, itemID)
			tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'SELESAI_PERBAIKAN', $2, 'system')`, itemID, "Servis selesai: "+assetName)
			tx.Exec(ctx, `INSERT INTO maintenance_logs (damage_report_id, inventory_item_id, service_date, vendor, cost, notes, performed_by) VALUES ($1, $2, NOW(), $3, $4, $5, 'system')`, id, itemID, vendor, cost, notes)
		}
	}
	return tx.Commit(ctx)
}

func (r *pgxMaintenanceRepository) GetAllSchedules(ctx context.Context, semesterID string) ([]*MaintenanceSchedule, error) {
	query := `
		SELECT id::text, schedule_code, asset_name, COALESCE(inventory_item_id, ''), COALESCE(scope_location, ''), 
		       frequency, last_service_date, next_service_date, is_active, COALESCE(semester_id, '')
		FROM maintenance_schedules
		WHERE ($1 = '' OR semester_id = $1) AND is_active = TRUE
		ORDER BY next_service_date ASC
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*MaintenanceSchedule
	for rows.Next() {
		var s MaintenanceSchedule
		err := rows.Scan(&s.ID, &s.ScheduleCode, &s.AssetName, &s.InventoryItemID, &s.ScopeLocation,
			&s.Frequency, &s.LastServiceDate, &s.NextServiceDate, &s.IsActive, &s.SemesterID)
		if err != nil { return nil, err }
		list = append(list, &s)
	}
	return list, nil
}

func (r *pgxMaintenanceRepository) CreateSchedule(ctx context.Context, s *MaintenanceSchedule) error {
	query := `
		INSERT INTO maintenance_schedules (schedule_code, asset_name, inventory_item_id, scope_location, frequency, last_service_date, next_service_date, semester_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`
	return r.db.QueryRow(ctx, query, s.ScheduleCode, s.AssetName, s.InventoryItemID, s.ScopeLocation, s.Frequency, s.LastServiceDate, s.NextServiceDate, s.SemesterID).Scan(&s.ID)
}

func (r *pgxMaintenanceRepository) CompleteSchedule(ctx context.Context, id string, serviceDate time.Time, cost float64, notes string, performedBy string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	var frequency, assetName, itemID string
	err = tx.QueryRow(ctx, `SELECT frequency, asset_name, COALESCE(inventory_item_id, '') FROM maintenance_schedules WHERE id = $1`, id).Scan(&frequency, &assetName, &itemID)
	if err != nil { return err }

	var nextDate time.Time
	switch frequency {
	case "Bulanan": nextDate = serviceDate.AddDate(0, 1, 0)
	case "3 Bulanan": nextDate = serviceDate.AddDate(0, 3, 0)
	case "6 Bulanan": nextDate = serviceDate.AddDate(0, 6, 0)
	case "Tahunan": nextDate = serviceDate.AddDate(1, 0, 0)
	default: nextDate = serviceDate.AddDate(0, 1, 0)
	}

	_, err = tx.Exec(ctx, `UPDATE maintenance_schedules SET last_service_date = $1, next_service_date = $2, updated_at = NOW() WHERE id = $3`, serviceDate, nextDate, id)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `INSERT INTO maintenance_logs (schedule_id, inventory_item_id, service_date, cost, notes, performed_by) VALUES ($1, $2, $3, $4, $5, $6)`, id, itemID, serviceDate, cost, notes, performedBy)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *pgxMaintenanceRepository) GetMaintenanceLogs(ctx context.Context, startDate, endDate string, semesterID string) ([]*MaintenanceLog, error) {
	query := `
		SELECT ml.id::text, COALESCE(ml.schedule_id::text, ''), COALESCE(ml.damage_report_id::text, ''), 
		       COALESCE(ml.inventory_item_id, ''), ml.service_date, COALESCE(ml.technician, ''), 
		       COALESCE(ml.vendor, ''), COALESCE(ml.cost, 0), COALESCE(ml.notes, ''), ml.performed_by, ml.created_at,
		       COALESCE(i.id, ''), COALESCE(c.name, '-')
		FROM maintenance_logs ml
		LEFT JOIN inventory_items i ON ml.inventory_item_id = i.id
		LEFT JOIN master_catalogs c ON i.catalog_id = c.id
		WHERE ($1 = '' OR ml.service_date >= $1::date)
		  AND ($2 = '' OR ml.service_date <= $2::date)
		  AND ($3 = '' OR ml.semester_id = $3)
		ORDER BY ml.service_date DESC
	`
	rows, err := r.db.Query(ctx, query, startDate, endDate, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*MaintenanceLog
	for rows.Next() {
		var l MaintenanceLog
		err := rows.Scan(&l.ID, &l.ScheduleID, &l.DamageReportID, &l.InventoryItemID, &l.ServiceDate, &l.Technician, &l.Vendor, &l.Cost, &l.Notes, &l.PerformedBy, &l.CreatedAt, &l.AssetCode, &l.AssetName)
		if err != nil { return nil, err }
		list = append(list, &l)
	}
	return list, nil
}
