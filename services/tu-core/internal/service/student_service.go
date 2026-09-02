package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	platformcontext "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/student/engine"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/student/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

// StudentService mengorchestrasi use case untuk domain Student.
type StudentService struct {
	repo   repository.StudentRepository
	engine interface {
		EnrollStudent(cmd commands.EnrollStudentCommand, opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
		GraduateStudent(opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
	}
	publisher EventPublisher
}

// NewStudentService membuat instance baru dari StudentService.
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

// CreateStudent membuat siswa baru.
func (s *StudentService) CreateStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID string, req CreateStudentRequest) (*StudentResponse, error) {
	// 1. Ambil Operational Context dari context (sudah di-inject dan divalidasi oleh middleware)
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 2. Validasi input dasar
	if req.Name == "" || req.NIS == "" {
		return nil, fmt.Errorf("%w: nama dan NIS wajib diisi", ErrInvalidInput)
	}

	// 3. Cek duplikasi NIS
	existing, err := s.repo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.NIS)
	if err != nil && !isNotFoundError(err) {
		return nil, NewServiceError("CreateStudent", "Student", "", "gagal cek duplikasi NIS", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("%w: NIS %s sudah terdaftar", ErrDuplicateEntry, req.NIS)
	}

	// 4. Siapkan StudentData
	studentData := models.StudentData{
		SchoolID:   opCtx.SchoolID,
		NISN:       req.NISN,
		NIS:        req.NIS,
		FirstName:  req.Name,
		MiddleName: "",
		LastName:   "",
		BirthDate:  req.BirthDate,
		Gender:     req.Gender,
		Address: models.Address{
			Street: req.Address,
		},
		ContactInfo: models.ContactInfo{
			Phone: req.Phone,
			Email: req.Email,
		},
		GuardianInfo: models.GuardianInfo{
			Name:  req.ParentName,
			Phone: req.ParentPhone,
		},
	}

	// 5. Bangun Domain Model
	domainStudent, err := models.NewStudent(studentData, opCtx)
	if err != nil {
		return nil, NewServiceError("CreateStudent", "Student", "", "gagal membuat domain model", err)
	}

	// 6. Siapkan Command untuk Engine (HANYA lakukan enroll jika ClassID disediakan)
	if req.ClassID != "" {
		enrollCmd := commands.EnrollStudentCommand{
			EnrollmentID: uuid.New().String(),
			ClassID:      req.ClassID,
		}

		// 7. Panggil Domain Engine
		result := s.engine.EnrollStudent(enrollCmd, opCtx, domainStudent)
		if result.Error != nil {
			return nil, NewServiceError("CreateStudent", "Student", "", "gagal enroll student", result.Error)
		}
	}

	// 8. Simpan ke Repository
	repoStudent := &repository.Student{
		ID:               domainStudent.StudentID(), // Gunakan ID yang sudah digenerate oleh domain
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		NIS:              req.NIS,
		NISN:             req.NISN,
		Name:             req.Name,
		Gender:           req.Gender,
		BirthDate:        req.BirthDate,
		BirthPlace:       req.BirthPlace,
		Address:          req.Address,
		Phone:            req.Phone,
		Email:            req.Email,
		ParentName:       req.ParentName,
		ParentPhone:      req.ParentPhone,
		ParentEmail:      req.ParentEmail,
		ClassName:        req.ClassID,
		Status:           "ACTIVE",
	}

	if err := s.repo.Create(ctx, repoStudent); err != nil {
		fmt.Printf("[StudentService] ❌ Database error: %v\n", err)
		fmt.Printf("[StudentService] Student data: %+v\n", repoStudent)
		return nil, NewServiceError("CreateStudent", "Student", repoStudent.ID, "gagal simpan ke database", err)
	}

	// 9. Publish Domain Event
	event := DomainEvent{
		EventType:     "StudentEnrolled", // Projection mendengarkan event ini
		AggregateType: "Student",
		AggregateID:   repoStudent.ID,
		Payload: map[string]interface{}{
			"studentId": repoStudent.ID,
			"classId":   req.ClassID,
			"name":      req.Name,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    requestID,
		Version:          "v1",
	}
	if err := s.publisher.Publish(ctx, event); err != nil {
		fmt.Printf("[StudentService] Gagal publish event StudentEnrolled: %v\n", err)
	}

	return toStudentResponse(repoStudent), nil
}

// GetStudent mengambil data siswa berdasarkan ID.
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

// ListStudents mengambil daftar siswa dengan filter dan pagination.
func (s *StudentService) ListStudents(ctx context.Context, schoolID, academicPeriodID string, filter StudentFilterRequest) (*StudentListResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	repoFilter := repository.StudentFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		ClassID:          filter.ClassID,
		Status:           filter.Status,
		Keyword:          filter.Keyword,
		Limit:            filter.Limit,
		Offset:           filter.Offset,
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
	if limit <= 0 {
		limit = 20
	}

	return &StudentListResponse{
		Data:    data,
		Total:   total,
		Limit:   limit,
		Offset:  filter.Offset,
		HasMore: filter.Offset+limit < total,
	}, nil
}

// UpdateStudent memperbarui data siswa.
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

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Gender != "" {
		existing.Gender = req.Gender
	}
	if !req.BirthDate.IsZero() {
		existing.BirthDate = req.BirthDate
	}
	if req.BirthPlace != "" {
		existing.BirthPlace = req.BirthPlace
	}
	if req.Address != "" {
		existing.Address = req.Address
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	if req.ParentName != "" {
		existing.ParentName = req.ParentName
	}
	if req.ParentPhone != "" {
		existing.ParentPhone = req.ParentPhone
	}
	if req.ParentEmail != "" {
		existing.ParentEmail = req.ParentEmail
	}

	if err := s.repo.Update(ctx, existing); err != nil {
		return nil, NewServiceError("UpdateStudent", "Student", studentID, "gagal update data", err)
	}

	event := DomainEvent{
		EventType:        "StudentUpdated",
		AggregateType:    "Student",
		AggregateID:      studentID,
		Payload:          map[string]interface{}{"studentId": studentID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toStudentResponse(existing), nil
}

// DeleteStudent menghapus siswa (soft delete).
func (s *StudentService) DeleteStudent(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if err := s.repo.Delete(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID); err != nil {
		if isNotFoundError(err) {
			return fmt.Errorf("%w: student %s", ErrNotFound, studentID)
		}
		return NewServiceError("DeleteStudent", "Student", studentID, "gagal hapus data", err)
	}

	event := DomainEvent{
		EventType:        "StudentDeleted",
		AggregateType:    "Student",
		AggregateID:      studentID,
		Payload:          map[string]interface{}{"studentId": studentID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// GraduateStudent menandai siswa sebagai lulus.
func (s *StudentService) GraduateStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		return NewServiceError("GraduateStudent", "Student", studentID, "gagal ambil data", err)
	}

	studentData := models.StudentData{
		StudentID:  existing.ID,
		SchoolID:   existing.SchoolID,
		NISN:       existing.NISN,
		NIS:        existing.NIS,
		FirstName:  existing.Name,
		MiddleName: "",
		LastName:   "",
		BirthDate:  existing.BirthDate,
		Gender:     existing.Gender,
		Address: models.Address{
			Street: existing.Address,
		},
		ContactInfo: models.ContactInfo{
			Phone: existing.Phone,
			Email: existing.Email,
		},
		GuardianInfo: models.GuardianInfo{
			Name:  existing.ParentName,
			Phone: existing.ParentPhone,
		},
	}

	domainStudent, err := models.NewStudent(studentData, opCtx)
	if err != nil {
		return NewServiceError("GraduateStudent", "Student", studentID, "gagal membuat domain model", err)
	}

	gradResult := s.engine.GraduateStudent(opCtx, domainStudent)
	if gradResult.Error != nil {
		return NewServiceError("GraduateStudent", "Student", studentID, "gagal graduate", gradResult.Error)
	}

	existing.Status = "GRADUATED"
	if err := s.repo.Update(ctx, existing); err != nil {
		return NewServiceError("GraduateStudent", "Student", studentID, "gagal update status", err)
	}

	event := DomainEvent{
		EventType:        "StudentGraduated",
		AggregateType:    "Student",
		AggregateID:      studentID,
		Payload:          map[string]interface{}{"studentId": studentID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    requestID,
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// EnrollStudent mendaftarkan siswa yang sudah ada ke dalam kelas tertentu.
func (s *StudentService) EnrollStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID, enrollmentID, classID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		return NewServiceError("EnrollStudent", "Student", studentID, "gagal ambil data", err)
	}

	existing.ClassName = classID
	existing.Status = "ACTIVE"
	if err := s.repo.Update(ctx, existing); err != nil {
		return NewServiceError("EnrollStudent", "Student", studentID, "gagal update kelas", err)
	}

	event := DomainEvent{
		EventType:     "StudentEnrolled",
		AggregateType: "Student",
		AggregateID:   studentID,
		Payload: map[string]interface{}{
			"studentId":    studentID,
			"enrollmentId": enrollmentID,
			"classId":      classID,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    requestID,
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// TransferStudent memindahkan siswa ke sekolah lain.
func (s *StudentService) TransferStudent(ctx context.Context, schoolID, academicPeriodID string, userID int, role, requestID, studentID, targetSchool, reason string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	existing, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		return NewServiceError("TransferStudent", "Student", studentID, "gagal ambil data", err)
	}

	existing.Status = "TRANSFERRED"
	if err := s.repo.Update(ctx, existing); err != nil {
		return NewServiceError("TransferStudent", "Student", studentID, "gagal update status", err)
	}

	event := DomainEvent{
		EventType:     "StudentTransferred",
		AggregateType: "Student",
		AggregateID:   studentID,
		Payload: map[string]interface{}{
			"studentId":    studentID,
			"targetSchool": targetSchool,
			"reason":       reason,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    requestID,
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

// CheckNISAvailability memeriksa apakah NIS sudah digunakan.
func (s *StudentService) CheckNISAvailability(ctx context.Context, schoolID, academicPeriodID, nis string) (bool, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return false, ErrMissingOperationalContext
	}

	student, err := s.repo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, nis)
	if err != nil {
		if isNotFoundError(err) {
			return true, nil
		}
		return false, NewServiceError("CheckNISAvailability", "Student", "", "gagal cek NIS", err)
	}

	if student.DeletedAt != nil {
		return true, nil
	}

	return false, nil
}

// CheckNISNAvailability memeriksa apakah NISN sudah digunakan.
func (s *StudentService) CheckNISNAvailability(ctx context.Context, schoolID, academicPeriodID, nisn string) (bool, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return false, ErrMissingOperationalContext
	}
	// TODO: Implementasi GetByNISN di repository jika diperlukan
	return true, nil
}

// === Helper Functions ===

func toStudentResponse(s *repository.Student) *StudentResponse {
	return &StudentResponse{
		ID:               s.ID,
		SchoolID:         s.SchoolID,
		AcademicPeriodID: s.AcademicPeriodID,
		NIS:              s.NIS,
		NISN:             s.NISN,
		Name:             s.Name,
		Gender:           s.Gender,
		BirthDate:        s.BirthDate,
		BirthPlace:       s.BirthPlace,
		Address:          s.Address,
		Phone:            s.Phone,
		Email:            s.Email,
		ParentName:       s.ParentName,
		ParentPhone:      s.ParentPhone,
		ParentEmail:      s.ParentEmail,
		ClassName:        s.ClassName,
		Status:           s.Status,
		CreatedAt:        s.CreatedAt,
		UpdatedAt:        s.UpdatedAt,
	}
}

func isNotFoundError(err error) bool {
	return err != nil && (err == repository.ErrNotFound || err.Error() == "data tidak ditemukan" || contains(err.Error(), "not found"))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

var _ = time.Now
