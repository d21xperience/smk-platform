// Package projection menyediakan Dashboard Projection Manager.
package projection

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/dashboard/models"
)

// DashboardProjection mengelola aggregated statistics untuk dashboard.
type DashboardProjection struct {
	db       *gorm.DB
	handlers map[string]EventHandler
}

// EventHandler adalah handler untuk event tertentu.
type EventHandler func(ctx context.Context, event types.DomainEvent) error

// NewDashboardProjection membuat DashboardProjection baru.
func NewDashboardProjection(db *gorm.DB) *DashboardProjection {
	dp := &DashboardProjection{
		db:       db,
		handlers: make(map[string]EventHandler),
	}

	// Register handlers
	dp.registerAttendanceHandlers()
	dp.registerAssessmentHandlers()
	dp.registerStudentHandlers()
	dp.registerTeacherHandlers()

	return dp
}

// Handle mendistribusikan event ke handler yang sesuai.
func (dp *DashboardProjection) Handle(ctx context.Context, event types.DomainEvent) error {
	eventName := event.GetEventName()

	handler, ok := dp.handlers[eventName]
	if !ok {
		// Event tidak punya handler, skip
		return nil
	}

	return handler(ctx, event)
}

// RegisterHandler mendaftarkan handler untuk event tertentu.
func (dp *DashboardProjection) RegisterHandler(eventName string, handler EventHandler) {
	dp.handlers[eventName] = handler
}

// === Attendance Handlers ===

func (dp *DashboardProjection) registerAttendanceHandlers() {
	dp.handlers["AttendanceSubmitted"] = dp.handleAttendanceSubmitted
	dp.handlers["AttendanceCorrected"] = dp.handleAttendanceCorrected
}

func (dp *DashboardProjection) handleAttendanceSubmitted(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})
	schoolID := event.GetSchoolID()
	periodID := event.GetAcademicPeriodID()
	classID := getStr(payload, "classId")
	teacherID := getStr(payload, "teacherId")

	stats := getMap(payload, "statistics")
	total := getInt(stats, "total")
	present := getInt(stats, "present")
	absent := getInt(stats, "absent")
	sick := getInt(stats, "sick")
	leave := getInt(stats, "leave")
	attendanceRate := getFloat(stats, "attendanceRate")

	// Update school stats
	return dp.db.Transaction(func(tx *gorm.DB) error {
		// Upsert school stats
		var schoolStats models.SchoolStats
		result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
		if result.Error == gorm.ErrRecordNotFound {
			schoolStats = models.SchoolStats{
				ID:        fmt.Sprintf("%s:%s", schoolID, periodID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				CreatedAt: time.Now().UTC(),
			}
		}

		schoolStats.TotalAttendanceRecords += total
		schoolStats.TotalPresent += present
		schoolStats.TotalAbsent += absent
		schoolStats.TotalSick += sick
		schoolStats.TotalLeave += leave

		// Recalculate attendance rate
		if schoolStats.TotalAttendanceRecords > 0 {
			schoolStats.AttendanceRate = float64(schoolStats.TotalPresent) / float64(schoolStats.TotalAttendanceRecords) * 100
		}
		schoolStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&schoolStats).Error; err != nil {
			return err
		}

		// Upsert class stats
		var classStats models.ClassStats
		result = tx.Where("school_id = ? AND period_id = ? AND class_id = ?", schoolID, periodID, classID).First(&classStats)
		if result.Error == gorm.ErrRecordNotFound {
			classStats = models.ClassStats{
				ID:        fmt.Sprintf("%s:%s:%s", schoolID, periodID, classID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				ClassID:   classID,
				CreatedAt: time.Now().UTC(),
			}
		}

		classStats.TotalAttendanceRecords += total
		classStats.TotalPresent += present
		if classStats.TotalAttendanceRecords > 0 {
			classStats.AttendanceRate = float64(classStats.TotalPresent) / float64(classStats.TotalAttendanceRecords) * 100
		}
		classStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&classStats).Error; err != nil {
			return err
		}

		// Upsert teacher stats
		var teacherStats models.TeacherStats
		result = tx.Where("school_id = ? AND period_id = ? AND teacher_id = ?", schoolID, periodID, teacherID).First(&teacherStats)
		if result.Error == gorm.ErrRecordNotFound {
			teacherStats = models.TeacherStats{
				ID:        fmt.Sprintf("%s:%s:%s", schoolID, periodID, teacherID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				TeacherID: teacherID,
				CreatedAt: time.Now().UTC(),
			}
		}

		teacherStats.AttendanceSubmissions++
		teacherStats.TotalStudentsAttended += total
		teacherStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&teacherStats).Error; err != nil {
			return err
		}

		_ = attendanceRate // Used for logging if needed

		return nil
	})
}

func (dp *DashboardProjection) handleAttendanceCorrected(ctx context.Context, event types.DomainEvent) error {
	// For simplicity, we don't recalculate stats on correction
	// In production, you might want to recalculate
	return nil
}

// === Assessment Handlers ===

func (dp *DashboardProjection) registerAssessmentHandlers() {
	dp.handlers["AssessmentSubmitted"] = dp.handleAssessmentSubmitted
	dp.handlers["ScoreUpdated"] = dp.handleScoreUpdated
}

func (dp *DashboardProjection) handleAssessmentSubmitted(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})
	schoolID := event.GetSchoolID()
	periodID := event.GetAcademicPeriodID()
	classID := getStr(payload, "classId")
	teacherID := getStr(payload, "teacherId")

	stats := getMap(payload, "statistics")
	total := getInt(stats, "total")
	// passing := getInt(stats, "passing")
	average := getFloat(stats, "average")
	passRate := getFloat(stats, "passRate")

	return dp.db.Transaction(func(tx *gorm.DB) error {
		// Update school stats
		var schoolStats models.SchoolStats
		result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
		if result.Error == gorm.ErrRecordNotFound {
			schoolStats = models.SchoolStats{
				ID:        fmt.Sprintf("%s:%s", schoolID, periodID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				CreatedAt: time.Now().UTC(),
			}
		}

		schoolStats.TotalAssessments++
		schoolStats.TotalScores += total

		// Recalculate average (weighted)
		if schoolStats.TotalScores > 0 {
			oldTotal := schoolStats.TotalScores - total
			oldSum := schoolStats.AverageScore * float64(oldTotal)
			newSum := average * float64(total)
			schoolStats.AverageScore = (oldSum + newSum) / float64(schoolStats.TotalScores)
		}

		// Recalculate pass rate
		if schoolStats.TotalScores > 0 {
			schoolStats.PassRate = passRate // Simplified: use latest pass rate
		}
		schoolStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&schoolStats).Error; err != nil {
			return err
		}

		// Update class stats
		var classStats models.ClassStats
		result = tx.Where("school_id = ? AND period_id = ? AND class_id = ?", schoolID, periodID, classID).First(&classStats)
		if result.Error == gorm.ErrRecordNotFound {
			classStats = models.ClassStats{
				ID:        fmt.Sprintf("%s:%s:%s", schoolID, periodID, classID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				ClassID:   classID,
				CreatedAt: time.Now().UTC(),
			}
		}

		classStats.TotalAssessments++
		classStats.AverageScore = average
		classStats.PassRate = passRate
		classStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&classStats).Error; err != nil {
			return err
		}

		// Update teacher stats
		var teacherStats models.TeacherStats
		result = tx.Where("school_id = ? AND period_id = ? AND teacher_id = ?", schoolID, periodID, teacherID).First(&teacherStats)
		if result.Error == gorm.ErrRecordNotFound {
			teacherStats = models.TeacherStats{
				ID:        fmt.Sprintf("%s:%s:%s", schoolID, periodID, teacherID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				TeacherID: teacherID,
				CreatedAt: time.Now().UTC(),
			}
		}

		teacherStats.AssessmentsCompleted++
		teacherStats.TotalScoresSubmitted += total
		if teacherStats.TotalScoresSubmitted > 0 {
			oldTotal := teacherStats.TotalScoresSubmitted - total
			oldSum := teacherStats.AverageStudentScore * float64(oldTotal)
			newSum := average * float64(total)
			teacherStats.AverageStudentScore = (oldSum + newSum) / float64(teacherStats.TotalScoresSubmitted)
		}
		teacherStats.LastUpdatedAt = time.Now().UTC()

		if err := tx.Save(&teacherStats).Error; err != nil {
			return err
		}

		return nil
	})
}

func (dp *DashboardProjection) handleScoreUpdated(ctx context.Context, event types.DomainEvent) error {
	// For simplicity, we don't recalculate stats on score update
	// In production, you might want to recalculate affected stats
	return nil
}

// === Student Handlers ===

func (dp *DashboardProjection) registerStudentHandlers() {
	dp.handlers["StudentCreated"] = dp.handleStudentCreated
	dp.handlers["StudentGraduated"] = dp.handleStudentStatusChange("GRADUATED")
	dp.handlers["StudentTransferred"] = dp.handleStudentStatusChange("TRANSFERRED")
}

func (dp *DashboardProjection) handleStudentCreated(ctx context.Context, event types.DomainEvent) error {
	schoolID := event.GetSchoolID()
	periodID := event.GetAcademicPeriodID()

	return dp.db.Transaction(func(tx *gorm.DB) error {
		var schoolStats models.SchoolStats
		result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
		if result.Error == gorm.ErrRecordNotFound {
			schoolStats = models.SchoolStats{
				ID:        fmt.Sprintf("%s:%s", schoolID, periodID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				CreatedAt: time.Now().UTC(),
			}
		}

		schoolStats.TotalStudents++
		schoolStats.ActiveStudents++
		schoolStats.LastUpdatedAt = time.Now().UTC()

		return tx.Save(&schoolStats).Error
	})
}

func (dp *DashboardProjection) handleStudentStatusChange(status string) EventHandler {
	return func(ctx context.Context, event types.DomainEvent) error {
		schoolID := event.GetSchoolID()
		periodID := event.GetAcademicPeriodID()

		return dp.db.Transaction(func(tx *gorm.DB) error {
			var schoolStats models.SchoolStats
			result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
			if result.Error != nil {
				return nil // Stats belum ada, skip
			}

			schoolStats.ActiveStudents--
			switch status {
			case "GRADUATED":
				schoolStats.GraduatedStudents++
			case "TRANSFERRED":
				schoolStats.TransferredStudents++
			}
			schoolStats.LastUpdatedAt = time.Now().UTC()

			return tx.Save(&schoolStats).Error
		})
	}
}

// === Teacher Handlers ===

func (dp *DashboardProjection) registerTeacherHandlers() {
	dp.handlers["TeacherCreated"] = dp.handleTeacherCreated
	dp.handlers["TeacherDeactivated"] = dp.handleTeacherDeactivated
}

func (dp *DashboardProjection) handleTeacherCreated(ctx context.Context, event types.DomainEvent) error {
	schoolID := event.GetSchoolID()
	periodID := event.GetAcademicPeriodID()

	return dp.db.Transaction(func(tx *gorm.DB) error {
		var schoolStats models.SchoolStats
		result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
		if result.Error == gorm.ErrRecordNotFound {
			schoolStats = models.SchoolStats{
				ID:        fmt.Sprintf("%s:%s", schoolID, periodID),
				SchoolID:  schoolID,
				PeriodID:  periodID,
				CreatedAt: time.Now().UTC(),
			}
		}

		schoolStats.TotalTeachers++
		schoolStats.ActiveTeachers++
		schoolStats.LastUpdatedAt = time.Now().UTC()

		return tx.Save(&schoolStats).Error
	})
}

func (dp *DashboardProjection) handleTeacherDeactivated(ctx context.Context, event types.DomainEvent) error {
	schoolID := event.GetSchoolID()
	periodID := event.GetAcademicPeriodID()

	return dp.db.Transaction(func(tx *gorm.DB) error {
		var schoolStats models.SchoolStats
		result := tx.Where("school_id = ? AND period_id = ?", schoolID, periodID).First(&schoolStats)
		if result.Error != nil {
			return nil
		}

		schoolStats.ActiveTeachers--
		schoolStats.LastUpdatedAt = time.Now().UTC()

		return tx.Save(&schoolStats).Error
	})
}

// === Helpers ===

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return int(n)
		case int:
			return n
		}
	}
	return 0
}

func getFloat(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch n := v.(type) {
		case float64:
			return n
		case int:
			return float64(n)
		}
	}
	return 0
}

func getMap(m map[string]interface{}, key string) map[string]interface{} {
	if v, ok := m[key]; ok {
		if subMap, ok := v.(map[string]interface{}); ok {
			return subMap
		}
	}
	return map[string]interface{}{}
}

// === Query Methods ===

// GetSchoolStats mengambil stats sekolah untuk periode tertentu.
func (dp *DashboardProjection) GetSchoolStats(ctx context.Context, schoolID, periodID string) (*models.SchoolStats, error) {
	var stats models.SchoolStats
	err := dp.db.WithContext(ctx).
		Where("school_id = ? AND period_id = ?", schoolID, periodID).
		First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &stats, nil
}

// GetClassStats mengambil stats kelas untuk periode tertentu.
func (dp *DashboardProjection) GetClassStats(ctx context.Context, schoolID, periodID, classID string) (*models.ClassStats, error) {
	var stats models.ClassStats
	err := dp.db.WithContext(ctx).
		Where("school_id = ? AND period_id = ? AND class_id = ?", schoolID, periodID, classID).
		First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &stats, nil
}

// ListClassStatsByPeriod mengambil semua stats kelas untuk periode tertentu.
func (dp *DashboardProjection) ListClassStatsByPeriod(ctx context.Context, schoolID, periodID string) ([]models.ClassStats, error) {
	var stats []models.ClassStats
	err := dp.db.WithContext(ctx).
		Where("school_id = ? AND period_id = ?", schoolID, periodID).
		Order("class_id ASC").
		Find(&stats).Error
	return stats, err
}

// GetTeacherStats mengambil stats guru untuk periode tertentu.
func (dp *DashboardProjection) GetTeacherStats(ctx context.Context, schoolID, periodID, teacherID string) (*models.TeacherStats, error) {
	var stats models.TeacherStats
	err := dp.db.WithContext(ctx).
		Where("school_id = ? AND period_id = ? AND teacher_id = ?", schoolID, periodID, teacherID).
		First(&stats).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &stats, nil
}

// HandlerCount mengembalikan jumlah handler yang terdaftar.
func (dp *DashboardProjection) HandlerCount() int {
	return len(dp.handlers)
}
