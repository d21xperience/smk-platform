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
// Konstanta COA default (dari seed migration 125).
// TODO FASE 3: refactor jadi lookup by kode via masterRepo.
// =====================================================================
const (
	coaPendapatanSPP      = "f3000000-0000-0000-0000-000000000032" // 4.1.1 Pendapatan SPP
	coaBelanjaOperasional = "f8000000-0000-0002-0000-000000000001" // 5.1.1 Belanja ATK & Operasional
)

// =====================================================================
// IuranService — pembayaran iuran siswa & manajemen tagihan
// =====================================================================
type IuranService struct {
	repo       repository.KeuanganTransaksiRepository
	masterRepo repository.KeuanganMasterRepository
}

func NewIuranService(
	repo repository.KeuanganTransaksiRepository,
	masterRepo repository.KeuanganMasterRepository,
) *IuranService {
	return &IuranService{repo: repo, masterRepo: masterRepo}
}

// ---------------------------------------------------------------------
// Request / Response DTO
// ---------------------------------------------------------------------
type TerimaPembayaranRequest struct {
	SchoolID         string  `json:"school_id"`
	AcademicPeriodID string  `json:"academic_period_id"`
	NIS              string  `json:"nis"`
	JenisIuranID     string  `json:"jenis_iuran_id"`
	PeriodeBulan     int     `json:"periode_bulan"`
	PeriodeTahun     int     `json:"periode_tahun"`
	Jumlah           float64 `json:"jumlah"`
	MetodeBayar      string  `json:"metode_bayar"`
	NoReferensi      string  `json:"no_referensi,omitempty"`
	Keterangan       string  `json:"keterangan,omitempty"`
	Tanggal          string  `json:"tanggal,omitempty"`
	CreatedBy        string  `json:"created_by,omitempty"`
}

type TerimaPembayaranResponse struct {
	Penerimaan *models.PenerimaanIuran `json:"penerimaan"`
	Tagihan    *models.TagihanSiswa    `json:"tagihan"`
	BuktiKas   *models.BukuKasUmum     `json:"bukti_kas"`
}

// ---------------------------------------------------------------------
// TerimaPembayaran — orkestrasi atomik:
//   1. Validasi input
//   2. Lookup siswa by NIS
//   3. Cari/create tagihan
//   4. BEGIN TX
//   5. Generate no_bkm
//   6. Insert penerimaan_iuran
//   7. Update tagihan status (LUNAS / SEBAGIAN)
//   8. Insert BKU (debit)
//   9. COMMIT
// ---------------------------------------------------------------------
func (s *IuranService) TerimaPembayaran(
	ctx context.Context, req TerimaPembayaranRequest,
) (*TerimaPembayaranResponse, error) {
	// 1. Validasi input
	if err := validateTerimaPembayaran(&req); err != nil {
		return nil, err
	}

	tgl, err := parseTanggalOrDefault(req.Tanggal)
	if err != nil {
		return nil, err
	}
	req.Tanggal = tgl.Format("2006-01-02")

	// 2. Lookup siswa
	siswa, err := s.repo.FindSiswaByNIS(ctx, req.SchoolID, req.NIS)
	if err != nil {
		return nil, err
	}

	// Lookup jenis iuran untuk nama di BKU uraian
	jenis, err := s.masterRepo.GetJenisIuranByID(ctx, req.SchoolID, req.JenisIuranID)
	if err != nil {
		return nil, fmt.Errorf("jenis iuran tidak ditemukan: %w", err)
	}

	// 3. Cari tagihan existing
	existingTagihan, err := s.repo.FindTagihanByNISPeriode(
		ctx, req.SchoolID, req.AcademicPeriodID,
		req.NIS, req.JenisIuranID, req.PeriodeBulan, req.PeriodeTahun,
	)
	if err != nil && !errors.Is(err, repository.ErrTransaksiKeuanganNotFound) {
		return nil, err
	}

	// 4. BEGIN TX
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("gagal begin transaction: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	// Cek duplikasi sebelum insert
	if existingTagihan != nil && existingTagihan.Status == models.TagihanLunas {
		return nil, fmt.Errorf("tagihan %s %s periode %d/%d untuk NIS %s sudah LUNAS",
			jenis.Nama, req.JenisIuranID, req.PeriodeBulan, req.PeriodeTahun, req.NIS)
	}

	// 5. Kalau tagihan belum ada → create dulu
	tagihan := existingTagihan
	if tagihan == nil {
		tagihan = &models.TagihanSiswa{
			SchoolID:          req.SchoolID,
			AcademicPeriodID:  req.AcademicPeriodID,
			PesertaDidikID:    siswa.PesertaDidikID,
			NIS:               req.NIS,
			KelasNama:         siswa.KelasNama,
			JenisIuranID:      req.JenisIuranID,
			JenisIuranNama:    jenis.Nama,
			PeriodeBulan:      req.PeriodeBulan,
			PeriodeTahun:      req.PeriodeTahun,
			Nominal:           req.Jumlah,
			TotalDibayar:      0,
			Status:            models.TagihanBelumLunas,
			TanggalJatuhTempo: tgl.AddDate(0, 1, 0).Format("2006-01-02"),
			CreatedBy:         req.CreatedBy,
		}
		if err := s.repo.CreateTagihan(ctx, tx, tagihan); err != nil {
			return nil, err
		}
	}

	// 6. Generate No BKM
	noBKM, err := s.repo.GenerateNomorBKM(ctx, tx, req.SchoolID, req.AcademicPeriodID, tgl)
	if err != nil {
		return nil, err
	}

	// 7. Insert penerimaan
	penerimaan := &models.PenerimaanIuran{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		NoBKM:            noBKM,
		Tanggal:          req.Tanggal,
		TagihanSiswaID:   tagihan.ID,
		PesertaDidikID:   siswa.PesertaDidikID,
		NIS:              req.NIS,
		JenisIuranID:     req.JenisIuranID,
		JenisIuranNama:   jenis.Nama,
		PeriodeBulan:     req.PeriodeBulan,
		PeriodeTahun:     req.PeriodeTahun,
		Jumlah:           req.Jumlah,
		MetodeBayar:      req.MetodeBayar,
		NoReferensi:      req.NoReferensi,
		Keterangan:       req.Keterangan,
		Status:           "LUNAS",
		CreatedBy:        req.CreatedBy,
	}
	if err := s.repo.CreatePenerimaan(ctx, tx, penerimaan); err != nil {
		return nil, err
	}

	// 8. Update tagihan status
	newTotalDibayar := tagihan.TotalDibayar + req.Jumlah
	newStatus := models.TagihanSebagian
	var tanggalBayar *time.Time
	if newTotalDibayar >= tagihan.Nominal {
		newStatus = models.TagihanLunas
		tanggalBayar = &tgl
	}
	if err := s.repo.UpdateTagihanStatus(ctx, tx, tagihan.ID, newStatus, newTotalDibayar, tanggalBayar); err != nil {
		return nil, err
	}
	tagihan.TotalDibayar = newTotalDibayar
	tagihan.Status = newStatus

	// 9. Insert BKU (debit)
	bku := &models.BukuKasUmum{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		Tanggal:          req.Tanggal,
		NoBukti:          noBKM,
		Uraian:           fmt.Sprintf("Penerimaan %s - NIS %s (%s) periode %d/%d", jenis.Nama, req.NIS, siswa.NamaSiswa, req.PeriodeBulan, req.PeriodeTahun),
		KodeRekeningID:   coaPendapatanSPP,
		Debet:            req.Jumlah,
		Kredit:           0,
		ReferensiTabel:   "penerimaan_iuran",
		ReferensiID:      penerimaan.ID,
		CreatedBy:        req.CreatedBy,
	}
	if err := s.repo.InsertBKU(ctx, tx, bku); err != nil {
		return nil, err
	}

	// 10. COMMIT
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit: %w", err)
	}
	committed = true

	return &TerimaPembayaranResponse{
		Penerimaan: penerimaan,
		Tagihan:    tagihan,
		BuktiKas:   bku,
	}, nil
}

// ---------------------------------------------------------------------
// Read methods
// ---------------------------------------------------------------------
func (s *IuranService) GetPenerimaan(ctx context.Context, schoolID, noBKM string) (*models.PenerimaanIuran, error) {
	return s.repo.GetPenerimaanByNoBKM(ctx, schoolID, noBKM)
}

func (s *IuranService) ListPenerimaan(ctx context.Context, filter models.PenerimaanFilter) ([]models.PenerimaanIuran, int, error) {
	return s.repo.ListPenerimaan(ctx, filter)
}

func (s *IuranService) GetTagihan(ctx context.Context, schoolID, id string) (*models.TagihanSiswa, error) {
	return s.repo.GetTagihanByID(ctx, schoolID, id)
}

func (s *IuranService) ListTagihan(ctx context.Context, filter models.TagihanFilter) ([]models.TagihanSiswa, int, error) {
	return s.repo.ListTagihan(ctx, filter)
}

// ---------------------------------------------------------------------
// Validators & helpers
// ---------------------------------------------------------------------
func validateTerimaPembayaran(r *TerimaPembayaranRequest) error {
	if r.SchoolID == "" {
		return errors.New("school_id wajib diisi")
	}
	if r.AcademicPeriodID == "" {
		return errors.New("academic_period_id wajib diisi")
	}
	if r.NIS == "" {
		return errors.New("nis wajib diisi")
	}
	if r.JenisIuranID == "" {
		return errors.New("jenis_iuran_id wajib diisi")
	}
	if r.PeriodeBulan < 1 || r.PeriodeBulan > 12 {
		return errors.New("periode_bulan harus 1-12")
	}
	if r.PeriodeTahun < 2000 || r.PeriodeTahun > 2100 {
		return errors.New("periode_tahun tidak valid")
	}
	if r.Jumlah <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	if r.MetodeBayar == "" {
		r.MetodeBayar = models.MetodeTunai
	}
	if !models.IsValidMetodeBayar(r.MetodeBayar) {
		return fmt.Errorf("metode_bayar tidak valid (TUNAI/TRANSFER/QRIS): %s", r.MetodeBayar)
	}
	return nil
}

func parseTanggalOrDefault(s string) (time.Time, error) {
	if s == "" {
		return time.Now(), nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, fmt.Errorf("format tanggal harus YYYY-MM-DD: %w", err)
	}
	return t, nil
}
