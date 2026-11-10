package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

type LaporanPDFGenerator struct{}

func NewLaporanPDFGenerator() *LaporanPDFGenerator {
	return &LaporanPDFGenerator{}
}

func (g *LaporanPDFGenerator) GenerateQR(data string) (string, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("qr_laporan_%d.png", time.Now().UnixNano()))
	err = qr.WriteFile(256, tempPath)
	return tempPath, err
}

func (g *LaporanPDFGenerator) drawKopSurat(pdf *gofpdf.Fpdf, school *repository.School) {
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "PEMERINTAH PROVINSI JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "DINAS PENDIDIKAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, school.Name, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s, %s", school.Address, school.Kelurahan, school.Kecamatan), "", 1, "C", false, 0, "")
	pdf.Line(15, 40, 195, 40)
	pdf.Line(15, 41, 195, 41)
	pdf.Ln(10)
}

func (g *LaporanPDFGenerator) GenerateLaporanPDF(school *repository.School, laporan *repository.LaporanBerkala) (string, error) {
	var data LaporanData
	_ = json.Unmarshal(laporan.DataJSON, &data)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	g.drawKopSurat(pdf, school)

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "LAPORAN KEPEGAWAIAN BERKALA", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 8, fmt.Sprintf("TRIWULAN %d TAHUN %d", laporan.Triwulan, laporan.Tahun), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, fmt.Sprintf("Periode: %s s.d. %s", laporan.TanggalMulai.Format("02 Jan 2006"), laporan.TanggalSelesai.Format("02 Jan 2006")), "", "L", false)
	pdf.Ln(5)

	// Tabel Data
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(80, 8, "Indikator", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 8, "Nilai", "1", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)

	rows := [][]string{
		{"Total Pegawai", fmt.Sprintf("%d", data.TotalPegawai)},
		{"Rata-rata Kehadiran", fmt.Sprintf("%.2f%%", data.RataKehadiran)},
		{"Total Alpha", fmt.Sprintf("%d", data.TotalAlpha)},
		{"Rata-rata SKP", fmt.Sprintf("%.2f", data.RataSkp)},
		{"Total Cuti", fmt.Sprintf("%d", data.TotalCuti)},
		{"Dokumen Kadaluarsa", fmt.Sprintf("%d", data.TotalDokumen)},
		{"Guru Beban < 24 JTM", fmt.Sprintf("%d", data.BebanKurang24)},
	}
	for _, row := range rows {
		pdf.CellFormat(80, 8, row[0], "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 8, row[1], "1", 1, "C", false, 0, "")
	}

	pdf.Ln(15)
	pdf.SetX(120)
	pdf.CellFormat(0, 6, fmt.Sprintf("%s, %s", school.KabKota, time.Now().Format("02 January 2006")), "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.CellFormat(0, 6, "Kepala Sekolah,", "", 1, "L", false, 0, "")
	pdf.Ln(15)
	pdf.SetX(120)
	pdf.SetFont("Arial", "BU", 11)
	pdf.CellFormat(0, 6, school.PrincipalName, "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("NIP: %s", school.PrincipalNIP), "", 1, "L", false, 0, "")

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("laporan_utama_%s.pdf", laporan.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}

func (g *LaporanPDFGenerator) GenerateRingkasanPDF(school *repository.School, laporan *repository.LaporanBerkala) (string, error) {
	var data LaporanData
	_ = json.Unmarshal(laporan.DataJSON, &data)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	g.drawKopSurat(pdf, school)

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "RINGKASAN EKSEKUTIF", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, fmt.Sprintf("Pada Triwulan %d Tahun %d, total pegawai aktif adalah %d orang. Rata-rata kehadiran mencapai %.2f%% dengan total alpha sebanyak %d hari.", laporan.Triwulan, laporan.Tahun, data.TotalPegawai, data.RataKehadiran, data.TotalAlpha), "", "J", false)
	pdf.Ln(5)
	pdf.MultiCell(0, 6, fmt.Sprintf("Evaluasi kinerja (SKP) menunjukkan rata-rata nilai %.2f. Terdapat %d guru dengan beban mengajar kurang dari 24 JTM yang perlu perhatian khusus.", data.RataSkp, data.BebanKurang24), "", "J", false)

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("ringkasan_%s.pdf", laporan.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}

func (g *LaporanPDFGenerator) GenerateSuratPengantar(school *repository.School, laporan *repository.LaporanBerkala) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	g.drawKopSurat(pdf, school)

	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 7, fmt.Sprintf("Nomor: 421.2/%d/Q%d/%d", time.Now().Unix()%1000, laporan.Triwulan, laporan.Tahun), "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "Perihal: Pengiriman Laporan Kepegawaian Berkala", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.CellFormat(0, 7, "Kepada Yth,", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "Kepala Dinas Pendidikan / Pengurus Yayasan", "", 1, "L", false, 0, "")
	pdf.CellFormat(0, 7, "Di Tempat", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.MultiCell(0, 7, "Dengan hormat,", "", "L", false)
	pdf.Ln(3)
	pdf.MultiCell(0, 7, fmt.Sprintf("Bersama surat ini kami sampaikan Laporan Kepegawaian Triwulan %d Tahun %d untuk lingkungan %s. Laporan ini mencakup data kehadiran, kinerja, dan beban mengajar pegawai.", laporan.Triwulan, laporan.Tahun, school.Name), "", "J", false)
	pdf.Ln(3)
	pdf.MultiCell(0, 7, "Demikian laporan ini kami sampaikan, atas perhatian dan kerja samanya kami ucapkan terima kasih.", "", "J", false)

	pdf.Ln(15)
	pdf.SetX(120)
	pdf.CellFormat(0, 6, "Kepala Sekolah,", "", 1, "L", false, 0, "")
	pdf.Ln(15)
	pdf.SetX(120)
	pdf.SetFont("Arial", "BU", 11)
	pdf.CellFormat(0, 6, school.PrincipalName, "", 1, "L", false, 0, "")
	pdf.SetX(120)
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 6, fmt.Sprintf("NIP: %s", school.PrincipalNIP), "", 1, "L", false, 0, "")

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("surat_pengantar_%s.pdf", laporan.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}
