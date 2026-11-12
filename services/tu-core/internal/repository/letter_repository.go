package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentLetter struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	LetterNumber     string
	ParentName       string
	Purpose          string
	Status           string
	FileURL          string
	StudentName      string // JOIN dari tabel_siswa
	ClassName        string // JOIN dari tabel_kelas
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type LetterFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Search           string
	Limit            int
	Offset           int
}

type LetterRepository interface {
	Create(ctx context.Context, l *StudentLetter) error
	GetByID(ctx context.Context, schoolID, periodID, letterID string) (*StudentLetter, error)
	List(ctx context.Context, filter LetterFilter) ([]StudentLetter, int, error)
	UpdateStatus(ctx context.Context, letterID, status string) error
}

type pgxLetterRepository struct {
	db *pgxpool.Pool
}

func NewLetterRepository(db *pgxpool.Pool) LetterRepository {
	return &pgxLetterRepository{db: db}
}

func (r *pgxLetterRepository) Create(ctx context.Context, l *StudentLetter) error {
	query := `
		INSERT INTO student_letters (
			id, school_id, academic_period_id, student_id, letter_number,
			parent_name, purpose, status, file_url, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		l.ID, l.SchoolID, l.AcademicPeriodID, l.StudentID, l.LetterNumber,
		l.ParentName, l.Purpose, l.Status, l.FileURL, l.CreatedAt, l.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxLetterRepository) GetByID(ctx context.Context, schoolID, periodID, letterID string) (*StudentLetter, error) {
	query := `
		SELECT 
			sl.id, sl.school_id, sl.academic_period_id, sl.student_id, sl.letter_number,
			COALESCE(sl.parent_name, ''), COALESCE(sl.purpose, ''), sl.status, COALESCE(sl.file_url, ''),
			ts.nm_siswa, COALESCE(tk.nm_kelas, ''), sl.created_at, sl.updated_at
		FROM student_letters sl
		JOIN tabel_siswa ts ON sl.student_id = ts.peserta_didik_id
		LEFT JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id AND tak.semester_id = sl.academic_period_id
		LEFT JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE sl.id = $1 AND sl.school_id = $2 AND sl.academic_period_id = $3
	`
	var l StudentLetter
	err := r.db.QueryRow(ctx, query, letterID, schoolID, periodID).Scan(
		&l.ID, &l.SchoolID, &l.AcademicPeriodID, &l.StudentID, &l.LetterNumber,
		&l.ParentName, &l.Purpose, &l.Status, &l.FileURL,
		&l.StudentName, &l.ClassName, &l.CreatedAt, &l.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &l, nil
}

func (r *pgxLetterRepository) List(ctx context.Context, filter LetterFilter) ([]StudentLetter, int, error) {
	whereClause := "WHERE sl.school_id = $1 AND sl.academic_period_id = $2"
	args := []any{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND sl.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}
	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (ts.nm_siswa ILIKE $%d OR sl.letter_number ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) 
		FROM student_letters sl
		JOIN tabel_siswa ts ON sl.student_id = ts.peserta_didik_id
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
			sl.id, sl.school_id, sl.academic_period_id, sl.student_id, sl.letter_number,
			COALESCE(sl.parent_name, ''), COALESCE(sl.purpose, ''), sl.status, COALESCE(sl.file_url, ''),
			ts.nm_siswa, COALESCE(tk.nm_kelas, ''), sl.created_at, sl.updated_at
		FROM student_letters sl
		JOIN tabel_siswa ts ON sl.student_id = ts.peserta_didik_id
		LEFT JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id AND tak.semester_id = sl.academic_period_id
		LEFT JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
		ORDER BY sl.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var letters []StudentLetter
	for rows.Next() {
		var l StudentLetter
		if err := rows.Scan(
			&l.ID, &l.SchoolID, &l.AcademicPeriodID, &l.StudentID, &l.LetterNumber,
			&l.ParentName, &l.Purpose, &l.Status, &l.FileURL,
			&l.StudentName, &l.ClassName, &l.CreatedAt, &l.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		letters = append(letters, l)
	}
	return letters, total, rows.Err()
}

func (r *pgxLetterRepository) UpdateStatus(ctx context.Context, letterID, status string) error {
	query := `UPDATE student_letters SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.Exec(ctx, query, status, letterID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
