package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Jurusan struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type JurusanRepository interface {
	List(ctx context.Context) ([]Jurusan, error)
}

type pgxJurusanRepository struct {
	db *pgxpool.Pool
}

func NewJurusanRepository(db *pgxpool.Pool) JurusanRepository {
	return &pgxJurusanRepository{db: db}
}

func (r *pgxJurusanRepository) List(ctx context.Context) ([]Jurusan, error) {
	query := `SELECT jurusan_id, nama_jurusan FROM tabel_jurusan ORDER BY nama_jurusan ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []Jurusan
	for rows.Next() {
		var j Jurusan
		if err := rows.Scan(&j.ID, &j.Name); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		list = append(list, j)
	}
	return list, rows.Err()
}
