package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Borrowing struct {
	ID                    string     `json:"id"`
	BorrowingNumber       string     `json:"borrowing_number"`
	BorrowerID            string     `json:"borrower_id"`
	BorrowerName          string     `json:"borrower_name"`
	InventoryItemID       string     `json:"inventory_item_id"`
	PurposeType           string     `json:"purpose_type"`
	PurposeDescription    string     `json:"purpose_description"`
	BorrowDate            time.Time  `json:"borrow_date"`
	ExpectedReturnDate    time.Time  `json:"expected_return_date"`
	UsageLocation         string     `json:"usage_location"`
	DurationDays          int        `json:"duration_days"`
	ApprovalLevel         int        `json:"approval_level"`
	Status                string     `json:"status"`
	CurrentApprovalStep   int        `json:"current_approval_step"`
	RejectionReason       string     `json:"rejection_reason"`
	SemesterID            string     `json:"semester_id"`
	
	// Pickup
	ActualPickupDate *time.Time `json:"actual_pickup_date"`
	PickupBy         string     `json:"pickup_by"`
	PickupCondition  string     `json:"pickup_condition"`
	PickupPhotos     string     `json:"pickup_photos"`
	PickupNotes      string     `json:"pickup_notes"`
	BastNumber       string     `json:"bast_number"`
	BastPdfUrl       string     `json:"bast_pdf_url"`
	
	// Return
	ActualReturnDate *time.Time `json:"actual_return_date"`
	ReturnBy         string     `json:"return_by"`
	ReturnCondition  string     `json:"return_condition"`
	ReturnPhotos     string     `json:"return_photos"`
	ReturnNotes      string     `json:"return_notes"`
	DaysLate         int        `json:"days_late"`
	FineAmount       float64    `json:"fine_amount"`
	FinePaid         bool       `json:"fine_paid"`
	BapNumber        string     `json:"bap_number"`
	BapPdfUrl        string     `json:"bap_pdf_url"`
	VerificationHash string     `json:"verification_hash"`
	
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Joined fields from inventory_items & catalogs
	AssetCode       string  `json:"asset_code"`
	AssetName       string  `json:"asset_name"`
	Merk            string  `json:"merk"`
	CategoryName    string  `json:"category_name"`
	RoomName        string  `json:"room_name"`
	AssetCondition  string  `json:"asset_condition"`
	Approvals       []BorrowingApproval `json:"approvals,omitempty"`
}

type BorrowingApproval struct {
	ID              string    `json:"id"`
	BorrowingID     string    `json:"borrowing_id"`
	ApprovalStep    int       `json:"approval_step"`
	ApproverID      string    `json:"approver_id"`
	Decision        string    `json:"decision"`
	Notes           string    `json:"notes"`
	ApprovalDate    time.Time `json:"approval_date"`
}

type BorrowingRepository interface {
	GenerateBorrowingNumber(ctx context.Context) (string, error)
	GenerateBASTNumber(ctx context.Context) (string, error)
	GenerateBAPNumber(ctx context.Context) (string, error)
	CountConflictingBorrowings(ctx context.Context, itemID string, borrowDate, expectedReturnDate time.Time) (int, error)
	
	Create(ctx context.Context, b *Borrowing) error
	GetAll(ctx context.Context, semesterID string) ([]*Borrowing, error)
	GetByID(ctx context.Context, id string) (*Borrowing, error)
	
	Approve(ctx context.Context, b *Borrowing, approval *BorrowingApproval) error
	Pickup(ctx context.Context, b *Borrowing) error
	Return(ctx context.Context, b *Borrowing, fine *BorrowingFine) error
}

type BorrowingFine struct {
	ID          string    `json:"id"`
	BorrowingID string    `json:"borrowing_id"`
	BorrowerID  string    `json:"borrower_id"`
	FineAmount  float64   `json:"fine_amount"`
	DaysLate    int       `json:"days_late"`
	FinePerDay  float64   `json:"fine_per_day"`
	FineDate    time.Time `json:"fine_date"`
	Status      string    `json:"status"`
}

type pgxBorrowingRepository struct {
	db *pgxpool.Pool
}

func NewBorrowingRepository(db *pgxpool.Pool) BorrowingRepository {
	return &pgxBorrowingRepository{db: db}
}

func (r *pgxBorrowingRepository) GenerateBorrowingNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM borrowings WHERE EXTRACT(YEAR FROM borrow_date) = $1`, year).Scan(&count)
	var schCode string
	r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schCode)
	return fmt.Sprintf("PJM/%s/%d/%03d", schCode, year, count+1), nil
}

func (r *pgxBorrowingRepository) GenerateBASTNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM borrowings WHERE bast_number IS NOT NULL AND EXTRACT(YEAR FROM actual_pickup_date) = $1`, year).Scan(&count)
	var schCode string
	r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schCode)
	return fmt.Sprintf("BAST-PJM/%s/%d/%03d", schCode, year, count+1), nil
}

func (r *pgxBorrowingRepository) GenerateBAPNumber(ctx context.Context) (string, error) {
	year := time.Now().Year()
	var count int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM borrowings WHERE bap_number IS NOT NULL AND EXTRACT(YEAR FROM actual_return_date) = $1`, year).Scan(&count)
	var schCode string
	r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schCode)
	return fmt.Sprintf("BAP-PJM/%s/%d/%03d", schCode, year, count+1), nil
}

func (r *pgxBorrowingRepository) CountConflictingBorrowings(ctx context.Context, itemID string, borrowDate, expectedReturnDate time.Time) (int, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM borrowings 
		WHERE inventory_item_id = $1 
		AND status IN ('PENDING_APPROVAL', 'APPROVED', 'APPROVED_STEP_1', 'APPROVED_STEP_2', 'APPROVED_STEP_3', 'PICKED_UP')
		AND borrow_date <= $3
		AND expected_return_date >= $2
	`
	err := r.db.QueryRow(ctx, query, itemID, borrowDate, expectedReturnDate).Scan(&count)
	return count, err
}

func (r *pgxBorrowingRepository) Create(ctx context.Context, b *Borrowing) error {
	query := `
		INSERT INTO borrowings 
		(borrowing_number, borrower_id, borrower_name, inventory_item_id, purpose_type, purpose_description, borrow_date, expected_return_date, usage_location, duration_days, approval_level, status, current_approval_step, semester_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, 'PENDING_APPROVAL', 0, $12)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query, b.BorrowingNumber, b.BorrowerID, b.BorrowerName, b.InventoryItemID, b.PurposeType, b.PurposeDescription, b.BorrowDate, b.ExpectedReturnDate, b.UsageLocation, b.DurationDays, b.ApprovalLevel, b.SemesterID).Scan(&b.ID, &b.CreatedAt, &b.UpdatedAt)
}

func (r *pgxBorrowingRepository) GetAll(ctx context.Context, semesterID string) ([]*Borrowing, error) {
	query := `
		SELECT b.id::text, b.borrowing_number, b.borrower_id, b.borrower_name, b.inventory_item_id, b.purpose_type, COALESCE(b.purpose_description, ''), 
		       b.borrow_date, b.expected_return_date, COALESCE(b.usage_location, ''), b.duration_days, b.approval_level, b.status, b.current_approval_step, 
		       COALESCE(b.rejection_reason, ''), COALESCE(b.semester_id, ''), b.actual_pickup_date, COALESCE(b.pickup_by, ''), COALESCE(b.pickup_condition, ''), COALESCE(b.bast_number, ''),
		       b.actual_return_date, COALESCE(b.bap_number, ''), b.days_late, b.fine_amount, b.fine_paid, b.created_at, b.updated_at,
		       i.id, c.name, COALESCE(c.brand, ''), cat.name, COALESCE(i.condition, ''), COALESCE(k.nm_kelas, '')
		FROM borrowings b
		JOIN inventory_items i ON b.inventory_item_id = i.id
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		LEFT JOIN tabel_kelas k ON i.room_id = k.rombongan_belajar_id
		WHERE b.semester_id = $1 OR $1 = ''
		ORDER BY b.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*Borrowing
	for rows.Next() {
		var b Borrowing
		err := rows.Scan(
			&b.ID, &b.BorrowingNumber, &b.BorrowerID, &b.BorrowerName, &b.InventoryItemID, &b.PurposeType, &b.PurposeDescription,
			&b.BorrowDate, &b.ExpectedReturnDate, &b.UsageLocation, &b.DurationDays, &b.ApprovalLevel, &b.Status, &b.CurrentApprovalStep,
			&b.RejectionReason, &b.SemesterID, &b.ActualPickupDate, &b.PickupBy, &b.PickupCondition, &b.BastNumber,
			&b.ActualReturnDate, &b.BapNumber, &b.DaysLate, &b.FineAmount, &b.FinePaid, &b.CreatedAt, &b.UpdatedAt,
			&b.AssetCode, &b.AssetName, &b.Merk, &b.CategoryName, &b.AssetCondition, &b.RoomName,
		)
		if err != nil { return nil, err }
		list = append(list, &b)
	}
	return list, nil
}

func (r *pgxBorrowingRepository) GetByID(ctx context.Context, id string) (*Borrowing, error) {
	query := `
		SELECT b.id::text, b.borrowing_number, b.borrower_id, b.borrower_name, b.inventory_item_id, b.purpose_type, COALESCE(b.purpose_description, ''), 
		       b.borrow_date, b.expected_return_date, COALESCE(b.usage_location, ''), b.duration_days, b.approval_level, b.status, b.current_approval_step, 
		       COALESCE(b.rejection_reason, ''), COALESCE(b.semester_id, ''), b.actual_pickup_date, COALESCE(b.pickup_by, ''), COALESCE(b.pickup_condition, ''), COALESCE(b.bast_number, ''),
		       b.actual_return_date, COALESCE(b.bap_number, ''), b.days_late, b.fine_amount, b.fine_paid, b.created_at, b.updated_at,
		       i.id, c.name, COALESCE(c.brand, ''), cat.name, COALESCE(i.condition, ''), COALESCE(k.nm_kelas, '')
		FROM borrowings b
		JOIN inventory_items i ON b.inventory_item_id = i.id
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		LEFT JOIN tabel_kelas k ON i.room_id = k.rombongan_belajar_id
		WHERE b.id = $1
	`
	var b Borrowing
	err := r.db.QueryRow(ctx, query, id).Scan(
		&b.ID, &b.BorrowingNumber, &b.BorrowerID, &b.BorrowerName, &b.InventoryItemID, &b.PurposeType, &b.PurposeDescription,
		&b.BorrowDate, &b.ExpectedReturnDate, &b.UsageLocation, &b.DurationDays, &b.ApprovalLevel, &b.Status, &b.CurrentApprovalStep,
		&b.RejectionReason, &b.SemesterID, &b.ActualPickupDate, &b.PickupBy, &b.PickupCondition, &b.BastNumber,
		&b.ActualReturnDate, &b.BapNumber, &b.DaysLate, &b.FineAmount, &b.FinePaid, &b.CreatedAt, &b.UpdatedAt,
		&b.AssetCode, &b.AssetName, &b.Merk, &b.CategoryName, &b.AssetCondition, &b.RoomName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, errors.New("peminjaman tidak ditemukan") }
		return nil, err
	}

	// Get Approvals
	appRows, _ := r.db.Query(ctx, `SELECT id::text, borrowing_id::text, approval_step, approver_id, decision, COALESCE(notes, ''), approval_date FROM borrowing_approvals WHERE borrowing_id = $1 ORDER BY approval_step ASC`, id)
	defer appRows.Close()
	for appRows.Next() {
		var a BorrowingApproval
		appRows.Scan(&a.ID, &a.BorrowingID, &a.ApprovalStep, &a.ApproverID, &a.Decision, &a.Notes, &a.ApprovalDate)
		b.Approvals = append(b.Approvals, a)
	}

	return &b, nil
}

func (r *pgxBorrowingRepository) Approve(ctx context.Context, b *Borrowing, approval *BorrowingApproval) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	// Insert approval record
	_, err = tx.Exec(ctx, `INSERT INTO borrowing_approvals (borrowing_id, approval_step, approver_id, decision, notes) VALUES ($1, $2, $3, $4, $5)`,
		b.ID, approval.ApprovalStep, approval.ApproverID, approval.Decision, approval.Notes)
	if err != nil { return err }

	// Determine new status
	newStatus := fmt.Sprintf("APPROVED_STEP_%d", approval.ApprovalStep)
	if approval.Decision == "REJECTED" {
		newStatus = "REJECTED"
	} else if approval.ApprovalStep == b.ApprovalLevel {
		newStatus = "APPROVED"
	}

	_, err = tx.Exec(ctx, `UPDATE borrowings SET status = $1, current_approval_step = $2, rejection_reason = $3, updated_at = NOW() WHERE id = $4`,
		newStatus, approval.ApprovalStep, approval.Notes, b.ID)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *pgxBorrowingRepository) Pickup(ctx context.Context, b *Borrowing) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	photos, _ := json.Marshal(b.PickupPhotos)
	_, err = tx.Exec(ctx, `UPDATE borrowings SET status = 'PICKED_UP', actual_pickup_date = $1, pickup_by = $2, pickup_condition = $3, pickup_photos = $4, pickup_notes = $5, bast_number = $6, updated_at = NOW() WHERE id = $7`,
		b.ActualPickupDate, b.PickupBy, b.PickupCondition, photos, b.PickupNotes, b.BastNumber, b.ID)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'DIPINJAM', current_borrowing_id = $1, updated_at = NOW() WHERE id = $2`, b.ID, b.InventoryItemID)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'PEMINJAMAN', $2, $3)`,
		b.InventoryItemID, "Dipinjam oleh "+b.BorrowerName+" ("+b.BorrowingNumber+")", b.PickupBy)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *pgxBorrowingRepository) Return(ctx context.Context, b *Borrowing, fine *BorrowingFine) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	photos, _ := json.Marshal(b.ReturnPhotos)
	_, err = tx.Exec(ctx, `UPDATE borrowings SET status = 'RETURNED', actual_return_date = $1, return_by = $2, return_condition = $3, return_photos = $4, return_notes = $5, days_late = $6, fine_amount = $7, fine_paid = $8, bap_number = $9, updated_at = NOW() WHERE id = $10`,
		b.ActualReturnDate, b.ReturnBy, b.ReturnCondition, photos, b.ReturnNotes, b.DaysLate, b.FineAmount, b.FinePaid, b.BapNumber, b.ID)
	if err != nil { return err }

	// Update asset: back to AKTIF, and set its condition based on return
	_, err = tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'AKTIF', current_borrowing_id = NULL, condition = $1, updated_at = NOW() WHERE id = $2`,
		b.ReturnCondition, b.InventoryItemID)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'PENGEMBALIAN', $2, $3)`,
		b.InventoryItemID, "Dikembalikan oleh "+b.BorrowerName+" ("+b.BorrowingNumber+")", b.ReturnBy)
	if err != nil { return err }

	// Insert fine record if any
	if fine != nil && fine.FineAmount > 0 {
		_, err = tx.Exec(ctx, `INSERT INTO borrowing_fines (borrowing_id, borrower_id, fine_amount, days_late, fine_per_day, fine_date, status) VALUES ($1, $2, $3, $4, $5, $6, 'UNPAID')`,
			fine.BorrowingID, fine.BorrowerID, fine.FineAmount, fine.DaysLate, fine.FinePerDay, fine.FineDate)
		if err != nil { return err }
	}

	return tx.Commit(ctx)
}
