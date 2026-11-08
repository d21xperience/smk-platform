package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Schedule struct {
	ID             string `json:"id"`
	PembelajaranID string `json:"pembelajaran_id"`
	Hari           string `json:"hari"`
	JamMulai       string `json:"jam_mulai"`
	JamSelesai     string `json:"jam_selesai"`
	ClassName      string `json:"class_name"`
	SubjectName    string `json:"subject_name"`
	TeacherName    string `json:"teacher_name"`
}

type ScheduleRepository interface {
	List(ctx context.Context, schoolID, periodID string) ([]Schedule, error)
	Create(ctx context.Context, s *Schedule, schoolID, periodID string) error
	Delete(ctx context.Context, id string) error
}

type pgxScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) ScheduleRepository {
	return &pgxScheduleRepository{db: db}
}

func (r *pgxScheduleRepository) List(ctx context.Context, schoolID, periodID string) ([]Schedule, error) {
	query := `
		SELECT
			tjp.id::text, tjp.pembelajaran_id::text, tjp.hari,
			TO_CHAR(tjp.jam_mulai, 'HH24:MI'), TO_CHAR(tjp.jam_selesai, 'HH24:MI'),
			tk.nm_kelas, tm.nm_mapel, COALESCE(tp.nama, '')
		FROM tabel_jadwal_pelajaran tjp
		JOIN tabel_pembelajaran tpm ON tjp.pembelajaran_id = tpm.pembelajaran_id
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		JOIN tabel_mapel tm ON tpm.mata_pelajaran_id = tm.mata_pelajaran_id
		LEFT JOIN tabel_ptk tp ON tpm.ptk_terdaftar_id = tp.ptk_id
		WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2
		ORDER BY
			CASE tjp.hari
				WHEN 'Senin' THEN 1 WHEN 'Selasa' THEN 2 WHEN 'Rabu' THEN 3
				WHEN 'Kamis' THEN 4 WHEN 'Jumat' THEN 5 WHEN 'Sabtu' THEN 6 ELSE 7
			END,
			tjp.jam_mulai ASC
	`
	rows, err := r.db.Query(ctx, query, schoolID, periodID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []Schedule
	for rows.Next() {
		var s Schedule
		if err := rows.Scan(&s.ID, &s.PembelajaranID, &s.Hari, &s.JamMulai, &s.JamSelesai, &s.ClassName, &s.SubjectName, &s.TeacherName); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, s)
	}
	return list, rows.Err()
}

func (r *pgxScheduleRepository) Create(ctx context.Context, s *Schedule, schoolID, periodID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer tx.Rollback(ctx)

	// Cari pembelajaran_id berdasarkan ClassID, SubjectID, TeacherID
	// (Kita asumsikan UI mengirim data ini, tapi karena struct Schedule belum punya field ID-nya, kita perlu parameter tambahan atau modifikasi request).
	// Untuk mempercepat, di handler kita akan passing "pembelajaran_id" langsung dari frontend jika ada.

	// Insert ke tabel_jadwal_pelajaran
	query := `
		INSERT INTO tabel_jadwal_pelajaran (id, pembelajaran_id, hari, jam_mulai, jam_selesai)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.Exec(ctx, query, uuid.New().String(), s.PembelajaranID, s.Hari, s.JamMulai, s.JamSelesai)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}

	return tx.Commit(ctx)
}

func (r *pgxScheduleRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tabel_jadwal_pelajaran WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
