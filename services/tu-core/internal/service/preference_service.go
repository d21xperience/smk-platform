package service

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type DBExecutor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Period struct {
	SemesterID       string `json:"semester_id"`
	TahunAjaranID    string `json:"tahun_ajaran_id"`
	TahunAjaranLabel string `json:"tahun_ajaran_label"`
	NamaSemester     string `json:"nama_semester"`
	Semester         int    `json:"semester"`
	PeriodeAktif     bool   `json:"periode_aktif"`
}

type PreferenceService struct {
	db DBExecutor
}

func NewPreferenceService(db DBExecutor) *PreferenceService {
	return &PreferenceService{db: db}
}

func (s *PreferenceService) ListPeriods(ctx context.Context) ([]Period, error) {
	query := `
		SELECT
			s.id AS semester_id,
			SUBSTRING(s.academic_year FROM 1 FOR 4) AS tahun_ajaran_id,
			s.academic_year AS tahun_ajaran_label,
			s.name AS nama_semester,
			s.semester,
			s.is_active
		FROM semesters s
		ORDER BY s.id DESC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list periods: %w", err)
	}
	defer rows.Close()

	var result []Period
	for rows.Next() {
		var p Period
		if err := rows.Scan(&p.SemesterID, &p.TahunAjaranID, &p.TahunAjaranLabel, &p.NamaSemester, &p.Semester, &p.PeriodeAktif); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		result = append(result, p)
	}
	if result == nil {
		result = []Period{}
	}
	return result, nil
}

func (s *PreferenceService) GetUserPreference(ctx context.Context, userID, schoolID string) (string, error) {
	var semesterID string
	query := `SELECT semester_id FROM user_period_preference WHERE user_id = $1 AND school_id = $2`
	err := s.db.QueryRowContext(ctx, query, userID, schoolID).Scan(&semesterID)

	if err == sql.ErrNoRows {
		queryActive := `SELECT id FROM semesters WHERE is_active = true LIMIT 1`
		if errActive := s.db.QueryRowContext(ctx, queryActive).Scan(&semesterID); errActive != nil {
			if errActive == sql.ErrNoRows {
				return "", fmt.Errorf("no active semester found")
			}
			return "", fmt.Errorf("get active semester: %w", errActive)
		}
		return semesterID, nil
	}
	if err != nil {
		return "", fmt.Errorf("get user preference: %w", err)
	}
	return semesterID, nil
}

func (s *PreferenceService) SaveUserPreference(ctx context.Context, userID, schoolID, semesterID string) error {
	query := `
		INSERT INTO user_period_preference (user_id, school_id, semester_id, updated_at)
		VALUES ($1, $2, $3, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, school_id) DO UPDATE SET
			semester_id = EXCLUDED.semester_id,
			updated_at = CURRENT_TIMESTAMP
	`
	_, err := s.db.ExecContext(ctx, query, userID, schoolID, semesterID)
	return err
}

func (s *PreferenceService) CreateAcademicYear(ctx context.Context, year string) error {
	if len(year) != 4 {
		return fmt.Errorf("format tahun harus 4 digit, contoh: 2027")
	}

	y, err := strconv.Atoi(year)
	if err != nil {
		return fmt.Errorf("tahun tidak valid")
	}

	nextYear := strconv.Itoa(y + 1)
	academicYearLabel := fmt.Sprintf("%s/%s", year, nextYear)

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Buat Semester 1
	sem1ID := fmt.Sprintf("%s1", year)
	querySem1 := `INSERT INTO semesters (id, name, academic_year, semester, is_active) VALUES ($1, $2, $3, 1, FALSE) ON CONFLICT (id) DO NOTHING`
	if _, err := tx.ExecContext(ctx, querySem1, sem1ID, "Semester 1", academicYearLabel); err != nil {
		return fmt.Errorf("gagal membuat semester 1: %w", err)
	}

	// 2. Buat Semester 2
	sem2ID := fmt.Sprintf("%s2", year)
	querySem2 := `INSERT INTO semesters (id, name, academic_year, semester, is_active) VALUES ($1, $2, $3, 2, FALSE) ON CONFLICT (id) DO NOTHING`
	if _, err := tx.ExecContext(ctx, querySem2, sem2ID, "Semester 2", academicYearLabel); err != nil {
		return fmt.Errorf("gagal membuat semester 2: %w", err)
	}

	return tx.Commit()
}

func (s *PreferenceService) ActivatePeriod(ctx context.Context, semesterID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `UPDATE semesters SET is_active = FALSE`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `UPDATE semesters SET is_active = TRUE WHERE id = $1`, semesterID); err != nil {
		return err
	}

	return tx.Commit()
}

// DeletePeriod menghapus periode, dengan validasi tidak boleh menghapus yang aktif
func (s *PreferenceService) DeletePeriod(ctx context.Context, semesterID string) error {
	var isActive bool
	queryCheck := `SELECT is_active FROM semesters WHERE id = $1`
	err := s.db.QueryRowContext(ctx, queryCheck, semesterID).Scan(&isActive)

	if err == sql.ErrNoRows {
		return fmt.Errorf("periode tidak ditemukan")
	}
	if err != nil {
		return fmt.Errorf("gagal memeriksa periode: %w", err)
	}

	if isActive {
		return fmt.Errorf("tidak dapat menghapus semester yang sedang aktif")
	}

	queryDelete := `DELETE FROM semesters WHERE id = $1`
	_, err = s.db.ExecContext(ctx, queryDelete, semesterID)
	if err != nil {
		return fmt.Errorf("gagal menghapus periode: %w", err)
	}

	return nil
}
