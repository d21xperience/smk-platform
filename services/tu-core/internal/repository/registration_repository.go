package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/registration/models"
	"sekolah-platform/services/tu-core/internal/domain/registration/valueobjects"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRegistrationNotFound = errors.New("registration not found")

type RegistrationFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Search           string
	Limit            int
	Offset           int
}

type RegistrationRepository interface {
	Create(ctx context.Context, r *models.Registration) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, registrationID string) (*models.Registration, error)
	List(ctx context.Context, filter RegistrationFilter) ([]models.Registration, int, error)
	Update(ctx context.Context, r *models.Registration) error
}

type pgxRegistrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) RegistrationRepository {
	return &pgxRegistrationRepository{db: db}
}

func (r *pgxRegistrationRepository) Create(ctx context.Context, reg *models.Registration) error {
	query := `
		INSERT INTO registrations (
			id, school_id, academic_period_id, calon_nisn, calon_nama, calon_jenis_kelamin,
			calon_tempat_lahir, calon_tanggal_lahir, calon_alamat, calon_telepon, calon_email,
			nama_ayah, nama_ibu, telepon_ortu, pekerjaan_ayah, pekerjaan_ibu,
			asal_sekolah, jurusan_dipilih, alasan_memilih, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22)
	`
	_, err := r.db.Exec(ctx, query,
		reg.ID, reg.SchoolID, reg.AcademicPeriodID, reg.CalonNISN, reg.CalonNama, reg.CalonJenisKelamin,
		reg.CalonTempatLahir, reg.CalonTanggalLahir, reg.CalonAlamat, reg.CalonTelepon, reg.CalonEmail,
		reg.NamaAyah, reg.NamaIbu, reg.TeleponOrtu, reg.PekerjaanAyah, reg.PekerjaanIbu,
		reg.AsalSekolah, reg.JurusanDipilih, reg.AlasanMemilih, reg.Status.String(), reg.CreatedAt, reg.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxRegistrationRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, registrationID string) (*models.Registration, error) {
	query := `
		SELECT id, school_id, academic_period_id, calon_nisn, calon_nama, calon_jenis_kelamin,
		       calon_tempat_lahir, calon_tanggal_lahir, calon_alamat, calon_telepon, calon_email,
		       nama_ayah, nama_ibu, telepon_ortu, pekerjaan_ayah, pekerjaan_ibu,
		       asal_sekolah, jurusan_dipilih, alasan_memilih, status,
		       COALESCE(catatan_verifikasi, '') as catatan_verifikasi,
		       COALESCE(verified_by, '') as verified_by,
		       verified_at, created_at, updated_at
		FROM registrations
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`

	var reg models.Registration
	var statusStr string
	var verifiedAt *time.Time // ✅ PERBAIKAN: Gunakan pointer untuk handle NULL

	err := r.db.QueryRow(ctx, query, registrationID, schoolID, academicPeriodID).Scan(
		&reg.ID, &reg.SchoolID, &reg.AcademicPeriodID, &reg.CalonNISN, &reg.CalonNama, &reg.CalonJenisKelamin,
		&reg.CalonTempatLahir, &reg.CalonTanggalLahir, &reg.CalonAlamat, &reg.CalonTelepon, &reg.CalonEmail,
		&reg.NamaAyah, &reg.NamaIbu, &reg.TeleponOrtu, &reg.PekerjaanAyah, &reg.PekerjaanIbu,
		&reg.AsalSekolah, &reg.JurusanDipilih, &reg.AlasanMemilih, &statusStr,
		&reg.CatatanVerifikasi, &reg.VerifiedBy, &verifiedAt, &reg.CreatedAt, &reg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRegistrationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	reg.Status, _ = valueobjects.NewRegistrationStatus(statusStr)

	// ✅ PERBAIKAN: Handle nil pointer untuk verified_at
	if verifiedAt != nil {
		reg.VerifiedAt = *verifiedAt
	} else {
		reg.VerifiedAt = time.Time{} // Zero value
	}

	return &reg, nil
}

func (r *pgxRegistrationRepository) List(ctx context.Context, filter RegistrationFilter) ([]models.Registration, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []interface{}{filter.SchoolID, filter.AcademicPeriodID}
	argIndex := 3

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Search != "" {
		whereClause += fmt.Sprintf(" AND (calon_nama ILIKE $%d OR calon_nisn ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Search+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM registrations %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, school_id, academic_period_id, calon_nisn, calon_nama, calon_jenis_kelamin,
		       calon_tempat_lahir, calon_tanggal_lahir, calon_alamat, calon_telepon, calon_email,
		       nama_ayah, nama_ibu, telepon_ortu, pekerjaan_ayah, pekerjaan_ibu,
		       asal_sekolah, jurusan_dipilih, alasan_memilih, status,
		       COALESCE(catatan_verifikasi, '') as catatan_verifikasi,
		       COALESCE(verified_by, '') as verified_by,
		       verified_at, created_at, updated_at
		FROM registrations %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var registrations []models.Registration
	for rows.Next() {
		var reg models.Registration
		var statusStr string
		var verifiedAt *time.Time // ✅ PERBAIKAN: Gunakan pointer untuk handle NULL

		if err := rows.Scan(
			&reg.ID, &reg.SchoolID, &reg.AcademicPeriodID, &reg.CalonNISN, &reg.CalonNama, &reg.CalonJenisKelamin,
			&reg.CalonTempatLahir, &reg.CalonTanggalLahir, &reg.CalonAlamat, &reg.CalonTelepon, &reg.CalonEmail,
			&reg.NamaAyah, &reg.NamaIbu, &reg.TeleponOrtu, &reg.PekerjaanAyah, &reg.PekerjaanIbu,
			&reg.AsalSekolah, &reg.JurusanDipilih, &reg.AlasanMemilih, &statusStr,
			&reg.CatatanVerifikasi, &reg.VerifiedBy, &verifiedAt, &reg.CreatedAt, &reg.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}

		reg.Status, _ = valueobjects.NewRegistrationStatus(statusStr)

		// ✅ PERBAIKAN: Handle nil pointer untuk verified_at
		if verifiedAt != nil {
			reg.VerifiedAt = *verifiedAt
		} else {
			reg.VerifiedAt = time.Time{} // Zero value
		}

		registrations = append(registrations, reg)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return registrations, total, nil
}

func (r *pgxRegistrationRepository) Update(ctx context.Context, reg *models.Registration) error {
	query := `
		UPDATE registrations
		SET calon_nisn = $1, calon_nama = $2, calon_jenis_kelamin = $3,
		    calon_tempat_lahir = $4, calon_tanggal_lahir = $5, calon_alamat = $6,
		    calon_telepon = $7, calon_email = $8, nama_ayah = $9, nama_ibu = $10,
		    telepon_ortu = $11, pekerjaan_ayah = $12, pekerjaan_ibu = $13,
		    asal_sekolah = $14, jurusan_dipilih = $15, alasan_memilih = $16,
		    status = $17, catatan_verifikasi = $18, verified_by = $19, verified_at = $20, updated_at = $21
		WHERE id = $22 AND school_id = $23 AND academic_period_id = $24
	`

	// ✅ PERBAIKAN: Handle zero time untuk verified_at
	var verifiedAt interface{}
	if reg.VerifiedAt.IsZero() {
		verifiedAt = nil
	} else {
		verifiedAt = reg.VerifiedAt
	}

	result, err := r.db.Exec(ctx, query,
		reg.CalonNISN, reg.CalonNama, reg.CalonJenisKelamin,
		reg.CalonTempatLahir, reg.CalonTanggalLahir, reg.CalonAlamat,
		reg.CalonTelepon, reg.CalonEmail, reg.NamaAyah, reg.NamaIbu,
		reg.TeleponOrtu, reg.PekerjaanAyah, reg.PekerjaanIbu,
		reg.AsalSekolah, reg.JurusanDipilih, reg.AlasanMemilih,
		reg.Status.String(), reg.CatatanVerifikasi, reg.VerifiedBy, verifiedAt, reg.UpdatedAt,
		reg.ID, reg.SchoolID, reg.AcademicPeriodID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrRegistrationNotFound
	}
	return nil
}
