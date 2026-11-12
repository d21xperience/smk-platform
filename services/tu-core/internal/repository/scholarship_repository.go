package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentScholarship struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	Category         string
	CardNumber       string
	Reason           string
	IsDocumentValid  bool
	Status           string
	BankName         string
	BankAccount      string
	DisbursedAmount  float64
	// JOIN fields
	StudentName string
	NISN        string
	ClassName   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ScholarshipFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Search           string
	Limit            int
	Offset           int
}

type ScholarshipRepository interface {
	Create(ctx context.Context, s *StudentScholarship) error
	GetByID(ctx context.Context, schoolID, periodID, scholarshipID string) (*StudentScholarship, error)
	List(ctx context.Context, filter ScholarshipFilter) ([]StudentScholarship, int, error)
	Update(ctx context.Context, s *StudentScholarship) error
	UpdateStatus(ctx context.Context, scholarshipID, status string) error
}

type pgxScholarshipRepository struct {
	db *pgxpool.Pool
}

func NewScholarshipRepository(db *pgxpool.Pool) ScholarshipRepository {
	return &pgxScholarshipRepository{db: db}
}

func (r *pgxScholarshipRepository) Create(ctx context.Context, s *StudentScholarship) error {
	query := `
		INSERT INTO student_scholarships (
			id, school_id, academic_period_id, student_id, category, card_number,
			reason, is_document_valid, status, bank_name, bank_account, disbursed_amount,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(ctx, query,
		s.ID, s.SchoolID, s.AcademicPeriodID, s.StudentID, s.Category, s.CardNumber,
		s.Reason, s.IsDocumentValid, s.Status, s.BankName, s.BankAccount, s.DisbursedAmount,
		s.CreatedAt, s.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxScholarshipRepository) GetByID(ctx context.Context, schoolID, periodID, scholarshipID string) (*StudentScholarship, error) {
	query := `
		SELECT 
			ss.id, ss.school_id, ss.academic_period_id, ss.student_id,
			ss.category, COALESCE(ss.card_number, ''), COALESCE(ss.reason, ''),
			ss.is_document_valid, ss.status, COALESCE(ss.bank_name, ''), COALESCE(ss.bank_account, ''),
			COALESCE(ss.disbursed_amount, 0),
			ts.nm_siswa, COALESCE(ts.nisn, ''), COALESCE(tk.nm_kelas, ''),
			ss.created_at, ss.updated_at
		FROM student_scholarships ss
		JOIN tabel_siswa ts ON ss.student_id = ts.peserta_didik_id
		LEFT JOIN tabel_anggotakelas tak ON ss.student_id = tak.peserta_didik_id AND tak.semester_id = ss.academic_period_id
		LEFT JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE ss.id = $1 AND ss.school_id = $2 AND ss.academic_period_id = $3
	`
	var s StudentScholarship
	err := r.db.QueryRow(ctx, query, scholarshipID, schoolID, periodID).Scan(
		&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.StudentID,
		&s.Category, &s.CardNumber, &s.Reason, &s.IsDocumentValid, &s.Status,
		&s.BankName, &s.BankAccount, &s.DisbursedAmount,
		&s.StudentName, &s.NISN, &s.ClassName, &s.CreatedAt, &s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &s, nil
}

func (r *pgxScholarshipRepository) List(ctx context.Context, filter ScholarshipFilter) ([]StudentScholarship, int, error) {
	whereClause := "WHERE ss.school_id = $1 AND ss.academic_period_id = $2"
	args := []any{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND ss.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (ts.nm_siswa ILIKE $%d OR ts.nisn ILIKE $%d OR COALESCE(ss.card_number, '') ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM student_scholarships ss
		JOIN tabel_siswa ts ON ss.student_id = ts.peserta_didik_id
		%s
	`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT 
			ss.id, ss.school_id, ss.academic_period_id, ss.student_id,
			ss.category, COALESCE(ss.card_number, ''), COALESCE(ss.reason, ''),
			ss.is_document_valid, ss.status, COALESCE(ss.bank_name, ''), COALESCE(ss.bank_account, ''),
			COALESCE(ss.disbursed_amount, 0),
			ts.nm_siswa, COALESCE(ts.nisn, ''), COALESCE(tk.nm_kelas, ''),
			ss.created_at, ss.updated_at
		FROM student_scholarships ss
		JOIN tabel_siswa ts ON ss.student_id = ts.peserta_didik_id
		LEFT JOIN tabel_anggotakelas tak ON ss.student_id = tak.peserta_didik_id AND tak.semester_id = ss.academic_period_id
		LEFT JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
		ORDER BY ts.nm_siswa ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var items []StudentScholarship
	for rows.Next() {
		var s StudentScholarship
		if err := rows.Scan(
			&s.ID, &s.SchoolID, &s.AcademicPeriodID, &s.StudentID,
			&s.Category, &s.CardNumber, &s.Reason, &s.IsDocumentValid, &s.Status,
			&s.BankName, &s.BankAccount, &s.DisbursedAmount,
			&s.StudentName, &s.NISN, &s.ClassName, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		items = append(items, s)
	}
	return items, total, rows.Err()
}

func (r *pgxScholarshipRepository) Update(ctx context.Context, s *StudentScholarship) error {
	query := `
		UPDATE student_scholarships SET
			category = $1, card_number = $2, reason = $3, is_document_valid = $4,
			status = $5, bank_name = $6, bank_account = $7, disbursed_amount = $8,
			updated_at = NOW()
		WHERE id = $9 AND school_id = $10 AND academic_period_id = $11
	`
	result, err := r.db.Exec(ctx, query,
		s.Category, s.CardNumber, s.Reason, s.IsDocumentValid,
		s.Status, s.BankName, s.BankAccount, s.DisbursedAmount,
		s.ID, s.SchoolID, s.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxScholarshipRepository) UpdateStatus(ctx context.Context, scholarshipID, status string) error {
	query := `UPDATE student_scholarships SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, status, scholarshipID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
