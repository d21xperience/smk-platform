package repository

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/expedition/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExpeditionRepository interface {
	Create(ctx context.Context, e *models.Expedition) error
	GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) ([]models.Expedition, error)
	UpdateBuktiURL(ctx context.Context, id, url string) error
	GetByVerifyToken(ctx context.Context, token string) (*models.Expedition, error)
}

type pgxExpeditionRepository struct {
	db *pgxpool.Pool
}

func NewExpeditionRepository(db *pgxpool.Pool) ExpeditionRepository {
	return &pgxExpeditionRepository{db: db}
}

func (r *pgxExpeditionRepository) Create(ctx context.Context, e *models.Expedition) error {
	query := `
		INSERT INTO expeditions (
			id, school_id, correspondence_id, penerima_nama, penerima_instansi, penerima_alamat,
			metode_kirim, nomor_resi, tanggal_kirim, jumlah_lembar, status, verify_token, created_by, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	e.CreatedAt = time.Now().UTC()

	_, err := r.db.Exec(ctx, query,
		e.ID, e.SchoolID, e.CorrespondenceID, e.PenerimaNama, e.PenerimaInstansi, e.PenerimaAlamat,
		e.MetodeKirim, e.NomorResi, e.TanggalKirim, e.JumlahLembar, e.Status, e.VerifyToken, e.CreatedBy, e.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxExpeditionRepository) GetByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) ([]models.Expedition, error) {
	query := `
		SELECT id, school_id, correspondence_id, penerima_nama, COALESCE(penerima_instansi,''), COALESCE(penerima_alamat,''),
		       metode_kirim, COALESCE(nomor_resi,''), tanggal_kirim, jumlah_lembar, status, COALESCE(bukti_pdf_url,''), verify_token, created_by, created_at
		FROM expeditions
		WHERE school_id = $1 AND correspondence_id = $2
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, schoolID, correspondenceID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	defer rows.Close()

	var list []models.Expedition
	for rows.Next() {
		var e models.Expedition
		if err := rows.Scan(
			&e.ID, &e.SchoolID, &e.CorrespondenceID, &e.PenerimaNama, &e.PenerimaInstansi, &e.PenerimaAlamat,
			&e.MetodeKirim, &e.NomorResi, &e.TanggalKirim, &e.JumlahLembar, &e.Status, &e.BuktiPdfURL, &e.VerifyToken, &e.CreatedBy, &e.CreatedAt,
		); err != nil {
			return nil, err
		}
		list = append(list, e)
	}
	return list, rows.Err()
}

func (r *pgxExpeditionRepository) UpdateBuktiURL(ctx context.Context, id, url string) error {
	query := `UPDATE expeditions SET bukti_pdf_url = $1 WHERE id = $2`
	_, err := r.db.Exec(ctx, query, url, id)
	return err
}

func (r *pgxExpeditionRepository) GetByVerifyToken(ctx context.Context, token string) (*models.Expedition, error) {
	query := `
		SELECT id, school_id, correspondence_id, penerima_nama, COALESCE(penerima_instansi,''), COALESCE(penerima_alamat,''),
		       metode_kirim, COALESCE(nomor_resi,''), tanggal_kirim, jumlah_lembar, status, COALESCE(bukti_pdf_url,''), verify_token, created_by, created_at
		FROM expeditions
		WHERE verify_token = $1
	`
	var e models.Expedition
	err := r.db.QueryRow(ctx, query, token).Scan(
		&e.ID, &e.SchoolID, &e.CorrespondenceID, &e.PenerimaNama, &e.PenerimaInstansi, &e.PenerimaAlamat,
		&e.MetodeKirim, &e.NomorResi, &e.TanggalKirim, &e.JumlahLembar, &e.Status, &e.BuktiPdfURL, &e.VerifyToken, &e.CreatedBy, &e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &e, nil
}
