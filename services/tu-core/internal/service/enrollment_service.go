// internal/service/enrollment_service.go
package service

import (
	"context"
	"errors"
	"fmt"
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

// EnrollmentService mengorchestrasi use case untuk Enrollment / Anggota Rombel.
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

// NewEnrollmentService membuat instance baru dari EnrollmentService.
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

// AddMemberRequest untuk menambahkan siswa ke rombel.
type AddMemberRequest struct {
	StudentID string `json:"studentId"`
}

// MoveStudentRequest untuk pindah rombel.
type MoveStudentRequest struct {
	StudentID string `json:"studentId"`
	ToClassID string `json:"toClassId"`
	Reason    string `json:"reason"`
}

// MemberResponse untuk data anggota rombel.
type MemberResponse struct {
	EnrollmentID string    `json:"enrollmentId"`
	StudentID    string    `json:"studentId"`
	NIS          string    `json:"nis"`
	Name         string    `json:"name"`
	Status       string    `json:"status"`
	EnrolledAt   time.Time `json:"enrolledAt"`
}

// AddMember menambahkan siswa ke rombel.
func (s *EnrollmentService) AddMember(ctx context.Context, classID string, req AddMemberRequest) (*repository.Enrollment, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	if classID == "" || req.StudentID == "" {
		return nil, fmt.Errorf("%w: classID dan studentID wajib diisi", ErrInvalidInput)
	}

	// 1. Cek kapasitas kelas
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
	if count >= class.Capacity {
		return nil, fmt.Errorf("%w: kapasitas kelas penuh", ErrCapacityFull)
	}

	// 2. Cek apakah siswa sudah punya enrollment aktif di periode ini
	active, err := s.enrollmentRepo.GetActiveByStudent(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.StudentID)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal cek enrollment aktif", err)
	}
	if active != nil {
		return nil, fmt.Errorf("siswa sudah memiliki enrollment aktif di periode ini")
	}

	// 3. Panggil Engine StudentEngine.EnrollStudent
	// Kita perlu mengambil data student untuk dibuat domain Student
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
	}
	domainStudent, err := models.NewStudent(studentData, opCtx)
	if err != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal membuat domain student", err)
	}

	enrollCmd := commands.EnrollStudentCommand{
		EnrollmentID: uuid.New().String(),
		ClassID:      classID,
	}
	result := s.engine.EnrollStudent(enrollCmd, opCtx, domainStudent)
	if result.Error != nil {
		return nil, NewServiceError("AddMember", "Enrollment", "", "gagal enroll student", result.Error)
	}

	// 4. Simpan enrollment ke repository
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

	// 5. Event
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

// RemoveMember mengakhiri keanggotaan siswa (update status menjadi completed).
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

	active.Status = "completed"
	active.Reason = "removed by admin"
	active.UpdatedAt = time.Now()
	if err := s.enrollmentRepo.Update(ctx, active); err != nil {
		return NewServiceError("RemoveMember", "Enrollment", active.ID, "gagal update status", err)
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

// MoveStudent memindahkan siswa ke rombel lain (dengan validasi grade & kapasitas).
// MoveStudent memindahkan siswa ke rombel lain (dengan validasi grade & kapasitas)
func (s *EnrollmentService) MoveStudent(ctx context.Context, schoolID, periodID string, req MoveStudentRequest) (*repository.Enrollment, error) {
	if schoolID == "" || periodID == "" || req.StudentID == "" || req.ToClassID == "" {
		return nil, repository.ErrInvalidInput
	}

	// 1. Ambil enrollment aktif siswa di periode ini
	current, err := s.enrollmentRepo.GetActiveByStudent(ctx, schoolID, periodID, req.StudentID)
	if err != nil {
		return nil, err
	}
	if current == nil {
		return nil, errors.New("student has no active enrollment in this period")
	}

	// 2. Ambil class asal dan tujuan
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

	// 3. Validasi grade (periode sama → grade harus sama)
	if fromClass.Grade != toClass.Grade {
		return nil, repository.ErrInvalidGradeTransition
	}

	// 4. Validasi kapasitas kelas tujuan
	count, err := s.enrollmentRepo.CountActiveByClass(ctx, schoolID, periodID, req.ToClassID)
	if err != nil {
		return nil, err
	}
	if count >= toClass.Capacity {
		return nil, repository.ErrCapacityFull
	}

	// 5. MULAI TRANSKAI MANUAL
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Kita perlu membuat context terpisah dengan transaksi, tetapi repository kita
	// tidak menerima parameter *pgxpool.Conn. Karena itu, kita harus melakukan query
	// UPDATE dan INSERT secara langsung di sini (bukan via repository).
	//
	// ATAU lebih baik: buat method di repository yang menerima *pgx.Tx.
	// Untuk sementara, saya akan tulis query langsung di service.

	// Update enrollment lama (langsung query)
	updateQuery := `
		UPDATE enrollments
		SET status = $1, reason = $2, updated_at = NOW()
		WHERE id = $3 AND school_id = $4 AND period_id = $5
	`
	_, err = tx.Exec(ctx, updateQuery, "transferred", req.Reason, current.ID, schoolID, periodID)
	if err != nil {
		return nil, err
	}

	// Buat enrollment baru (langsung query)
	newEnrollment := &repository.Enrollment{
		ID:        generateID(),
		StudentID: req.StudentID,
		SchoolID:  schoolID,
		PeriodID:  periodID,
		ClassID:   req.ToClassID,
		Status:    "active",
		Reason:    "moved from " + fromClass.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	insertQuery := `
		INSERT INTO enrollments (id, student_id, school_id, period_id, class_id, status, reason)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err = tx.Exec(ctx, insertQuery,
		newEnrollment.ID, newEnrollment.StudentID, newEnrollment.SchoolID,
		newEnrollment.PeriodID, newEnrollment.ClassID, newEnrollment.Status, newEnrollment.Reason,
	)
	if err != nil {
		return nil, err
	}

	// Commit transaksi
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return newEnrollment, nil
}

// GetClassMembers mengambil daftar anggota aktif suatu rombel (dengan data siswa).
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
			continue // log error
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

// GetStudentHistory mengambil histori enrollment seorang siswa.
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
