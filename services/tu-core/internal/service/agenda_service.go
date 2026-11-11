package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/agenda/models"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/xuri/excelize/v2"
)

type AgendaService struct {
	repo repository.AgendaRepository
}

func NewAgendaService(repo repository.AgendaRepository) *AgendaService {
	return &AgendaService{repo: repo}
}

func (s *AgendaService) GenerateAndCreateAgenda(ctx context.Context, correspondenceID, jenis string) (*models.AgendaEntry, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	tahun := time.Now().Year()
	nomorAgenda, nomorUrut, err := s.repo.GenerateNomorAgenda(ctx, opCtx.SchoolID, jenis, tahun)
	if err != nil {
		return nil, NewServiceError("GenerateAndCreateAgenda", "Agenda", "", "failed to generate number", err)
	}

	entry := &models.AgendaEntry{
		SchoolID:         opCtx.SchoolID,
		CorrespondenceID: correspondenceID,
		Jenis:            jenis,
		NomorAgenda:      nomorAgenda,
		NomorUrut:        nomorUrut,
		Tahun:            tahun,
		PetugasID:        opCtx.UserID,
	}

	if err := s.repo.CreateAgendaEntry(ctx, entry); err != nil {
		return nil, NewServiceError("GenerateAndCreateAgenda", "Agenda", "", "failed to save", err)
	}

	return entry, nil
}

func (s *AgendaService) GetByCorrespondenceID(ctx context.Context, correspondenceID string) (*models.AgendaEntry, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.GetByCorrespondenceID(ctx, opCtx.SchoolID, correspondenceID)
}

func (s *AgendaService) ListAgenda(ctx context.Context, jenis string, bulan, tahun int) ([]map[string]interface{}, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}
	return s.repo.ListAgenda(ctx, repository.AgendaFilter{
		SchoolID: opCtx.SchoolID,
		Jenis:    jenis,
		Bulan:    bulan,
		Tahun:    tahun,
	})
}

func (s *AgendaService) GenerateExcel(ctx context.Context, jenis string, bulan, tahun int) (string, error) {
	list, err := s.ListAgenda(ctx, jenis, bulan, tahun)
	if err != nil {
		return "", err
	}

	f := excelize.NewFile()
	sheet := "Buku Agenda"
	index, _ := f.NewSheet(sheet)
	f.SetActiveSheet(index)
	f.DeleteSheet("Sheet1")

	// Header
	title := fmt.Sprintf("BUKU AGENDA SURAT %s - %d/%d", jenis, bulan, tahun)
	f.SetCellValue(sheet, "A1", title)
	f.MergeCell(sheet, "A1", "I1")

	headers := []string{"No", "No. Agenda", "Tgl Catat", "Jenis", "No. Surat", "Tgl Surat", "Perihal", "Pengirim/Penerima", "Klasifikasi"}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 3)
		f.SetCellValue(sheet, cell, h)
	}

	// Style Header
	style, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D9E1F2"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	f.SetCellStyle(sheet, "A3", "I3", style)

	// Data
	for i, row := range list {
		r := i + 4
		f.SetCellValue(sheet, fmt.Sprintf("A%d", r), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", r), row["nomor_agenda"])
		f.SetCellValue(sheet, fmt.Sprintf("C%d", r), row["tanggal_catat"])
		f.SetCellValue(sheet, fmt.Sprintf("D%d", r), row["jenis"])
		f.SetCellValue(sheet, fmt.Sprintf("E%d", r), row["number"])
		f.SetCellValue(sheet, fmt.Sprintf("F%d", r), row["date"])
		f.SetCellValue(sheet, fmt.Sprintf("G%d", r), row["subject"])
		f.SetCellValue(sheet, fmt.Sprintf("H%d", r), fmt.Sprintf("%s / %s", row["from_party"], row["to_party"]))
		f.SetCellValue(sheet, fmt.Sprintf("I%d", r), row["kode_klasifikasi"])
	}

	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "C", 20)
	f.SetColWidth(sheet, "D", "D", 10)
	f.SetColWidth(sheet, "E", "F", 15)
	f.SetColWidth(sheet, "G", "G", 30)
	f.SetColWidth(sheet, "H", "H", 30)

	saveDir := "./storage/reports/agenda"
	_ = os.MkdirAll(saveDir, os.ModePerm)
	fileName := fmt.Sprintf("Buku_Agenda_%s_%d_%02d.xlsx", jenis, tahun, bulan)
	filePath := filepath.Join(saveDir, fileName)

	if err := f.SaveAs(filePath); err != nil {
		return "", err
	}

	return fmt.Sprintf("/storage/reports/agenda/%s", fileName), nil
}
