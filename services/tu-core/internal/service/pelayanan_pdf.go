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

func GeneratePelayananPDF(school *repository.School, pengajuan *repository.PengajuanLayanan, nomorSurat, qrPath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetMargins(20, 20, 20)

	// Kop Surat
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "PEMERINTAH PROVINSI JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "DINAS PENDIDIKAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, school.Name, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s, %s", school.Address, school.Kelurahan, school.Kecamatan), "", 1, "C", false, 0, "")
	pdf.Line(20, 45, 190, 45)
	pdf.Line(20, 46, 190, 46)
	pdf.Ln(10)

	// Judul
	pdf.SetFont("Arial", "BU", 12)
	pdf.CellFormat(0, 8, pengajuan.LayananNama, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(0, 7, "Nomor: "+nomorSurat, "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Isi
	pdf.SetFont("Arial", "", 11)
	pdf.MultiCell(0, 6, "Yang bertanda tangan di bawah ini Kepala "+school.Name+", dengan ini menerangkan bahwa:", "", "L", false)
	pdf.Ln(5)

	// Data Pegawai
	pdf.CellFormat(50, 6, "Nama", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 6, pengajuan.NamaPegawai, "", 1, "L", false, 0, "")
	pdf.CellFormat(50, 6, "NIP", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 6, pengajuan.NIP, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Isi Spesifik
	if pengajuan.TujuanSurat != "" {
		pdf.MultiCell(0, 6, "Surat keterangan ini diterbitkan untuk keperluan: "+pengajuan.TujuanSurat, "", "L", false)
	}
	if pengajuan.Kepada != "" {
		pdf.MultiCell(0, 6, "Ditujukan kepada: "+pengajuan.Kepada, "", "L", false)
	}
	if pengajuan.Keterangan != "" {
		pdf.MultiCell(0, 6, "Keterangan: "+pengajuan.Keterangan, "", "L", false)
	}
	pdf.Ln(10)

	// Tanda Tangan
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

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("surat_pelayanan_%s.pdf", pengajuan.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}

func GenerateQRPelayanan(nomorSurat string) (string, error) {
	qr, err := qrcode.New(nomorSurat, qrcode.Medium)
	if err != nil { return "", err }
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("qr_pelayanan_%d.png", time.Now().UnixNano()))
	err = qr.WriteFile(256, tempPath)
	return tempPath, err
}
