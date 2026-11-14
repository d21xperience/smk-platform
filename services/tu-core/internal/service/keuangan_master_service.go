package service

import (
    "context"
    "errors"
    "strings"

    "sekolah-platform/services/tu-core/internal/domain/finance/models"
    "sekolah-platform/services/tu-core/internal/repository"
)

// =====================================================================
// KeuanganMasterService — orkestrasi master data keuangan
// =====================================================================
type KeuanganMasterService struct {
    repo repository.KeuanganMasterRepository
}

func NewKeuanganMasterService(repo repository.KeuanganMasterRepository) *KeuanganMasterService {
    return &KeuanganMasterService{repo: repo}
}

// =====================================================================
// Jenis Iuran
// =====================================================================
func (s *KeuanganMasterService) CreateJenisIuran(ctx context.Context, j *models.JenisIuran) error {
    if err := validateJenisIuran(j); err != nil {
        return err
    }
    return s.repo.CreateJenisIuran(ctx, j)
}

func (s *KeuanganMasterService) UpdateJenisIuran(ctx context.Context, j *models.JenisIuran) error {
    if j.ID == "" {
        return errors.New("id wajib diisi")
    }
    if err := validateJenisIuran(j); err != nil {
        return err
    }
    return s.repo.UpdateJenisIuran(ctx, j)
}

func (s *KeuanganMasterService) DeleteJenisIuran(ctx context.Context, schoolID, id string) error {
    return s.repo.DeleteJenisIuran(ctx, schoolID, id)
}

func (s *KeuanganMasterService) GetJenisIuran(ctx context.Context, schoolID, id string) (*models.JenisIuran, error) {
    return s.repo.GetJenisIuranByID(ctx, schoolID, id)
}

func (s *KeuanganMasterService) ListJenisIuran(ctx context.Context, filter models.JenisIuranFilter) ([]models.JenisIuran, int, error) {
    return s.repo.ListJenisIuran(ctx, filter)
}

// =====================================================================
// Tarif Iuran
// =====================================================================
func (s *KeuanganMasterService) CreateTarifIuran(ctx context.Context, t *models.TarifIuran) error {
    if err := validateTarifIuran(t); err != nil {
        return err
    }
    return s.repo.CreateTarifIuran(ctx, t)
}

func (s *KeuanganMasterService) UpdateTarifIuran(ctx context.Context, t *models.TarifIuran) error {
    if t.ID == "" {
        return errors.New("id wajib diisi")
    }
    if err := validateTarifIuran(t); err != nil {
        return err
    }
    return s.repo.UpdateTarifIuran(ctx, t)
}

func (s *KeuanganMasterService) DeleteTarifIuran(ctx context.Context, schoolID, id string) error {
    return s.repo.DeleteTarifIuran(ctx, schoolID, id)
}

func (s *KeuanganMasterService) GetTarifIuran(ctx context.Context, schoolID, id string) (*models.TarifIuran, error) {
    return s.repo.GetTarifIuranByID(ctx, schoolID, id)
}

func (s *KeuanganMasterService) ListTarifIuran(ctx context.Context, filter models.TarifIuranFilter) ([]models.TarifIuran, int, error) {
    return s.repo.ListTarifIuran(ctx, filter)
}

// =====================================================================
// Kode Rekening
// =====================================================================
func (s *KeuanganMasterService) CreateKodeRekening(ctx context.Context, k *models.KodeRekening) error {
    if err := validateKodeRekening(k); err != nil {
        return err
    }
    return s.repo.CreateKodeRekening(ctx, k)
}

func (s *KeuanganMasterService) UpdateKodeRekening(ctx context.Context, k *models.KodeRekening) error {
    if k.ID == "" {
        return errors.New("id wajib diisi")
    }
    if err := validateKodeRekening(k); err != nil {
        return err
    }
    return s.repo.UpdateKodeRekening(ctx, k)
}

func (s *KeuanganMasterService) DeleteKodeRekening(ctx context.Context, schoolID, id string) error {
    return s.repo.DeleteKodeRekening(ctx, schoolID, id)
}

func (s *KeuanganMasterService) GetKodeRekening(ctx context.Context, schoolID, id string) (*models.KodeRekening, error) {
    return s.repo.GetKodeRekeningByID(ctx, schoolID, id)
}

func (s *KeuanganMasterService) ListKodeRekening(ctx context.Context, filter models.KodeRekeningFilter) ([]models.KodeRekening, int, error) {
    return s.repo.ListKodeRekening(ctx, filter)
}

// =====================================================================
// Rekening Bank
// =====================================================================
func (s *KeuanganMasterService) CreateRekeningBank(ctx context.Context, rb *models.RekeningBank) error {
    if err := validateRekeningBank(rb); err != nil {
        return err
    }
    return s.repo.CreateRekeningBank(ctx, rb)
}

func (s *KeuanganMasterService) UpdateRekeningBank(ctx context.Context, rb *models.RekeningBank) error {
    if rb.ID == "" {
        return errors.New("id wajib diisi")
    }
    if err := validateRekeningBank(rb); err != nil {
        return err
    }
    return s.repo.UpdateRekeningBank(ctx, rb)
}

func (s *KeuanganMasterService) DeleteRekeningBank(ctx context.Context, schoolID, id string) error {
    return s.repo.DeleteRekeningBank(ctx, schoolID, id)
}

func (s *KeuanganMasterService) GetRekeningBank(ctx context.Context, schoolID, id string) (*models.RekeningBank, error) {
    return s.repo.GetRekeningBankByID(ctx, schoolID, id)
}

func (s *KeuanganMasterService) ListRekeningBank(ctx context.Context, filter models.RekeningBankFilter) ([]models.RekeningBank, int, error) {
    return s.repo.ListRekeningBank(ctx, filter)
}

// =====================================================================
// Validators
// =====================================================================
func validateJenisIuran(j *models.JenisIuran) error {
    if j == nil {
        return errors.New("payload kosong")
    }
    if j.SchoolID == "" {
        return errors.New("school_id wajib diisi")
    }
    if strings.TrimSpace(j.Kode) == "" {
        return errors.New("kode wajib diisi")
    }
    if strings.TrimSpace(j.Nama) == "" {
        return errors.New("nama wajib diisi")
    }
    if !models.IsValidKategoriIuran(j.Kategori) {
        return errors.New("kategori tidak valid (SPP/GEDUNG/EKSTRA/SERAGAM/KEGIATAN/LAINNYA)")
    }
    return nil
}

func validateTarifIuran(t *models.TarifIuran) error {
    if t == nil {
        return errors.New("payload kosong")
    }
    if t.SchoolID == "" {
        return errors.New("school_id wajib diisi")
    }
    if t.AcademicPeriodID == "" {
        return errors.New("academic_period_id wajib diisi")
    }
    if t.JenisIuranID == "" {
        return errors.New("jenis_iuran_id wajib diisi")
    }
    if t.Nominal < 0 {
        return errors.New("nominal tidak boleh negatif")
    }
    if t.TahunBerlaku < 2000 || t.TahunBerlaku > 2100 {
        return errors.New("tahun_berlaku tidak valid")
    }
    return nil
}

func validateKodeRekening(k *models.KodeRekening) error {
    if k == nil {
        return errors.New("payload kosong")
    }
    if k.SchoolID == "" {
        return errors.New("school_id wajib diisi")
    }
    if strings.TrimSpace(k.Kode) == "" {
        return errors.New("kode wajib diisi")
    }
    if strings.TrimSpace(k.Nama) == "" {
        return errors.New("nama wajib diisi")
    }
    if !models.IsValidTipeKodeRekening(k.Tipe) {
        return errors.New("tipe tidak valid (ASET/KEWAJIBAN/EKUITAS/PENDAPATAN/BELANJA)")
    }
    if k.Level < 1 {
        k.Level = 1
    }
    return nil
}

func validateRekeningBank(rb *models.RekeningBank) error {
    if rb == nil {
        return errors.New("payload kosong")
    }
    if rb.SchoolID == "" {
        return errors.New("school_id wajib diisi")
    }
    if strings.TrimSpace(rb.NamaBank) == "" {
        return errors.New("nama_bank wajib diisi")
    }
    if strings.TrimSpace(rb.NomorRekening) == "" {
        return errors.New("nomor_rekening wajib diisi")
    }
    if strings.TrimSpace(rb.AtasNama) == "" {
        return errors.New("atas_nama wajib diisi")
    }
    return nil
}
