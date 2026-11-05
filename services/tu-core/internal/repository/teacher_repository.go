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

type Teacher struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	NIP             string    `json:"nip"`
	NUPTK           string    `json:"nuptk"`
	Gender          string    `json:"gender"`
	BirthPlace      string    `json:"birthPlace"`
	BirthDate       time.Time `json:"birthDate"`
	Address         string    `json:"address"`
	StatusKeaktifan int       `json:"statusKeaktifan"`
	SoftDelete      int       `json:"softDelete"`
	GelarDepan      string    `json:"gelarDepan"`
	GelarBelakang   string    `json:"gelarBelakang"`
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
}

type pgxTeacherRepository struct {
	db *pgxpool.Pool
}

func NewTeacherRepository(db *pgxpool.Pool) TeacherRepository {
	return &pgxTeacherRepository{db: db}
}

const teacherSelectColumns = `
	tp.ptk_id, 
	COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
	COALESCE(tp.nip, ''), COALESCE(tp.nuptk, ''), tp.jenis_kelamin, 
	tp.tempat_lahir, tp.tanggal_lahir, tp.alamat_jalan, tp.status_keaktifan_id, tp.soft_delete,
	COALESCE(tpp.gelar_depan, ''), COALESCE(tpp.gelar_belakang, '')
`

const teacherFromClause = `
	FROM tabel_ptk tp
	LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
	JOIN tabel_pembelajaran tpm ON tp.ptk_id = tpm.ptk_terdaftar_id
	JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
`

func scanTeacher(row pgx.Row, t *Teacher) error {
	var birthDate *time.Time
	if err := row.Scan(
		&t.ID, &t.Name, &t.NIP, &t.NUPTK, &t.Gender,
		&t.BirthPlace, &birthDate, &t.Address, &t.StatusKeaktifan, &t.SoftDelete,
		&t.GelarDepan, &t.GelarBelakang,
	); err != nil {
		return err
	}
	if birthDate != nil {
		t.BirthDate = *birthDate
	}
	return nil
}

func (r *pgxTeacherRepository) GetByID(ctx context.Context, teacherID string) (*Teacher, error) {
	query := `SELECT ` + teacherSelectColumns + teacherFromClause + ` WHERE tp.ptk_id = $1 AND tp.soft_delete = 0 LIMIT 1`
	var t Teacher
	if err := scanTeacher(r.db.QueryRow(ctx, query, teacherID), &t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return &t, nil
}

func (r *pgxTeacherRepository) List(ctx context.Context, filter TeacherFilter) ([]Teacher, int, error) {
	if filter.SchoolID == "" || filter.PeriodID == "" {
		return nil, 0, ErrMissingOperationalContext
	}

	whereClause := "WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2 AND tp.soft_delete = 0"
	args := []interface{}{filter.SchoolID, filter.PeriodID}
	argIndex := 3

	if filter.Keyword != "" {
		whereClause += fmt.Sprintf(" AND (tp.nama ILIKE $%d OR tp.nip ILIKE $%d)", argIndex, argIndex)
		args = append(args, "%"+filter.Keyword+"%")
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(DISTINCT tp.ptk_id) %s %s", teacherFromClause, whereClause)
	var total int
	if err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	if filter.Limit <= 0 { filter.Limit = 50 }
	if filter.Offset < 0 { filter.Offset = 0 }

	dataQuery := fmt.Sprintf(`
	SELECT DISTINCT `+teacherSelectColumns+`
	%s %s 
	ORDER BY nama_lengkap ASC 
	LIMIT $%d OFFSET $%d
	`, teacherFromClause, whereClause, argIndex, argIndex+1)
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

func (r *pgxTeacherRepository) Create(ctx context.Context, t *Teacher) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	// 1. Insert ke tabel_ptk
	queryPtk := `
		INSERT INTO tabel_ptk (
			ptk_id, nama, nip, jenis_ptk_id, jenis_kelamin, tempat_lahir, tanggal_lahir, 
			nuptk, alamat_jalan, status_keaktifan_id, soft_delete
		) VALUES ($1, $2, $3, 1, $4, $5, $6, $7, $8, 1, 0)
	`
	var birthDate any = nil
	if !t.BirthDate.IsZero() {
		birthDate = t.BirthDate
	}

	_, err = tx.Exec(ctx, queryPtk,
		t.ID, t.Name, nullableString(t.NIP), t.Gender, nullableString(t.BirthPlace),
		birthDate, nullableString(t.NUPTK), t.Address,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	// 2. Insert ke tabel_ptk_pelengkap
	queryPelengkap := `
		INSERT INTO tabel_ptk_pelengkap (
			ptk_pelengkap_id, ptk_id, gelar_depan, gelar_belakang, nip_niy
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, queryPelengkap,
		uuid.New().String(), t.ID, nullableString(t.GelarDepan), nullableString(t.GelarBelakang), nullableString(t.NIP),
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return tx.Commit(ctx)
}

func (r *pgxTeacherRepository) Update(ctx context.Context, t *Teacher) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	// 1. Update tabel_ptk
	queryPtk := `
		UPDATE tabel_ptk SET
			nama = $1, nip = $2, jenis_kelamin = $3, tempat_lahir = $4, 
			tanggal_lahir = $5, nuptk = $6, alamat_jalan = $7
		WHERE ptk_id = $8 AND soft_delete = 0
	`
	var birthDate any = nil
	if !t.BirthDate.IsZero() {
		birthDate = t.BirthDate
	}

	result, err := tx.Exec(ctx, queryPtk,
		t.Name, nullableString(t.NIP), t.Gender, nullableString(t.BirthPlace),
		birthDate, nullableString(t.NUPTK), t.Address, t.ID,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	// 2. Upsert tabel_ptk_pelengkap
	queryPelengkap := `
		INSERT INTO tabel_ptk_pelengkap (ptk_pelengkap_id, ptk_id, gelar_depan, gelar_belakang, nip_niy)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (ptk_id) DO UPDATE SET
			gelar_depan = EXCLUDED.gelar_depan,
			gelar_belakang = EXCLUDED.gelar_belakang,
			nip_niy = EXCLUDED.nip_niy
	`
	_, err = tx.Exec(ctx, queryPelengkap,
		uuid.New().String(), t.ID, nullableString(t.GelarDepan), nullableString(t.GelarBelakang), nullableString(t.NIP),
	)
	if err != nil {
		// Jika constraint ptk_id belum ada di tabel_pelengkap, kita handle fallback (update manual)
		queryUpdate := `UPDATE tabel_ptk_pelengkap SET gelar_depan = $1, gelar_belakang = $2, nip_niy = $3 WHERE ptk_id = $4`
		_, errUpdate := tx.Exec(ctx, queryUpdate, nullableString(t.GelarDepan), nullableString(t.GelarBelakang), nullableString(t.NIP), t.ID)
		if errUpdate != nil {
			return fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *pgxTeacherRepository) Delete(ctx context.Context, teacherID string) error {
	query := `UPDATE tabel_ptk SET soft_delete = 1 WHERE ptk_id = $1 AND soft_delete = 0`
	result, err := r.db.Exec(ctx, query, teacherID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
