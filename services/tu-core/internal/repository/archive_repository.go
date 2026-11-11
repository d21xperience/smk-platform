package repository

import (
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/archive/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArchiveRepository interface {
	CreateArchive(ctx context.Context, a *models.Archive) error
	GetArchiveByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.Archive, error)
	GetArchiveByID(ctx context.Context, schoolID, id string) (*models.Archive, error)
	UpdateArchiveStatus(ctx context.Context, id, status, holderID string) error
	UpdateArchiveLocation(ctx context.Context, id, lokasiBaru string) error
	
	CreateMovement(ctx context.Context, m *models.ArchiveMovement) error
	ListMovements(ctx context.Context, archiveID string) ([]models.ArchiveMovement, error)
	
	GetRetentionSchedule(ctx context.Context, kategori string) (*models.RetentionSchedule, error)
}

type pgxArchiveRepository struct {
	db *pgxpool.Pool
}

func NewArchiveRepository(db *pgxpool.Pool) ArchiveRepository {
	return &pgxArchiveRepository{db: db}
}

func (r *pgxArchiveRepository) CreateArchive(ctx context.Context, a *models.Archive) error {
	query := `
		INSERT INTO archives (
			id, school_id, correspondence_id, lokasi_kode, kategori, tanggal_arsip, 
			tanggal_musnah, retensi_tahun, status, keterangan, archived_by, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`
	if a.ID == "" {
		a.ID = uuid.New().String()
	}
	now := time.Now()
	a.CreatedAt = now
	a.UpdatedAt = now

	_, err := r.db.Exec(ctx, query,
		a.ID, a.SchoolID, a.CorrespondenceID, a.LokasiKode, a.Kategori, a.TanggalArsip,
		a.TanggalMusnah, a.RetensiTahun, a.Status, a.Keterangan, a.ArchivedBy, a.CreatedAt, a.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDatabaseError, err)
	}
	return nil
}

func (r *pgxArchiveRepository) GetArchiveByCorrespondenceID(ctx context.Context, schoolID, correspondenceID string) (*models.Archive, error) {
	query := `
		SELECT id, school_id, correspondence_id, lokasi_kode, kategori, tanggal_arsip, 
		       tanggal_musnah, retensi_tahun, status, COALESCE(current_holder_id::text, ''), 
		       COALESCE(keterangan, ''), archived_by, created_at, updated_at
		FROM archives WHERE school_id = $1 AND correspondence_id = $2
	`
	var a models.Archive
	err := r.db.QueryRow(ctx, query, schoolID, correspondenceID).Scan(
		&a.ID, &a.SchoolID, &a.CorrespondenceID, &a.LokasiKode, &a.Kategori, &a.TanggalArsip,
		&a.TanggalMusnah, &a.RetensiTahun, &a.Status, &a.CurrentHolderID,
		&a.Keterangan, &a.ArchivedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err // Biarkan nil jika tidak ditemukan
	}
	return &a, nil
}

func (r *pgxArchiveRepository) GetArchiveByID(ctx context.Context, schoolID, id string) (*models.Archive, error) {
	query := `
		SELECT id, school_id, correspondence_id, lokasi_kode, kategori, tanggal_arsip, 
		       tanggal_musnah, retensi_tahun, status, COALESCE(current_holder_id::text, ''), 
		       COALESCE(keterangan, ''), archived_by, created_at, updated_at
		FROM archives WHERE school_id = $1 AND id = $2
	`
	var a models.Archive
	err := r.db.QueryRow(ctx, query, schoolID, id).Scan(
		&a.ID, &a.SchoolID, &a.CorrespondenceID, &a.LokasiKode, &a.Kategori, &a.TanggalArsip,
		&a.TanggalMusnah, &a.RetensiTahun, &a.Status, &a.CurrentHolderID,
		&a.Keterangan, &a.ArchivedBy, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *pgxArchiveRepository) UpdateArchiveStatus(ctx context.Context, id, status, holderID string) error {
	query := `UPDATE archives SET status = $1, current_holder_id = NULLIF($2, '')::uuid, updated_at = NOW() WHERE id = $3`
	_, err := r.db.Exec(ctx, query, status, holderID, id)
	return err
}

func (r *pgxArchiveRepository) UpdateArchiveLocation(ctx context.Context, id, lokasiBaru string) error {
	query := `UPDATE archives SET lokasi_kode = $1, updated_at = NOW() WHERE id = $2`
	_, err := r.db.Exec(ctx, query, lokasiBaru, id)
	return err
}

func (r *pgxArchiveRepository) CreateMovement(ctx context.Context, m *models.ArchiveMovement) error {
	query := `
		INSERT INTO archive_movements (
			id, archive_id, type, petugas_id, peminjam_id, tanggal, 
			tanggal_kembali_estimasi, tanggal_kembali_aktual, tujuan, keterangan, status
		) VALUES ($1, $2, $3, $4, NULLIF($5, '')::uuid, $6, $7, $8, $9, $10, $11)
	`
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	if m.Tanggal.IsZero() {
		m.Tanggal = time.Now()
	}
	_, err := r.db.Exec(ctx, query,
		m.ID, m.ArchiveID, m.Type, m.PetugasID, m.PeminjamID, m.Tanggal,
		m.TanggalKembaliEstimasi, m.TanggalKembaliAktual, m.Tujuan, m.Keterangan, m.Status,
	)
	return err
}

func (r *pgxArchiveRepository) ListMovements(ctx context.Context, archiveID string) ([]models.ArchiveMovement, error) {
	query := `
		SELECT id, archive_id, type, petugas_id, COALESCE(peminjam_id::text, ''), tanggal, 
		       tanggal_kembali_estimasi, tanggal_kembali_aktual, COALESCE(tujuan, ''), COALESCE(keterangan, ''), status
		FROM archive_movements WHERE archive_id = $1 ORDER BY tanggal DESC
	`
	rows, err := r.db.Query(ctx, query, archiveID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.ArchiveMovement
	for rows.Next() {
		var m models.ArchiveMovement
		if err := rows.Scan(
			&m.ID, &m.ArchiveID, &m.Type, &m.PetugasID, &m.PeminjamID, &m.Tanggal,
			&m.TanggalKembaliEstimasi, &m.TanggalKembaliAktual, &m.Tujuan, &m.Keterangan, &m.Status,
		); err != nil {
			return nil, err
		}
		list = append(list, m)
	}
	return list, rows.Err()
}

func (r *pgxArchiveRepository) GetRetentionSchedule(ctx context.Context, kategori string) (*models.RetentionSchedule, error) {
	query := `SELECT id, kategori, retensi_tahun, COALESCE(keterangan, ''), created_at FROM retention_schedules WHERE kategori = $1`
	var rs models.RetentionSchedule
	err := r.db.QueryRow(ctx, query, kategori).Scan(&rs.ID, &rs.Kategori, &rs.RetensiTahun, &rs.Keterangan, &rs.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &rs, nil
}
