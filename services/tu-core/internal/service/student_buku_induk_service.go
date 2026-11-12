package service

import (
	"context"

	"sekolah-platform/services/tu-core/internal/repository"
)

// StudentBukuIndukService menangani logika Buku Induk Siswa
type StudentBukuIndukService struct {
	repo       repository.StudentBukuIndukRepository
	schoolRepo repository.SchoolRepository
}

func NewStudentBukuIndukService(repo repository.StudentBukuIndukRepository, schoolRepo repository.SchoolRepository) *StudentBukuIndukService {
	return &StudentBukuIndukService{repo: repo, schoolRepo: schoolRepo}
}

func (s *StudentBukuIndukService) ListBukuInduk(ctx context.Context, schoolID, academicPeriodID string) ([]repository.StudentBukuIndukDetail, error) {
	return s.repo.ListBukuInduk(ctx, schoolID, academicPeriodID)
}

func (s *StudentBukuIndukService) GetBukuIndukDetail(ctx context.Context, studentID, schoolID, academicPeriodID string) (*repository.StudentBukuIndukDetail, error) {
	detail, err := s.repo.GetBukuIndukDetail(ctx, studentID, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	calculateStudentCompleteness(detail)
	return detail, nil
}

// calculateStudentCompleteness menghitung skor kelengkapan data siswa (0-100)
// - Identitas (40): Name, NIS, NISN, BirthPlace, BirthDate, Gender, Agama
// - Kontak & Alamat (20): Address, Phone
// - Ortu (20): NamaAyah, NamaIbu, PekerjaanAyah, PekerjaanIbu
// - Pelengkap (10): SekolahAsal, FotoURL
// - Akademik (10): punya minimal 1 nilai
func calculateStudentCompleteness(d *repository.StudentBukuIndukDetail) {
	skor := 0
	// Identitas (40)
	if d.Name != "" {
		skor += 10
	}
	if d.NIS != "" {
		skor += 5
	}
	if d.NISN != "" {
		skor += 5
	}
	if d.BirthPlace != "" {
		skor += 5
	}
	if d.BirthDate != nil {
		skor += 5
	}
	if d.Gender != "" {
		skor += 5
	}
	if d.Agama != "" {
		skor += 5
	}

	// Kontak & Alamat (20)
	if d.Address != "" {
		skor += 10
	}
	if d.Phone != "" {
		skor += 10
	}

	// Ortu (20)
	if d.NamaAyah != "" {
		skor += 5
	}
	if d.NamaIbu != "" {
		skor += 5
	}
	if d.PekerjaanAyah != "" {
		skor += 5
	}
	if d.PekerjaanIbu != "" {
		skor += 5
	}

	// Pelengkap (10)
	if d.SekolahAsal != "" {
		skor += 5
	}
	if d.FotoURL != "" {
		skor += 5
	}

	// Akademik (10)
	if len(d.AcademicScores) > 0 {
		skor += 10
	}

	d.SkorKelengkapan = skor
	switch {
	case skor >= 81:
		d.StatusKelengkapan = "Lengkap"
	case skor >= 51:
		d.StatusKelengkapan = "Cukup"
	default:
		d.StatusKelengkapan = "Belum Lengkap"
	}
}
