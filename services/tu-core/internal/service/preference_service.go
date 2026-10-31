package service

import (
	"context"
	"database/sql"
	"fmt"
)

type DBExecutor interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
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
			s.semester_id,
			SUBSTRING(ay.year FROM 1 FOR 4) AS tahun_ajaran_id,
			ay.year AS tahun_ajaran_label,
			s.nama_semester,
			s.semester,
			s.is_active
		FROM semester s
		INNER JOIN academic_years ay ON ay.id = s.academic_year_id
		ORDER BY s.semester_id ASC
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list periods: %w", err)
	}
	defer rows.Close()

	var result []Period
	for rows.Next() {
		var p Period
		if err := rows.Scan(
			&p.SemesterID,
			&p.TahunAjaranID,
			&p.TahunAjaranLabel,
			&p.NamaSemester,
			&p.Semester,
			&p.PeriodeAktif,
		); err != nil {
			return nil, fmt.Errorf("scan period: %w", err)
		}
		result = append(result, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate periods: %w", err)
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
		queryActive := `SELECT semester_id FROM semester WHERE is_active = true LIMIT 1`
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
	if err != nil {
		return fmt.Errorf("save user preference: %w", err)
	}
	return nil
}
