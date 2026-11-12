package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type KesiswaanReportService struct {
	repo       repository.KesiswaanReportRepository
	schoolRepo repository.SchoolRepository
}

func NewKesiswaanReportService(repo repository.KesiswaanReportRepository, schoolRepo repository.SchoolRepository) *KesiswaanReportService {
	return &KesiswaanReportService{repo: repo, schoolRepo: schoolRepo}
}

type ReportData struct {
	Stats        []repository.JurusanStat `json:"stats"`
	TotalMasuk   int                      `json:"totalMasuk"`
	TotalKeluar  int                      `json:"totalKeluar"`
	TotalAkhir   int                      `json:"totalAkhir"`
	AcademicYear string                   `json:"academicYear"`
}

func (s *KesiswaanReportService) GetStatistik(ctx context.Context, schoolID, academicPeriodID string) (*ReportData, error) {
	stats, err := s.repo.GetStatistikJurusan(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}
	masuk, keluar, err := s.repo.GetTotalMutasi(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, err
	}

	totalAkhir := 0
	for _, st := range stats {
		totalAkhir += st.SiswaAkhir
	}

	return &ReportData{
		Stats:       stats,
		TotalMasuk:  masuk,
		TotalKeluar: keluar,
		TotalAkhir:  totalAkhir,
	}, nil
}

func (s *KesiswaanReportService) GenerateExcel(ctx context.Context, schoolID, academicPeriodID string) (*bytes.Buffer, string, error) {
	data, err := s.GetStatistik(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, "", err
	}
	school, _ := s.schoolRepo.GetByID(ctx, schoolID)
	schoolName := "Sekolah"
	if school != nil {
		schoolName = school.Name
	}

	f := excelize.NewFile()
	sheet := "Rekap Bulanan"
	f.SetSheetName("Sheet1", sheet)

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1},
		},
	})

	f.SetCellValue(sheet, "A1", "LAPORAN BULANAN KESISWAAN")
	f.SetCellValue(sheet, "A2", fmt.Sprintf("Sekolah: %s", schoolName))
	f.SetCellValue(sheet, "A3", fmt.Sprintf("Periode: %s", academicPeriodID))
	f.MergeCell(sheet, "A1", "E1")
	f.MergeCell(sheet, "A2", "E2")
	f.MergeCell(sheet, "A3", "E3")

	headers := []string{"Kompetensi Keahlian", "Siswa Awal", "Mutasi Masuk", "Mutasi Keluar", "Siswa Akhir"}
	for i, h := range headers {
		cell := fmt.Sprintf("%c5", 'A'+i)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}

	for i, st := range data.Stats {
		row := i + 6
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), st.Jurusan)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), st.SiswaAwal)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), st.MutasiMasuk)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), st.MutasiKeluar)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), st.SiswaAkhir)
	}

	f.SetColWidth(sheet, "A", "A", 40)
	f.SetColWidth(sheet, "B", "E", 15)

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("Rekap_Kesiswaan_%s.xlsx", academicPeriodID)
	return &buf, filename, nil
}

func (s *KesiswaanReportService) GeneratePDF(ctx context.Context, schoolID, academicPeriodID string) (*bytes.Buffer, string, error) {
	data, err := s.GetStatistik(ctx, schoolID, academicPeriodID)
	if err != nil {
		return nil, "", err
	}
	school, _ := s.schoolRepo.GetByID(ctx, schoolID)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	drawKopSekolahBukuInduk(pdf, school)

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "LAPORAN BULANAN KESISWAAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("Periode Akademik: %s", academicPeriodID), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Table Header
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(220, 220, 220)
	pdf.CellFormat(70, 8, "Kompetensi Keahlian", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, "Siswa Awal", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, "Mutasi (+)", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, "Mutasi (-)", "1", 0, "C", true, 0, "")
	pdf.CellFormat(26, 8, "Siswa Akhir", "1", 1, "C", true, 0, "")

	// Table Body
	pdf.SetFont("Arial", "", 10)
	for _, st := range data.Stats {
		pdf.CellFormat(70, 7, st.Jurusan, "1", 0, "L", false, 0, "")
		pdf.CellFormat(22, 7, fmt.Sprintf("%d", st.SiswaAwal), "1", 0, "C", false, 0, "")
		pdf.CellFormat(22, 7, fmt.Sprintf("%d", st.MutasiMasuk), "1", 0, "C", false, 0, "")
		pdf.CellFormat(22, 7, fmt.Sprintf("%d", st.MutasiKeluar), "1", 0, "C", false, 0, "")
		pdf.CellFormat(26, 7, fmt.Sprintf("%d", st.SiswaAkhir), "1", 1, "C", false, 0, "")
	}

	// Total Row
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(70, 8, "TOTAL KESELURUHAN", "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, "", "1", 0, "C", true, 0, "") // Awal ditangani terpisah atau dihitung
	pdf.CellFormat(22, 8, fmt.Sprintf("%d", data.TotalMasuk), "1", 0, "C", true, 0, "")
	pdf.CellFormat(22, 8, fmt.Sprintf("%d", data.TotalKeluar), "1", 0, "C", true, 0, "")
	pdf.CellFormat(26, 8, fmt.Sprintf("%d", data.TotalAkhir), "1", 1, "C", true, 0, "")

	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("Dicetak pada: %s", time.Now().Format("02 January 2006 15:04")), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", err
	}
	filename := fmt.Sprintf("Rekap_Kesiswaan_%s.pdf", academicPeriodID)
	return &buf, filename, nil
}
