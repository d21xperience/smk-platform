package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

type MaintenanceService struct {
	repo repository.MaintenanceRepository
}

func NewMaintenanceService(repo repository.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{repo: repo}
}

func (s *MaintenanceService) CreateDamageReport(ctx context.Context, d *repository.DamageReport) (*repository.DamageReport, error) {
	if d.Description == "" || d.ReporterName == "" || d.AssetName == "" {
		return nil, errors.New("nama aset, pelapor, dan deskripsi wajib diisi")
	}
	reqNo, _ := s.repo.GenerateDamageNumber(ctx)
	d.ReportNumber = reqNo
	d.Status = "Pending"
	return d, s.repo.CreateDamageReport(ctx, d)
}

func (s *MaintenanceService) GetAllDamageReports(ctx context.Context, semesterID string) ([]*repository.DamageReport, error) {
	return s.repo.GetAllDamageReports(ctx, semesterID)
}

func (s *MaintenanceService) UpdateDamageReport(ctx context.Context, id string, status, vendor string, cost float64, notes string) error {
	if status == "" {
		return errors.New("status wajib diisi")
	}
	return s.repo.UpdateDamageReport(ctx, id, status, vendor, cost, notes)
}

func (s *MaintenanceService) CreateSchedule(ctx context.Context, sch *repository.MaintenanceSchedule) (*repository.MaintenanceSchedule, error) {
	if sch.AssetName == "" || sch.Frequency == "" {
		return nil, errors.New("nama aset dan frekuensi wajib diisi")
	}
	code, _ := s.repo.GenerateScheduleCode(ctx)
	sch.ScheduleCode = code
	sch.IsActive = true
	if sch.NextServiceDate.IsZero() {
		sch.NextServiceDate = time.Now()
	}
	return sch, s.repo.CreateSchedule(ctx, sch)
}

func (s *MaintenanceService) GetAllSchedules(ctx context.Context, semesterID string) ([]*repository.MaintenanceSchedule, error) {
	return s.repo.GetAllSchedules(ctx, semesterID)
}

func (s *MaintenanceService) CompleteSchedule(ctx context.Context, id string, serviceDate time.Time, cost float64, notes string, performedBy string) error {
	if serviceDate.IsZero() { serviceDate = time.Now() }
	return s.repo.CompleteSchedule(ctx, id, serviceDate, cost, notes, performedBy)
}

func (s *MaintenanceService) GetMaintenanceLogs(ctx context.Context, startDate, endDate string) ([]*repository.MaintenanceLog, error) {
	return s.repo.GetMaintenanceLogs(ctx, startDate, endDate)
}

func (s *MaintenanceService) GenerateMaintenanceBookPDF(ctx context.Context, startDate, endDate string) (*bytes.Buffer, error) {
	logs, err := s.repo.GetMaintenanceLogs(ctx, startDate, endDate)
	if err != nil { return nil, err }

	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "BUKU PERAWATAN & PERBAIKAN ASET", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "SMK Pasundan Jatinangor", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Periode: %s s.d. %s", startDate, endDate), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 7, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 7, "Kode Aset", "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, 7, "Nama Aset", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Vendor/Teknisi", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Biaya (Rp)", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 7, "Catatan", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	var totalCost float64
	for i, l := range logs {
		pdf.CellFormat(10, 6, fmt.Sprint(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 6, l.ServiceDate.Format("02/01/2006"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(50, 6, l.AssetCode, "1", 0, "L", false, 0, "")
		pdf.CellFormat(70, 6, l.AssetName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 6, l.Vendor, "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("%.0f", l.Cost), "1", 0, "R", false, 0, "")
		pdf.CellFormat(50, 6, l.Notes, "1", 1, "L", false, 0, "")
		totalCost += l.Cost
	}

	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(230, 7, "TOTAL BIAYA PERAWATAN", "1", 0, "R", false, 0, "")
	pdf.CellFormat(30, 7, fmt.Sprintf("%.0f", totalCost), "1", 1, "R", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil { return nil, err }
	return &buf, nil
}
