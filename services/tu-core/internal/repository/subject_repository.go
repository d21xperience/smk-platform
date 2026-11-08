package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Subject struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Kelompok  string `json:"kelompok"`
	Semester  int    `json:"semester"`
	JurusanID string `json:"jurusanId"`
	NmRingkas string `json:"nmRingkas"`
}

type SubjectRepository interface {
	List(ctx context.Context) ([]Subject, error)
	Create(ctx context.Context, s *Subject) error
	Update(ctx context.Context, s *Subject) error
	Delete(ctx context.Context, id int) error
}

type pgxSubjectRepository struct {
	db *pgxpool.Pool
}

func NewSubjectRepository(db *pgxpool.Pool) SubjectRepository {
	return &pgxSubjectRepository{db: db}
}

func (r *pgxSubjectRepository) List(ctx context.Context) ([]Subject, error) {
	query := `
		SELECT mata_pelajaran_id, COALESCE(nm_mapel, ''), COALESCE(kelompok, ''),
		       COALESCE(semester, 0), COALESCE(jurusan_id, ''), COALESCE(nm_ringkas, '')
		FROM tabel_mapel
		ORDER BY nm_mapel ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var subjects []Subject
	for rows.Next() {
		var s Subject
		if err := rows.Scan(&s.ID, &s.Name, &s.Kelompok, &s.Semester, &s.JurusanID, &s.NmRingkas); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		subjects = append(subjects, s)
	}
	return subjects, rows.Err()
}

func (r *pgxSubjectRepository) Create(ctx context.Context, s *Subject) error {
	query := `
		INSERT INTO tabel_mapel (mata_pelajaran_id, nm_mapel, kelompok, semester, jurusan_id, nm_ringkas)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	// Jika ID 0, gunakan sequence (asumsi ada auto-increment/sequence di DB, jika tidak pakai MAX(id)+1)
	var nextID int
	if s.ID == 0 {
		_ = r.db.QueryRow(ctx, "SELECT COALESCE(MAX(mata_pelajaran_id), 0) + 1 FROM tabel_mapel").Scan(&nextID)
		s.ID = nextID
	}

	_, err := r.db.Exec(ctx, query, s.ID, s.Name, nullableString(s.Kelompok), s.Semester, nullableString(s.JurusanID), nullableString(s.NmRingkas))
	if err != nil {
		if isDuplicateKeyError(err) {
			return fmt.Errorf("%w: %v", ErrDuplicateEntry, err)
		}
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxSubjectRepository) Update(ctx context.Context, s *Subject) error {
	query := `
		UPDATE tabel_mapel SET
			nm_mapel = $1, kelompok = $2, semester = $3, jurusan_id = $4, nm_ringkas = $5
		WHERE mata_pelajaran_id = $6
	`
	result, err := r.db.Exec(ctx, query, s.Name, nullableString(s.Kelompok), s.Semester, nullableString(s.JurusanID), nullableString(s.NmRingkas), s.ID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *pgxSubjectRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM tabel_mapel WHERE mata_pelajaran_id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
