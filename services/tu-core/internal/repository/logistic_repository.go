package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Logistic struct {
	ID          string    `json:"id"`
	Kode        string    `json:"kode"`
	Nama        string    `json:"nama"`
	Kategori    string    `json:"kategori"`
	Satuan      string    `json:"satuan"`
	Stok        int       `json:"stok"`
	StokMinimum int       `json:"stok_minimum"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LogisticMutation struct {
	ID         string    `json:"id"`
	LogisticID string    `json:"logistic_id"`
	Tipe       string    `json:"tipe"`
	Jumlah     int       `json:"jumlah"`
	Keterangan string    `json:"keterangan"`
	CreatedAt  time.Time `json:"created_at"`
}

type LogisticRepository interface {
	Create(ctx context.Context, req *Logistic) error
	FindAll(ctx context.Context) ([]Logistic, error)
	FindByID(ctx context.Context, id string) (*Logistic, error)
	UpdateStock(ctx context.Context, id string, newStock int) error
	RecordMutation(ctx context.Context, mutation *LogisticMutation) error
}

type pgxLogisticRepository struct {
	db *pgxpool.Pool
}

func NewLogisticRepository(db *pgxpool.Pool) LogisticRepository {
	return &pgxLogisticRepository{db: db}
}

func (r *pgxLogisticRepository) Create(ctx context.Context, req *Logistic) error {
	if req.ID == "" {
		req.ID = uuid.New().String()
	}
	query := `INSERT INTO sarpras_logistics (id, kode, nama, kategori, satuan, stok, stok_minimum) 
	          VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, req.ID, req.Kode, req.Nama, req.Kategori, req.Satuan, req.Stok, req.StokMinimum)
	return err
}

func (r *pgxLogisticRepository) FindAll(ctx context.Context) ([]Logistic, error) {
	query := `SELECT id, kode, nama, kategori, satuan, stok, stok_minimum, created_at, updated_at 
	          FROM sarpras_logistics ORDER BY kode ASC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Logistic
	for rows.Next() {
		var l Logistic
		if err := rows.Scan(&l.ID, &l.Kode, &l.Nama, &l.Kategori, &l.Satuan, &l.Stok, &l.StokMinimum, &l.CreatedAt, &l.UpdatedAt); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *pgxLogisticRepository) FindByID(ctx context.Context, id string) (*Logistic, error) {
	query := `SELECT id, kode, nama, kategori, satuan, stok, stok_minimum, created_at, updated_at 
	          FROM sarpras_logistics WHERE id = $1`
	var l Logistic
	err := r.db.QueryRow(ctx, query, id).
		Scan(&l.ID, &l.Kode, &l.Nama, &l.Kategori, &l.Satuan, &l.Stok, &l.StokMinimum, &l.CreatedAt, &l.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("logistik tidak ditemukan")
		}
		return nil, err
	}
	return &l, nil
}

func (r *pgxLogisticRepository) UpdateStock(ctx context.Context, id string, newStock int) error {
	query := `UPDATE sarpras_logistics SET stok = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, newStock, time.Now(), id)
	return err
}

func (r *pgxLogisticRepository) RecordMutation(ctx context.Context, mutation *LogisticMutation) error {
	if mutation.ID == "" {
		mutation.ID = uuid.New().String()
	}
	query := `INSERT INTO sarpras_logistic_mutations (id, logistic_id, tipe, jumlah, keterangan) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, mutation.ID, mutation.LogisticID, mutation.Tipe, mutation.Jumlah, mutation.Keterangan)
	return err
}
