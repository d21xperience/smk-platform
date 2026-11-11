package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRegistrationNotFound = errors.New("registration not found")

type Registration struct {
	ID                string
	SchoolID          string
	AcademicPeriodID  string
	CalonNISN         string
	NIS               string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir *time.Time
	AsalSekolah       string
	JurusanDipilih    string
	BiayaStandar      float64
	BiayaPendaftaran  float64
	BerkasIjazah      bool
	BerkasKK          bool
	BerkasAkta        bool
	Status            string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type RegistrationFilter struct {
	SchoolID         string
	AcademicPeriodID string
	Status           string
	Search           string
	Limit            int
	Offset           int
}

type RegistrationRepository interface {
	Create(ctx context.Context, reg *Registration) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, regID string) (*Registration, error)
	List(ctx context.Context, filter RegistrationFilter) ([]Registration, int, error)
	Update(ctx context.Context, reg *Registration) error
	GenerateNISMassal(ctx context.Context, schoolID, academicPeriodID, prefix string) (int, error)
}

type pgxRegistrationRepository struct {
	db *pgxpool.Pool
}

func NewRegistrationRepository(db *pgxpool.Pool) RegistrationRepository {
	return &pgxRegistrationRepository{db: db}
}

func (r *pgxRegistrationRepository) Create(ctx context.Context, reg *Registration) error {
	query := `
		INSERT INTO registrations (
			id, school_id, academic_period_id, calon_nisn, nis, calon_nama, calon_jenis_kelamin,
			calon_tempat_lahir, calon_tanggal_lahir, asal_sekolah, jurusan_dipilih,
			biaya_standar, biaya_pendaftaran, berkas_ijazah, berkas_kk, berkas_akta,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
	`
	_, err := r.db.Exec(ctx, query,
		reg.ID, reg.SchoolID, reg.AcademicPeriodID, reg.CalonNISN, reg.NIS, reg.CalonNama, reg.CalonJenisKelamin,
		reg.CalonTempatLahir, reg.CalonTanggalLahir, reg.AsalSekolah, reg.JurusanDipilih,
		reg.BiayaStandar, reg.BiayaPendaftaran, reg.BerkasIjazah, reg.BerkasKK, reg.BerkasAkta,
		reg.Status, reg.CreatedAt, reg.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxRegistrationRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, regID string) (*Registration, error) {
	query := `
		SELECT 
			id, school_id, academic_period_id, calon_nisn, COALESCE(nis, ''), calon_nama, calon_jenis_kelamin,
			COALESCE(calon_tempat_lahir, ''), calon_tanggal_lahir, COALESCE(asal_sekolah, ''), jurusan_dipilih,
			biaya_standar, biaya_pendaftaran, berkas_ijazah, berkas_kk, berkas_akta,
			status, created_at, updated_at
		FROM registrations
		WHERE id = $1 AND school_id = $2 AND academic_period_id = $3
	`
	var reg Registration
	err := r.db.QueryRow(ctx, query, regID, schoolID, academicPeriodID).Scan(
		&reg.ID, &reg.SchoolID, &reg.AcademicPeriodID, &reg.CalonNISN, &reg.NIS, &reg.CalonNama, &reg.CalonJenisKelamin,
		&reg.CalonTempatLahir, &reg.CalonTanggalLahir, &reg.AsalSekolah, &reg.JurusanDipilih,
		&reg.BiayaStandar, &reg.BiayaPendaftaran, &reg.BerkasIjazah, &reg.BerkasKK, &reg.BerkasAkta,
		&reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRegistrationNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &reg, nil
}

func (r *pgxRegistrationRepository) List(ctx context.Context, filter RegistrationFilter) ([]Registration, int, error) {
	whereClause := "WHERE school_id = $1 AND academic_period_id = $2"
	args := []any{filter.SchoolID, filter.AcademicPeriodID}
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

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT 
			id, school_id, academic_period_id, calon_nisn, COALESCE(nis, ''), calon_nama, calon_jenis_kelamin,
			COALESCE(calon_tempat_lahir, ''), calon_tanggal_lahir, COALESCE(asal_sekolah, ''), jurusan_dipilih,
			biaya_standar, biaya_pendaftaran, berkas_ijazah, berkas_kk, berkas_akta,
			status, created_at, updated_at
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

	var registrations []Registration
	for rows.Next() {
		var reg Registration
		if err := rows.Scan(
			&reg.ID, &reg.SchoolID, &reg.AcademicPeriodID, &reg.CalonNISN, &reg.NIS, &reg.CalonNama, &reg.CalonJenisKelamin,
			&reg.CalonTempatLahir, &reg.CalonTanggalLahir, &reg.AsalSekolah, &reg.JurusanDipilih,
			&reg.BiayaStandar, &reg.BiayaPendaftaran, &reg.BerkasIjazah, &reg.BerkasKK, &reg.BerkasAkta,
			&reg.Status, &reg.CreatedAt, &reg.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		registrations = append(registrations, reg)
	}
	return registrations, total, rows.Err()
}

func (r *pgxRegistrationRepository) Update(ctx context.Context, reg *Registration) error {
	query := `
		UPDATE registrations SET
			calon_nisn = $1, nis = $2, calon_nama = $3, calon_jenis_kelamin = $4,
			calon_tempat_lahir = $5, calon_tanggal_lahir = $6, asal_sekolah = $7,
			jurusan_dipilih = $8, biaya_standar = $9, biaya_pendaftaran = $10,
			berkas_ijazah = $11, berkas_kk = $12, berkas_akta = $13,
			status = $14, updated_at = $15
		WHERE id = $16 AND school_id = $17 AND academic_period_id = $18
	`
	result, err := r.db.Exec(ctx, query,
		reg.CalonNISN, reg.NIS, reg.CalonNama, reg.CalonJenisKelamin,
		reg.CalonTempatLahir, reg.CalonTanggalLahir, reg.AsalSekolah,
		reg.JurusanDipilih, reg.BiayaStandar, reg.BiayaPendaftaran,
		reg.BerkasIjazah, reg.BerkasKK, reg.BerkasAkta,
		reg.Status, reg.UpdatedAt,
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

// GenerateNISMassal membuat NIS secara berurutan untuk semua siswa yang belum memiliki NIS.
// Menggunakan Row Locking (FOR UPDATE) pada nis_sequences agar aman dari race condition.
func (r *pgxRegistrationRepository) GenerateNISMassal(ctx context.Context, schoolID, academicPeriodID, prefix string) (int, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	// 1. Dapatkan atau buat sequence untuk periode & prefix ini
	var currentSeq int
	querySeq := `SELECT last_sequence FROM nis_sequences WHERE academic_period_id = $1 AND prefix = $2 FOR UPDATE`
	errSeq := tx.QueryRow(ctx, querySeq, academicPeriodID, prefix).Scan(&currentSeq)
	
	if errors.Is(errSeq, pgx.ErrNoRows) {
		// Jika belum ada, buat baru dengan sequence 0
		_, errInsertSeq := tx.Exec(ctx, `INSERT INTO nis_sequences (academic_period_id, prefix, last_sequence) VALUES ($1, $2, 0)`, academicPeriodID, prefix)
		if errInsertSeq != nil {
			return 0, fmt.Errorf("%w: %v", ErrDatabaseError, errInsertSeq)
		}
		currentSeq = 0
	} else if errSeq != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, errSeq)
	}

	// 2. Ambil semua siswa yang belum punya NIS
	queryGet := `
		SELECT id FROM registrations 
		WHERE school_id = $1 AND academic_period_id = $2 AND (nis IS NULL OR nis = '')
		ORDER BY created_at ASC
	`
	rows, err := tx.Query(ctx, queryGet, schoolID, academicPeriodID)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	
	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		ids = append(ids, id)
	}
	rows.Close()

	if len(ids) == 0 {
		return 0, nil // Tidak ada yang perlu di-generate
	}

	// 3. Update NIS untuk setiap siswa
	updatedCount := 0
	for _, id := range ids {
		currentSeq++
		newNIS := fmt.Sprintf("%s%05d", prefix, currentSeq) // Contoh: 2627 + 00001 = 262700001
		
		_, errUpdate := tx.Exec(ctx, `UPDATE registrations SET nis = $1, updated_at = NOW() WHERE id = $2`, newNIS, id)
		if errUpdate != nil {
			return 0, fmt.Errorf("%w: %v", ErrDatabaseError, errUpdate)
		}
		updatedCount++
	}

	// 4. Update last_sequence di tabel nis_sequences
	_, errUpdateSeq := tx.Exec(ctx, `UPDATE nis_sequences SET last_sequence = $1 WHERE academic_period_id = $2 AND prefix = $3`, currentSeq, academicPeriodID, prefix)
	if errUpdateSeq != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, errUpdateSeq)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return updatedCount, nil
}
