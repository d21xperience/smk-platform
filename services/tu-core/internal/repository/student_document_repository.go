package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type StudentDocument struct {
	ID          string    `json:"id"`
	StudentID   string    `json:"studentId"`
	JenisBerkas string    `json:"jenisBerkas"`
	NamaFile    string    `json:"namaFile"`
	URLFile     string    `json:"urlFile"`
	CreatedAt   time.Time `json:"createdAt"`
}

type StudentDocumentRepository interface {
	Create(ctx context.Context, doc *StudentDocument) error
	GetByStudentID(ctx context.Context, studentID string) ([]StudentDocument, error)
	Delete(ctx context.Context, id string) error
}

type pgxStudentDocumentRepository struct {
	db *pgxpool.Pool
}

func NewStudentDocumentRepository(db *pgxpool.Pool) StudentDocumentRepository {
	return &pgxStudentDocumentRepository{db: db}
}

func (r *pgxStudentDocumentRepository) Create(ctx context.Context, doc *StudentDocument) error {
	query := `
		INSERT INTO tabel_berkas_siswa (
			id, peserta_didik_id, jenis_berkas, nama_file, url_file
		) VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(ctx, query, doc.ID, doc.StudentID, doc.JenisBerkas, doc.NamaFile, doc.URLFile)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxStudentDocumentRepository) GetByStudentID(ctx context.Context, studentID string) ([]StudentDocument, error) {
	query := `
		SELECT id, peserta_didik_id, jenis_berkas, nama_file, url_file, created_at
		FROM tabel_berkas_siswa
		WHERE peserta_didik_id = $1
		ORDER BY created_at ASC
	`
	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var docs []StudentDocument
	for rows.Next() {
		var d StudentDocument
		if err := rows.Scan(&d.ID, &d.StudentID, &d.JenisBerkas, &d.NamaFile, &d.URLFile, &d.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
		}
		docs = append(docs, d)
	}
	return docs, rows.Err()
}

func (r *pgxStudentDocumentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tabel_berkas_siswa WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	if result.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
