package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Pembelajaran struct {
	ID          string `json:"id"`
	ClassID     string `json:"classId"`
	ClassName   string `json:"className"`
	SubjectID   int    `json:"subjectId"`
	SubjectName string `json:"subjectName"`
	TeacherID   string `json:"teacherId"`
	TeacherName string `json:"teacherName"`
	SemesterID  string `json:"semesterId"`
}

type PembelajaranRepository interface {
	List(ctx context.Context, schoolID, periodID, classID string) ([]Pembelajaran, error)
	Create(ctx context.Context, p *Pembelajaran) error
	Delete(ctx context.Context, id string) error
}

type pgxPembelajaranRepository struct {
	db *pgxpool.Pool
}

func NewPembelajaranRepository(db *pgxpool.Pool) PembelajaranRepository {
	return &pgxPembelajaranRepository{db: db}
}

func (r *pgxPembelajaranRepository) List(ctx context.Context, schoolID, periodID, classID string) ([]Pembelajaran, error) {
	query := `
		SELECT 
			tpm.pembelajaran_id, tk.rombongan_belajar_id, tk.nm_kelas,
			tpm.mata_pelajaran_id, tm.nm_mapel,
			COALESCE(tpm.ptk_terdaftar_id::text, ''), COALESCE(tp.nama, ''),
			tpm.semester_id
		FROM tabel_pembelajaran tpm
		JOIN tabel_kelas tk ON tpm.rombongan_belajar_id = tk.rombongan_belajar_id
		JOIN tabel_mapel tm ON tpm.mata_pelajaran_id = tm.mata_pelajaran_id
		LEFT JOIN tabel_ptk tp ON tpm.ptk_terdaftar_id = tp.ptk_id
		WHERE tk.sekolah_id::text = $1 AND tpm.semester_id = $2
	`
	args := []interface{}{schoolID, periodID}
	
	if classID != "" {
		query += " AND tpm.rombongan_belajar_id = $3"
		args = append(args, classID)
	}
	query += " ORDER BY tk.nm_kelas ASC, tm.nm_mapel ASC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []Pembelajaran
	for rows.Next() {
		var p Pembelajaran
		if err := rows.Scan(&p.ID, &p.ClassID, &p.ClassName, &p.SubjectID, &p.SubjectName, &p.TeacherID, &p.TeacherName, &p.SemesterID); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, p)
	}
	return list, rows.Err()
}

func (r *pgxPembelajaranRepository) Create(ctx context.Context, p *Pembelajaran) error {
	query := `
		INSERT INTO tabel_pembelajaran (
			pembelajaran_id, rombongan_belajar_id, mata_pelajaran_id, semester_id, ptk_terdaftar_id
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, p.ID, p.ClassID, p.SubjectID, p.SemesterID, nullableUUID(p.TeacherID))
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxPembelajaranRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tabel_pembelajaran WHERE pembelajaran_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
