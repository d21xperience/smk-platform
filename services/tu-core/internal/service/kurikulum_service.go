package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"sekolah-platform/services/tu-core/internal/repository"
)

type KurikulumService struct {
	subjectRepo      repository.SubjectRepository
	pembelajaranRepo repository.PembelajaranRepository
	statRepo         repository.KurikulumStatRepository
}

func NewKurikulumService(sr repository.SubjectRepository, pr repository.PembelajaranRepository, str repository.KurikulumStatRepository) *KurikulumService {
	return &KurikulumService{subjectRepo: sr, pembelajaranRepo: pr, statRepo: str}
}

// === MAPEL ===
func (s *KurikulumService) ListSubjects(ctx context.Context) ([]repository.Subject, error) {
	return s.subjectRepo.List(ctx)
}
func (s *KurikulumService) CreateSubject(ctx context.Context, sub *repository.Subject) error {
	return s.subjectRepo.Create(ctx, sub)
}
func (s *KurikulumService) UpdateSubject(ctx context.Context, sub *repository.Subject) error {
	return s.subjectRepo.Update(ctx, sub)
}
func (s *KurikulumService) DeleteSubject(ctx context.Context, id int) error {
	return s.subjectRepo.Delete(ctx, id)
}

// === PEMBELAJARAN ===
func (s *KurikulumService) ListPembelajaran(ctx context.Context, schoolID, periodID, classID string) ([]repository.Pembelajaran, error) {
	return s.pembelajaranRepo.List(ctx, schoolID, periodID, classID)
}

func (s *KurikulumService) CreatePembelajaran(ctx context.Context, classID, subjectID, teacherID, periodID string) error {
	var subID int
	fmt.Sscanf(subjectID, "%d", &subID)
	p := &repository.Pembelajaran{
		ID:         uuid.New().String(),
		ClassID:    classID,
		SubjectID:  subID,
		TeacherID:  teacherID,
		SemesterID: periodID,
	}
	return s.pembelajaranRepo.Create(ctx, p)
}
func (s *KurikulumService) DeletePembelajaran(ctx context.Context, id string) error {
	return s.pembelajaranRepo.Delete(ctx, id)
}

// === STATISTIK ===
func (s *KurikulumService) GetDashboardStats(ctx context.Context, schoolID, periodID string) (*repository.KurikulumStats, error) {
	return s.statRepo.GetStats(ctx, schoolID, periodID)
}
