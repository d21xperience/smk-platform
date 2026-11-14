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
// PengeluaranService — workflow 3 tingkat (Opsi A):
//   Pemohon → Bendahara → Kepala Sekolah → auto BKK + auto BKU
// =====================================================================
type PengeluaranService struct {
	repo repository.KeuanganTransaksiRepository
}

func NewPengeluaranService(repo repository.KeuanganTransaksiRepository) *PengeluaranService {
	return &PengeluaranService{repo: repo}
}

// ---------------------------------------------------------------------
// DTO
// ---------------------------------------------------------------------
type AjukanPengeluaranRequest struct {
	SchoolID         string  `json:"school_id"`
	AcademicPeriodID string  `json:"academic_period_id"`
	PemohonUserID    string  `json:"pemohon_user_id"`
	Keperluan        string  `json:"keperluan"`
	Jumlah           float64 `json:"jumlah"`
	Prioritas        string  `json:"prioritas"`
	CatatanPemohon   string  `json:"catatan_pemohon,omitempty"`
	Tanggal          string  `json:"tanggal,omitempty"`
}

type VerifikasiBendaharaRequest struct {
	SchoolID    string `json:"school_id"`
	PengajuanID string `json:"pengajuan_id"`
	Aksi        string `json:"aksi"` // SETUJUI | TOLAK | REVISI
	BendaharaID string `json:"bendahara_id"`
	Catatan     string `json:"catatan,omitempty"`
}

type ApprovalKepalaRequest struct {
	SchoolID    string `json:"school_id"`
	PengajuanID string `json:"pengajuan_id"`
	Aksi        string `json:"aksi"` // SETUJUI | TOLAK
	KepalaID    string `json:"kepala_id"`
	Catatan     string `json:"catatan,omitempty"`
}

type ApprovalKepalaResponse struct {
	Pengajuan   *models.PengajuanPengeluaran `json:"pengajuan"`
	Pengeluaran *models.Pengeluaran          `json:"pengeluaran,omitempty"`
	BuktiKas    *models.BukuKasUmum          `json:"bukti_kas,omitempty"`
}

// ---------------------------------------------------------------------
// AjukanPengeluaran — status: MENUNGGU_VERIFIKASI_BENDAHARA
// ---------------------------------------------------------------------
func (s *PengeluaranService) AjukanPengeluaran(
	ctx context.Context, req AjukanPengeluaranRequest,
) (*models.PengajuanPengeluaran, error) {
	if err := validateAjukanPengeluaran(&req); err != nil {
		return nil, err
	}

	tgl, err := parseTanggalOrDefault(req.Tanggal)
	if err != nil {
		return nil, err
	}

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

	nomor, err := s.repo.GenerateNomorPengajuan(ctx, tx, req.SchoolID, req.AcademicPeriodID, tgl)
	if err != nil {
		return nil, err
	}

	p := &models.PengajuanPengeluaran{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		NomorPengajuan:   nomor,
		Tanggal:          tgl.Format("2006-01-02"),
		PemohonUserID:    req.PemohonUserID,
		Keperluan:        req.Keperluan,
		Jumlah:           req.Jumlah,
		Prioritas:        req.Prioritas,
		Status:           models.PengajuanMenungguBendahara,
		CatatanPemohon:   req.CatatanPemohon,
	}
	if err := s.repo.CreatePengajuan(ctx, tx, p); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("gagal commit: %w", err)
	}
	committed = true
	return p, nil
}

// ---------------------------------------------------------------------
// VerifikasiBendahara — SETUJUI → MENUNGGU_PERSETUJUAN_KEPALA
//                       TOLAK   → DITOLAK_BENDAHARA
//                       REVISI  → PERLU_REVISI
// ---------------------------------------------------------------------
func (s *PengeluaranService) VerifikasiBendahara(
	ctx context.Context, req VerifikasiBendaharaRequest,
) (*models.PengajuanPengeluaran, error) {
	if req.SchoolID == "" || req.PengajuanID == "" || req.BendaharaID == "" {
		return nil, errors.New("school_id, pengajuan_id, bendahara_id wajib diisi")
	}

	p, err := s.repo.GetPengajuanByID(ctx, req.SchoolID, req.PengajuanID)
	if err != nil {
		return nil, err
	}
	if p.Status != models.PengajuanMenungguBendahara {
		return nil, fmt.Errorf("status pengajuan bukan MENUNGGU_VERIFIKASI_BENDAHARA (saat ini: %s)", p.Status)
	}

	var newStatus string
	switch req.Aksi {
	case "SETUJUI":
		newStatus = models.PengajuanMenungguKepala
	case "TOLAK":
		newStatus = models.PengajuanDitolakBendahara
	case "REVISI":
		newStatus = models.PengajuanPerluRevisi
	default:
		return nil, fmt.Errorf("aksi tidak valid: %s (SETUJUI/TOLAK/REVISI)", req.Aksi)
	}

	if !models.IsValidTransisiPengajuan(p.Status, newStatus) {
		return nil, fmt.Errorf("transisi tidak diizinkan: %s → %s", p.Status, newStatus)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var catatanBendahara *string
	if req.Catatan != "" {
		catatanBendahara = &req.Catatan
	}
	if err := s.repo.UpdatePengajuanStatus(
		ctx, tx, p.ID, newStatus,
		&req.BendaharaID, nil, catatanBendahara, nil, nil,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	committed = true

	return s.repo.GetPengajuanByID(ctx, req.SchoolID, req.PengajuanID)
}

// ---------------------------------------------------------------------
// ApprovalKepala — SETUJUI → DISETUJUI + auto BKK + auto BKU
//                   TOLAK   → DITOLAK_KEPALA
// ---------------------------------------------------------------------
func (s *PengeluaranService) ApprovalKepala(
	ctx context.Context, req ApprovalKepalaRequest,
) (*ApprovalKepalaResponse, error) {
	if req.SchoolID == "" || req.PengajuanID == "" || req.KepalaID == "" {
		return nil, errors.New("school_id, pengajuan_id, kepala_id wajib diisi")
	}

	p, err := s.repo.GetPengajuanByID(ctx, req.SchoolID, req.PengajuanID)
	if err != nil {
		return nil, err
	}
	if p.Status != models.PengajuanMenungguKepala {
		return nil, fmt.Errorf("status pengajuan bukan MENUNGGU_PERSETUJUAN_KEPALA (saat ini: %s)", p.Status)
	}

	if req.Aksi == "TOLAK" {
		return s.approvalTolak(ctx, req, p)
	}
	if req.Aksi != "SETUJUI" {
		return nil, fmt.Errorf("aksi tidak valid: %s (SETUJUI/TOLAK)", req.Aksi)
	}

	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	tgl := time.Now()
	noBKK, err := s.repo.GenerateNomorBKK(ctx, tx, req.SchoolID, p.AcademicPeriodID, tgl)
	if err != nil {
		return nil, err
	}

	// Update pengajuan → DISETUJUI + set no_bkk (atomic)
	var catatanKepala *string
	if req.Catatan != "" {
		catatanKepala = &req.Catatan
	}
	if err := s.repo.UpdatePengajuanStatus(
		ctx, tx, p.ID, models.PengajuanDisetujui,
		nil, &req.KepalaID, nil, catatanKepala, &noBKK,
	); err != nil {
		return nil, err
	}

	// Insert pengeluaran
	pe := &models.Pengeluaran{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: p.AcademicPeriodID,
		NoBKK:            noBKK,
		Tanggal:          tgl.Format("2006-01-02"),
		PengajuanID:      p.ID,
		Jumlah:           p.Jumlah,
		Keperluan:        p.Keperluan,
		Status:           models.PengeluaranSiapCair,
		CreatedBy:        req.KepalaID,
	}
	if err := s.repo.CreatePengeluaran(ctx, tx, pe); err != nil {
		return nil, err
	}

	// Insert BKU kredit
	bku := &models.BukuKasUmum{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: p.AcademicPeriodID,
		Tanggal:          pe.Tanggal,
		NoBukti:          noBKK,
		Uraian:           "Pengeluaran: " + p.Keperluan,
		KodeRekeningID:   coaBelanjaOperasional,
		Debet:            0,
		Kredit:           p.Jumlah,
		ReferensiTabel:   "pengeluaran",
		ReferensiID:      pe.ID,
		CreatedBy:        req.KepalaID,
	}
	if err := s.repo.InsertBKU(ctx, tx, bku); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	committed = true

	updated, err := s.repo.GetPengajuanByID(ctx, req.SchoolID, req.PengajuanID)
	if err != nil {
		return nil, err
	}
	return &ApprovalKepalaResponse{
		Pengajuan:   updated,
		Pengeluaran: pe,
		BuktiKas:    bku,
	}, nil
}

func (s *PengeluaranService) approvalTolak(
	ctx context.Context, req ApprovalKepalaRequest, p *models.PengajuanPengeluaran,
) (*ApprovalKepalaResponse, error) {
	tx, err := s.repo.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback(ctx)
		}
	}()

	var catatanKepala *string
	if req.Catatan != "" {
		catatanKepala = &req.Catatan
	}
	if err := s.repo.UpdatePengajuanStatus(
		ctx, tx, p.ID, models.PengajuanDitolakKepala,
		nil, &req.KepalaID, nil, catatanKepala, nil,
	); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	committed = true

	updated, err := s.repo.GetPengajuanByID(ctx, req.SchoolID, req.PengajuanID)
	if err != nil {
		return nil, err
	}
	return &ApprovalKepalaResponse{Pengajuan: updated}, nil
}

// ---------------------------------------------------------------------
// Read methods
// ---------------------------------------------------------------------
func (s *PengeluaranService) GetPengajuan(ctx context.Context, schoolID, id string) (*models.PengajuanPengeluaran, error) {
	return s.repo.GetPengajuanByID(ctx, schoolID, id)
}

func (s *PengeluaranService) ListPengajuan(ctx context.Context, filter models.PengajuanFilter) ([]models.PengajuanPengeluaran, int, error) {
	return s.repo.ListPengajuan(ctx, filter)
}

func (s *PengeluaranService) GetPengeluaran(ctx context.Context, schoolID, noBKK string) (*models.Pengeluaran, error) {
	return s.repo.GetPengeluaranByNoBKK(ctx, schoolID, noBKK)
}

func (s *PengeluaranService) ListPengeluaran(ctx context.Context, filter models.PengeluaranFilter) ([]models.Pengeluaran, int, error) {
	return s.repo.ListPengeluaran(ctx, filter)
}

// ---------------------------------------------------------------------
// Validator
// ---------------------------------------------------------------------
func validateAjukanPengeluaran(r *AjukanPengeluaranRequest) error {
	if r.SchoolID == "" {
		return errors.New("school_id wajib diisi")
	}
	if r.AcademicPeriodID == "" {
		return errors.New("academic_period_id wajib diisi")
	}
	if r.PemohonUserID == "" {
		return errors.New("pemohon_user_id wajib diisi")
	}
	if r.Keperluan == "" {
		return errors.New("keperluan wajib diisi")
	}
	if r.Jumlah <= 0 {
		return errors.New("jumlah harus lebih dari 0")
	}
	if r.Prioritas == "" {
		r.Prioritas = models.PrioritasSedang
	}
	if !models.IsValidPrioritas(r.Prioritas) {
		return fmt.Errorf("prioritas tidak valid (RENDAH/SEDANG/TINGGI/DARURAT): %s", r.Prioritas)
	}
	return nil
}
