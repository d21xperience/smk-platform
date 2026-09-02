// Package projection menyediakan projection handlers untuk update read model.
package projection

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/edge-bff/internal/repository"
)

// ProjectionManager mengelola semua projection handlers.
type ProjectionManager struct {
	repo     *repository.SQLiteRepository
	handlers map[string]ProjectionHandler
}

// ProjectionHandler adalah handler untuk event tertentu.
type ProjectionHandler func(ctx context.Context, event types.DomainEvent) error

// NewProjectionManager membuat ProjectionManager baru.
func NewProjectionManager(repo *repository.SQLiteRepository) *ProjectionManager {
	pm := &ProjectionManager{
		repo:     repo,
		handlers: make(map[string]ProjectionHandler),
	}

	// Register handlers
	pm.registerStudentHandlers()
	pm.registerTeacherHandlers()
	pm.registerAcademicHandlers()

	return pm
}

// Handle mendistribusikan event ke handler yang sesuai.
func (pm *ProjectionManager) Handle(ctx context.Context, event types.DomainEvent) error {
	// Idempotency check
	eventID := event.GetEventID()
	processed, err := pm.repo.HasProcessedEvent(ctx, eventID)
	if err != nil {
		return err
	}
	if processed {
		return nil // Skip, sudah diproses
	}

	// Find handler
	handler, ok := pm.handlers[event.GetEventName()]
	if !ok {
		// Event tidak punya handler, skip
		return nil
	}

	return handler(ctx, event)
}

// RegisterHandler mendaftarkan handler untuk event tertentu.
func (pm *ProjectionManager) RegisterHandler(eventName string, handler ProjectionHandler) {
	pm.handlers[eventName] = handler
}

// === Student Handlers ===

func (pm *ProjectionManager) registerStudentHandlers() {
	pm.handlers["StudentCreated"] = pm.handleStudentCreated
	pm.handlers["StudentEnrolled"] = pm.handleStudentEnrolled
	pm.handlers["StudentGraduated"] = pm.handleStudentStatusChange("GRADUATED")
	pm.handlers["StudentTransferred"] = pm.handleStudentStatusChange("TRANSFERRED")
	pm.handlers["StudentProfileUpdated"] = pm.handleStudentProfileUpdated
}

func (pm *ProjectionManager) handleStudentCreated(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})

	projection := &repository.StudentProjection{
		StudentID:   getStr(payload, "studentId"),
		SchoolID:    event.GetSchoolID(),
		NISN:        getStr(payload, "nisn"),
		NIS:         getStr(payload, "nis"),
		FirstName:   getStr(payload, "firstName"),
		LastName:    getStr(payload, "lastName"),
		Gender:      getStr(payload, "gender"),
		Status:      "ACTIVE",
		CreatedAt:   event.GetOccurredAt(),
		UpdatedAt:   event.GetOccurredAt(),
		LastEventID: event.GetEventID(),
	}

	// Extract birthDate jika ada
	if bd, ok := payload["birthDate"].(string); ok {
		if t, err := time.Parse(time.RFC3339, bd); err == nil {
			projection.BirthDate = t
		}
	}

	return pm.repo.UpsertStudent(ctx, projection)
}

func (pm *ProjectionManager) handleStudentEnrolled(ctx context.Context, event types.DomainEvent) error {
	// Update enrollment info (bisa ditambahkan table enrollments)
	// Untuk simplicity, kita update timestamp
	payload := event.GetPayload().(map[string]interface{})
	studentID := getStr(payload, "studentId")

	existing, err := pm.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil // Student belum ada di projection
	}

	existing.UpdatedAt = event.GetOccurredAt()
	existing.LastEventID = event.GetEventID()
	return pm.repo.UpsertStudent(ctx, existing)
}

func (pm *ProjectionManager) handleStudentStatusChange(status string) ProjectionHandler {
	return func(ctx context.Context, event types.DomainEvent) error {
		payload := event.GetPayload().(map[string]interface{})
		studentID := getStr(payload, "studentId")
		return pm.repo.UpdateStudentStatus(ctx, studentID, status, event.GetEventID())
	}
}

func (pm *ProjectionManager) handleStudentProfileUpdated(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})
	studentID := getStr(payload, "studentId")

	existing, err := pm.repo.FindStudentByID(ctx, studentID)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}

	existing.UpdatedAt = event.GetOccurredAt()
	existing.LastEventID = event.GetEventID()
	return pm.repo.UpsertStudent(ctx, existing)
}

// === Teacher Handlers ===

func (pm *ProjectionManager) registerTeacherHandlers() {
	pm.handlers["TeacherCreated"] = pm.handleTeacherCreated
	pm.handlers["TeacherProfileUpdated"] = pm.handleTeacherProfileUpdated
	pm.handlers["TeacherDeactivated"] = pm.handleTeacherStatusChange("INACTIVE")
}

func (pm *ProjectionManager) handleTeacherCreated(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})

	projection := &repository.TeacherProjection{
		TeacherID:   getStr(payload, "teacherId"),
		SchoolID:    event.GetSchoolID(),
		NIP:         getStr(payload, "nip"),
		FirstName:   getStr(payload, "fullName"), // Simplified
		Subject:     getStr(payload, "subject"),
		Gender:      getStr(payload, "gender"),
		Status:      "ACTIVE",
		UpdatedAt:   event.GetOccurredAt(),
		LastEventID: event.GetEventID(),
	}

	return pm.repo.UpsertTeacher(ctx, projection)
}

func (pm *ProjectionManager) handleTeacherProfileUpdated(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})
	teacherID := getStr(payload, "teacherId")

	existing, err := pm.repo.FindTeacherByID(ctx, teacherID)
	if err != nil {
		return err
	}
	if existing == nil {
		return nil
	}

	existing.UpdatedAt = event.GetOccurredAt()
	existing.LastEventID = event.GetEventID()
	return pm.repo.UpsertTeacher(ctx, existing)
}

func (pm *ProjectionManager) handleTeacherStatusChange(status string) ProjectionHandler {
	return func(ctx context.Context, event types.DomainEvent) error {
		payload := event.GetPayload().(map[string]interface{})
		teacherID := getStr(payload, "teacherId")

		existing, err := pm.repo.FindTeacherByID(ctx, teacherID)
		if err != nil {
			return err
		}
		if existing == nil {
			return nil
		}

		existing.Status = status
		existing.UpdatedAt = event.GetOccurredAt()
		existing.LastEventID = event.GetEventID()
		return pm.repo.UpsertTeacher(ctx, existing)
	}
}

// === Academic Handlers ===

func (pm *ProjectionManager) registerAcademicHandlers() {
	pm.handlers["ClassCreated"] = pm.handleClassCreated
	pm.handlers["ClassHomeroomChanged"] = pm.handleClassHomeroomChanged
}

func (pm *ProjectionManager) handleClassCreated(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})

	projection := &repository.ClassProjection{
		ClassID:     getStr(payload, "classId"),
		SchoolID:    event.GetSchoolID(),
		PeriodID:    getStr(payload, "periodId"),
		Name:        getStr(payload, "name"),
		Grade:       getInt(payload, "grade"),
		Level:       getStr(payload, "level"),
		Major:       getStr(payload, "major"),
		HomeroomID:  getStr(payload, "homeroomId"),
		Capacity:    getInt(payload, "capacity"),
		UpdatedAt:   event.GetOccurredAt(),
		LastEventID: event.GetEventID(),
	}

	return pm.repo.UpsertClass(ctx, projection)
}

func (pm *ProjectionManager) handleClassHomeroomChanged(ctx context.Context, event types.DomainEvent) error {
	payload := event.GetPayload().(map[string]interface{})
	classID := getStr(payload, "classId")
	homeroomID := getStr(payload, "homeroomId")

	// Find and update class
	var existing repository.ClassProjection
	if err := pm.repo.DB().Where("class_id = ?", classID).First(&existing).Error; err != nil {
		return err
	}

	existing.HomeroomID = homeroomID
	existing.UpdatedAt = event.GetOccurredAt()
	existing.LastEventID = event.GetEventID()
	return pm.repo.UpsertClass(ctx, &existing)
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

// RegisteredEvents mengembalikan daftar event yang terdaftar.
func (pm *ProjectionManager) RegisteredEvents() []string {
	events := make([]string, 0, len(pm.handlers))
	for name := range pm.handlers {
		events = append(events, name)
	}
	return events
}

// HandlerCount mengembalikan jumlah handler yang terdaftar.
func (pm *ProjectionManager) HandlerCount() int {
	return len(pm.handlers)
}

// Placeholder untuk fmt
var _ = fmt.Sprintf
