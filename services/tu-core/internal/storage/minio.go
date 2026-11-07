package storage

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"

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
	if accessKey == "" {
		accessKey = "minioadmin"
	}

	secretKey := os.Getenv("MINIO_SECRET_KEY")
	if secretKey == "" {
		secretKey = "minioadmin"
	}

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

func (m *MinioClient) UploadFile(ctx context.Context, bucketName, objectName string, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// Pastikan bucket ada
	exists, err := m.client.BucketExists(ctx, bucketName)
	if err != nil {
		return "", err
	}
	
	if !exists {
		err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return "", err
		}
	}

	// FIX: Set bucket policy menjadi public read agar browser bisa mengakses foto
	policyJSON := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, bucketName)
	_ = m.client.SetBucketPolicy(ctx, bucketName, policyJSON) // Abaikan error jika policy sudah ada

	contentType := fileHeader.Header.Get("Content-Type")
	_, err = m.client.PutObject(ctx, bucketName, objectName, file, fileHeader.Size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", err
	}

	endpoint := os.Getenv("MINIO_ENDPOINT")
	if endpoint == "" {
		endpoint = "localhost:9000"
	}
	return fmt.Sprintf("http://%s/%s/%s", endpoint, bucketName, objectName), nil
}
