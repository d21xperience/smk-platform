package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/xuri/excelize/v2"
)

type LaporanExcelGenerator struct{}

func NewLaporanExcelGenerator() *LaporanExcelGenerator {
	return &LaporanExcelGenerator{}
}

func (g *LaporanExcelGenerator) GenerateExcel(tipe string, tahun, triwulan int, dataJson json.RawMessage) (string, error) {
	f := excelize.NewFile()
	sheet := "Rekap Laporan"
	f.SetSheetName("Sheet1", sheet)

	// Style
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1}, {Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1}, {Type: "bottom", Color: "000000", Style: 1},
		},
	})

	// Header
	f.MergeCell(sheet, "A1", "B1")
	f.SetCellValue(sheet, "A1", fmt.Sprintf("LAPORAN KEPEGAWAIAN TRIWULAN %d TAHUN %d - %s", triwulan, tahun, tipe))
	f.SetCellStyle(sheet, "A1", "B1", headerStyle)

	f.SetCellValue(sheet, "A3", "Indikator")
	f.SetCellValue(sheet, "B3", "Nilai")
	f.SetCellStyle(sheet, "A3", "B3", headerStyle)

	// Parse JSON
	var data LaporanData
	if err := json.Unmarshal(dataJson, &data); err != nil {
		return "", err
	}

	// Isi Data
	rows := [][]interface{}{
		{"Total Pegawai", data.TotalPegawai},
		{"Rata-rata Kehadiran (%)", data.RataKehadiran},
		{"Total Alpha", data.TotalAlpha},
		{"Rata-rata SKP", data.RataSkp},
		{"Total Cuti", data.TotalCuti},
		{"Dokumen Kadaluarsa", data.TotalDokumen},
		{"Guru Beban < 24 JTM", data.BebanKurang24},
		{"Catatan", data.CatatanTambahan},
	}

	for i, row := range rows {
		cellA := fmt.Sprintf("A%d", i+4)
		cellB := fmt.Sprintf("B%d", i+4)
		f.SetCellValue(sheet, cellA, row[0])
		f.SetCellValue(sheet, cellB, row[1])
	}

	f.SetColWidth(sheet, "A", "A", 30)
	f.SetColWidth(sheet, "B", "B", 20)

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("laporan_%s_Q%d_%d_%d.xlsx", tipe, triwulan, tahun, time.Now().Unix()))
	if err := f.SaveAs(tempPath); err != nil {
		return "", err
	}
	return tempPath, nil
}
