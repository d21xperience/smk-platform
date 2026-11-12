package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Jurusan struct {
	JurusanID   string `json:"jurusan_id"`
	NamaJurusan string `json:"nama_jurusan"`
	LevelBidang string `json:"level_bidang_id"`
}

type JurusanRepository interface {
	List(ctx context.Context) ([]Jurusan, error)
	Create(ctx context.Context, j *Jurusan) error
}

type jurusanRepository struct {
	db *pgxpool.Pool
}

func NewJurusanRepository(db *pgxpool.Pool) JurusanRepository {
	return &jurusanRepository{db: db}
}

func (r *jurusanRepository) List(ctx context.Context) ([]Jurusan, error) {
	query := `SELECT jurusan_id, nama_jurusan, level_bidang_id FROM tabel_jurusan ORDER BY nama_jurusan ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Jurusan
	for rows.Next() {
		var j Jurusan
		if err := rows.Scan(&j.JurusanID, &j.NamaJurusan, &j.LevelBidang); err != nil {
			return nil, err
		}
		result = append(result, j)
	}
	if result == nil {
		result = []Jurusan{}
	}
	return result, nil
}

func (r *jurusanRepository) Create(ctx context.Context, j *Jurusan) error {
	query := `INSERT INTO tabel_jurusan (jurusan_id, nama_jurusan, level_bidang_id) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, query, j.JurusanID, j.NamaJurusan, j.LevelBidang)
	return err
}
