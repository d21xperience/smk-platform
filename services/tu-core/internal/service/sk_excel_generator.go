package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/xuri/excelize/v2"
)

type SKExcelGenerator struct{}

func NewSKExcelGenerator() *SKExcelGenerator {
	return &SKExcelGenerator{}
}

// GenerateRekapExcel membuat file Excel rekap pembagian tugas dan mengembalikan path-nya
func (g *SKExcelGenerator) GenerateRekapExcel(sk *repository.SKDetail) (string, error) {
	f := excelize.NewFile()
	sheet := "Rekap_Tugas"
	f.SetSheetName("Sheet1", sheet)

	// Style Definisi
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 12},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#D3D3D3"}, Pattern: 1},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	cellStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
		},
	})
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	// Judul
	f.MergeCell(sheet, "A1", "F1")
	f.SetCellValue(sheet, "A1", "REKAP PEMBAGIAN TUGAS GURU DAN KARYAWAN")
	f.SetCellStyle(sheet, "A1", "F1", titleStyle)

	f.SetCellValue(sheet, "A2", fmt.Sprintf("Nomor SK: %s", sk.NomorSK))
	semesterLabel := fmt.Sprintf("Semester %d", sk.Semester); if sk.Semester == 0 { semesterLabel = "Semester Ganjil & Genap (Tahun Penuh)" }; f.SetCellValue(sheet, "A3", fmt.Sprintf("%s | Tahun Ajaran: %s", semesterLabel, sk.TahunAjaran))
	f.SetCellValue(sheet, "A4", fmt.Sprintf("Tanggal Cetak: %s", time.Now().Format("02 January 2006")))
	f.SetRowHeight(sheet, 1, 25)

	// Header Tabel
	headers := []string{"No", "Nama Guru", "NIP", "Mata Pelajaran", "Kelas", "Jam/Minggu"}
	for i, h := range headers {
		col := string(rune(65 + i)) // A, B, C, D, E, F
		f.SetCellValue(sheet, col+"6", h)
		f.SetCellStyle(sheet, col+"6", col+"6", headerStyle)
	}

	// Body Tabel
	row := 7
	for i, a := range sk.Assignments {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), a.GuruName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), a.GuruNIP)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), a.MataPelajaran)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), a.Kelas)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), a.JamPerMinggu)

		for col := 65; col <= 70; col++ {
			c := string(rune(col))
			f.SetCellStyle(sheet, fmt.Sprintf("%s%d", c, row), fmt.Sprintf("%s%d", c, row), cellStyle)
		}
		row++
	}

	// Tugas Tambahan (jika ada)
	if len(sk.TugasTambahan) > 0 {
		row += 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), "TUGAS TAMBAHAN")
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("A%d", row), titleStyle)
		row++

		for i, t := range sk.TugasTambahan {
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t.GuruName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t.JenisTugas)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), t.Keterangan)
			for col := 65; col <= 68; col++ {
				c := string(rune(col))
				f.SetCellStyle(sheet, fmt.Sprintf("%s%d", c, row), fmt.Sprintf("%s%d", c, row), cellStyle)
			}
			row++
		}
	}

	// Auto-fit Column Widths
	f.SetColWidth(sheet, "A", "A", 5)
	f.SetColWidth(sheet, "B", "B", 30)
	f.SetColWidth(sheet, "C", "C", 20)
	f.SetColWidth(sheet, "D", "D", 25)
	f.SetColWidth(sheet, "E", "E", 15)
	f.SetColWidth(sheet, "F", "F", 15)

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("rekap_sk_%s_%d.xlsx", sk.ID, time.Now().UnixNano()))
	if err := f.SaveAs(tempPath); err != nil {
		return "", err
	}
	return tempPath, nil
}
