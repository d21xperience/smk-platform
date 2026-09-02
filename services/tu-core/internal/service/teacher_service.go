package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"sekolah-platform/services/tu-core/internal/repository"
)

// TeacherService mengorchestrasi use case untuk domain Teacher.
type TeacherService struct {
	repo      repository.TeacherRepository
	publisher EventPublisher
}

// NewTeacherService membuat instance baru.
func NewTeacherService(repo repository.TeacherRepository, publisher EventPublisher) *TeacherService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &TeacherService{repo: repo, publisher: publisher}
}

// CreateTeacher membuat guru baru.
func (s *TeacherService) CreateTeacher(ctx context.Context, schoolID, academicPeriodID string, req CreateTeacherRequest) (*TeacherResponse, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	if req.Name == "" || req.NIP == "" {
		return nil, fmt.Errorf("%w: nama dan NIP wajib diisi", ErrInvalidInput)
	}

	// Cek duplikasi NIP
	existing, err := s.repo.GetByNIP(ctx, schoolID, academicPeriodID, req.NIP)
	if err != nil && !isNotFoundError(err) {
		return nil, NewServiceError("CreateTeacher", "Teacher", "", "gagal cek duplikasi NIP", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: NIP %s sudah terdaftar", ErrDuplicateEntry, req.NIP)
	}

	teacher := &repository.Teacher{
		ID:               uuid.New().String(),
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		NIP:              req.NIP,
		Name:             req.Name,
		Gender:           req.Gender,
		BirthDate:        req.BirthDate,
		Email:            req.Email,
		Phone:            req.Phone,
		Address:          req.Address,
		SubjectID:        req.SubjectID,
		Status:           "ACTIVE",
	}

	if err := s.repo.Create(ctx, teacher); err != nil {
		return nil, NewServiceError("CreateTeacher", "Teacher", teacher.ID, "gagal simpan", err)
	}

	// Publish event
	event := DomainEvent{
		EventType:        "TeacherCreated",
		AggregateType:    "Teacher",
		AggregateID:      teacher.ID,
		Payload:          map[string]interface{}{"teacherId": teacher.ID, "name": teacher.Name},
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toTeacherResponse(teacher), nil
}

// GetTeacher mengambil data guru berdasarkan ID.
func (s *TeacherService) GetTeacher(ctx context.Context, schoolID, academicPeriodID, teacherID string) (*TeacherResponse, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}

	teacher, err := s.repo.GetByID(ctx, schoolID, academicPeriodID, teacherID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: teacher %s", ErrNotFound, teacherID)
		}
		return nil, NewServiceError("GetTeacher", "Teacher", teacherID, "gagal ambil data", err)
	}

	return toTeacherResponse(teacher), nil
}

// ListTeachers mengambil daftar guru.
func (s *TeacherService) ListTeachers(ctx context.Context, schoolID, academicPeriodID, subjectID, status, keyword string, limit, offset int) ([]TeacherResponse, int, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	filter := repository.TeacherFilter{
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		SubjectID:        subjectID,
		Status:           status,
		Keyword:          keyword,
		Limit:            limit,
		Offset:           offset,
	}

	teachers, total, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, 0, NewServiceError("ListTeachers", "Teacher", "", "gagal ambil daftar", err)
	}

	data := make([]TeacherResponse, 0, len(teachers))
	for _, t := range teachers {
		data = append(data, *toTeacherResponse(&t))
	}

	return data, total, nil
}

// DeleteTeacher menghapus guru (soft delete).
func (s *TeacherService) DeleteTeacher(ctx context.Context, schoolID, academicPeriodID, teacherID string) error {
	if schoolID == "" || academicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	if err := s.repo.Delete(ctx, schoolID, academicPeriodID, teacherID); err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("%w: teacher %s", ErrNotFound, teacherID)
		}
		return NewServiceError("DeleteTeacher", "Teacher", teacherID, "gagal hapus", err)
	}

	event := DomainEvent{
		EventType:        "TeacherDeleted",
		AggregateType:    "Teacher",
		AggregateID:      teacherID,
		Payload:          map[string]interface{}{"teacherId": teacherID},
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// toTeacherResponse mengkonversi repository.Teacher ke TeacherResponse DTO.
func toTeacherResponse(t *repository.Teacher) *TeacherResponse {
	return &TeacherResponse{
		ID:               t.ID,
		SchoolID:         t.SchoolID,
		AcademicPeriodID: t.AcademicPeriodID,
		NIP:              t.NIP,
		Name:             t.Name,
		Gender:           t.Gender,
		BirthDate:        t.BirthDate,
		Email:            t.Email,
		Phone:            t.Phone,
		Address:          t.Address,
		SubjectID:        t.SubjectID,
		Status:           t.Status,
		CreatedAt:        t.CreatedAt,
	}
}
