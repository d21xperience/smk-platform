package projections

import (
	"context"
	"encoding/json"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// AttendanceProjectionHandler mengupdate read model laporan kehadiran harian.
type AttendanceProjectionHandler struct {
	db *pgxpool.Pool
}

func NewAttendanceProjectionHandler(db *pgxpool.Pool) *AttendanceProjectionHandler {
	return &AttendanceProjectionHandler{db: db}
}

// HandleAttendanceSubmitted mengagregasi data saat event kehadiran masuk.
func (h *AttendanceProjectionHandler) HandleAttendanceSubmitted(ctx context.Context, msg []byte) error {
	var event struct {
		SchoolID         string `json:"schoolId"`
		AcademicPeriodID string `json:"academicPeriodId"`
		ClassID          string `json:"classId"`
		Date             string `json:"date"` // Format YYYY-MM-DD
		TotalStudents    int    `json:"totalStudents"`
		PresentCount     int    `json:"presentCount"`
		AbsentCount      int    `json:"absentCount"`
	}

	if err := json.Unmarshal(msg, &event); err != nil {
		return err
	}

	// Hitung rate kehadiran
	var rate float64
	if event.TotalStudents > 0 {
		rate = float64(event.PresentCount) / float64(event.TotalStudents) * 100
	}

	// UPSERT ke tabel reporting (Idempotent)
	query := `
		INSERT INTO reporting_attendance_daily
			(school_id, academic_period_id, class_id, date, total_students, present_count, absent_count, attendance_rate)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (school_id, academic_period_id, class_id, date) DO UPDATE SET
			total_students = EXCLUDED.total_students,
			present_count = EXCLUDED.present_count,
			absent_count = EXCLUDED.absent_count,
			attendance_rate = EXCLUDED.attendance_rate
	`
	_, err := h.db.Exec(ctx, query,
		event.SchoolID, event.AcademicPeriodID, event.ClassID, event.Date,
		event.TotalStudents, event.PresentCount, event.AbsentCount, rate,
	)

	if err != nil {
		log.Printf("[Reporting Projection] Gagal update attendance report: %v", err)
		return err
	}

	return nil
}
