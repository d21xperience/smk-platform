// internal/service/class_service.go
package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	platformcontext "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/academic/engine"
	"sekolah-platform/services/tu-core/internal/domain/academic/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/academic/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

// ClassService mengorchestrasi use case untuk domain Class/Rombel.
type ClassService struct {
	repo   repository.ClassRepository
	engine interface {
		CreateClass(cmd commands.CreateClassCommand, opCtx *platformcontext.OperationalContext) engine.ClassResult
		UpdateClass(cmd commands.UpdateClassCommand, current *models.Class, opCtx *platformcontext.OperationalContext) engine.ClassResult
		SetHomeroom(classID, teacherID string, current *models.Class, opCtx *platformcontext.OperationalContext) engine.ClassResult
	}
	publisher EventPublisher
}

// NewClassService membuat instance baru dari ClassService.
func NewClassService(
	repo repository.ClassRepository,
	eng interface {
		CreateClass(cmd commands.CreateClassCommand, opCtx *platformcontext.OperationalContext) engine.ClassResult
		UpdateClass(cmd commands.UpdateClassCommand, current *models.Class, opCtx *platformcontext.OperationalContext) engine.ClassResult
		SetHomeroom(classID, teacherID string, current *models.Class, opCtx *platformcontext.OperationalContext) engine.ClassResult
	},
	publisher EventPublisher,
) *ClassService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &ClassService{
		repo:      repo,
		engine:    eng,
		publisher: publisher,
	}
}

// ClassRequest untuk create & update
type CreateClassRequest struct {
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
}

type UpdateClassRequest struct {
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
	Status     string `json:"status"` // active / inactive
}

type ClassResponse struct {
	ID         string    `json:"id"`
	SchoolID   string    `json:"schoolId"`
	PeriodID   string    `json:"periodId"`
	Name       string    `json:"name"`
	Grade      int       `json:"grade"`
	Level      string    `json:"level"`
	Major      string    `json:"major"`
	HomeroomID string    `json:"homeroomId"`
	Capacity   int       `json:"capacity"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

// CreateClass membuat rombel baru.
func (s *ClassService) CreateClass(ctx context.Context, req CreateClassRequest) (*ClassResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	if req.Name == "" {
		return nil, fmt.Errorf("%w: nama kelas wajib diisi", ErrInvalidInput)
	}
	if req.Capacity <= 0 {
		return nil, fmt.Errorf("%w: kapasitas harus lebih dari 0", ErrInvalidInput)
	}

	// 1. Siapkan command
	cmd := commands.CreateClassCommand{
		ClassID:    uuid.New().String(),
		Name:       req.Name,
		Grade:      req.Grade,
		Level:      req.Level,
		Major:      req.Major,
		HomeroomID: req.HomeroomID,
		Capacity:   req.Capacity,
	}

	// 2. Panggil Engine
	result := s.engine.CreateClass(cmd, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateClass", "Class", "", "gagal membuat class", result.Error)
	}

	domainClass := result.Class

	// 3. Simpan ke Repository
	repoClass := &repository.Class{
		ID:               domainClass.ID(),
		SchoolID:         domainClass.SchoolID(),
		AcademicPeriodID: domainClass.PeriodID(),
		Name:             domainClass.Name(),
		Grade:            domainClass.Grade(),
		Level:            mapLevelToInt(domainClass.Level()),
		Major:            domainClass.Major(),
		HomeroomID:       domainClass.HomeroomID(),
		Capacity:         domainClass.Capacity(),
		Status:           domainClass.Status(),
		CreatedAt:        domainClass.CreatedAt(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.Create(ctx, repoClass); err != nil {
		return nil, NewServiceError("CreateClass", "Class", repoClass.ID, "gagal simpan ke database", err)
	}

	// 4. Publish Event (opsional)
	event := DomainEvent{
		EventType:        "ClassCreated",
		AggregateType:    "Class",
		AggregateID:      repoClass.ID,
		Payload:          map[string]interface{}{"classId": repoClass.ID, "name": repoClass.Name},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toClassResponse(repoClass), nil
}

// GetClass mengambil data rombel berdasarkan ID.
func (s *ClassService) GetClass(ctx context.Context, classID string) (*ClassResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	class, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: class %s", ErrNotFound, classID)
		}
		return nil, NewServiceError("GetClass", "Class", classID, "gagal ambil data", err)
	}
	return toClassResponse(class), nil
}

// ListClasses mengambil daftar rombel dengan filter.
func (s *ClassService) ListClasses(ctx context.Context, filter repository.ClassFilter) (*ClassListResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter.SchoolID = opCtx.SchoolID
	filter.AcademicPeriodID = opCtx.AcademicPeriodID

	classes, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListClasses", "Class", "", "gagal ambil daftar", err)
	}

	data := make([]ClassResponse, 0, len(classes))
	for _, c := range classes {
		data = append(data, *toClassResponse(&c))
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	return &ClassListResponse{
		Data:    data,
		Total:   total,
		Limit:   limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+limit < total,
	}, nil
}

// UpdateClass memperbarui data rombel.
func (s *ClassService) UpdateClass(ctx context.Context, classID string, req UpdateClassRequest) (*ClassResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Ambil existing dari repository
	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: class %s", ErrNotFound, classID)
		}
		return nil, NewServiceError("UpdateClass", "Class", classID, "gagal ambil data", err)
	}

	// 2. Rekonstruksi domain Class dari data repository
	domainClass := models.LoadClass(
		existing.ID,
		existing.SchoolID,
		existing.AcademicPeriodID,
		existing.Name,
		existing.Grade,
		mapLevelToString(existing.Level),
		existing.Major,
		existing.HomeroomID,
		existing.Capacity,
		existing.Status,
		existing.CreatedAt,
	)

	// 3. Siapkan command
	cmd := commands.UpdateClassCommand{
		ClassID:    classID,
		Name:       req.Name,
		Grade:      req.Grade,
		Level:      req.Level,
		Major:      req.Major,
		HomeroomID: req.HomeroomID,
		Capacity:   req.Capacity,
		Status:     req.Status,
	}

	// 4. Panggil Engine
	result := s.engine.UpdateClass(cmd, domainClass, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("UpdateClass", "Class", classID, "gagal update class", result.Error)
	}

	updated := result.Class

	// 5. Simpan ke Repository
	repoUpdated := &repository.Class{
		ID:               updated.ID(),
		SchoolID:         updated.SchoolID(),
		AcademicPeriodID: updated.PeriodID(),
		Name:             updated.Name(),
		Grade:            updated.Grade(),
		Level:            mapLevelToInt(updated.Level()),
		Major:            updated.Major(),
		HomeroomID:       updated.HomeroomID(),
		Capacity:         updated.Capacity(),
		Status:           updated.Status(),
		CreatedAt:        updated.CreatedAt(),
		UpdatedAt:        updated.UpdatedAt,
	}

	if err := s.repo.Update(ctx, repoUpdated); err != nil {
		return nil, NewServiceError("UpdateClass", "Class", classID, "gagal update database", err)
	}

	// 6. Event
	event := DomainEvent{
		EventType:        "ClassUpdated",
		AggregateType:    "Class",
		AggregateID:      classID,
		Payload:          map[string]interface{}{"classId": classID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toClassResponse(repoUpdated), nil
}

// DeleteClass menghapus rombel (soft delete).
func (s *ClassService) DeleteClass(ctx context.Context, classID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if err := s.repo.Delete(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID); err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("%w: class %s", ErrNotFound, classID)
		}
		return NewServiceError("DeleteClass", "Class", classID, "gagal hapus data", err)
	}

	event := DomainEvent{
		EventType:        "ClassDeleted",
		AggregateType:    "Class",
		AggregateID:      classID,
		Payload:          map[string]interface{}{"classId": classID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// SetHomeroom mengatur wali kelas.
func (s *ClassService) SetHomeroom(ctx context.Context, classID, teacherID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if classID == "" || teacherID == "" {
		return fmt.Errorf("%w: classID dan teacherID wajib diisi", ErrInvalidInput)
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		return NewServiceError("SetHomeroom", "Class", classID, "gagal ambil data", err)
	}

	domainClass := models.LoadClass(
		existing.ID,
		existing.SchoolID,
		existing.AcademicPeriodID,
		existing.Name,
		existing.Grade,
		mapLevelToString(existing.Level),
		existing.Major,
		existing.HomeroomID,
		existing.Capacity,
		existing.Status,
		existing.CreatedAt,
	)

	result := s.engine.SetHomeroom(classID, teacherID, domainClass, opCtx)
	if result.Error != nil {
		return NewServiceError("SetHomeroom", "Class", classID, "gagal set wali kelas", result.Error)
	}

	updated := result.Class
	repoUpdated := &repository.Class{
		ID:               updated.ID(),
		SchoolID:         updated.SchoolID(),
		AcademicPeriodID: updated.PeriodID(),
		Name:             updated.Name(),
		Grade:            updated.Grade(),
		Level:            mapLevelToInt(updated.Level()),
		Major:            updated.Major(),
		HomeroomID:       updated.HomeroomID(),
		Capacity:         updated.Capacity(),
		Status:           updated.Status(),
		CreatedAt:        updated.CreatedAt(),
		UpdatedAt:        updated.UpdatedAt,
	}
	if err := s.repo.Update(ctx, repoUpdated); err != nil {
		return NewServiceError("SetHomeroom", "Class", classID, "gagal update database", err)
	}

	event := DomainEvent{
		EventType:        "ClassHomeroomChanged",
		AggregateType:    "Class",
		AggregateID:      classID,
		Payload:          map[string]interface{}{"classId": classID, "homeroomId": teacherID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// --- Response DTOs ---

type ClassListResponse struct {
	Data    []ClassResponse `json:"data"`
	Total   int             `json:"total"`
	Limit   int             `json:"limit"`
	Offset  int             `json:"offset"`
	HasMore bool            `json:"hasMore"`
}

// Helper functions

func toClassResponse(c *repository.Class) *ClassResponse {
	return &ClassResponse{
		ID:         c.ID,
		SchoolID:   c.SchoolID,
		PeriodID:   c.AcademicPeriodID,
		Name:       c.Name,
		Grade:      c.Grade,
		Level:      mapLevelToString(c.Level),
		Major:      c.Major,
		HomeroomID: c.HomeroomID,
		Capacity:   c.Capacity,
		Status:     c.Status,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

// mapLevelToInt mengubah level string ke integer (untuk repository)
func mapLevelToInt(level string) int {
	switch level {
	case "X":
		return 10
	case "XI":
		return 11
	case "XII":
		return 12
	default:
		return 0
	}
}

// mapLevelToString mengubah integer ke string (untuk response)
func mapLevelToString(level int) string {
	switch level {
	case 10:
		return "X"
	case 11:
		return "XI"
	case 12:
		return "XII"
	default:
		return ""
	}
}
