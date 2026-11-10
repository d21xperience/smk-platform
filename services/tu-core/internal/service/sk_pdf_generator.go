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

type SKPDFGenerator struct{}

func NewSKPDFGenerator() *SKPDFGenerator {
	return &SKPDFGenerator{}
}

// GenerateQR membuat QR Code TTE sementara dan mengembalikan path-nya
func (g *SKPDFGenerator) GenerateQR(data string) (string, error) {
	qr, err := qrcode.New(data, qrcode.Medium)
	if err != nil {
		return "", err
	}
	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("qr_tte_%d.png", time.Now().UnixNano()))
	err = qr.WriteFile(256, tempPath)
	return tempPath, err
}

// GenerateSKPDF membuat file PDF utama (SK) dan mengembalikan path sementaranya
func (g *SKPDFGenerator) GenerateSKPDF(school *repository.School, sk *repository.SKDetail, qrPath string) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// --- KOP SURAT ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "PEMERINTAH PROVINSI JAWA TIMUR", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "DINAS PENDIDIKAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, school.Name, "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("%s, %s, %s", school.Address, school.Kelurahan, school.Kecamatan), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Telp: %s | Email: %s | NPSN: %s", school.Phone, school.Email, school.NPSN), "", 1, "C", false, 0, "")
	
	// Garis Kop
	pdf.Line(15, 40, 195, 40)
	pdf.Line(15, 41, 195, 41)
	pdf.Ln(10)

	// --- JUDUL SK ---
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 7, "KEPUTUSAN KEPALA SEKOLAH", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, fmt.Sprintf("Nomor: %s", sk.NomorSK), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "TENTANG", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 7, "PEMBAGIAN TUGAS GURU DAN KARYAWAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(0, 7, fmt.Sprintf("SEMESTER %d TAHUN AJARAN %s", sk.Semester, sk.TahunAjaran), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// --- KONSIDERANS ---
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(20, 6, "Menimbang", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.MultiCell(0, 6, "Bahwa dalam rangka kelancaran pelaksanaan proses belajar mengajar, perlu menetapkan pembagian tugas guru dan karyawan.", "", "J", false)
	pdf.CellFormat(20, 6, "Mengingat", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.MultiCell(0, 6, "1. Undang-Undang Nomor 20 Tahun 2003 tentang Sistem Pendidikan Nasional\n2. Peraturan Pemerintah Nomor 19 Tahun 2005\n3. Kalender Pendidikan Tahun Ajaran "+sk.TahunAjaran, "", "J", false)
	pdf.Ln(5)

	// --- DIKTUM ---
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 6, "MEMUTUSKAN:", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 11)
	pdf.CellFormat(20, 6, "Menetapkan", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.MultiCell(0, 6, "Pembagian Tugas Guru dan Karyawan sebagaimana tercantum dalam lampiran keputusan ini.", "", "J", false)
	pdf.CellFormat(20, 6, "Pertama", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.MultiCell(0, 6, "Menugaskan nama-nama tersebut dalam lampiran untuk melaksanakan tugas mengajar.", "", "J", false)
	pdf.CellFormat(20, 6, "Kedua", "", 0, "L", false, 0, "")
	pdf.CellFormat(5, 6, ":", "", 0, "C", false, 0, "")
	pdf.MultiCell(0, 6, "Keputusan ini berlaku sejak tanggal ditetapkan dan apabila di kemudian hari terdapat kekeliruan akan diadakan perbaikan sebagaimana mestinya.", "", "J", false)
	pdf.Ln(15)

	// --- TANDA TANGAN & QR ---
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

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("sk_utama_%s.pdf", sk.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}

// GenerateLampiranPDF membuat file PDF Lampiran matriks
func (g *SKPDFGenerator) GenerateLampiranPDF(school *repository.School, sk *repository.SKDetail) (string, error) {
	pdf := gofpdf.New("L", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 10, "LAMPIRAN KEPUTUSAN KEPALA SEKOLAH", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Nomor: %s", sk.NomorSK), "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, "TENTANG PEMBAGIAN TUGAS GURU DAN KARYAWAN", "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Table Header
	pdf.SetFillColor(230, 230, 230)
	pdf.SetFont("Arial", "B", 10)
	headers := []string{"No", "Nama Guru", "NIP", "Mata Pelajaran", "Kelas", "Jam/Minggu"}
	widths := []float64{10, 60, 40, 60, 40, 30}

	for i, h := range headers {
		pdf.CellFormat(widths[i], 8, h, "1", 0, "C", true, 0, "")
	}
	pdf.Ln(-1)

	// Table Body
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for i, a := range sk.Assignments {
		pdf.CellFormat(widths[0], 8, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(widths[1], 8, a.GuruName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[2], 8, a.GuruNIP, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[3], 8, a.MataPelajaran, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[4], 8, a.Kelas, "1", 0, "L", false, 0, "")
		pdf.CellFormat(widths[5], 8, fmt.Sprintf("%d", a.JamPerMinggu), "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	// Tugas Tambahan (jika ada)
	if len(sk.TugasTambahan) > 0 {
		pdf.Ln(10)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(0, 8, "TUGAS TAMBAHAN:", "", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		for _, t := range sk.TugasTambahan {
			pdf.CellFormat(10, 8, "-", "", 0, "C", false, 0, "")
			pdf.CellFormat(60, 8, t.GuruName, "", 0, "L", false, 0, "")
			pdf.CellFormat(100, 8, fmt.Sprintf("%s (%s)", t.JenisTugas, t.Keterangan), "", 1, "L", false, 0, "")
		}
	}

	tempPath := filepath.Join(os.TempDir(), fmt.Sprintf("lampiran_sk_%s.pdf", sk.ID))
	return tempPath, pdf.OutputFileAndClose(tempPath)
}
