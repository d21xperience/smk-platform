package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Teacher struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	NIP              string    `json:"nip"`
	NUPTK            string    `json:"nuptk"`
	Gender           string    `json:"gender"`
	BirthPlace       string    `json:"birthPlace"`
	BirthDate        time.Time `json:"birthDate"`
	Address          string    `json:"address"`
	JenisPtkID       int       `json:"jenisPtkId"`
	StatusKeaktifan  int       `json:"statusKeaktifan"`
	SoftDelete       int       `json:"softDelete"`
	GelarDepan       string    `json:"gelarDepan"`
	GelarBelakang    string    `json:"gelarBelakang"`
	IsActiveInPeriod bool      `json:"isActiveInPeriod"`

	// Field Baru Fase 1
	NIK               string `json:"nik"`
	AgamaID           int    `json:"agamaId"`
	StatusKepegawaian int    `json:"statusKepegawaian"`
	FotoURL           string `json:"fotoUrl"`
	NoHP              string `json:"noHp"`
	Email             string `json:"email"`
	StatusPerkawinan  int    `json:"statusPerkawinan"`
	NamaIbuKandung    string `json:"namaIbuKandung"`
}

type ScheduleDetail struct {
	ID          string `json:"id"`
	Day         string `json:"day"`
	StartTime   string `json:"startTime"`
	EndTime     string `json:"endTime"`
	ClassName   string `json:"className"`
	SubjectName string `json:"subjectName"`
}

// FIX: Renamed to TeacherWorkloadSummary to avoid conflict with kurikulum_stat_repository.go
type TeacherWorkloadSummary struct {
	TeacherID     string           `json:"teacherId"`
	TeacherName   string           `json:"teacherName"`
	TotalClasses  int              `json:"totalClasses"`
	TotalSubjects int              `json:"totalSubjects"`
	TotalSessions int              `json:"totalSessions"`
	Schedules     []ScheduleDetail `json:"schedules"`
}

type TeacherFilter struct {
	SchoolID string
	PeriodID string
	Keyword  string
	Limit    int
	Offset   int
}

type TeacherRepository interface {
	GetByID(ctx context.Context, teacherID string) (*Teacher, error)
	List(ctx context.Context, filter TeacherFilter) ([]Teacher, int, error)
	Create(ctx context.Context, t *Teacher) error
	Update(ctx context.Context, t *Teacher) error
	Delete(ctx context.Context, teacherID string) error
	UpdatePhoto(ctx context.Context, teacherID, photoURL string) error
	ListWorkloads(ctx context.Context, schoolID, periodID string) ([]TeacherWorkloadSummary, error)
}

type pgxTeacherRepository struct {
	db *pgxpool.Pool
}

func NewTeacherRepository(db *pgxpool.Pool) TeacherRepository {
	return &pgxTeacherRepository{db: db}
}

const teacherSelectColumns = `
	tp.ptk_id::text, 
	COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
	COALESCE(tp.nip, ''), COALESCE(tp.nuptk, ''), tp.jenis_kelamin, 
	tp.tempat_lahir, tp.tanggal_lahir, tp.alamat_jalan, tp.jenis_ptk_id, tp.status_keaktifan_id, tp.soft_delete,
	COALESCE(tpp.gelar_depan, ''), COALESCE(tpp.gelar_belakang, ''),
	EXISTS (
		SELECT 1 FROM tabel_pembelajaran tpm
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		WHERE tpm.ptk_terdaftar_id = tp.ptk_id 
		  AND tpm.semester_id = $1
		  AND tk.sekolah_id::text = $2
	) AS is_active,
	COALESCE(tp.nik, ''), COALESCE(tp.agama_id, 0), COALESCE(tp.status_kepegawaian_id, 0), 
	COALESCE(tp.foto_url, ''), COALESCE(tp.no_hp, ''), COALESCE(tp.email, ''), 
	COALESCE(tp.status_perkawinan, 0), COALESCE(tp.nama_ibu_kandung, '')
`

func scanTeacher(row pgx.Row, t *Teacher) error {
	var birthDate *time.Time
	if err := row.Scan(
		&t.ID, &t.Name, &t.NIP, &t.NUPTK, &t.Gender,
		&t.BirthPlace, &birthDate, &t.Address, &t.JenisPtkID, &t.StatusKeaktifan, &t.SoftDelete,
		&t.GelarDepan, &t.GelarBelakang, &t.IsActiveInPeriod,
		&t.NIK, &t.AgamaID, &t.StatusKepegawaian, &t.FotoURL, &t.NoHP, &t.Email, &t.StatusPerkawinan, &t.NamaIbuKandung,
	); err != nil {
		return err
	}
	if birthDate != nil {
		t.BirthDate = *birthDate
	}
	return nil
}

func (r *pgxTeacherRepository) GetByID(ctx context.Context, teacherID string) (*Teacher, error) {
	query := `SELECT ` + teacherSelectColumns + ` FROM tabel_ptk tp 
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id 
		WHERE tp.ptk_id = $3 AND tp.soft_delete = 0 LIMIT 1`

	var t Teacher
	if err := scanTeacher(r.db.QueryRow(ctx, query, "", "", teacherID), &t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &t, nil
}

func (r *pgxTeacherRepository) List(ctx context.Context, filter TeacherFilter) ([]Teacher, int, error) {
	if filter.PeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE tp.soft_delete = 0"
	args := []interface{}{filter.PeriodID, filter.SchoolID}
	argIndex := 3

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (tp.nama ILIKE $%d OR tp.nip ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tabel_ptk tp %s", whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args[2:]...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 50 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
	SELECT `+teacherSelectColumns+`
	FROM tabel_ptk tp
	LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
	%s 
	ORDER BY nama_lengkap ASC 
	LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var teachers []Teacher
	for rows.Next() {
		var t Teacher
		if err := scanTeacher(rows, &t); err != nil {
			return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		teachers = append(teachers, t)
	}
	return teachers, total, rows.Err()
}

func (r *pgxTeacherRepository) ListWorkloads(ctx context.Context, schoolID, periodID string) ([]TeacherWorkloadSummary, error) {
	query := `
		SELECT 
			tp.ptk_id::text, 
			COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
			COUNT(DISTINCT tk.rombongan_belajar_id) AS total_classes,
			COUNT(DISTINCT tm.mata_pelajaran_id) AS total_subjects,
			COUNT(jp.id) AS total_sessions,
			COALESCE(
				json_agg(
					json_build_object(
						'id', jp.id,
						'day', jp.hari,
						'startTime', jp.jam_mulai::text,
						'endTime', jp.jam_selesai::text,
						'className', tk.nm_kelas,
						'subjectName', tm.nm_mapel
					) ORDER BY jp.hari, jp.jam_mulai
				) FILTER (WHERE jp.id IS NOT NULL), 
				'[]'::json
			) AS schedules
		FROM tabel_ptk tp
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		JOIN tabel_pembelajaran tpm ON tp.ptk_id = tpm.ptk_terdaftar_id
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		JOIN tabel_mapel tm ON tpm.mata_pelajaran_id = tm.mata_pelajaran_id
		LEFT JOIN tabel_jadwal_pelajaran jp ON tpm.pembelajaran_id = jp.pembelajaran_id
		WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2 AND tp.soft_delete = 0
		GROUP BY tp.ptk_id, tpp.gelar_depan, tp.nama, tpp.gelar_belakang
		ORDER BY nama_lengkap;
	`
	rows, err := r.db.Query(ctx, query, schoolID, periodID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workloads []TeacherWorkloadSummary
	for rows.Next() {
		var w TeacherWorkloadSummary
		var schedulesJSON []byte
		if err := rows.Scan(&w.TeacherID, &w.TeacherName, &w.TotalClasses, &w.TotalSubjects, &w.TotalSessions, &schedulesJSON); err != nil {
			return nil, err
		}
		if len(schedulesJSON) > 0 {
			if err := json.Unmarshal(schedulesJSON, &w.Schedules); err != nil {
				return nil, err
			}
		}
		workloads = append(workloads, w)
	}
	return workloads, nil
}

func (r *pgxTeacherRepository) Create(ctx context.Context, t *Teacher) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }
	defer tx.Rollback(ctx)

	if t.ID == "" { t.ID = uuid.New().String() }

	queryPtk := `
		INSERT INTO tabel_ptk (
			ptk_id, nama, nip, jenis_ptk_id, jenis_kelamin, tempat_lahir, tanggal_lahir, 
			nuptk, alamat_jalan, status_keaktifan_id, soft_delete,
			nik, agama_id, status_kepegawaian_id, no_hp, email, status_perkawinan, nama_ibu_kandung
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 0, $11, $12, $13, $14, $15, $16, $17)
	`
	var birthDate any = nil
	if !t.BirthDate.IsZero() { birthDate = t.BirthDate }

	_, err = tx.Exec(ctx, queryPtk,
		t.ID, t.Name, nullableString(t.NIP), t.JenisPtkID, t.Gender, nullableString(t.BirthPlace),
		birthDate, nullableString(t.NUPTK), t.Address, t.StatusKeaktifan,
		nullableString(t.NIK), t.AgamaID, t.StatusKepegawaian, nullableString(t.NoHP), nullableString(t.Email), t.StatusPerkawinan, nullableString(t.NamaIbuKandung),
	)
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }

	queryPelengkap := `
		INSERT INTO tabel_ptk_pelengkap (ptk_pelengkap_id, ptk_id, gelar_depan, gelar_belakang, nip_niy)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryPelengkap, uuid.New().String(), t.ID, nullableString(t.GelarDepan), nullableString(t.GelarBelakang), nullableString(t.NIP))
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }

	return tx.Commit(ctx)
}

func (r *pgxTeacherRepository) Update(ctx context.Context, t *Teacher) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }
	defer tx.Rollback(ctx)

	queryPtk := `
		UPDATE tabel_ptk SET
			nama = $1, nip = $2, jenis_kelamin = $3, tempat_lahir = $4, tanggal_lahir = $5, 
			nuptk = $6, alamat_jalan = $7, jenis_ptk_id = $8, status_keaktifan_id = $9,
			nik = $10, agama_id = $11, status_kepegawaian_id = $12, no_hp = $13, email = $14, status_perkawinan = $15, nama_ibu_kandung = $16
		WHERE ptk_id = $17 AND soft_delete = 0
	`
	var birthDate any = nil
	if !t.BirthDate.IsZero() { birthDate = t.BirthDate }

	result, err := tx.Exec(ctx, queryPtk,
		t.Name, nullableString(t.NIP), t.Gender, nullableString(t.BirthPlace),
		birthDate, nullableString(t.NUPTK), t.Address, t.JenisPtkID, t.StatusKeaktifan,
		nullableString(t.NIK), t.AgamaID, t.StatusKepegawaian, nullableString(t.NoHP), nullableString(t.Email), t.StatusPerkawinan, nullableString(t.NamaIbuKandung),
		t.ID,
	)
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }
	if result.RowsAffected() == 0 { return ErrNotFound }

	queryPelengkap := `
		INSERT INTO tabel_ptk_pelengkap (ptk_pelengkap_id, ptk_id, gelar_depan, gelar_belakang, nip_niy)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ptk_id) DO UPDATE SET
			gelar_depan = EXCLUDED.gelar_depan, gelar_belakang = EXCLUDED.gelar_belakang, nip_niy = EXCLUDED.nip_niy
	`
	_, err = tx.Exec(ctx, queryPelengkap, uuid.New().String(), t.ID, nullableString(t.GelarDepan), nullableString(t.GelarBelakang), nullableString(t.NIP))
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }

	return tx.Commit(ctx)
}

func (r *pgxTeacherRepository) Delete(ctx context.Context, teacherID string) error {
	query := `UPDATE tabel_ptk SET soft_delete = 1 WHERE ptk_id = $1 AND soft_delete = 0`
	result, err := r.db.Exec(ctx, query, teacherID)
	if err != nil { return fmt.Errorf("%w: %v", ErrDatabaseError, err) }
	if result.RowsAffected() == 0 { return ErrNotFound }
	return nil
}

func (r *pgxTeacherRepository) UpdatePhoto(ctx context.Context, teacherID, photoURL string) error {
	query := `UPDATE tabel_ptk SET foto_url = $1 WHERE ptk_id = $2 AND soft_delete = 0`
	_, err := r.db.Exec(ctx, query, photoURL, teacherID)
	return err
}
