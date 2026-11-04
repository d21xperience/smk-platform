package service

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"

	"sekolah-platform/services/tu-core/internal/repository"
)

type AuthService struct {
	userRepo     repository.UserRepository
	academicRepo repository.AcademicRepository
}

// NewAuthService — signature berubah: sekarang butuh academicRepo untuk
// me-resolve academic_period_id aktif saat login.
func NewAuthService(userRepo repository.UserRepository, academicRepo repository.AcademicRepository) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		academicRepo: academicRepo,
	}
}

type LoginResult struct {
	Token            string
	User             *repository.User
	AcademicPeriodID string
	AcademicYear     string
	Semester         int
}

func (s *AuthService) Login(ctx context.Context, username, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, errors.New("username atau password salah")
	}

	// 1. Verifikasi password via bcrypt.
	// 2. Fallback plaintext (khusus development / seed awal).
	if bcryptErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); bcryptErr != nil {
		if user.Password != password {
			return nil, errors.New("username atau password salah")
		}
	}

	// 3. Generate JWT token (dummy, untuk development).
	token := "jwt-token-dummy-" + user.UserID

	result := &LoginResult{
		Token: token,
		User:  user,
	}

	// 4. Resolve academic period aktif berdasarkan sekolah user.
	if user.SchoolID != "" {
		if period, err := s.academicRepo.GetActivePeriod(ctx, user.SchoolID); err == nil && period != nil {
			result.AcademicPeriodID = period.ID
			result.AcademicYear = period.AcademicYear
			result.Semester = period.Semester
		}
	}

	return result, nil
}
