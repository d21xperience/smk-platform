package service

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type SKPPDFGenerator struct{}

func NewSKPPDFGenerator() *SKPPDFGenerator {
	return &SKPPDFGenerator{}
}

func (g *SKPPDFGenerator) GenerateQR(data string) (string, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("skp_qr_%d.png", time.Now().UnixNano()))
	err = qr.WriteFile(256, tempPath)
	return tempPath, err
}

func (g *SKPPDFGenerator) GenerateSKPPDF(school *repository.School, penilaian *repository.SKPPenilaian, hasil *repository.SKPHasil, detail []repository.SKPDetailItem, qrPath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	
	// --- PAGE 1: COVER ---
	pdf.AddPage()
	
	// Kop Surat
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "PEMERINTAH PROVINSI JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "DINAS PENDIDIKAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, school.Name, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s, %s", school.Address, school.Kelurahan, school.Kecamatan), "", 1, "C", false, 0, "")
	pdf.Line(15, 40, 195, 40)
	pdf.Line(15, 41, 195, 41)
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, "SASARAN KINERJA PEGAWAI (SKP)", "", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 8, fmt.Sprintf("Periode: %s", penilaian.PeriodeID), "", 1, "C", false, 0, "")
	pdf.Ln(15)

	// Identitas
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(40, 7, "Nama Pegawai", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 7, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 7, penilaian.TeacherName, "", 1, "L", false, 0, "")
	
	pdf.CellFormat(40, 7, "Status Penilaian", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 7, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 7, penilaian.Status, "", 1, "L", false, 0, "")
	pdf.Ln(15)

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "HASIL AKHIR", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 15, fmt.Sprintf("%.2f", hasil.NilaiAkhir), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 10, hasil.PredikatAkhir, "", 1, "C", false, 0, "")

	// --- PAGE 2: DETAIL PENILAIAN ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "A. PERILAKU KERJA (Bobot 50%)", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(10, 8, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(100, 8, "Aspek Perilaku", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Nilai", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Keterangan", "1", 1, "C", true, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	no := 1
	for _, d := range detail {
		if d.Kelompok == "PERILAKU" {
			pdf.CellFormat(10, 8, fmt.Sprintf("%d", no), "1", 0, "C", false, 0, "")
			pdf.CellFormat(100, 8, d.Aspek, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 8, fmt.Sprintf("%d", d.Nilai), "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 8, "", "1", 1, "C", false, 0, "")
			no++
		}
	}
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, fmt.Sprintf("Rata-rata Perilaku: %.2f (%s)", hasil.NilaiPerilaku, hasil.PredikatPerilaku), "", 1, "R", false, 0, "")

	// Hasil Kerja
	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "B. HASIL KERJA", "", 1, "L", false, 0, "")
	pdf.Ln(3)
	
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(10, 8, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(100, 8, "Aspek Hasil Kerja", "1", 0, "L", true, 0, "")
	pdf.CellFormat(30, 8, "Nilai", "1", 0, "C", true, 0, "")
	pdf.CellFormat(40, 8, "Keterangan", "1", 1, "C", true, 0, "")
	
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	no = 1
	for _, d := range detail {
		if d.Kelompok == "HASIL_UTAMA" || d.Kelompok == "HASIL_TAMBAHAN" {
			pdf.CellFormat(10, 8, fmt.Sprintf("%d", no), "1", 0, "C", false, 0, "")
			pdf.CellFormat(100, 8, d.Aspek, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 8, fmt.Sprintf("%d", d.Nilai), "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 8, "", "1", 1, "C", false, 0, "")
			no++
		}
	}

	// --- PAGE 3: TANDA TANGAN ---
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "C. PENGESAHAN", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(90, 7, "Catatan Penilai:", "", 1, "L", false, 0, "")
	pdf.MultiCell(0, 6, hasil.CatatanPenilai, "", "L", false)
	pdf.Ln(5)
	pdf.CellFormat(90, 7, "Catatan Atasan Berikat:", "", 1, "L", false, 0, "")
	pdf.MultiCell(0, 6, hasil.CatatanBerikat, "", "L", false)
	
	pdf.Ln(15)
	pdf.SetX(120)
	pdf.CellFormat(0, 6, fmt.Sprintf("Ditetapkan di: %s", school.KabKota), "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(0, 6, fmt.Sprintf("Pada tanggal: %s", time.Now().Format("02 January 2006")), "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(0, 6, "Kepala Sekolah,", "", 1, "L", false, 0, "")
	
	if qrPath != "" {
		pdf.ImageOptions(qrPath, 125, pdf.GetY()+2, 25, 25, false, gofpdf.ImageOptions{ImageType: "PNG"}, 0, "")
		pdf.Ln(28)
	} else {
		pdf.Ln(20)
	}

	pdf.SetX(120)
	pdf.SetFont("Arial", "BU", 11)
	pdf.CellFormat(0, 6, school.PrincipalName, "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("NIP: %s", school.PrincipalNIP), "", 1, "L", false, 0, "")

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("skp_formulir_%s.pdf", penilaian.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}
