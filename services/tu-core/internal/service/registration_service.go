package service

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/registration/engine"
	"sekolah-platform/services/tu-core/internal/domain/registration/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type RegistrationService struct {
	repo        repository.RegistrationRepository
	studentRepo repository.StudentRepository
	regEngine   *engine.RegistrationEngine
	publisher   EventPublisher
}

func NewRegistrationService(
	repo repository.RegistrationRepository,
	studentRepo repository.StudentRepository,
	regEngine *engine.RegistrationEngine,
	publisher EventPublisher,
) *RegistrationService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &RegistrationService{
		repo:        repo,
		studentRepo: studentRepo,
		regEngine:   regEngine,
		publisher:   publisher,
	}
}

type CreateRegistrationRequest struct {
	CalonNISN         string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir time.Time
	CalonAlamat       string
	CalonTelepon      string
	CalonEmail        string
	NamaAyah          string
	NamaIbu           string
	TeleponOrtu       string
	PekerjaanAyah     string
	PekerjaanIbu      string
	AsalSekolah       string
	JurusanDipilih    string
	AlasanMemilih     string
}

func (s *RegistrationService) CreateRegistration(ctx context.Context, req CreateRegistrationRequest) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	data := models.RegistrationData{
		CalonNISN:         req.CalonNISN,
		CalonNama:         req.CalonNama,
		CalonJenisKelamin: req.CalonJenisKelamin,
		CalonTempatLahir:  req.CalonTempatLahir,
		CalonTanggalLahir: req.CalonTanggalLahir,
		CalonAlamat:       req.CalonAlamat,
		CalonTelepon:      req.CalonTelepon,
		CalonEmail:        req.CalonEmail,
		NamaAyah:          req.NamaAyah,
		NamaIbu:           req.NamaIbu,
		TeleponOrtu:       req.TeleponOrtu,
		PekerjaanAyah:     req.PekerjaanAyah,
		PekerjaanIbu:      req.PekerjaanIbu,
		AsalSekolah:       req.AsalSekolah,
		JurusanDipilih:    req.JurusanDipilih,
		AlasanMemilih:     req.AlasanMemilih,
	}

	result := s.regEngine.CreateRegistration(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateRegistration", "Registration", "", "invalid registration data", result.Error)
	}

	// ✅ TAMBAHKAN LOGGING SEBELUM SAVE
	fmt.Printf("[RegistrationService] 📝 Saving registration: ID=%s, NISN=%s, Nama=%s\n",
		result.Registration.ID, result.Registration.CalonNISN, result.Registration.CalonNama)

	if err := s.repo.Create(ctx, result.Registration); err != nil {
		// ✅ TAMBAHKAN LOGGING ERROR DETAIL
		fmt.Printf("[RegistrationService] ❌ Database error: %v\n", err)
		return nil, NewServiceError("CreateRegistration", "Registration", result.Registration.ID,
			fmt.Sprintf("failed to save to database: %v", err), err)
	}

	fmt.Printf("[RegistrationService] ✅ Registration saved successfully\n")

	event := DomainEvent{
		EventType:     "RegistrationCreated",
		AggregateType: "Registration",
		AggregateID:   result.Registration.ID,
		Payload: map[string]interface{}{
			"registrationId": result.Registration.ID,
			"calonNisn":      result.Registration.CalonNISN,
			"calonNama":      result.Registration.CalonNama,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toRegistrationResponse(result.Registration), nil
}

func (s *RegistrationService) SubmitRegistration(ctx context.Context, registrationID string) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, registrationID)
	if err != nil {
		return nil, NewServiceError("SubmitRegistration", "Registration", registrationID, "failed to get registration", err)
	}

	result := s.regEngine.SubmitRegistration(reg, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("SubmitRegistration", "Registration", registrationID, "failed to submit", result.Error)
	}

	if err := s.repo.Update(ctx, result.Registration); err != nil {
		return nil, NewServiceError("SubmitRegistration", "Registration", registrationID, "failed to update database", err)
	}

	event := DomainEvent{
		EventType:     "RegistrationSubmitted",
		AggregateType: "Registration",
		AggregateID:   result.Registration.ID,
		Payload: map[string]interface{}{
			"registrationId": result.Registration.ID,
			"calonNisn":      result.Registration.CalonNISN,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toRegistrationResponse(result.Registration), nil
}

func (s *RegistrationService) VerifyRegistration(ctx context.Context, registrationID, catatan string) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, registrationID)
	if err != nil {
		return nil, NewServiceError("VerifyRegistration", "Registration", registrationID, "failed to get registration", err)
	}

	result := s.regEngine.VerifyRegistration(reg, catatan, opCtx.UserID, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("VerifyRegistration", "Registration", registrationID, "failed to verify", result.Error)
	}

	if err := s.repo.Update(ctx, result.Registration); err != nil {
		return nil, NewServiceError("VerifyRegistration", "Registration", registrationID, "failed to update database", err)
	}

	event := DomainEvent{
		EventType:     "RegistrationVerified",
		AggregateType: "Registration",
		AggregateID:   result.Registration.ID,
		Payload: map[string]interface{}{
			"registrationId": result.Registration.ID,
			"verifiedBy":     opCtx.UserID,
			"catatan":        catatan,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toRegistrationResponse(result.Registration), nil
}

func (s *RegistrationService) ApproveRegistration(ctx context.Context, registrationID, catatan string) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, registrationID)
	if err != nil {
		return nil, NewServiceError("ApproveRegistration", "Registration", registrationID, "failed to get registration", err)
	}

	result := s.regEngine.ApproveRegistration(reg, catatan, opCtx.UserID, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("ApproveRegistration", "Registration", registrationID, "failed to approve", result.Error)
	}

	if err := s.repo.Update(ctx, result.Registration); err != nil {
		return nil, NewServiceError("ApproveRegistration", "Registration", registrationID, "failed to update database", err)
	}

	// ============================================================
	// ✅ SIDE-EFFECT: Buat Student baru dari Registration yang disetujui
	// ============================================================
	studentID := uuid.New().String()

	// Cek duplikasi NISN (safety check)
	existingStudent, _ := s.studentRepo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, reg.CalonNISN)
	if existingStudent != nil {
		return nil, NewServiceError("ApproveRegistration", "Registration", registrationID,
			"failed to create student: NISN already exists",
			fmt.Errorf("NISN %s sudah terdaftar", reg.CalonNISN))
	}

	// Buat Student object langsung
	newStudent := &repository.Student{
		ID:               studentID,
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		NIS:              reg.CalonNISN, // Gunakan NISN sebagai NIS awal
		NISN:             reg.CalonNISN,
		Name:             reg.CalonNama,
		Gender:           reg.CalonJenisKelamin,
		BirthDate:        reg.CalonTanggalLahir,
		BirthPlace:       reg.CalonTempatLahir,
		Address:          reg.CalonAlamat,
		Phone:            reg.CalonTelepon,
		Email:            reg.CalonEmail,
		ParentName:       reg.NamaAyah,
		ParentPhone:      reg.TeleponOrtu,
		ParentEmail:      "", // Tidak ada di registration
		ClassName:        "", // Belum ada kelas
		Status:           "ACTIVE",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	if err := s.studentRepo.Create(ctx, newStudent); err != nil {
		return nil, NewServiceError("ApproveRegistration", "Registration", registrationID,
			"failed to create student from registration", err)
	}

	// Publish event StudentCreated
	studentEvent := DomainEvent{
		EventType:     "StudentCreated",
		AggregateType: "Student",
		AggregateID:   studentID,
		Payload: map[string]interface{}{
			"studentId":               studentID,
			"registrationId":          registrationID,
			"nisn":                    reg.CalonNISN,
			"name":                    reg.CalonNama,
			"createdFromRegistration": true,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, studentEvent)

	// Publish event RegistrationApproved
	event := DomainEvent{
		EventType:     "RegistrationApproved",
		AggregateType: "Registration",
		AggregateID:   result.Registration.ID,
		Payload: map[string]interface{}{
			"registrationId": result.Registration.ID,
			"studentId":      studentID,
			"approvedBy":     opCtx.UserID,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toRegistrationResponse(result.Registration), nil
}

func (s *RegistrationService) RejectRegistration(ctx context.Context, registrationID, alasan string) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, registrationID)
	if err != nil {
		return nil, NewServiceError("RejectRegistration", "Registration", registrationID, "failed to get registration", err)
	}

	result := s.regEngine.RejectRegistration(reg, alasan, opCtx.UserID, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("RejectRegistration", "Registration", registrationID, "failed to reject", result.Error)
	}

	if err := s.repo.Update(ctx, result.Registration); err != nil {
		return nil, NewServiceError("RejectRegistration", "Registration", registrationID, "failed to update database", err)
	}

	event := DomainEvent{
		EventType:     "RegistrationRejected",
		AggregateType: "Registration",
		AggregateID:   result.Registration.ID,
		Payload: map[string]interface{}{
			"registrationId": result.Registration.ID,
			"alasan":         alasan,
			"rejectedBy":     opCtx.UserID,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toRegistrationResponse(result.Registration), nil
}

func (s *RegistrationService) GetRegistration(ctx context.Context, registrationID string) (*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, registrationID)
	if err != nil {
		return nil, NewServiceError("GetRegistration", "Registration", registrationID, "failed to get registration", err)
	}

	return toRegistrationResponse(reg), nil
}

func (s *RegistrationService) ListRegistrations(ctx context.Context, status, search string) ([]*RegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.RegistrationFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		Status:           status,
		Search:           search,
		Limit:            100,
		Offset:           0,
	}

	registrations, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListRegistrations", "Registration", "", "failed to list registrations", err)
	}

	var responses []*RegistrationResponse
	for _, r := range registrations {
		responses = append(responses, toRegistrationResponse(&r))
	}

	return responses, nil
}

type RegistrationResponse struct {
	ID                string
	SchoolID          string
	AcademicPeriodID  string
	CalonNISN         string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir time.Time
	CalonAlamat       string
	CalonTelepon      string
	CalonEmail        string
	NamaAyah          string
	NamaIbu           string
	TeleponOrtu       string
	PekerjaanAyah     string
	PekerjaanIbu      string
	AsalSekolah       string
	JurusanDipilih    string
	AlasanMemilih     string
	Status            string
	CatatanVerifikasi string
	VerifiedBy        string
	VerifiedAt        time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func toRegistrationResponse(r *models.Registration) *RegistrationResponse {
	if r == nil {
		return nil
	}
	return &RegistrationResponse{
		ID:                r.ID,
		SchoolID:          r.SchoolID,
		AcademicPeriodID:  r.AcademicPeriodID,
		CalonNISN:         r.CalonNISN,
		CalonNama:         r.CalonNama,
		CalonJenisKelamin: r.CalonJenisKelamin,
		CalonTempatLahir:  r.CalonTempatLahir,
		CalonTanggalLahir: r.CalonTanggalLahir,
		CalonAlamat:       r.CalonAlamat,
		CalonTelepon:      r.CalonTelepon,
		CalonEmail:        r.CalonEmail,
		NamaAyah:          r.NamaAyah,
		NamaIbu:           r.NamaIbu,
		TeleponOrtu:       r.TeleponOrtu,
		PekerjaanAyah:     r.PekerjaanAyah,
		PekerjaanIbu:      r.PekerjaanIbu,
		AsalSekolah:       r.AsalSekolah,
		JurusanDipilih:    r.JurusanDipilih,
		AlasanMemilih:     r.AlasanMemilih,
		Status:            r.Status.String(),
		CatatanVerifikasi: r.CatatanVerifikasi,
		VerifiedBy:        r.VerifiedBy,
		VerifiedAt:        r.VerifiedAt,
		CreatedAt:         r.CreatedAt,
		UpdatedAt:         r.UpdatedAt,
	}
}
