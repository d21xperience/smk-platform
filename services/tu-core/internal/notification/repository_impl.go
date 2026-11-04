package notification

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/notification/handlers"
)

// NotificationRepositoryImpl mengimplementasikan semua interface repository
// yang dibutuhkan handlers notifikasi, menembak skema DAPODIK.
type NotificationRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewNotificationRepositoryImpl(db *pgxpool.Pool) *NotificationRepositoryImpl {
	return &NotificationRepositoryImpl{db: db}
}

// GetStudentContactInfo mengambil informasi kontak siswa + orang tua/wali
// dari tabel_siswa dan tabel_siswa_pelengkap.
//
// Pemetaan (DAPODIK):
//   - siswa.peserta_didik_id  -> StudentID
//   - siswa.nm_siswa          -> StudentName
//   - siswa.nm_ibu/nm_ayah/nm_wali -> ParentName (prioritas ibu)
//   - pelengkap.telepon_ortu/telepon_wali/siswa.telepon_siswa -> ParentPhone
//   - email ortu: TIDAK ADA di DAPODIK -> string kosong
func (r *NotificationRepositoryImpl) GetStudentContactInfo(ctx context.Context, studentID string) (*handlers.StudentContactInfo, error) {
	query := `
		SELECT
			s.peserta_didik_id,
			s.nm_siswa,
			COALESCE(NULLIF(s.nm_ibu, ''), NULLIF(s.nm_ayah, ''), NULLIF(s.nm_wali, ''), '')  AS parent_name,
			COALESCE(NULLIF(sp.telepon_ortu, ''), NULLIF(sp.telepon_wali, ''), NULLIF(s.telepon_siswa, ''), '') AS parent_phone,
			''::text AS parent_email
		FROM tabel_siswa s
		LEFT JOIN tabel_siswa_pelengkap sp
			ON sp.peserta_didik_id = s.peserta_didik_id
		WHERE s.peserta_didik_id = $1
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("student %s tidak ditemukan di tabel_siswa: %w", studentID, err)
		}
		return nil, fmt.Errorf("gagal fetch contact info untuk student %s: %w", studentID, err)
	}

	return &info, nil
}

// GetClassName mengambil nama rombel (nm_kelas) dari tabel_kelas.
func (r *NotificationRepositoryImpl) GetClassName(ctx context.Context, classID string) (string, error) {
	query := `
		SELECT COALESCE(nm_kelas, '')
		FROM tabel_kelas
		WHERE rombongan_belajar_id = $1
	`

	var className string
	err := r.db.QueryRow(ctx, query, classID).Scan(&className)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", fmt.Errorf("class %s tidak ditemukan di tabel_kelas: %w", classID, err)
		}
		return "", fmt.Errorf("gagal fetch class name untuk class %s: %w", classID, err)
	}

	return className, nil
}
