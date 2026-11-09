package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcurementRequest struct {
	ID              string     `json:"id"`
	RequestNumber   string     `json:"request_number"`
	RequesterID     string     `json:"requester_id"`
	RequesterName   string     `json:"requester_name"`
	ItemName        string     `json:"item_name"`
	Specification   string     `json:"specification"`
	Quantity        int        `json:"quantity"`
	Unit            string     `json:"unit"`
	EstimatedPrice  float64    `json:"estimated_price"`
	TotalEstimated  float64    `json:"total_estimated"`
	Priority        string     `json:"priority"`
	Justification   string     `json:"justification"`
	Status          string     `json:"status"`
	SemesterID      string     `json:"semester_id"` // <-- NEW
	ApprovedBy      *string    `json:"approved_by"`
	ApprovedAt      *time.Time `json:"approved_at"`
	RejectionReason *string    `json:"rejection_reason"`
	Notes           *string    `json:"notes"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type ProcurementRepository interface {
	Create(ctx context.Context, req *ProcurementRequest) error
	GetAll(ctx context.Context, semesterID string) ([]*ProcurementRequest, error)
	GetByID(ctx context.Context, id string) (*ProcurementRequest, error)
	Update(ctx context.Context, req *ProcurementRequest) error
	UpdateStatus(ctx context.Context, id string, status string, approvedBy string, rejectionReason *string, notes *string) error
	Delete(ctx context.Context, id string) error
	GenerateRequestNumber(ctx context.Context, semesterID string) (string, error)
}

type pgxProcurementRepository struct {
	db *pgxpool.Pool
}

func NewProcurementRepository(db *pgxpool.Pool) ProcurementRepository {
	return &pgxProcurementRepository{db: db}
}

func (r *pgxProcurementRepository) GenerateRequestNumber(ctx context.Context, semesterID string) (string, error) {
	year := time.Now().Year()
	var count int
	// Hitung berdasarkan tahun dan semester
	query := `SELECT COUNT(*) FROM procurement_requests WHERE EXTRACT(YEAR FROM created_at) = $1 AND semester_id = $2`
	err := r.db.QueryRow(ctx, query, year, semesterID).Scan(&count)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("REQ/%d/%03d", year, count+1), nil
}

func (r *pgxProcurementRepository) Create(ctx context.Context, req *ProcurementRequest) error {
	query := `
		INSERT INTO procurement_requests 
		(request_number, requester_id, requester_name, item_name, specification, quantity, unit, estimated_price, total_estimated, priority, justification, status, semester_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		req.RequestNumber, req.RequesterID, req.RequesterName, req.ItemName, req.Specification,
		req.Quantity, req.Unit, req.EstimatedPrice, req.TotalEstimated, req.Priority, req.Justification, req.Status, req.SemesterID,
	).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	return err
}

func (r *pgxProcurementRepository) GetAll(ctx context.Context, semesterID string) ([]*ProcurementRequest, error) {
	query := `
		SELECT id::text, request_number, requester_id, requester_name, item_name, COALESCE(specification, ''),
		       quantity, unit, COALESCE(estimated_price, 0), COALESCE(total_estimated, 0), priority, justification,
		       status, COALESCE(semester_id, ''), approved_by, approved_at, rejection_reason, notes, created_at, updated_at
		FROM procurement_requests
		WHERE semester_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*ProcurementRequest
	for rows.Next() {
		var req ProcurementRequest
		err := rows.Scan(
			&req.ID, &req.RequestNumber, &req.RequesterID, &req.RequesterName, &req.ItemName, &req.Specification,
			&req.Quantity, &req.Unit, &req.EstimatedPrice, &req.TotalEstimated, &req.Priority, &req.Justification,
			&req.Status, &req.SemesterID, &req.ApprovedBy, &req.ApprovedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, &req)
	}
	return list, nil
}

func (r *pgxProcurementRepository) GetByID(ctx context.Context, id string) (*ProcurementRequest, error) {
	query := `
		SELECT id::text, request_number, requester_id, requester_name, item_name, COALESCE(specification, ''),
		       quantity, unit, COALESCE(estimated_price, 0), COALESCE(total_estimated, 0), priority, justification,
		       status, COALESCE(semester_id, ''), approved_by, approved_at, rejection_reason, notes, created_at, updated_at
		FROM procurement_requests WHERE id = $1
	`
	var req ProcurementRequest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&req.ID, &req.RequestNumber, &req.RequesterID, &req.RequesterName, &req.ItemName, &req.Specification,
		&req.Quantity, &req.Unit, &req.EstimatedPrice, &req.TotalEstimated, &req.Priority, &req.Justification,
		&req.Status, &req.SemesterID, &req.ApprovedBy, &req.ApprovedAt, &req.RejectionReason, &req.Notes, &req.CreatedAt, &req.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("pengajuan tidak ditemukan")
		}
		return nil, err
	}
	return &req, nil
}

func (r *pgxProcurementRepository) Update(ctx context.Context, req *ProcurementRequest) error {
	query := `
		UPDATE procurement_requests 
		SET item_name = $1, specification = $2, quantity = $3, unit = $4, estimated_price = $5, total_estimated = $6, priority = $7, justification = $8
		WHERE id = $9 AND status = 'PENDING'
	`
	_, err := r.db.Exec(ctx, query, req.ItemName, req.Specification, req.Quantity, req.Unit, req.EstimatedPrice, req.TotalEstimated, req.Priority, req.Justification, req.ID)
	return err
}

func (r *pgxProcurementRepository) UpdateStatus(ctx context.Context, id string, status string, approvedBy string, rejectionReason *string, notes *string) error {
	query := `
		UPDATE procurement_requests 
		SET status = $1, approved_by = $2, approved_at = NOW(), rejection_reason = $3, notes = $4
		WHERE id = $5
	`
	_, err := r.db.Exec(ctx, query, status, approvedBy, rejectionReason, notes, id)
	return err
}

func (r *pgxProcurementRepository) Delete(ctx context.Context, id string) error {
	// Hanya bisa hapus jika status PENDING
	query := `DELETE FROM procurement_requests WHERE id = $1 AND status = 'PENDING'`
	_, err := r.db.Exec(ctx, query, id)
	return err
}
