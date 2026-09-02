package projection

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DashboardProjection menangani projection untuk dashboard.
type DashboardProjection struct {
	db *pgxpool.Pool
}

// NewDashboardProjection membuat instance baru.
func NewDashboardProjection(db *pgxpool.Pool) *DashboardProjection {
	return &DashboardProjection{db: db}
}

// Name mengembalikan nama projection.
func (p *DashboardProjection) Name() string {
	return "dashboard_projection"
}

// SupportedEvents mengembalikan daftar event yang didukung.
func (p *DashboardProjection) SupportedEvents() []string {
	return []string{
		"StudentEnrolled",
		"TeacherAssigned",
		"AttendanceSubmitted",
	}
}

// HandleEvent memproses event untuk dashboard.
func (p *DashboardProjection) HandleEvent(ctx context.Context, event DomainEvent) error {
	switch event.EventType {
	case "StudentEnrolled":
		return p.updateStudentCount(ctx, event, 1)
	case "TeacherAssigned":
		return p.updateTeacherCount(ctx, event, 1)
	case "AttendanceSubmitted":
		return p.updateAttendanceStats(ctx, event)
	default:
		return fmt.Errorf("event type '%s' tidak didukung", event.EventType)
	}
}

// updateStudentCount mengupdate counter siswa di dashboard.
func (p *DashboardProjection) updateStudentCount(ctx context.Context, event DomainEvent, delta int) error {
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		INSERT INTO projection_dashboard_stats
			(school_id, academic_period_id, stat_type, stat_value, updated_at)
		VALUES ($1, $2, 'total_students', $3, NOW())
		ON CONFLICT (school_id, academic_period_id, stat_type) DO UPDATE SET
			stat_value = projection_dashboard_stats.stat_value + $3,
			updated_at = NOW()
	`

	_, err := p.db.Exec(ctx, query, event.SchoolID, periodID, delta)
	if err != nil {
		return fmt.Errorf("gagal update student count: %w", err)
	}

	log.Printf("[DashboardProjection] Updated student count by %d (CorrelationID: %s)",
		delta, event.CorrelationID)
	return nil
}

// updateTeacherCount mengupdate counter guru di dashboard.
func (p *DashboardProjection) updateTeacherCount(ctx context.Context, event DomainEvent, delta int) error {
	periodID := event.Payload["academicPeriodId"].(string)

	query := `
		INSERT INTO projection_dashboard_stats
			(school_id, academic_period_id, stat_type, stat_value, updated_at)
		VALUES ($1, $2, 'total_teachers', $3, NOW())
		ON CONFLICT (school_id, academic_period_id, stat_type) DO UPDATE SET
			stat_value = projection_dashboard_stats.stat_value + $3,
			updated_at = NOW()
	`

	_, err := p.db.Exec(ctx, query, event.SchoolID, periodID, delta)
	if err != nil {
		return fmt.Errorf("gagal update teacher count: %w", err)
	}

	log.Printf("[DashboardProjection] Updated teacher count by %d (CorrelationID: %s)",
		delta, event.CorrelationID)
	return nil
}

// updateAttendanceStats mengupdate statistik kehadiran di dashboard.
func (p *DashboardProjection) updateAttendanceStats(ctx context.Context, event DomainEvent) error {
	periodID := event.Payload["academicPeriodId"].(string)
	totalStudents := int(event.Payload["totalStudents"].(float64))
	presentCount := int(event.Payload["presentCount"].(float64))

	var rate float64
	if totalStudents > 0 {
		rate = float64(presentCount) / float64(totalStudents) * 100
	}

	query := `
		INSERT INTO projection_dashboard_stats
			(school_id, academic_period_id, stat_type, stat_value, updated_at)
		VALUES ($1, $2, 'avg_attendance_rate', $3, NOW())
		ON CONFLICT (school_id, academic_period_id, stat_type) DO UPDATE SET
			stat_value = (projection_dashboard_stats.stat_value + $3) / 2,
			updated_at = NOW()
	`

	_, err := p.db.Exec(ctx, query, event.SchoolID, periodID, rate)
	if err != nil {
		return fmt.Errorf("gagal update attendance stats: %w", err)
	}

	log.Printf("[DashboardProjection] Updated attendance rate: %.2f%% (CorrelationID: %s)",
		rate, event.CorrelationID)
	return nil
}
