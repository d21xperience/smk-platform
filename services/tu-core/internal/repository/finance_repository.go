package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/domain/finance/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrInvoiceNotFound = errors.New("invoice not found")
	ErrPaymentNotFound = errors.New("payment not found")
)

type InvoiceFilter struct {
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	ComponentType    string
	Status           string
	Month            int
	Year             int
	Limit            int
	Offset           int
}

type PaymentFilter struct {
	SchoolID         string
	AcademicPeriodID string
	InvoiceID        string
	PaymentMethod    string
	StartDate        time.Time
	EndDate          time.Time
	Limit            int
	Offset           int
}

type FinanceRepository interface {
	CreateInvoice(ctx context.Context, i *models.Invoice) error
	GetInvoiceByID(ctx context.Context, schoolID, academicPeriodID, invoiceID string) (*models.Invoice, error)
	ListInvoices(ctx context.Context, filter InvoiceFilter) ([]models.Invoice, int, error)
	UpdateInvoice(ctx context.Context, i *models.Invoice) error
	CreatePayment(ctx context.Context, p *models.Payment) error
	ListPayments(ctx context.Context, filter PaymentFilter) ([]models.Payment, int, error)
	GetPaymentSummary(ctx context.Context, schoolID, academicPeriodID string, month, year int, componentType string) (*PaymentSummary, error)
}

type PaymentSummary struct {
	TotalInvoices     int
	PaidInvoices      int
	UnpaidInvoices    int
	TotalAmount       float64
	TotalPaid         float64
	TotalOutstanding  float64
	PaymentPercentage float64
}

type pgxFinanceRepository struct {
	db *pgxpool.Pool
}

func NewFinanceRepository(db *pgxpool.Pool) FinanceRepository {
	return &pgxFinanceRepository{db: db}
}

func (r *pgxFinanceRepository) CreateInvoice(ctx context.Context, i *models.Invoice) error {
	query := `
		INSERT INTO invoices (
			id, school_id, academic_period_id, student_id, component_type,
			month, year, amount, due_date, status, paid_amount, notes,
			created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	`
	_, err := r.db.Exec(ctx, query,
		i.ID, i.SchoolID, i.AcademicPeriodID, i.StudentID, i.ComponentType.String(),
		i.Month, i.Year, i.Amount, i.DueDate, i.Status.String(), i.PaidAmount, i.Notes,
		i.CreatedBy, i.CreatedAt, i.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxFinanceRepository) GetInvoiceByID(ctx context.Context, schoolID, academicPeriodID, invoiceID string) (*models.Invoice, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, component_type,
		       month, year, amount, due_date, status, paid_amount,
		       paid_at, COALESCE(notes, '') as notes, COALESCE(created_by, '') as created_by,
		       created_at, updated_at
		FROM invoices
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var i models.Invoice
	var componentTypeStr, statusStr string
	var paidAt *time.Time

	err := r.db.QueryRow(ctx, query, invoiceID, schoolID, academicPeriodID).Scan(
		&i.ID, &i.SchoolID, &i.AcademicPeriodID, &i.StudentID, &componentTypeStr,
		&i.Month, &i.Year, &i.Amount, &i.DueDate, &statusStr, &i.PaidAmount,
		&paidAt, &i.Notes, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvoiceNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	i.ComponentType, _ = valueobjects.NewComponentType(componentTypeStr)
	i.Status, _ = valueobjects.NewInvoiceStatus(statusStr)
	if paidAt != nil {
		i.PaidAt = *paidAt
	}

	return &i, nil
}

func (r *pgxFinanceRepository) ListInvoices(ctx context.Context, filter InvoiceFilter) ([]models.Invoice, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.StudentID != "" {
		whereClause += fmt.Sprintf(" AND student_id = $%d", argIndex)
		args = append(args, filter.StudentID)
		argIndex++
	}
	if filter.ComponentType != "" {
		whereClause += fmt.Sprintf(" AND component_type = $%d", argIndex)
		args = append(args, filter.ComponentType)
		argIndex++
	}
	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.Month > 0 {
		whereClause += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, filter.Month)
		argIndex++
	}
	if filter.Year > 0 {
		whereClause += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, filter.Year)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM invoices %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, student_id, component_type,
		       month, year, amount, due_date, status, paid_amount,
		       paid_at, COALESCE(notes, '') as notes, COALESCE(created_by, '') as created_by,
		       created_at, updated_at
		FROM invoices %s
		ORDER BY due_date ASC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var invoices []models.Invoice
	for rows.Next() {
		var i models.Invoice
		var componentTypeStr, statusStr string
		var paidAt *time.Time

		if err := rows.Scan(
			&i.ID, &i.SchoolID, &i.AcademicPeriodID, &i.StudentID, &componentTypeStr,
			&i.Month, &i.Year, &i.Amount, &i.DueDate, &statusStr, &i.PaidAmount,
			&paidAt, &i.Notes, &i.CreatedBy, &i.CreatedAt, &i.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		i.ComponentType, _ = valueobjects.NewComponentType(componentTypeStr)
		i.Status, _ = valueobjects.NewInvoiceStatus(statusStr)
		if paidAt != nil {
			i.PaidAt = *paidAt
		}

		invoices = append(invoices, i)
	}

	return invoices, total, rows.Err()
}

func (r *pgxFinanceRepository) UpdateInvoice(ctx context.Context, i *models.Invoice) error {
	query := `
		UPDATE invoices
		SET status = $1, paid_amount = $2, paid_at = $3, notes = $4, updated_at = $5
		WHERE id = $6 AND school_id = $7 AND academic_period_id = $8
	`

	var paidAt interface{}
	if i.PaidAt.IsZero() {
		paidAt = nil
	} else {
		paidAt = i.PaidAt
	}

	result, err := r.db.Exec(ctx, query,
		i.Status.String(), i.PaidAmount, paidAt, i.Notes, i.UpdatedAt,
		i.ID, i.SchoolID, i.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrInvoiceNotFound
	}
	return nil
}

func (r *pgxFinanceRepository) CreatePayment(ctx context.Context, p *models.Payment) error {
	query := `
		INSERT INTO payments (
			id, school_id, academic_period_id, invoice_id, amount, payment_method,
			payment_date, reference_number, notes, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.Exec(ctx, query,
		p.ID, p.SchoolID, p.AcademicPeriodID, p.InvoiceID, p.Amount, p.PaymentMethod.String(),
		p.PaymentDate, p.ReferenceNumber, p.Notes, p.CreatedBy, p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxFinanceRepository) ListPayments(ctx context.Context, filter PaymentFilter) ([]models.Payment, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.InvoiceID != "" {
		whereClause += fmt.Sprintf(" AND invoice_id = $%d", argIndex)
		args = append(args, filter.InvoiceID)
		argIndex++
	}
	if filter.PaymentMethod != "" {
		whereClause += fmt.Sprintf(" AND payment_method = $%d", argIndex)
		args = append(args, filter.PaymentMethod)
		argIndex++
	}
	if !filter.StartDate.IsZero() {
		whereClause += fmt.Sprintf(" AND payment_date >= $%d", argIndex)
		args = append(args, filter.StartDate)
		argIndex++
	}
	if !filter.EndDate.IsZero() {
		whereClause += fmt.Sprintf(" AND payment_date <= $%d", argIndex)
		args = append(args, filter.EndDate)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM payments %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, invoice_id, amount, payment_method,
		       payment_date, COALESCE(reference_number, '') as reference_number,
		       COALESCE(notes, '') as notes, COALESCE(created_by, '') as created_by,
		       created_at, updated_at
		FROM payments %s
		ORDER BY payment_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var p models.Payment
		var paymentMethodStr string

		if err := rows.Scan(
			&p.ID, &p.SchoolID, &p.AcademicPeriodID, &p.InvoiceID, &p.Amount, &paymentMethodStr,
			&p.PaymentDate, &p.ReferenceNumber, &p.Notes, &p.CreatedBy, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		p.PaymentMethod, _ = valueobjects.NewPaymentMethod(paymentMethodStr)
		payments = append(payments, p)
	}

	return payments, total, rows.Err()
}

func (r *pgxFinanceRepository) GetPaymentSummary(ctx context.Context, schoolID, academicPeriodID string, month, year int, componentType string) (*PaymentSummary, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{schoolID, academicPeriodID}
	argIndex := 3

	if month > 0 {
		whereClause += fmt.Sprintf(" AND month = $%d", argIndex)
		args = append(args, month)
		argIndex++
	}
	if year > 0 {
		whereClause += fmt.Sprintf(" AND year = $%d", argIndex)
		args = append(args, year)
		argIndex++
	}
	if componentType != "" {
		whereClause += fmt.Sprintf(" AND component_type = $%d", argIndex)
		args = append(args, componentType)
		argIndex++
	}

	query := fmt.Sprintf(`
		SELECT
			COUNT(*) as total_invoices,
			COUNT(*) FILTER (WHERE status = 'PAID') as paid_invoices,
			COUNT(*) FILTER (WHERE status != 'PAID' AND status != 'CANCELLED') as unpaid_invoices,
			COALESCE(SUM(amount), 0) as total_amount,
			COALESCE(SUM(paid_amount), 0) as total_paid,
			COALESCE(SUM(amount - paid_amount) FILTER (WHERE status != 'CANCELLED'), 0) as total_outstanding
		FROM invoices %s
	`, whereClause)

	var summary PaymentSummary
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&summary.TotalInvoices,
		&summary.PaidInvoices,
		&summary.UnpaidInvoices,
		&summary.TotalAmount,
		&summary.TotalPaid,
		&summary.TotalOutstanding,
	)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if summary.TotalAmount > 0 {
		summary.PaymentPercentage = (summary.TotalPaid / summary.TotalAmount) * 100
	}

	return &summary, nil
}
