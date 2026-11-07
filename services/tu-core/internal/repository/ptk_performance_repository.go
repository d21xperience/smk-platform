package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Performance struct {
	ID          string    `json:"id"`
	PtkID       string    `json:"ptkId"`
	TeacherName string    `json:"teacherName"`
	Tahun       string    `json:"tahun"`
	Semester    string    `json:"semester"`
	NilaiSkp    float64   `json:"nilaiSkp"`
	Predikat    string    `json:"predikat"`
	Catatan     string    `json:"catatan"`
	Evaluator   string    `json:"evaluator"`
	CreatedAt   time.Time `json:"createdAt"`
}

type PtkPerformanceRepository interface {
	ListPerformances(ctx context.Context) ([]Performance, error)
	CreatePerformance(ctx context.Context, req *Performance) error
	DeletePerformance(ctx context.Context, id string) error
}

type pgxPtkPerformanceRepository struct {
	db *pgxpool.Pool
}

func NewPtkPerformanceRepository(db *pgxpool.Pool) PtkPerformanceRepository {
	return &pgxPtkPerformanceRepository{db: db}
}

func (r *pgxPtkPerformanceRepository) ListPerformances(ctx context.Context) ([]Performance, error) {
	query := `
		SELECT p.id::text, p.ptk_id::text, 
		       COALESCE(tpp.gelar_depan || ' ', '') || tp.nama || COALESCE(', ' || tpp.gelar_belakang, '') AS nama_lengkap,
		       p.tahun, p.semester, p.nilai_skp, p.predikat, COALESCE(p.catatan, ''), COALESCE(p.evaluator, ''), p.created_at
		FROM ptk_performances p
		JOIN tabel_ptk tp ON p.ptk_id = tp.ptk_id
		LEFT JOIN tabel_ptk_pelengkap tpp ON tp.ptk_id = tpp.ptk_id
		ORDER BY p.created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Performance
	for rows.Next() {
		var p Performance
		if err := rows.Scan(&p.ID, &p.PtkID, &p.TeacherName, &p.Tahun, &p.Semester, &p.NilaiSkp, &p.Predikat, &p.Catatan, &p.Evaluator, &p.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, p)
	}
	return list, nil
}

func (r *pgxPtkPerformanceRepository) CreatePerformance(ctx context.Context, req *Performance) error {
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	query := `
		INSERT INTO ptk_performances (id, ptk_id, tahun, semester, nilai_skp, predikat, catatan, evaluator)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(ctx, query, req.ID, req.PtkID, req.Tahun, req.Semester, req.NilaiSkp, req.Predikat, req.Catatan, req.Evaluator)
	return err
}

func (r *pgxPtkPerformanceRepository) DeletePerformance(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM ptk_performances WHERE id = $1`, id)
	return err
}
