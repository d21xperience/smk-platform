package repository

import (
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/domain/assessment/models"
	"sekolah-platform/services/tu-core/internal/domain/assessment/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrAssessmentNotFound = errors.New("assessment not found")

type AssessmentFilter struct {
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	SubjectID        string
	AssessmentType   string
	Semester         int
	Limit            int
	Offset           int
}

type AssessmentRepository interface {
	Create(ctx context.Context, a *models.Assessment) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, assessmentID string) (*models.Assessment, error)
	List(ctx context.Context, filter AssessmentFilter) ([]models.Assessment, int, error)
	Update(ctx context.Context, a *models.Assessment) error
	GetByStudent(ctx context.Context, schoolID, academicPeriodID, studentID string, semester int, subjectID string) ([]models.Assessment, error)
	GetBySubject(ctx context.Context, schoolID, academicPeriodID, subjectID, classID string, semester int) ([]models.Assessment, error)
}

type pgxAssessmentRepository struct {
	db *pgxpool.Pool
}

func NewAssessmentRepository(db *pgxpool.Pool) AssessmentRepository {
	return &pgxAssessmentRepository{db: db}
}

func (r *pgxAssessmentRepository) Create(ctx context.Context, a *models.Assessment) error {
	query := `
		INSERT INTO assessments (
			id, school_id, academic_period_id, student_id, subject_id, assessment_type,
			score, max_score, assessment_date, semester, notes, created_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.Exec(ctx, query,
		a.ID, a.SchoolID, a.AcademicPeriodID, a.StudentID, a.SubjectID, a.AssessmentType.String(),
		a.Score, a.MaxScore, a.AssessmentDate, a.Semester, a.Notes, a.CreatedBy, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxAssessmentRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, assessmentID string) (*models.Assessment, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, subject_id, assessment_type,
		       score, max_score, assessment_date, semester, COALESCE(notes, '') as notes,
		       COALESCE(created_by, '') as created_by, created_at, updated_at
		FROM assessments
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var a models.Assessment
	var assessmentTypeStr string

	err := r.db.QueryRow(ctx, query, assessmentID, schoolID, academicPeriodID).Scan(
		&a.ID, &a.SchoolID, &a.AcademicPeriodID, &a.StudentID, &a.SubjectID, &assessmentTypeStr,
		&a.Score, &a.MaxScore, &a.AssessmentDate, &a.Semester, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrAssessmentNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	a.AssessmentType, _ = valueobjects.NewAssessmentType(assessmentTypeStr)

	return &a, nil
}

func (r *pgxAssessmentRepository) List(ctx context.Context, filter AssessmentFilter) ([]models.Assessment, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.StudentID != "" {
		whereClause += fmt.Sprintf(" AND student_id = $%d", argIndex)
		args = append(args, filter.StudentID)
		argIndex++
	}

	if filter.SubjectID != "" {
		whereClause += fmt.Sprintf(" AND subject_id = $%d", argIndex)
		args = append(args, filter.SubjectID)
		argIndex++
	}

	if filter.AssessmentType != "" {
		whereClause += fmt.Sprintf(" AND assessment_type = $%d", argIndex)
		args = append(args, filter.AssessmentType)
		argIndex++
	}

	if filter.Semester > 0 {
		whereClause += fmt.Sprintf(" AND semester = $%d", argIndex)
		args = append(args, filter.Semester)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM assessments %s", whereClause)
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
		SELECT id, school_id, academic_period_id, student_id, subject_id, assessment_type,
		       score, max_score, assessment_date, semester, COALESCE(notes, '') as notes,
		       COALESCE(created_by, '') as created_by, created_at, updated_at
		FROM assessments %s
		ORDER BY assessment_date DESC, created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		var assessmentTypeStr string

		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.AcademicPeriodID, &a.StudentID, &a.SubjectID, &assessmentTypeStr,
			&a.Score, &a.MaxScore, &a.AssessmentDate, &a.Semester, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		a.AssessmentType, _ = valueobjects.NewAssessmentType(assessmentTypeStr)
		assessments = append(assessments, a)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return assessments, total, nil
}

func (r *pgxAssessmentRepository) Update(ctx context.Context, a *models.Assessment) error {
	query := `
		UPDATE assessments
		SET score = $1, max_score = $2, notes = $3, updated_at = $4
		WHERE id = $5 AND school_id = $6 AND academic_period_id = $7
	`
	result, err := r.db.Exec(ctx, query,
		a.Score, a.MaxScore, a.Notes, a.UpdatedAt,
		a.ID, a.SchoolID, a.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrAssessmentNotFound
	}
	return nil
}

func (r *pgxAssessmentRepository) GetByStudent(ctx context.Context, schoolID, academicPeriodID, studentID string, semester int, subjectID string) ([]models.Assessment, error) {
	query := `
		SELECT id, school_id, academic_period_id, student_id, subject_id, assessment_type,
		       score, max_score, assessment_date, semester, COALESCE(notes, '') as notes,
		       COALESCE(created_by, '') as created_by, created_at, updated_at
		FROM assessments
		WHERE school_id = $1 AND academic_period_id = $2 AND student_id = $3
	`
	args := []interface{}{schoolID, academicPeriodID, studentID}
	argIndex := 4

	if semester > 0 {
		query += fmt.Sprintf(" AND semester = $%d", argIndex)
		args = append(args, semester)
		argIndex++
	}

	if subjectID != "" {
		query += fmt.Sprintf(" AND subject_id = $%d", argIndex)
		args = append(args, subjectID)
		argIndex++
	}

	query += " ORDER BY assessment_date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		var assessmentTypeStr string

		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.AcademicPeriodID, &a.StudentID, &a.SubjectID, &assessmentTypeStr,
			&a.Score, &a.MaxScore, &a.AssessmentDate, &a.Semester, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		a.AssessmentType, _ = valueobjects.NewAssessmentType(assessmentTypeStr)
		assessments = append(assessments, a)
	}

	return assessments, rows.Err()
}

func (r *pgxAssessmentRepository) GetBySubject(ctx context.Context, schoolID, academicPeriodID, subjectID, classID string, semester int) ([]models.Assessment, error) {
	query := `
		SELECT a.id, a.school_id, a.academic_period_id, a.student_id, a.subject_id, a.assessment_type,
		       a.score, a.max_score, a.assessment_date, a.semester, COALESCE(a.notes, '') as notes,
		       COALESCE(a.created_by, '') as created_by, a.created_at, a.updated_at
		FROM assessments a
		JOIN students s ON a.student_id = s.id
		WHERE a.school_id = $1 AND a.academic_period_id = $2 AND a.subject_id = $3
	`
	args := []interface{}{schoolID, academicPeriodID, subjectID}
	argIndex := 4

	if classID != "" {
		query += fmt.Sprintf(" AND s.class_name = $%d", argIndex)
		args = append(args, classID)
		argIndex++
	}

	if semester > 0 {
		query += fmt.Sprintf(" AND a.semester = $%d", argIndex)
		args = append(args, semester)
		argIndex++
	}

	query += " ORDER BY s.name ASC, a.assessment_date DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var assessments []models.Assessment
	for rows.Next() {
		var a models.Assessment
		var assessmentTypeStr string

		if err := rows.Scan(
			&a.ID, &a.SchoolID, &a.AcademicPeriodID, &a.StudentID, &a.SubjectID, &assessmentTypeStr,
			&a.Score, &a.MaxScore, &a.AssessmentDate, &a.Semester, &a.Notes, &a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		a.AssessmentType, _ = valueobjects.NewAssessmentType(assessmentTypeStr)
		assessments = append(assessments, a)
	}

	return assessments, rows.Err()
}
