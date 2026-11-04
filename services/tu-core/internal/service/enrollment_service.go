// internal/service/enrollment_service.go
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	platformcontext "sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/student/engine"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/student/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type EnrollmentService struct {
	db             *pgxpool.Pool
	enrollmentRepo repository.EnrollmentRepository
	classRepo      repository.ClassRepository
	studentRepo    repository.StudentRepository
	engine         interface {
		EnrollStudent(cmd commands.EnrollStudentCommand, opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
	}
	publisher EventPublisher
}

func NewEnrollmentService(
	db *pgxpool.Pool,
	enrollmentRepo repository.EnrollmentRepository,
	classRepo repository.ClassRepository,
	studentRepo repository.StudentRepository,
	eng interface {
		EnrollStudent(cmd commands.EnrollStudentCommand, opCtx *platformcontext.OperationalContext, student *models.Student) engine.Result
	},
	publisher EventPublisher,
) *EnrollmentService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &EnrollmentService{
		db:             db,
		enrollmentRepo: enrollmentRepo,
		classRepo:      classRepo,
		studentRepo:    studentRepo,
		engine:         eng,
		publisher:      publisher,
	}
}

// === DTOs ===
type AddMemberRequest struct {
	StudentID string `json:"studentId"`
}

type MoveStudentRequest struct {
	StudentID string `json:"studentId"`
	ToClassID string `json:"toClassId"`
	Reason    string `json:"reason"`
}

type MemberResponse struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID    string    `json:"studentId"`
	NIS          string    `json:"nis"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	EnrolledAt   time.Time `json:"enrolledAt"`
}

// === METHODS ===

func (s *EnrollmentService) AddMember(ctx context.Context, classID string, req AddMemberRequest) (*repository.Enrollment, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	if classID == "" || req.StudentID == "" {
		return nil, fmt.Errorf("%w: classID dan studentID wajib diisi", ErrInvalidInput)
	}

	count, err := s.enrollmentRepo.CountActiveByClass(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal cek kapasitas", err)
	}
	class, err := s.classRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal ambil data kelas", err)
	}
	if class == nil {
		return nil, fmt.Errorf("%w: class %s", ErrNotFound, classID)
	}
	if class.Capacity > 0 && count >= class.Capacity {
		return nil, fmt.Errorf("%w: kapasitas kelas penuh", ErrCapacityFull)
	}

	active, err := s.enrollmentRepo.GetActiveByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.StudentID)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal cek enrollment aktif", err)
	}
	if active != nil {
		return nil, fmt.Errorf("siswa sudah memiliki enrollment aktif di periode ini")
	}

	studentRepo, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.StudentID)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal ambil data siswa", err)
	}
	if studentRepo == nil {
		return nil, fmt.Errorf("%w: student %s", ErrNotFound, req.StudentID)
	}

	studentData := models.StudentData{
		StudentID: studentRepo.ID,
		SchoolID:  studentRepo.SchoolID,
		NISN:      studentRepo.NISN,
		NIS:       studentRepo.NIS,
		FirstName: studentRepo.Name,
		BirthDate: studentRepo.BirthDate,
		Gender:    studentRepo.Gender,
	}
	domainStudent, err := models.NewStudent(studentData, opCtx)
	if err != nil {
		log.Printf("❌ [Domain Error] %v", err)
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal membuat domain student", err)
	}

	enrollCmd := commands.EnrollStudentCommand{
		StudentID:    req.StudentID, // Wajib ada agar validasi command lolos
		EnrollmentID: uuid.New().String(),
		ClassID:      classID,
	}
	result := s.engine.EnrollStudent(enrollCmd, opCtx, domainStudent)
	if result.Error != nil {
		log.Printf("❌ [Enrollment Engine Error] %v", result.Error)
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal enroll student", result.Error)
	}

	enrollment := &repository.Enrollment{
		ID:        enrollCmd.EnrollmentID,
		StudentID: req.StudentID,
		SchoolID:  opCtx.SchoolID,
		PeriodID:  opCtx.AcademicPeriodID,
		ClassID:   classID,
		Status:    "active",
		Reason:    "enrolled",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.enrollmentRepo.Create(ctx, enrollment); err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", enrollment.ID, "gagal simpan ke database", err)
	}

	event := DomainEvent{
		EventType:        "StudentEnrolled",
		AggregateType:    "Enrollment",
		AggregateID:      enrollment.ID,
		Payload:          map[string]interface{}{"studentId": req.StudentID, "classId": classID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return enrollment, nil
}

func (s *EnrollmentService) RemoveMember(ctx context.Context, classID, studentID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	if classID == "" || studentID == "" {
		return fmt.Errorf("%w: classID dan studentID wajib diisi", ErrInvalidInput)
	}

	active, err := s.enrollmentRepo.GetActiveByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, studentID)
	if err != nil {
		return NewServiceError("RemoveMember", "Enrollment", "", "gagal cek enrollment aktif", err)
	}
	if active == nil {
		return fmt.Errorf("%w: tidak ada enrollment aktif untuk siswa ini", ErrNotFound)
	}
	if active.ClassID != classID {
		return fmt.Errorf("siswa tidak terdaftar di kelas ini")
	}

	// Hapus dari tabel_anggotakelas
	if err := s.enrollmentRepo.Delete(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, active.ID); err != nil {
		return NewServiceError("RemoveMember", "Enrollment", active.ID, "gagal hapus data", err)
	}

	event := DomainEvent{
		EventType:        "StudentRemoved",
		AggregateType:    "Enrollment",
		AggregateID:      active.ID,
		Payload:          map[string]interface{}{"studentId": studentID, "classId": classID},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return nil
}

func (s *EnrollmentService) MoveStudent(ctx context.Context, schoolID, periodID string, req MoveStudentRequest) (*repository.Enrollment, error) {
	if schoolID == "" || periodID == "" || req.StudentID == "" || req.ToClassID == "" {
		return nil, repository.ErrInvalidInput
	}

	current, err := s.enrollmentRepo.GetActiveByStudent(ctx, schoolID, periodID, req.StudentID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, fmt.Errorf("student has no active enrollment in this period")
	}

	fromClass, err := s.classRepo.GetByID(ctx, schoolID, periodID, current.ClassID)
	if err != nil {
		return nil, err
	}
	toClass, err := s.classRepo.GetByID(ctx, schoolID, periodID, req.ToClassID)
	if err != nil {
		return nil, err
	}
	if toClass == nil {
		return nil, repository.ErrNotFound
	}

	if fromClass.Grade != toClass.Grade {
		return nil, repository.ErrInvalidGradeTransition
	}

	count, err := s.enrollmentRepo.CountActiveByClass(ctx, schoolID, periodID, req.ToClassID)
	if err != nil {
		return nil, err
	}
	if toClass.Capacity > 0 && count >= toClass.Capacity {
		return nil, repository.ErrCapacityFull
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	updateQuery := `
		UPDATE tabel_anggotakelas
		SET rombongan_belajar_id = $1
		WHERE anggota_rombel_id = $2 AND semester_id = $3
	`
	_, err = tx.Exec(ctx, updateQuery, req.ToClassID, current.ID, periodID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	current.ClassID = req.ToClassID
	current.Reason = "moved from " + fromClass.Name
	current.UpdatedAt = time.Now()
	return current, nil
}

func (s *EnrollmentService) GetClassMembers(ctx context.Context, classID string) ([]MemberResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	if classID == "" {
		return nil, fmt.Errorf("%w: classID wajib diisi", ErrInvalidInput)
	}

	enrollments, err := s.enrollmentRepo.GetMembersByClass(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, classID)
	if err != nil {
		return nil, NewServiceError("GetClassMembers", "Enrollment", "", "gagal ambil data anggota", err)
	}

	var result []MemberResponse
	for _, e := range enrollments {
		student, err := s.studentRepo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, e.StudentID)
		if err != nil {
			continue
		}
		if student == nil {
			continue
		}
		result = append(result, MemberResponse{
			EnrollmentID: e.ID,
			StudentID:    e.StudentID,
			NIS:          student.NIS,
			Name:         student.Name,
			Status:       e.Status,
			EnrolledAt:   e.CreatedAt,
		})
	}
	return result, nil
}

func (s *EnrollmentService) GetStudentHistory(ctx context.Context, studentID string) ([]repository.Enrollment, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	if studentID == "" {
		return nil, fmt.Errorf("%w: studentID wajib diisi", ErrInvalidInput)
	}

	history, err := s.enrollmentRepo.GetHistoryByStudent(ctx, opCtx.SchoolID, studentID)
	if err != nil {
		return nil, NewServiceError("GetStudentHistory", "Enrollment", "", "gagal ambil histori", err)
	}
	return history, nil
}
