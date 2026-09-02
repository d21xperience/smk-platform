package notification

import (
	"context"
	"fmt"

	"sekolah-platform/services/tu-core/internal/notification/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
)

// NotificationRepositoryImpl mengimplementasikan semua interface repository yang dibutuhkan handlers.
type NotificationRepositoryImpl struct {
	db *pgxpool.Pool
}

// NewNotificationRepositoryImpl membuat instance baru.
func NewNotificationRepositoryImpl(db *pgxpool.Pool) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{db: db}
}

// GetStudentContactInfo mengambil informasi kontak siswa dan orang tua.
func (r *NotificationRepositoryImpl) GetStudentContactInfo(ctx context.Context, studentID string) (*handlers.StudentContactInfo, error) {
	query := `
		SELECT
			s.id,
			s.name,
			COALESCE(s.parent_name, ''),
			COALESCE(s.parent_phone, ''),
			COALESCE(s.parent_email, '')
		FROM students s
		WHERE s.id = $1 AND s.deleted_at IS NULL
	`

	var info handlers.StudentContactInfo
	err := r.db.QueryRow(ctx, query, studentID).Scan(
		&info.StudentID,
		&info.StudentName,
		&info.ParentName,
		&info.ParentPhone,
		&info.ParentEmail,
	)

	if err != nil {
		return nil, fmt.Errorf("gagal fetch contact info untuk student %s: %w", studentID, err)
	}

	return &info, nil
}

// GetClassName mengambil nama kelas berdasarkan classID.
func (r *NotificationRepositoryImpl) GetClassName(ctx context.Context, classID string) (string, error) {
	query := `
		SELECT name
		FROM classes
		WHERE id = $1 AND deleted_at IS NULL
	`

	var className string
	err := r.db.QueryRow(ctx, query, classID).Scan(&className)
	if err != nil {
		return "", fmt.Errorf("gagal fetch class name untuk class %s: %w", classID, err)
	}

	return className, nil
}
