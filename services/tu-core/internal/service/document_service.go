package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/storage"
)

type DocumentService struct {
	repo  repository.DocumentRepository
	minio *storage.MinioClient
}

func NewDocumentService(repo repository.DocumentRepository, minio *storage.MinioClient) *DocumentService {
	return &DocumentService{repo: repo, minio: minio}
}

func (s *DocumentService) UploadAndLink(ctx context.Context, doc *repository.InventoryDocument, entityType, entityID string, file multipart.File, fileHeader *multipart.FileHeader, userID string) error {
	// 1. Hitung checksum
	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return err
	}
	doc.Checksum = hex.EncodeToString(hasher.Sum(nil))
	file.Seek(0, io.SeekStart) // Reset pointer file

	// 2. Upload ke MinIO (Bucket documents)
	objectKey := fmt.Sprintf("inventory/%s/%s/%d_%s", entityType, entityID, time.Now().Unix(), fileHeader.Filename)
	_, err := s.minio.UploadFile(ctx, "documents", objectKey, file, fileHeader)
	if err != nil {
		return errors.New("gagal upload ke storage: " + err.Error())
	}

	doc.BucketName = "documents"
	doc.ObjectKey = objectKey
	doc.OriginalFilename = fileHeader.Filename
	doc.MimeType = fileHeader.Header.Get("Content-Type")
	doc.FileSize = fileHeader.Size
	doc.Status = "UPLOADED"
	doc.UploadedBy = userID

	// 3. Simpan Metadata
	if err := s.repo.CreateDocument(ctx, doc); err != nil {
		return err
	}

	// 4. Buat Relasi
	return s.repo.LinkDocument(ctx, doc.ID, entityType, entityID, "PRIMARY", userID)
}

func (s *DocumentService) GetDocuments(ctx context.Context, entityType, entityID string) ([]*repository.InventoryDocument, error) {
	return s.repo.GetDocumentsByEntity(ctx, entityType, entityID)
}

func (s *DocumentService) VerifyDocument(ctx context.Context, id string, userID string) error {
	return s.repo.UpdateStatus(ctx, id, "VERIFIED", userID, "")
}

func (s *DocumentService) RejectDocument(ctx context.Context, id string, userID string, reason string) error {
	if reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}
	return s.repo.UpdateStatus(ctx, id, "REJECTED", userID, reason)
}

func (s *DocumentService) GetDownloadURL(ctx context.Context, id string) (string, error) {
	doc, err := s.repo.GetDocumentByID(ctx, id)
	if err != nil {
		return "", errors.New("dokumen tidak ditemukan")
	}

	// Jika PUBLIC, bisa langsung return URL (karena settingan bucket saat ini public)
	// Tapi untuk keamanan, kita paksa pakai Presigned URL (expiry 15 menit)
	return s.minio.GetPresignedURL(ctx, doc.BucketName, doc.ObjectKey, 15*time.Minute)
}
