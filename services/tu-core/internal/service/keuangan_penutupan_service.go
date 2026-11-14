package service

import (
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
	"time"
	"sekolah-platform/services/tu-core/internal/ttd"
)

// =====================================================================
// PenutupanService — orkestrasi penutupan buku kas
// =====================================================================
type PenutupanService struct {
	repo         repository.PenutupanRepository
	rekonRepo    repository.RekonsiliasiRepository
}

func NewPenutupanService(
	repo repository.PenutupanRepository,
	rekonRepo repository.RekonsiliasiRepository,
) *PenutupanService {
	return &PenutupanService{repo: repo, rekonRepo: rekonRepo}
}

// =====================================================================
// TutupBuku — alur utama:
//   1. Validasi input
//   2. Cek duplikasi (periode sudah ditutup?)
//   3. Cek ada rekon BALANCE di bulan tsb (design decision 4B.1.A)
//   4. Hitung snapshot BKU (saldo_awal, debet, kredit, saldo_akhir)
//   5. Hitung snapshot posisi keuangan (aset, kewajiban, ekuitas)
//   6. Hitung selisih = saldo_akhir_buku - saldo_fisik
//   7. Insert atomik + lock BKU via repo
//   8. Reload & return
// =====================================================================
func (s *PenutupanService) TutupBuku(
	ctx context.Context, req models.CreatePenutupanRequest,
) (*models.PenutupanBukuKas, error) {
	if err := validateCreatePenutupan(&req); err != nil {
		return nil, err
	}

	// 1. Cek duplikasi
	exists, err := s.repo.ExistsByPeriode(
		ctx, req.SchoolID, req.AcademicPeriodID, req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf(
			"penutupan periode %d/%d untuk periode akademik %s sudah ada",
			req.PeriodeBulan, req.PeriodeTahun, req.AcademicPeriodID,
		)
	}

	// 2. Cek ada rekon BALANCE di bulan tsb
	hasBalance, err := s.rekonRepo.HasBalanceRekonsiliasi(
		ctx, req.SchoolID, req.AcademicPeriodID, req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}
	if !hasBalance {
		return nil, fmt.Errorf(
			"tidak dapat menutup buku: belum ada rekonsiliasi bank dengan status BALANCE untuk periode %d/%d",
			req.PeriodeBulan, req.PeriodeTahun,
		)
	}

	// 3. Snapshot BKU
	snapshotBku, err := s.repo.GetSnapshotBku(
		ctx, req.SchoolID, req.AcademicPeriodID, req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}

	// 4. Snapshot posisi keuangan
	snapshotPosisi, err := s.repo.GetSnapshotPosisi(ctx, req.SchoolID, req.AcademicPeriodID)
	if err != nil {
		return nil, err
	}

	// 5. Validasi balance posisi keuangan
	if err := validatePosisiKeuangan(snapshotPosisi); err != nil {
		return nil, err
	}

	// 6. Hitung selisih
	selisih := snapshotBku.SaldoAkhirBuku - req.SaldoFisikKas

	// 7. Build & insert
	p := &models.PenutupanBukuKas{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		PeriodeBulan:     req.PeriodeBulan,
		PeriodeTahun:     req.PeriodeTahun,
		SaldoAwal:        snapshotBku.SaldoAwal,
		TotalDebet:       snapshotBku.TotalDebet,
		TotalKredit:      snapshotBku.TotalKredit,
		SaldoAkhirBuku:   snapshotBku.SaldoAkhirBuku,
		SaldoFisikKas:    req.SaldoFisikKas,
		SelisihKas:       selisih,
		TotalAset:        snapshotPosisi.TotalAset,
		TotalKewajiban:   snapshotPosisi.TotalKewajiban,
		TotalEkuitas:     snapshotPosisi.TotalEkuitas,
		Status:           models.PenutupanClosed,
		Catatan:          req.Catatan,
		ClosedBy:         req.ClosedBy,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}

	// 8. Reload untuk response lengkap
	return s.repo.GetByID(ctx, req.SchoolID, p.ID)
}

// =====================================================================
// Read wrappers
// =====================================================================
func (s *PenutupanService) GetPenutupan(
	ctx context.Context, schoolID, id string,
) (*models.PenutupanBukuKas, error) {
	return s.repo.GetByID(ctx, schoolID, id)
}

func (s *PenutupanService) ListPenutupan(
	ctx context.Context, filter models.PenutupanFilter,
) ([]models.PenutupanBukuKas, int, error) {
	return s.repo.List(ctx, filter)
}

// =====================================================================
// Preview endpoints — supaya frontend bisa tampilkan data sebelum confirm
// =====================================================================
func (s *PenutupanService) GetSnapshotBku(
	ctx context.Context, schoolID, periodID string, bulan, tahun int,
) (*models.SnapshotBku, error) {
	return s.repo.GetSnapshotBku(ctx, schoolID, periodID, bulan, tahun)
}

func (s *PenutupanService) GetSnapshotPosisi(
	ctx context.Context, schoolID, periodID string,
) (*models.SnapshotPosisi, error) {
	return s.repo.GetSnapshotPosisi(ctx, schoolID, periodID)
}

// =====================================================================
// Validators
// =====================================================================
func validateCreatePenutupan(r *models.CreatePenutupanRequest) error {
	if r.SchoolID == "" {
		return errors.New("school_id wajib diisi")
	}
	if r.AcademicPeriodID == "" {
		return errors.New("academic_period_id wajib diisi")
	}
	if r.PeriodeBulan < 1 || r.PeriodeBulan > 12 {
		return errors.New("periode_bulan harus 1-12")
	}
	if r.PeriodeTahun < 2000 || r.PeriodeTahun > 2100 {
		return errors.New("periode_tahun tidak valid (2000-2100)")
	}
	// Catatan: SaldoFisikKas TIDAK divalidasi > 0, karena bisa 0 atau
	// bahkan negatif (kalau kas defisit). Yang penting user input
	// nilai yang benar sesuai cash count fisik.
	return nil
}

func validatePosisiKeuangan(p *models.SnapshotPosisi) error {
	// Aset = Kewajiban + Ekuitas (toleransi 0.01 untuk float)
	selisih := p.TotalAset - p.TotalKewajiban - p.TotalEkuitas
	if selisih > 0.01 || selisih < -0.01 {
		return fmt.Errorf(
			"posisi keuangan tidak balance: Aset (%.2f) ≠ Kewajiban (%.2f) + Ekuitas (%.2f), selisih %.2f",
			p.TotalAset, p.TotalKewajiban, p.TotalEkuitas, selisih,
		)
	}
	return nil
}

// =====================================================================
// SignPenutupan — TTD oleh Kepala Sekolah (FASE 5B-2)
//
// Asumsi: PIN verification sudah dilakukan di handler (separation of concern).
// Service ini hanya:
//   1. Load & validasi status dokumen (harus CLOSED)
//   2. Generate hash SHA-256
//   3. Persist hash + signer + timestamp + ubah status
//   4. Return dokumen fresh
// =====================================================================
func (s *PenutupanService) SignPenutupan(
	ctx context.Context, schoolID, id, signerID string,
) (*models.PenutupanBukuKas, error) {
	if schoolID == "" || id == "" || signerID == "" {
		return nil, errors.New("school_id, id, signer_id wajib diisi")
	}

	// 1. Load & cek status
	p, err := s.repo.GetByID(ctx, schoolID, id)
	if err != nil {
		return nil, err
	}
	if p.Status != models.PenutupanClosed {
		return nil, fmt.Errorf("status penutupan bukan CLOSED (saat ini: %s)", p.Status)
	}

	// 2. Generate hash
	signedAt := time.Now().UTC()
	hash := ttd.GenerateSignature("PENUTUPAN", id, signerID, signedAt)

	// 3. Persist
	if err := s.repo.SignPenutupan(ctx, schoolID, id, hash, signerID, signedAt); err != nil {
		return nil, err
	}

	// 4. Reload
	return s.repo.GetByID(ctx, schoolID, id)
}

// =====================================================================
// VerifyByTTDHash — lookup dokumen by hash TTD (endpoint publik)
// =====================================================================
func (s *PenutupanService) VerifyByTTDHash(
	ctx context.Context, hash string,
) (*models.PenutupanBukuKas, error) {
	if hash == "" {
		return nil, errors.New("hash wajib diisi")
	}
	return s.repo.GetPenutupanByTTDHash(ctx, hash)
}

