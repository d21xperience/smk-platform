package projection

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// TeacherProjection menangani projection untuk domain Teacher.
type TeacherProjection struct {
	db *pgxpool.Pool
}

// NewTeacherProjection membuat instance baru.
func NewTeacherProjection(db *pgxpool.Pool) *TeacherProjection {
	return &TeacherProjection{db: db}
}

// Name mengembalikan nama projection.
func (p *TeacherProjection) Name() string {
	return "teacher_projection"
}

// SupportedEvents mengembalikan daftar event yang didukung.
func (p *TeacherProjection) SupportedEvents() []string {
	return []string{
		"TeacherAssigned",
		"TeacherUnassigned",
	}
}

// HandleEvent memproses event teacher.
func (p *TeacherProjection) HandleEvent(ctx context.Context, event DomainEvent) error {
	switch event.EventType {
	case "TeacherAssigned":
		return p.handleTeacherAssigned(ctx, event)
	case "TeacherUnassigned":
		return p.handleTeacherUnassigned(ctx, event)
	default:
		return fmt.Errorf("event type '%s' tidak didukung", event.EventType)
	}
}

// handleTeacherAssigned mengupdate read model saat guru ditugaskan.
func (p *TeacherProjection) handleTeacherAssigned(ctx context.Context, event DomainEvent) error {
	teacherID := event.Payload["teacherId"].(string)
	classID := event.Payload["classId"].(string)
	subjectID := event.Payload["subjectId"].(string)
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		INSERT INTO projection_teacher_assignment
			(school_id, academic_period_id, teacher_id, class_id, subject_id, assigned_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		ON CONFLICT (school_id, academic_period_id, teacher_id, class_id, subject_id) DO UPDATE SET
			assigned_at = NOW(),
			status = 'ACTIVE'
	`

	_, err := p.db.Exec(ctx, query,
		event.SchoolID,
		periodID,
		teacherID,
		classID,
		subjectID,
	)
	if err != nil {
		return fmt.Errorf("gagal upsert teacher assignment: %w", err)
	}

	log.Printf("[TeacherProjection] Teacher %s assigned ke class %s, subject %s (CorrelationID: %s)",
		teacherID, classID, subjectID, event.CorrelationID)
	return nil
}

// handleTeacherUnassigned menangani pencabutan tugas guru.
func (p *TeacherProjection) handleTeacherUnassigned(ctx context.Context, event DomainEvent) error {
	teacherID := event.Payload["teacherId"].(string)
	classID := event.Payload["classId"].(string)
	subjectID := event.Payload["subjectId"].(string)
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		UPDATE projection_teacher_assignment
		SET status = 'INACTIVE', unassigned_at = NOW()
		WHERE school_id = $1 AND academic_period_id = $2
		  AND teacher_id = $3 AND class_id = $4 AND subject_id = $5
	`

	_, err := p.db.Exec(ctx, query,
		event.SchoolID,
		periodID,
		teacherID,
		classID,
		subjectID,
	)
	if err != nil {
		return fmt.Errorf("gagal update teacher unassigned: %w", err)
	}

	log.Printf("[TeacherProjection] Teacher %s unassigned dari class %s (CorrelationID: %s)",
		teacherID, classID, event.CorrelationID)
	return nil
}
