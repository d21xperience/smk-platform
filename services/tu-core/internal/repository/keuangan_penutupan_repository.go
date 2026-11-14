package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// =====================================================================
// PenutupanRepository — interface CRUD + snapshot + lock/unlock
// =====================================================================
type PenutupanRepository interface {
	Create(ctx context.Context, p *models.PenutupanBukuKas) error
	GetByID(ctx context.Context, schoolID, id string) (*models.PenutupanBukuKas, error)
	List(ctx context.Context, filter models.PenutupanFilter) ([]models.PenutupanBukuKas, int, error)
	ExistsByPeriode(ctx context.Context, schoolID, periodID string, bulan, tahun int) (bool, error)
	GetSnapshotBku(ctx context.Context, schoolID, periodID string, bulan, tahun int) (*models.SnapshotBku, error)
	GetSnapshotPosisi(ctx context.Context, schoolID, periodID string) (*models.SnapshotPosisi, error)
	UnlockBku(ctx context.Context, schoolID, periodID string, bulan, tahun int) error

	// FASE 5B-2: E-Sign
	SignPenutupan(ctx context.Context, schoolID, id, hash, signerID string, signedAt time.Time) error
	GetPenutupanByTTDHash(ctx context.Context, hash string) (*models.PenutupanBukuKas, error)
}

type pgxPenutupanRepository struct {
	db *pgxpool.Pool
}

func NewPenutupanRepository(db *pgxpool.Pool) PenutupanRepository {
	return &pgxPenutupanRepository{db: db}
}

// =====================================================================
// Create — atomik: insert penutupan + lock semua BKU di periode tsb
// =====================================================================
func (r *pgxPenutupanRepository) Create(
	ctx context.Context, p *models.PenutupanBukuKas,
) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("gagal begin tx: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// 1. Insert penutupan
	err = tx.QueryRow(ctx, `
		INSERT INTO penutupan_buku_kas (
			id, school_id, academic_period_id, periode_bulan, periode_tahun,
			saldo_awal, total_debet, total_kredit, saldo_akhir_buku,
			saldo_fisik_kas, selisih_kas,
			total_aset, total_kewajiban, total_ekuitas,
			status, catatan, closed_by, closed_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4,
			$5, $6, $7, $8,
			$9, $10,
			$11, $12, $13,
			$14, NULLIF($15,''), NULLIF($16,''), NOW()
		)
		RETURNING id, closed_at, created_at, updated_at
	`, p.SchoolID, p.AcademicPeriodID, p.PeriodeBulan, p.PeriodeTahun,
		p.SaldoAwal, p.TotalDebet, p.TotalKredit, p.SaldoAkhirBuku,
		p.SaldoFisikKas, p.SelisihKas,
		p.TotalAset, p.TotalKewajiban, p.TotalEkuitas,
		p.Status, p.Catatan, p.ClosedBy,
	).Scan(&p.ID, &p.ClosedAt, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("gagal insert penutupan_buku_kas: %w", err)
	}

	// 2. Lock semua BKU di periode tsb (bulan + tahun)
	_, err = tx.Exec(ctx, `
		UPDATE buku_kas_umum
		SET is_locked = TRUE
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND EXTRACT(MONTH FROM tanggal) = $3
		  AND EXTRACT(YEAR FROM tanggal) = $4
		  AND is_locked = FALSE
	`, p.SchoolID, p.AcademicPeriodID, p.PeriodeBulan, p.PeriodeTahun)
	if err != nil {
		return fmt.Errorf("gagal lock BKU: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("gagal commit: %w", err)
	}
	committed = true
	return nil
}

// =====================================================================
// UnlockBku — untuk rollback / buka kembali penutupan
// (dipakai FASE 5 saat implementasi "buka kembali")
// =====================================================================
func (r *pgxPenutupanRepository) UnlockBku(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) error {
	_, err := r.db.Exec(ctx, `
		UPDATE buku_kas_umum
		SET is_locked = FALSE
		WHERE school_id = $1
		  AND academic_period_id = $2
		  AND EXTRACT(MONTH FROM tanggal) = $3
		  AND EXTRACT(YEAR FROM tanggal) = $4
	`, schoolID, periodID, bulan, tahun)
	if err != nil {
		return fmt.Errorf("gagal unlock BKU: %w", err)
	}
	return nil
}
