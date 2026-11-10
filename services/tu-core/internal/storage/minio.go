package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioClient struct {
	client *minio.Client
}

func NewMinioClient() (*MinioClient, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	accessKey := os.Getenv("MINIO_ACCESS_KEY")
	if accessKey == "" { accessKey = "minioadmin" }
	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" { secretKey = "minioadmin" }
	useSSL := os.Getenv("MINIO_USE_SSL") == "true"

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("gagal inisialisasi minio: %w", err)
	}

	return &MinioClient{client: client}, nil
}

func (m *MinioClient) ensureBucket(ctx context.Context, bucketName string) error {
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil { return err }
	
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil { return err }
	}

	policyJSON := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [{"Effect": "Allow","Principal": {"AWS": ["*"]},"Action": ["s3:GetObject"],"Resource": ["arn:aws:s3:::%s/*"]}]
	}`, bucketName)
	_ = m.client.SetBucketPolicy(ctx, bucketName, policyJSON)
	return nil
}

func (m *MinioClient) UploadFile(ctx context.Context, bucketName, objectName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	if err := m.ensureBucket(ctx, bucketName); err != nil { return "", err }

	contentType := fileHeader.Header.Get("Content-Type")
	_, err := m.client.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil { return "", err }

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" { endpoint = "localhost:9000" }
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName), nil
}

// Fungsi baru untuk upload raw file (seperti hasil generate PDF)
func (m *MinioClient) UploadRawFile(ctx context.Context, bucketName, objectName string, file io.Reader, contentType string) (string, error) {
	if err := m.ensureBucket(ctx, bucketName); err != nil { return "", err }

	_, err := m.client.PutObject(ctx, bucketName, objectName, file, -1, minio.PutObjectOptions{ContentType: contentType})
	if err != nil { return "", err }

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" { endpoint = "localhost:9000" }
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName), nil
}

// Fungsi yang sempat hilang untuk Presigned URL
func (m *MinioClient) GetPresignedURL(ctx context.Context, bucketName, objectName string, expires time.Duration) (string, error) {
	// Pastikan bucket ada
	if err := m.ensureBucket(ctx, bucketName); err != nil { return "", err }
	
	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expires, nil)
	if err != nil {
		return "", err
	}
	return url.String(), nil
}
