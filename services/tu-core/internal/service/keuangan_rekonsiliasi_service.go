package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// RekonsiliasiService — orkestrasi auto-match & saldo adjustment
// =====================================================================
type RekonsiliasiService struct {
	repo repository.RekonsiliasiRepository
}

func NewRekonsiliasiService(repo repository.RekonsiliasiRepository) *RekonsiliasiService {
	return &RekonsiliasiService{repo: repo}
}

// =====================================================================
// JalankanRekonsiliasi — main flow
//
// 1. Validasi input
// 2. Cek duplikasi (per school×period×rekening×bulan×tahun)
// 3. Parse mutasi items ke MutasiParsed
// 4. Load BKU untuk periode & rekening
// 5. Auto-match
// 6. Compute saldo adjustment
// 7. Build header + details + unmatched
// 8. Insert atomik via repo
// 9. Reload & return
// =====================================================================
func (s *RekonsiliasiService) JalankanRekonsiliasi(
	ctx context.Context, req models.CreateRekonsiliasiRequest,
) (*models.RekonsiliasiBank, error) {
	if err := validateRekonsiliasiRequest(&req); err != nil {
		return nil, err
	}

	// Cek duplikasi
	exists, err := s.repo.ExistsRekonsiliasi(
		ctx, req.SchoolID, req.AcademicPeriodID, req.RekeningBankID,
		req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf(
			"rekonsiliasi periode %d/%d untuk rekening %s sudah ada",
			req.PeriodeBulan, req.PeriodeTahun, req.RekeningBankID,
		)
	}

	// Parse mutasi
	mutasi, err := parseMutasiItems(req.MutasiItems)
	if err != nil {
		return nil, err
	}

	// Load BKU & saldo awal
	bkuItems, err := s.repo.GetBkuForPeriod(
		ctx, req.SchoolID, req.AcademicPeriodID, req.RekeningBankID,
		req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}
	saldoBuku, err := s.repo.GetSaldoBkuByPeriode(
		ctx, req.SchoolID, req.AcademicPeriodID, req.RekeningBankID,
		req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil {
		return nil, err
	}

	// Auto-match
	matched, unmatchedBank, unmatchedBuku := autoMatchMutasi(mutasi, bkuItems)

	// Compute saldo adjustment
	saldoBank := lastMutasiSaldo(mutasi)
	saldoRekBank, saldoRekBuku, isBalance := computeSaldoRekonsiliasi(
		saldoBuku, saldoBank, unmatchedBank, unmatchedBuku,
	)
	status := models.RekonUnbalance
	if isBalance {
		status = models.RekonBalance
	}

	// Build header
	h := &models.RekonsiliasiBank{
		SchoolID:              req.SchoolID,
		AcademicPeriodID:      req.AcademicPeriodID,
		RekeningBankID:        req.RekeningBankID,
		PeriodeBulan:          req.PeriodeBulan,
		PeriodeTahun:          req.PeriodeTahun,
		SaldoBuku:             saldoBuku,
		SaldoBank:             saldoBank,
		Selisih:               saldoBank - saldoBuku,
		JumlahMatched:         len(matched),
		JumlahUnmatchedBank:   len(unmatchedBank),
		JumlahUnmatchedBuku:   len(unmatchedBuku),
		SaldoRekonsiliasiBank: saldoRekBank,
		SaldoRekonsiliasiBuku: saldoRekBuku,
		Status:                status,
		Catatan:               req.Catatan,
		CreatedBy:             req.CreatedBy,
	}

	// Build detail & unmatched rows
	details := buildRekonsiliasiDetails(matched)
	unmatched := buildRekonsiliasiUnmatched(unmatchedBank, unmatchedBuku)

	// Insert atomik
	if err := s.repo.CreateRekonsiliasi(ctx, h, details, unmatched); err != nil {
		return nil, err
	}

	// Reload untuk response lengkap
	return s.repo.GetRekonsiliasiByID(ctx, req.SchoolID, h.ID)
}

// =====================================================================
// Read wrappers
// =====================================================================
func (s *RekonsiliasiService) GetRekonsiliasi(
	ctx context.Context, schoolID, id string,
) (*models.RekonsiliasiBank, error) {
	return s.repo.GetRekonsiliasiByID(ctx, schoolID, id)
}

func (s *RekonsiliasiService) ListRekonsiliasi(
	ctx context.Context, filter models.RekonsiliasiFilter,
) ([]models.RekonsiliasiBank, int, error) {
	return s.repo.ListRekonsiliasi(ctx, filter)
}

func (s *RekonsiliasiService) DeleteRekonsiliasi(
	ctx context.Context, schoolID, id string,
) error {
	return s.repo.DeleteRekonsiliasi(ctx, schoolID, id)
}

// =====================================================================
// Validator
// =====================================================================
func validateRekonsiliasiRequest(r *models.CreateRekonsiliasiRequest) error {
	if r.SchoolID == "" {
		return errors.New("school_id wajib diisi")
	}
	if r.AcademicPeriodID == "" {
		return errors.New("academic_period_id wajib diisi")
	}
	if r.RekeningBankID == "" {
		return errors.New("rekening_bank_id wajib diisi")
	}
	if r.PeriodeBulan < 1 || r.PeriodeBulan > 12 {
		return errors.New("periode_bulan harus 1-12")
	}
	if r.PeriodeTahun < 2000 || r.PeriodeTahun > 2100 {
		return errors.New("periode_tahun tidak valid")
	}
	if len(r.MutasiItems) == 0 {
		return errors.New("mutasi_items tidak boleh kosong")
	}
	return nil
}

// =====================================================================
// Helper — parse MutasiInput ke MutasiParsed
// =====================================================================
func parseMutasiItems(items []models.MutasiInput) ([]models.MutasiParsed, error) {
	out := make([]models.MutasiParsed, 0, len(items))
	for i, in := range items {
		if in.Tanggal == "" {
			return nil, fmt.Errorf("mutasi #%d: tanggal wajib diisi", i+1)
		}
		tgl, err := time.Parse("2006-01-02", in.Tanggal)
		if err != nil {
			return nil, fmt.Errorf("mutasi #%d: tanggal harus YYYY-MM-DD: %w", i+1, err)
		}
		if in.Debet < 0 || in.Kredit < 0 {
			return nil, fmt.Errorf("mutasi #%d: debet/kredit tidak boleh negatif", i+1)
		}
		if in.Debet == 0 && in.Kredit == 0 {
			return nil, fmt.Errorf("mutasi #%d: debet atau kredit harus > 0", i+1)
		}
		out = append(out, models.MutasiParsed{
			Tanggal:    tgl,
			Keterangan: in.Keterangan,
			Debet:      in.Debet,
			Kredit:     in.Kredit,
			Saldo:      in.Saldo,
			Referensi:  in.Referensi,
		})
	}
	return out, nil
}

// =====================================================================
// Helper — ambil saldo dari mutasi terakhir (berdasarkan tanggal)
// =====================================================================
func lastMutasiSaldo(mutasi []models.MutasiParsed) float64 {
	if len(mutasi) == 0 {
		return 0
	}
	last := mutasi[0]
	for _, m := range mutasi[1:] {
		if m.Tanggal.After(last.Tanggal) {
			last = m
		}
	}
	return last.Saldo
}
