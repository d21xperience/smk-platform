package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Student struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	NIS              string
	NISN             string
	Name             string
	Gender           string
	BirthDate        time.Time
	BirthPlace       string
	Address          string
	Phone            string
	Email            string
	ParentName       string
	ParentPhone      string
	ParentEmail      string
	ClassName        string
	Status           string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	DeletedAt        *time.Time
}

type StudentFilter struct {
	SchoolID         string
	AcademicPeriodID string
	ClassID          string
	Status           string
	Keyword          string
	Limit            int
	Offset           int
}

type StudentUIDTO struct {
	ID           string     `json:"id"`
	NIS          string     `json:"nis"`
	NISN         string     `json:"nisn"`
	Name         string     `json:"name"`
	BirthPlace   *string    `json:"birthPlace,omitempty"`
	BirthDate    *time.Time `json:"birthDate,omitempty"`
	Gender       *string    `json:"gender,omitempty"`
	Agama        *string    `json:"agama,omitempty"`
	Address      *string    `json:"address,omitempty"`
	Phone        *string    `json:"phone,omitempty"`
	ParentName   *string    `json:"parentName,omitempty"`
	ParentNameIbu *string   `json:"nm_ibu,omitempty"`
	ClassName    string     `json:"className"`
	Status       string     `json:"status"`
}

type StudentRepository interface {
	Create(ctx context.Context, student *Student) error
	GetByID(ctx context.Context, schoolID, academicPeriodID, studentID string) (*Student, error)
	GetByNIS(ctx context.Context, schoolID, academicPeriodID, nis string) (*Student, error)
	GetByNISN(ctx context.Context, schoolID, academicPeriodID, nisn string) (*Student, error)
	List(ctx context.Context, filter StudentFilter) ([]Student, int, error)
	Update(ctx context.Context, student *Student) error
	Delete(ctx context.Context, schoolID, academicPeriodID, studentID string) error
	CountByClass(ctx context.Context, schoolID, academicPeriodID, classID string) (int, error)

	// Area UI
	ListForUI(ctx context.Context, filter StudentFilter) ([]StudentUIDTO, int, error)
	GetByIDForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) (*StudentUIDTO, error)
	UpdateForUI(ctx context.Context, schoolID, academicPeriodID, studentID string, name, birthPlace, gender, address, phone, className string) error
	DeleteForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) error
}

type pgxStudentRepository struct {
	db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) StudentRepository {
	return &pgxStudentRepository{db: db}
}

func (r *pgxStudentRepository) Create(ctx context.Context, student *Student) error {
	if student.SchoolID == "" || student.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	querySiswa := `
		INSERT INTO tabel_siswa (
			peserta_didik_id, nis, nisn, nm_siswa, tempat_lahir, tanggal_lahir,
			jenis_kelamin, alamat_siswa, telepon_siswa, nm_ayah
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	var birthDate any = nil
	if !student.BirthDate.IsZero() {
		birthDate = student.BirthDate
	}

	_, err = tx.Exec(ctx, querySiswa,
		student.ID, student.NIS, student.NISN, student.Name, student.BirthPlace,
		birthDate, student.Gender, student.Address, student.Phone, student.ParentName,
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	queryPelengkap := `
		INSERT INTO tabel_siswa_pelengkap (
			pelengkap_siswa_id, peserta_didik_id, telepon_ortu
		) VALUES ($1, $2, $3)
	`
	_, err = tx.Exec(ctx, queryPelengkap, uuid.New().String(), student.ID, student.ParentPhone)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if student.ClassName != "" {
		var rombelID string
		queryFindRombel := `
			SELECT rombongan_belajar_id::text
			FROM tabel_kelas
			WHERE (nm_kelas = $1 OR rombongan_belajar_id::text = $1)
			  AND semester_id = $2
			  AND sekolah_id::text = $3
			LIMIT 1
		`
		errFind := tx.QueryRow(ctx, queryFindRombel, student.ClassName, student.AcademicPeriodID, student.SchoolID).Scan(&rombelID)

		if errFind == nil {
			queryAnggota := `
				INSERT INTO tabel_anggotakelas (
					anggota_rombel_id, peserta_didik_id, rombongan_belajar_id, semester_id
				) VALUES ($1, $2, $3, $4)
			`
			_, err = tx.Exec(ctx, queryAnggota, uuid.New().String(), student.ID, rombelID, student.AcademicPeriodID)
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDatabaseError, err)
			}
		} else {
			fmt.Printf("[Repo] ⚠️ Peringatan: Kelas '%s' tidak ditemukan di sekolah %s. Siswa disimpan tanpa kelas.\n", student.ClassName, student.SchoolID)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxStudentRepository) GetByNIS(ctx context.Context, schoolID, academicPeriodID, nis string) (*Student, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}
	query := `
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.alamat_siswa, ts.telepon_siswa
		FROM tabel_siswa ts
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE ts.nis = $1 AND tak.semester_id = $2 AND tk.sekolah_id::text = $3
	`
	var s Student
	var birthDate *time.Time
	err := r.db.QueryRow(ctx, query, nis, academicPeriodID, schoolID).Scan(
		&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &birthDate,
		&s.Gender, &s.Address, &s.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if birthDate != nil { s.BirthDate = *birthDate }
	s.SchoolID = schoolID
	s.AcademicPeriodID = academicPeriodID
	s.Status = "Aktif"
	return &s, nil
}

func (r *pgxStudentRepository) GetByNISN(ctx context.Context, schoolID, academicPeriodID, nisn string) (*Student, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}
	query := `
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.alamat_siswa, ts.telepon_siswa
		FROM tabel_siswa ts
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE ts.nisn = $1 AND tak.semester_id = $2 AND tk.sekolah_id::text = $3
	`
	var s Student
	var birthDate *time.Time
	err := r.db.QueryRow(ctx, query, nisn, academicPeriodID, schoolID).Scan(
		&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &birthDate,
		&s.Gender, &s.Address, &s.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if birthDate != nil { s.BirthDate = *birthDate }
	s.SchoolID = schoolID
	s.AcademicPeriodID = academicPeriodID
	s.Status = "Aktif"
	return &s, nil
}

// GetByID mengambil profil siswa murni dari tabel_siswa.
// Menggunakan LEFT JOIN agar siswa tetap ditemukan meskipun belum terdaftar di kelas manapun.
func (r *pgxStudentRepository) GetByID(ctx context.Context, schoolID, academicPeriodID, studentID string) (*Student, error) {
	if schoolID == "" || academicPeriodID == "" {
		return nil, ErrMissingOperationalContext
	}
	query := `
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.alamat_siswa, ts.telepon_siswa,
			COALESCE(ts.nm_ayah, '') as nm_ayah, COALESCE(tsp.telepon_ortu, '') as telepon_ortu,
			COALESCE(tk.nm_kelas, '') as nm_kelas
		FROM tabel_siswa ts
		LEFT JOIN tabel_siswa_pelengkap tsp ON ts.peserta_didik_id = tsp.peserta_didik_id
		LEFT JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id AND tak.semester_id = $2
		LEFT JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE ts.peserta_didik_id = $1
	`
	var s Student
	var birthDate *time.Time
	err := r.db.QueryRow(ctx, query, studentID, academicPeriodID).Scan(
		&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &birthDate,
		&s.Gender, &s.Address, &s.Phone, &s.ParentName, &s.ParentPhone, &s.ClassName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if birthDate != nil { s.BirthDate = *birthDate }
	s.SchoolID = schoolID
	s.AcademicPeriodID = academicPeriodID
	s.Status = "Aktif"
	return &s, nil
}

func (r *pgxStudentRepository) List(ctx context.Context, filter StudentFilter) ([]Student, int, error) {
	if filter.SchoolID == "" || filter.AcademicPeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE tak.semester_id = $1 AND tk.sekolah_id::text = $2"
	args := []any{filter.AcademicPeriodID, filter.SchoolID}
	argIndex := 3

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (ts.nm_siswa ILIKE $%d OR ts.nis ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM tabel_siswa ts
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
	`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.alamat_siswa, ts.telepon_siswa,
			COALESCE(ts.nm_ayah, '') as nm_ayah, COALESCE(tsp.telepon_ortu, '') as telepon_ortu,
			COALESCE(tk.nm_kelas, '') as nm_kelas
		FROM tabel_siswa ts
		LEFT JOIN tabel_siswa_pelengkap tsp ON ts.peserta_didik_id = tsp.peserta_didik_id
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
		ORDER BY ts.nm_siswa ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var s Student
		var birthDate *time.Time
		if err := rows.Scan(
			&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &birthDate,
			&s.Gender, &s.Address, &s.Phone, &s.ParentName, &s.ParentPhone, &s.ClassName,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		if birthDate != nil { s.BirthDate = *birthDate }
		s.SchoolID = filter.SchoolID
		s.AcademicPeriodID = filter.AcademicPeriodID
		s.Status = "Aktif"
		students = append(students, s)
	}
	return students, total, rows.Err()
}

func (r *pgxStudentRepository) Update(ctx context.Context, student *Student) error {
	if student.SchoolID == "" || student.AcademicPeriodID == "" {
		return ErrMissingOperationalContext
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	querySiswa := `
		UPDATE tabel_siswa SET
			nis = $1, nisn = $2, nm_siswa = $3, jenis_kelamin = $4, tanggal_lahir = $5,
			tempat_lahir = $6, alamat_siswa = $7, telepon_siswa = $8, nm_ayah = $9
		WHERE peserta_didik_id = $10
	`
	var birthDate any = nil
	if !student.BirthDate.IsZero() {
		birthDate = student.BirthDate
	}

	result, err := tx.Exec(ctx, querySiswa,
		student.NIS, student.NISN, student.Name, student.Gender, birthDate,
		student.BirthPlace, student.Address, student.Phone, student.ParentName, student.ID,
	)
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	queryPelengkap := `
		UPDATE tabel_siswa_pelengkap SET
			telepon_ortu = $1
		WHERE peserta_didik_id = $2
	`
	_, err = tx.Exec(ctx, queryPelengkap, student.ParentPhone, student.ID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxStudentRepository) Delete(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	if schoolID == "" || academicPeriodID == "" {
		return ErrMissingOperationalContext
	}
	query := `
		DELETE FROM tabel_anggotakelas
		WHERE peserta_didik_id = $1 AND semester_id = $2
	`
	result, err := r.db.Exec(ctx, query, studentID, academicPeriodID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxStudentRepository) CountByClass(ctx context.Context, schoolID, academicPeriodID, classID string) (int, error) {
	if schoolID == "" || academicPeriodID == "" {
		return 0, ErrMissingOperationalContext
	}
	query := `
		SELECT COUNT(*)
		FROM tabel_anggotakelas tak
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tak.rombongan_belajar_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id::text = $3
	`
	var count int
	err := r.db.QueryRow(ctx, query, classID, academicPeriodID, schoolID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return count, nil
}

// === AREA UI ===
func (r *pgxStudentRepository) ListForUI(ctx context.Context, filter StudentFilter) ([]StudentUIDTO, int, error) {
	whereClause := "WHERE tak.semester_id = $1 AND tk.sekolah_id::text = $2"
	args := []any{filter.AcademicPeriodID, filter.SchoolID}
	argIndex := 3

	if filter.ClassID != "" {
		whereClause += fmt.Sprintf(" AND tak.rombongan_belajar_id::text = $%d", argIndex)
		args = append(args, filter.ClassID)
		argIndex++
	}
	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (ts.nm_siswa ILIKE $%d OR ts.nis ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf(`
		SELECT COUNT(*)
		FROM tabel_siswa ts
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
	`, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 20 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.agama, ts.alamat_siswa, ts.telepon_siswa,
			ts.nm_ayah, ts.nm_ibu, COALESCE(tk.nm_kelas, '') as nm_kelas
		FROM tabel_siswa ts
		LEFT JOIN tabel_siswa_pelengkap tsp ON ts.peserta_didik_id = tsp.peserta_didik_id
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		%s
		ORDER BY ts.nm_siswa ASC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var students []StudentUIDTO
	for rows.Next() {
		var s StudentUIDTO
		if err := rows.Scan(
			&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &s.BirthDate,
			&s.Gender, &s.Agama, &s.Address, &s.Phone,
			&s.ParentName, &s.ParentNameIbu, &s.ClassName,
		); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		s.Status = "Aktif"
		students = append(students, s)
	}
	return students, total, rows.Err()
}

func (r *pgxStudentRepository) GetByIDForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) (*StudentUIDTO, error) {
	query := `
		SELECT
			ts.peserta_didik_id, ts.nis, ts.nisn, ts.nm_siswa, ts.tempat_lahir, ts.tanggal_lahir,
			ts.jenis_kelamin, ts.agama, ts.alamat_siswa, ts.telepon_siswa,
			ts.nm_ayah, ts.nm_ibu, COALESCE(tk.nm_kelas, '') as nm_kelas
		FROM tabel_siswa ts
		LEFT JOIN tabel_siswa_pelengkap tsp ON ts.peserta_didik_id = tsp.peserta_didik_id
		JOIN tabel_anggotakelas tak ON ts.peserta_didik_id = tak.peserta_didik_id
		JOIN tabel_kelas tk ON tak.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE ts.peserta_didik_id = $1 AND tak.semester_id = $2 AND tk.sekolah_id::text = $3
	`
	var s StudentUIDTO
	err := r.db.QueryRow(ctx, query, studentID, academicPeriodID, schoolID).Scan(
		&s.ID, &s.NIS, &s.NISN, &s.Name, &s.BirthPlace, &s.BirthDate,
		&s.Gender, &s.Agama, &s.Address, &s.Phone,
		&s.ParentName, &s.ParentNameIbu, &s.ClassName,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	s.Status = "Aktif"
	return &s, nil
}

func (r *pgxStudentRepository) UpdateForUI(ctx context.Context, schoolID, academicPeriodID, studentID string, name, birthPlace, gender, address, phone, className string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	querySiswa := `
		UPDATE tabel_siswa SET
			nm_siswa = $1, tempat_lahir = $2, jenis_kelamin = $3,
			alamat_siswa = $4, telepon_siswa = $5
		WHERE peserta_didik_id = $6
	`
	result, err := tx.Exec(ctx, querySiswa, name, birthPlace, gender, address, phone, studentID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	if className != "" {
		var rombelID string
		queryFindRombel := `
			SELECT rombongan_belajar_id::text
			FROM tabel_kelas
			WHERE (nm_kelas = $1 OR rombongan_belajar_id::text = $1)
			  AND semester_id = $2
			  AND sekolah_id::text = $3
			LIMIT 1
		`
		errFind := tx.QueryRow(ctx, queryFindRombel, className, academicPeriodID, schoolID).Scan(&rombelID)

		if errFind == nil {
			var count int
			tx.QueryRow(ctx, `SELECT COUNT(*) FROM tabel_anggotakelas WHERE peserta_didik_id = $1 AND semester_id = $2`, studentID, academicPeriodID).Scan(&count)

			if count > 0 {
				queryAnggota := `
					UPDATE tabel_anggotakelas SET rombongan_belajar_id = $1
					WHERE peserta_didik_id = $2 AND semester_id = $3
				`
				_, err = tx.Exec(ctx, queryAnggota, rombelID, studentID, academicPeriodID)
			} else {
				queryAnggota := `
					INSERT INTO tabel_anggotakelas (anggota_rombel_id, peserta_didik_id, rombongan_belajar_id, semester_id)
					VALUES ($1, $2, $3, $4)
				`
				_, err = tx.Exec(ctx, queryAnggota, uuid.New().String(), studentID, rombelID, academicPeriodID)
			}
			if err != nil {
				return fmt.Errorf("%w: %v", ErrDatabaseError, err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxStudentRepository) DeleteForUI(ctx context.Context, schoolID, academicPeriodID, studentID string) error {
	return r.Delete(ctx, schoolID, academicPeriodID, studentID)
}

func isDuplicateKeyError(err error) bool {
	return err != nil && err.Error() != "" &&
		(contains(err.Error(), "duplicate key") || contains(err.Error(), "23505"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
