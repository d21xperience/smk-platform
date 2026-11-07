package service

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type TeacherService struct {
	repo repository.TeacherRepository
}

func NewTeacherService(repo repository.TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

type TeacherRequest struct {
	Name            string `json:"name"`
	NIP             string `json:"nip"`
	NUPTK           string `json:"nuptk"`
	Gender          string `json:"gender"`
	BirthPlace      string `json:"birthPlace"`
	BirthDate       string `json:"birthDate"` // Format: YYYY-MM-DD
	Address         string `json:"address"`
	GelarDepan      string `json:"gelarDepan"`
	GelarBelakang   string `json:"gelarBelakang"`
	JenisPtkID      int    `json:"jenisPtkId"`
	StatusKeaktifan int    `json:"statusKeaktifan"`
}

func (s *TeacherService) ListTeachers(ctx context.Context, schoolID, periodID, keyword string) ([]repository.Teacher, int, error) {
	filter := repository.TeacherFilter{
		SchoolID: schoolID,
		PeriodID: periodID,
		Keyword:  keyword,
		Limit:    100,
	}
	return s.repo.List(ctx, filter)
}

func (s *TeacherService) CreateTeacher(ctx context.Context, req TeacherRequest) error {
	t := &repository.Teacher{
		ID:              uuid.New().String(),
		Name:            req.Name,
		NIP:             req.NIP,
		NUPTK:           req.NUPTK,
		Gender:          req.Gender,
		BirthPlace:      req.BirthPlace,
		Address:         req.Address,
		GelarDepan:      req.GelarDepan,
		GelarBelakang:   req.GelarBelakang,
		JenisPtkID:      req.JenisPtkID,
		StatusKeaktifan: req.StatusKeaktifan,
	}
	if req.BirthDate != "" {
		parsedTime, err := parseDate(req.BirthDate)
		if err == nil {
			t.BirthDate = parsedTime
		}
	}
	return s.repo.Create(ctx, t)
}

func (s *TeacherService) UpdateTeacher(ctx context.Context, id string, req TeacherRequest) error {
	t := &repository.Teacher{
		ID:              id,
		Name:            req.Name,
		NIP:             req.NIP,
		NUPTK:           req.NUPTK,
		Gender:          req.Gender,
		BirthPlace:      req.BirthPlace,
		Address:         req.Address,
		GelarDepan:      req.GelarDepan,
		GelarBelakang:   req.GelarBelakang,
		JenisPtkID:      req.JenisPtkID,
		StatusKeaktifan: req.StatusKeaktifan,
	}
	if req.BirthDate != "" {
		parsedTime, err := parseDate(req.BirthDate)
		if err == nil {
			t.BirthDate = parsedTime
		}
	}
	return s.repo.Update(ctx, t)
}

func (s *TeacherService) DeleteTeacher(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("%w: ID wajib diisi", ErrInvalidInput)
	}
	return s.repo.Delete(ctx, id)
}

func parseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}
