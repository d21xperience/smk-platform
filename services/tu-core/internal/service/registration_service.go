package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
)

type RegistrationService struct {
	repo        repository.RegistrationRepository
	studentRepo repository.StudentRepository
	publisher   EventPublisher
}

func NewRegistrationService(
	repo repository.RegistrationRepository,
	studentRepo repository.StudentRepository,
	publisher EventPublisher,
) *RegistrationService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &RegistrationService{
		repo:        repo,
		studentRepo: studentRepo,
		publisher:   publisher,
	}
}

type CreateRegistrationRequest struct {
	CalonNISN         string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir *time.Time
	AsalSekolah       string
	JurusanDipilih    string
	BiayaStandar      float64
	BiayaPendaftaran  float64
	BerkasIjazah      bool
	BerkasKK          bool
	BerkasAkta        bool
}

func (s *RegistrationService) CreateRegistration(ctx context.Context, req CreateRegistrationRequest) (*repository.Registration, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	reg := &repository.Registration{
		ID:                uuid.New().String(),
		SchoolID:          opCtx.SchoolID,
		AcademicPeriodID:  opCtx.AcademicPeriodID,
		CalonNISN:         req.CalonNISN,
		NIS:               "",
		CalonNama:         req.CalonNama,
		CalonJenisKelamin: req.CalonJenisKelamin,
		CalonTempatLahir:  req.CalonTempatLahir,
		CalonTanggalLahir: req.CalonTanggalLahir,
		AsalSekolah:       req.AsalSekolah,
		JurusanDipilih:    req.JurusanDipilih,
		BiayaStandar:      req.BiayaStandar,
		BiayaPendaftaran:  req.BiayaPendaftaran,
		BerkasIjazah:      req.BerkasIjazah,
		BerkasKK:          req.BerkasKK,
		BerkasAkta:        req.BerkasAkta,
		Status:            "REGISTERED",
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, reg); err != nil {
		return nil, NewServiceError("CreateRegistration", "Registration", reg.ID, "gagal simpan ke database", err)
	}

	return reg, nil
}

func (s *RegistrationService) ListRegistrations(ctx context.Context, status, search string) ([]repository.Registration, error) {
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

	regs, _, err := s.repo.List(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListRegistrations", "Registration", "", "gagal ambil daftar", err)
	}
	return regs, nil
}

// GetRegistrationByID (PERBAIKAN: Fungsi baru untuk memenuhi panggilan gRPC)
func (s *RegistrationService) GetRegistrationByID(ctx context.Context, regID string) (*repository.Registration, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, regID)
}

func (s *RegistrationService) GenerateNISMassal(ctx context.Context, prefix string) (int, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return 0, ErrMissingOperationalContext
	}

	count, err := s.repo.GenerateNISMassal(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, prefix)
	if err != nil {
		return 0, NewServiceError("GenerateNISMassal", "Registration", "", "gagal generate NIS massal", err)
	}
	return count, nil
}

func (s *RegistrationService) ActivateRegistration(ctx context.Context, regID string) error {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return ErrMissingOperationalContext
	}

	reg, err := s.repo.GetByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, regID)
	if err != nil {
		return NewServiceError("ActivateRegistration", "Registration", regID, "registrasi tidak ditemukan", err)
	}

	if reg.NIS == "" {
		return NewServiceError("ActivateRegistration", "Registration", regID, "NIS belum di-generate", nil)
	}
	if reg.BiayaPendaftaran < reg.BiayaStandar {
		return NewServiceError("ActivateRegistration", "Registration", regID, "biaya pendaftaran belum lunas", nil)
	}
	if !reg.BerkasIjazah || !reg.BerkasKK || !reg.BerkasAkta {
		return NewServiceError("ActivateRegistration", "Registration", regID, "berkas fisik belum lengkap", nil)
	}

	existing, _ := s.studentRepo.GetByNIS(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, reg.NIS)
	if existing != nil {
		return NewServiceError("ActivateRegistration", "Registration", regID, "NIS sudah terdaftar di tabel siswa", nil)
	}

	newStudent := &repository.Student{
		ID:               uuid.New().String(),
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		NIS:              reg.NIS,
		NISN:             reg.CalonNISN,
		Name:             reg.CalonNama,
		Gender:           reg.CalonJenisKelamin,
		BirthDate:        time.Time{},
		BirthPlace:       reg.CalonTempatLahir,
		Address:          "",
		Phone:            "",
		Email:            "",
		ParentName:       "",
		ParentPhone:      "",
		ParentEmail:      "",
		ClassName:        "",
		Status:           "ACTIVE",
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}
	if reg.CalonTanggalLahir != nil {
		newStudent.BirthDate = *reg.CalonTanggalLahir
	}

	if err := s.studentRepo.Create(ctx, newStudent); err != nil {
		return NewServiceError("ActivateRegistration", "Registration", regID, "gagal memindahkan ke siswa aktif", err)
	}

	reg.Status = "ACTIVE"
	reg.UpdatedAt = time.Now().UTC()
	if err := s.repo.Update(ctx, reg); err != nil {
		return NewServiceError("ActivateRegistration", "Registration", regID, "gagal update status registrasi", err)
	}

	return nil
}
