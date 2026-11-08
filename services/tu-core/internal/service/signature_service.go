package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sekolah-platform/services/tu-core/internal/domain/signature/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type SignatureService struct {
	repo repository.SignatureRepository
}

func NewSignatureService(repo repository.SignatureRepository) *SignatureService {
	return &SignatureService{repo: repo}
}

func (s *SignatureService) SignCorrespondence(ctx context.Context, correspondenceID string, req SignCorrespondenceRequest) (*SignatureResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Parse Base64 image
	base64Data := req.ImageBase64
	// Hapus prefix jika ada (e.g., "data:image/png;base64,")
	if idx := strings.Index(base64Data, ","); idx != -1 {
		base64Data = base64Data[idx+1:]
	}

	imgBytes, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return nil, NewServiceError("SignCorrespondence", "Signature", "", "invalid base64 image", err)
	}

	// 2. Simpan file PNG ke direktori storage
	saveDir := "./storage/ttd"
	_ = os.MkdirAll(saveDir, os.ModePerm)
	filename := fmt.Sprintf("%s_%s.png", correspondenceID, uuid.New().String()[:8])
	filePath := filepath.Join(saveDir, filename)

	if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
		return nil, NewServiceError("SignCorrespondence", "Signature", "", "failed to save image", err)
	}

	// 3. Generate kode verifikasi unik (8 karakter)
	verificationCode := uuid.New().String()[:8]

	sig := &models.DigitalSignature{
		ID:               uuid.New().String(),
		CorrespondenceID: correspondenceID,
		SignerUserID:     opCtx.UserID,
		SignerRole:       req.SignerRole,
		ImageURL:         filePath,
		VerificationCode: verificationCode,
	}

	if err := s.repo.Create(ctx, sig); err != nil {
		return nil, NewServiceError("SignCorrespondence", "Signature", "", "failed to save to db", err)
	}

	return toSignatureResponse(sig), nil
}

func (s *SignatureService) VerifyCode(ctx context.Context, code string) (*SignatureResponse, error) {
	sig, err := s.repo.GetByVerificationCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return toSignatureResponse(sig), nil
}

// Method internal untuk diakses oleh GeneratePDF di CorrespondenceService
func (s *SignatureService) GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.DigitalSignature, error) {
	return s.repo.GetByCorrespondenceID(ctx, schoolID, correspondenceID)
}

func toSignatureResponse(s *models.DigitalSignature) *SignatureResponse {
	return &SignatureResponse{
		ID:               s.ID,
		CorrespondenceID: s.CorrespondenceID,
		SignerUserID:     s.SignerUserID,
		SignerRole:       s.SignerRole,
		ImageURL:         s.ImageURL,
		VerificationCode: s.VerificationCode,
		SignedAt:         s.SignedAt,
	}
}
