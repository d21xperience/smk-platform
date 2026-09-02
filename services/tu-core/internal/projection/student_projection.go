package projection

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// StudentProjection menangani projection untuk domain Student.
type StudentProjection struct {
	db *pgxpool.Pool
}

// NewStudentProjection membuat instance baru.
func NewStudentProjection(db *pgxpool.Pool) *StudentProjection {
	return &StudentProjection{db: db}
}

// Name mengembalikan nama projection.
func (p *StudentProjection) Name() string {
	return "student_projection"
}

// SupportedEvents mengembalikan daftar event yang didukung.
func (p *StudentProjection) SupportedEvents() []string {
	return []string{
		"StudentEnrolled",
		"StudentGraduated",
		"StudentTransferred",
	}
}

// HandleEvent memproses event student.
func (p *StudentProjection) HandleEvent(ctx context.Context, event DomainEvent) error {
	switch event.EventType {
	case "StudentEnrolled":
		return p.handleStudentEnrolled(ctx, event)
	case "StudentGraduated":
		return p.handleStudentGraduated(ctx, event)
	case "StudentTransferred":
		return p.handleStudentTransferred(ctx, event)
	default:
		return fmt.Errorf("event type '%s' tidak didukung", event.EventType)
	}
}

// handleStudentEnrolled mengupdate read model saat siswa terdaftar.
func (p *StudentProjection) handleStudentEnrolled(ctx context.Context, event DomainEvent) error {
	studentID := event.Payload["studentId"].(string)
	classID := event.Payload["classId"].(string)
	periodID := event.Payload["academicPeriodId"].(string)

	// UPSERT ke read model (misalnya: dashboard_student_summary)
	query := `
		INSERT INTO projection_student_summary
			(school_id, academic_period_id, class_id, student_id, status, enrolled_at)
		VALUES ($1, $2, $3, $4, 'ACTIVE', NOW())
		ON CONFLICT (school_id, academic_period_id, student_id) DO UPDATE SET
			class_id = EXCLUDED.class_id,
			status = 'ACTIVE',
			updated_at = NOW()
	`

	_, err := p.db.Exec(ctx, query,
		event.SchoolID,
		periodID,
		classID,
		studentID,
	)
	if err != nil {
		return fmt.Errorf("gagal upsert student summary: %w", err)
	}

	log.Printf("[StudentProjection] Student %s enrolled di class %s (CorrelationID: %s)",
		studentID, classID, event.CorrelationID)
	return nil
}

// handleStudentGraduated mengupdate status siswa menjadi lulus.
func (p *StudentProjection) handleStudentGraduated(ctx context.Context, event DomainEvent) error {
	studentID := event.Payload["studentId"].(string)
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		UPDATE projection_student_summary
		SET status = 'GRADUATED', graduated_at = NOW(), updated_at = NOW()
		WHERE school_id = $1 AND academic_period_id = $2 AND student_id = $3
	`

	_, err := p.db.Exec(ctx, query, event.SchoolID, periodID, studentID)
	if err != nil {
		return fmt.Errorf("gagal update student graduated: %w", err)
	}

	log.Printf("[StudentProjection] Student %s graduated (CorrelationID: %s)",
		studentID, event.CorrelationID)
	return nil
}

// handleStudentTransferred menangani mutasi siswa.
func (p *StudentProjection) handleStudentTransferred(ctx context.Context, event DomainEvent) error {
	studentID := event.Payload["studentId"].(string)
	newClassID := event.Payload["newClassId"].(string)
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		UPDATE projection_student_summary
		SET class_id = $1, status = 'TRANSFERRED', updated_at = NOW()
		WHERE school_id = $2 AND academic_period_id = $3 AND student_id = $4
	`

	_, err := p.db.Exec(ctx, query, newClassID, event.SchoolID, periodID, studentID)
	if err != nil {
		return fmt.Errorf("gagal update student transferred: %w", err)
	}

	log.Printf("[StudentProjection] Student %s transferred ke class %s (CorrelationID: %s)",
		studentID, newClassID, event.CorrelationID)
	return nil
}
