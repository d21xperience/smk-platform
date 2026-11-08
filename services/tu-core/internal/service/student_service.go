package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	platformcontext "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/student/engine"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/student/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type StudentService struct {
	repo   repository.StudentRepository
	engine interface {
		EnrollStudent(cmd commands.EnrollStudentCommand, opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
		GraduateStudent(opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
	}
	publisher EventPublisher
}

func NewStudentService(
	repo repository.StudentRepository,
	engine interface {
		EnrollStudent(cmd commands.EnrollStudentCommand, opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
		GraduateStudent(opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
	},
	publisher EventPublisher,
) *StudentService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &StudentService{
		repo:      repo,
		engine:    engine,
		publisher: publisher,
	}
}

func (s *StudentService) CreateStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID string, req CreateStudentRequest) (*StudentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	nisStr := string(req.NIS)
	nisnStr := string(req.NISN)
	phoneStr := string(req.Phone)
	parentPhoneStr := string(req.ParentPhone)

	if req.Name == "" || nisStr == "" {
		return nil, fmt.Errorf("%w: nama dan NIS wajib diisi", ErrInvalidInput)
	}
	if len(nisStr) != 8 {
		return nil, fmt.Errorf("%w: NIS harus tepat 8 digit (yang dikirim: %d digit)", ErrInvalidInput, len(nisStr))
	}

	existing, err := s.repo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, nisStr)
	if err != nil && !isNotFoundError(err) {
		return nil, NewServiceError("CreateStudent", "Student", "", "gagal cek duplikasi NIS", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: NIS %s sudah terdaftar", ErrDuplicateEntry, nisStr)
	}

	if nisnStr != "" {
		existingNISN, err := s.repo.GetByNISN(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, nisnStr)
		if err != nil && !isNotFoundError(err) {
			return nil, NewServiceError("CreateStudent", "Student", "", "gagal cek duplikasi NISN", err)
		}
		if existingNISN != nil {
			return nil, fmt.Errorf("%w: NISN %s sudah terdaftar", ErrDuplicateEntry, nisnStr)
		}
	}

	var birthDate time.Time
	if req.BirthDate != nil {
		birthDate = *req.BirthDate
	}

	genderCode := "L"
	gUpper := strings.ToUpper(req.Gender)
	if strings.Contains(gUpper, "P") || strings.Contains(gUpper, "FEMALE") || strings.Contains(gUpper, "WANITA") || strings.Contains(gUpper, "PEREMPUAN") {
		genderCode = "P"
	}

	studentUUID := uuid.New().String()

	repoStudent := &repository.Student{
		ID:               studentUUID,
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		NIS:              nisStr,
		NISN:             nisnStr,
		Name:             req.Name,
		Gender:           genderCode,
		BirthDate:        birthDate,
		BirthPlace:       req.BirthPlace,
		Address:          req.Address,
		Phone:            phoneStr,
		Email:            req.Email,
		ParentName:       req.ParentName,
		ParentPhone:      parentPhoneStr,
		ParentEmail:      req.ParentEmail,
		ClassName:        req.ClassID,
		Status:           "ACTIVE",
	}

	if err := s.repo.Create(ctx, repoStudent); err != nil {
		return nil, NewServiceError("CreateStudent", "Student", repoStudent.ID, "gagal simpan ke database", err)
	}

	event := DomainEvent{
		EventType: "StudentEnrolled", AggregateType: "Student", AggregateID: repoStudent.ID,
		Payload: map[string]interface{}{"studentId": repoStudent.ID, "classId": req.ClassID, "name": req.Name},
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID, CorrelationID: requestID, Version: "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toStudentResponse(repoStudent), nil
}

func (s *StudentService) GetStudent(ctx context.Context, schoolID, academicPeriodID, studentID string) (*StudentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	student, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: student %s", ErrNotFound, studentID)
		}
		return nil, NewServiceError("GetStudent", "Student", studentID, "gagal ambil data", err)
	}
	return toStudentResponse(student), nil
}

func (s *StudentService) ListStudents(ctx context.Context, schoolID, academicPeriodID string, filter StudentFilterRequest) (*StudentListResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	repoFilter := repository.StudentFilter{
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID: filter.ClassID, Status: filter.Status, Keyword: filter.Keyword,
		Limit: filter.Limit, Offset: filter.Offset,
	}

	students, total, err := s.repo.List(ctx, repoFilter)
	if err != nil {
		return nil, NewServiceError("ListStudents", "Student", "", "gagal ambil daftar", err)
	}

	data := make([]StudentResponse, 0, len(students))
	for _, st := range students {
		data = append(data, *toStudentResponse(&st))
	}

	limit := filter.Limit
	if limit <= 0 { limit = 20 }

	return &StudentListResponse{
		Data: data, Total: total, Limit: limit,
		Offset: filter.Offset, HasMore: filter.Offset+limit < total,
	}, nil
}

func (s *StudentService) UpdateStudent(ctx context.Context, schoolID, academicPeriodID, studentID string, req UpdateStudentRequest) (*StudentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		if isNotFoundError(err) {
			return nil, fmt.Errorf("%w: student %s", ErrNotFound, studentID)
		}
		return nil, NewServiceError("UpdateStudent", "Student", studentID, "gagal ambil data", err)
	}

	if req.Name != "" { existing.Name = req.Name }
	if req.Gender != "" {
		gUpper := strings.ToUpper(req.Gender)
		if strings.Contains(gUpper, "P") || strings.Contains(gUpper, "FEMALE") || strings.Contains(gUpper, "WANITA") || strings.Contains(gUpper, "PEREMPUAN") {
			existing.Gender = "P"
		} else {
			existing.Gender = "L"
		}
	}
	if req.BirthDate != nil { existing.BirthDate = *req.BirthDate }
	if req.BirthPlace != "" { existing.BirthPlace = req.BirthPlace }
	if req.Address != "" { existing.Address = req.Address }
	if req.Phone != "" { existing.Phone = string(req.Phone) }
	if req.Email != "" { existing.Email = req.Email }
	if req.ParentName != "" { existing.ParentName = req.ParentName }
	if req.ParentPhone != "" { existing.ParentPhone = string(req.ParentPhone) }
	if req.ParentEmail != "" { existing.ParentEmail = req.ParentEmail }
	if req.ClassID != "" { existing.ClassName = req.ClassID }

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, NewServiceError("UpdateStudent", "Student", studentID, "gagal update data", err)
	}

	event := DomainEvent{
		EventType: "StudentUpdated", AggregateType: "Student", AggregateID: studentID,
		Payload: map[string]interface{}{"studentId": studentID},
		SchoolID: opCtx.SchoolID, AcademicPeriodID: opCtx.AcademicPeriodID, CorrelationID: uuid.New().String(), Version: "v1",
	}
	_ = s.publisher.Publish(ctx, event)
	return toStudentResponse(existing), nil
}

func (s *StudentService) DeleteStudent(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return ErrMissingOperationalContext }

	if err := s.repo.Delete(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID); err != nil {
		if isNotFoundError(err) { return fmt.Errorf("%w: student %s", ErrNotFound, studentID) }
		return NewServiceError("DeleteStudent", "Student", studentID, "gagal hapus data", err)
	}
	return nil
}

func (s *StudentService) GraduateStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return ErrMissingOperationalContext }
	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil { return NewServiceError("GraduateStudent", "Student", studentID, "gagal ambil data", err) }
	existing.Status = "GRADUATED"
	return s.repo.Update(ctx, existing)
}

func (s *StudentService) EnrollStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID, enrollmentID, classID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return ErrMissingOperationalContext }
	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil { return NewServiceError("EnrollStudent", "Student", studentID, "gagal ambil data", err) }
	existing.ClassName = classID
	existing.Status = "ACTIVE"
	return s.repo.Update(ctx, existing)
}

func (s *StudentService) TransferStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID, targetSchool, reason string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return ErrMissingOperationalContext }
	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil { return NewServiceError("TransferStudent", "Student", studentID, "gagal ambil data", err) }
	existing.Status = "TRANSFERRED"
	return s.repo.Update(ctx, existing)
}

func (s *StudentService) CheckNISAvailability(ctx context.Context, schoolID, academicPeriodID, nis string) (bool, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return false, ErrMissingOperationalContext }
	student, err := s.repo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, nis)
	if err != nil {
		if isNotFoundError(err) { return true, nil }
		return false, NewServiceError("CheckNISAvailability", "Student", "", "gagal cek NIS", err)
	}
	return student == nil, nil
}

func (s *StudentService) CheckNISNAvailability(ctx context.Context, schoolID, academicPeriodID, nisn string) (bool, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil { return false, ErrMissingOperationalContext }
	student, err := s.repo.GetByNISN(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, nisn)
	if err != nil {
		if isNotFoundError(err) { return true, nil }
		return false, NewServiceError("CheckNISNAvailability", "Student", "", "gagal cek NISN", err)
	}
	return student == nil, nil
}

func toStudentResponse(s *repository.Student) *StudentResponse {
	return &StudentResponse{
		ID: s.ID, SchoolID: s.SchoolID, AcademicPeriodID: s.AcademicPeriodID,
		NIS: s.NIS, NISN: s.NISN, Name: s.Name, Gender: s.Gender,
		BirthDate: s.BirthDate, BirthPlace: s.BirthPlace, Address: s.Address,
		Phone: s.Phone, Email: s.Email, ParentName: s.ParentName,
		ParentPhone: s.ParentPhone, ParentEmail: s.ParentEmail,
		ClassName: s.ClassName, Status: s.Status,
		CreatedAt: s.CreatedAt, UpdatedAt: s.UpdatedAt,
	}
}

func isNotFoundError(err error) bool {
	return err != nil && (err == repository.ErrNotFound || err.Error() == "data tidak ditemukan" || contains(err.Error(), "not found"))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr { return true }
	}
	return false
}

// =====================================================================
// AREA BARU: Service Khusus UI
// =====================================================================

type StudentUIResponse struct {
	Data       []repository.StudentUIDTO `json:"data"`
	Total      int                       `json:"total"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
	TotalPages int                       `json:"totalPages"`
}

func (s *StudentService) ListStudentsForUI(ctx context.Context, schoolID, academicPeriodID string, filter repository.StudentFilter) (*StudentUIResponse, error) {
	filter.SchoolID = schoolID
	filter.AcademicPeriodID = academicPeriodID

	students, total, err := s.repo.ListForUI(ctx, filter)
	if err != nil { return nil, err }

	page := filter.Offset / filter.Limit + 1
	if filter.Limit == 0 { page = 1; filter.Limit = 20 }

	return &StudentUIResponse{
		Data: students, Total: total, Page: page, Limit: filter.Limit,
		TotalPages: (total + filter.Limit - 1) / filter.Limit,
	}, nil
}

func (s *StudentService) GetStudentForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) (*repository.StudentUIDTO, error) {
	return s.repo.GetByIDForUI(ctx, schoolID, academicPeriodID, studentID)
}

func (s *StudentService) UpdateStudentForUI(ctx context.Context, schoolID, academicPeriodID, studentID string, req struct {
	Name        string `json:"name"`
	NIS         string `json:"nis"`
	NISN        string `json:"nisn"`
	BirthPlace  string `json:"birthPlace"`
	BirthDate   string `json:"birthDate"`
	Gender      string `json:"gender"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	ClassName   string `json:"classId"`
	ParentName  string `json:"parentName"`
	ParentPhone string `json:"parentPhone"`
}) (*repository.StudentUIDTO, error) {
	gender := req.Gender
	if gender == "MALE" { gender = "L" }
	if gender == "FEMALE" { gender = "P" }

	err := s.repo.UpdateForUI(ctx, schoolID, academicPeriodID, studentID, req.Name, req.BirthPlace, req.BirthDate, gender, req.Address, req.Phone, req.ClassName, req.ParentName, req.ParentPhone)
	if err != nil { return nil, err }
	return s.repo.GetByIDForUI(ctx, schoolID, academicPeriodID, studentID)
}

func (s *StudentService) DeleteStudentForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	return s.repo.DeleteForUI(ctx, schoolID, academicPeriodID, studentID)
}

func (s *StudentService) ListAvailableStudentsForUI(ctx context.Context, schoolID, academicPeriodID string, filter repository.StudentFilter) (*StudentUIResponse, error) {
	filter.SchoolID = schoolID
	filter.AcademicPeriodID = academicPeriodID

	students, total, err := s.repo.ListAvailableForUI(ctx, filter)
	if err != nil { return nil, err }

	page := filter.Offset / filter.Limit + 1
	if filter.Limit == 0 { page = 1; filter.Limit = 50 }

	return &StudentUIResponse{
		Data: students, Total: total, Page: page, Limit: filter.Limit,
		TotalPages: (total + filter.Limit - 1) / filter.Limit,
	}, nil
}
