package service

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

// GenerateBukuIndukPDF menghasilkan PDF Buku Induk Siswa (A4 portrait).
// Return: buffer PDF, NIS (untuk nama file), error.
func (s *StudentBukuIndukService) GenerateBukuIndukPDF(ctx context.Context, studentID, schoolID, academicPeriodID string) (*bytes.Buffer, string, error) {
	detail, err := s.GetBukuIndukDetail(ctx, studentID, schoolID, academicPeriodID)
	if err != nil {
		return nil, "", err
	}
	school, err := s.schoolRepo.GetByID(ctx, schoolID)
	if err != nil {
		return nil, "", err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.AddPage()

	drawKopSekolahBukuInduk(pdf, school)

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 8, "BUKU INDUK SISWA", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, fmt.Sprintf("No. Buku Induk: BI/%s/SMK", detail.NIS), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	// ---- A. DATA PRIBADI ----
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "A. DATA PRIBADI", "", 1, "L", false, 0, "")
	writeRow := func(label, value string) {
		if value == "" {
			value = "-"
		}
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(45, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(2, 6, ":", "", 0, "L", false, 0, "")
		pdf.CellFormat(0, 6, value, "", 1, "L", false, 0, "")
	}
	writeLongRow := func(label, value string) {
		if value == "" {
			value = "-"
		}
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(45, 6, label, "", 0, "L", false, 0, "")
		pdf.CellFormat(2, 6, ":", "", 0, "L", false, 0, "")
		pdf.MultiCell(0, 6, value, "", "L", false)
	}
	birthDate := "-"
	if detail.BirthDate != nil {
		birthDate = detail.BirthDate.Format("02 January 2006")
	}
	writeRow("Nama Lengkap", detail.Name)
	writeRow("NIS", detail.NIS)
	writeRow("NISN", detail.NISN)
	writeRow("Tempat Lahir", detail.BirthPlace)
	writeRow("Tanggal Lahir", birthDate)
	writeRow("Jenis Kelamin", detail.Gender)
	writeRow("Agama", detail.Agama)
	writeLongRow("Alamat", detail.Address)
	writeRow("Telepon", detail.Phone)
	pdf.Ln(3)

	// ---- B. DATA ORANG TUA / WALI ----
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "B. DATA ORANG TUA / WALI", "", 1, "L", false, 0, "")
	writeRow("Nama Ayah", detail.NamaAyah)
	writeRow("Pekerjaan Ayah", detail.PekerjaanAyah)
	writeRow("Nama Ibu", detail.NamaIbu)
	writeRow("Pekerjaan Ibu", detail.PekerjaanIbu)
	writeRow("Nama Wali", detail.NamaWali)
	writeRow("Pekerjaan Wali", detail.PekerjaanWali)
	writeLongRow("Alamat Ortu", detail.AlamatOrtu)
	writeRow("Telepon Ortu", detail.TeleponOrtu)
	writeLongRow("Alamat Wali", detail.AlamatWali)
	writeRow("Telepon Wali", detail.TeleponWali)
	pdf.Ln(3)

	// ---- C. RIWAYAT PENDIDIKAN ----
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "C. RIWAYAT PENDIDIKAN", "", 1, "L", false, 0, "")
	diterimaTgl := "-"
	if detail.DiterimaTanggal != nil {
		diterimaTgl = detail.DiterimaTanggal.Format("02 January 2006")
	}
	writeRow("Sekolah Asal", detail.SekolahAsal)
	writeRow("Diterima di Kelas", detail.DiterimaKelas)
	writeRow("Tanggal Diterima", diterimaTgl)
	pdf.Ln(3)

	// ---- D. NILAI RAPOR ----
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 7, "D. NILAI RAPOR", "", 1, "L", false, 0, "")
	pdf.Ln(1)

	// Header tabel nilai
	pdf.SetFont("Arial", "B", 9)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(8, 7, "No", "1", 0, "C", true, 0, "")
	pdf.CellFormat(60, 7, "Mata Pelajaran", "1", 0, "C", true, 0, "")
	pdf.CellFormat(12, 7, "Smt", "1", 0, "C", true, 0, "")
	pdf.CellFormat(18, 7, "Harian", "1", 0, "C", true, 0, "")
	pdf.CellFormat(18, 7, "UTS", "1", 0, "C", true, 0, "")
	pdf.CellFormat(18, 7, "UAS", "1", 0, "C", true, 0, "")
	pdf.CellFormat(18, 7, "Akhir", "1", 0, "C", true, 0, "")
	pdf.CellFormat(20, 7, "Predikat", "1", 1, "C", true, 0, "")

	// Kelompokkan nilai per (semester, mapel)
	type subjectKey struct {
		Semester int
		MapelID  int
		Mapel    string
	}
	grouped := make(map[subjectKey]map[string]float64)
	var order []subjectKey
	for _, sc := range detail.AcademicScores {
		k := subjectKey{Semester: sc.Semester, MapelID: sc.MataPelajaranID, Mapel: sc.NamaMataPelajaran}
		if _, ok := grouped[k]; !ok {
			grouped[k] = make(map[string]float64)
			order = append(order, k)
		}
		grouped[k][sc.AssessmentType] = sc.Score
	}

	pdf.SetFont("Arial", "", 9)
	if len(order) == 0 {
		pdf.CellFormat(172, 6, "Belum ada data nilai.", "1", 1, "C", false, 0, "")
	} else {
		for i, k := range order {
			v := grouped[k]
			daily := v["DAILY"]
			midterm := v["MIDTERM"]
			final := v["FINAL"]
			sum, cnt := 0.0, 0
			if daily > 0 {
				sum += daily
				cnt++
			}
			if midterm > 0 {
				sum += midterm
				cnt++
			}
			if final > 0 {
				sum += final
				cnt++
			}
			avg := 0.0
			if cnt > 0 {
				avg = sum / float64(cnt)
			}
			predikat := "-"
			switch {
			case avg >= 90:
				predikat = "A"
			case avg >= 80:
				predikat = "B"
			case avg >= 70:
				predikat = "C"
			case avg > 0:
				predikat = "D"
			}
			pdf.CellFormat(8, 6, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
			pdf.CellFormat(60, 6, k.Mapel, "1", 0, "L", false, 0, "")
			pdf.CellFormat(12, 6, fmt.Sprintf("%d", k.Semester), "1", 0, "C", false, 0, "")
			pdf.CellFormat(18, 6, studentBukuIndukFormatScore(daily), "1", 0, "C", false, 0, "")
			pdf.CellFormat(18, 6, studentBukuIndukFormatScore(midterm), "1", 0, "C", false, 0, "")
			pdf.CellFormat(18, 6, studentBukuIndukFormatScore(final), "1", 0, "C", false, 0, "")
			pdf.CellFormat(18, 6, studentBukuIndukFormatScore(avg), "1", 0, "C", false, 0, "")
			pdf.CellFormat(20, 6, predikat, "1", 1, "C", false, 0, "")
		}
	}
	pdf.Ln(8)

	// ---- Tanda Tangan ----
	kabKota := school.KabKota
	if kabKota == "" {
		kabKota = school.Kecamatan
	}
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("%s, %s", kabKota, time.Now().Format("02 January 2006")), "", 1, "R", false, 0, "")
	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(0, 6, "Kepala Sekolah", "", 1, "R", false, 0, "")
	pdf.Ln(18)
	namaKepsek := school.PrincipalName
	if namaKepsek == "" {
		namaKepsek = "............................"
	}
	nipKepsek := school.PrincipalNIP
	if nipKepsek == "" {
		nipKepsek = "............................"
	}
	pdf.CellFormat(0, 6, namaKepsek, "", 1, "R", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, fmt.Sprintf("NIP. %s", nipKepsek), "", 1, "R", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, "", err
	}
	return &buf, detail.NIS, nil
}

// drawKopSekolahBukuInduk menggambar kop sekolah generik (dari DB).
// Berbeda dari LaporanPDFGenerator.drawKopSurat yang hardcoded Pemprov Jatim.
func drawKopSekolahBukuInduk(pdf *gofpdf.Fpdf, school *repository.School) {
	if school == nil {
		return
	}
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(0, 6, "YAYASAN PENDIDIKAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 8, strings.ToUpper(school.Name), "", 1, "C", false, 0, "")

	alamat := school.Address
	for _, part := range []string{school.Kelurahan, school.Kecamatan, school.KabKota, school.Propinsi} {
		if part != "" {
			alamat += ", " + part
		}
	}
	pdf.SetFont("Arial", "", 9)
	pdf.CellFormat(0, 5, alamat, "", 1, "C", false, 0, "")

	kontak := ""
	if school.Phone != "" {
		kontak = "Telp. " + school.Phone
	}
	if school.Email != "" {
		if kontak != "" {
			kontak += " | "
		}
		kontak += "Email: " + school.Email
	}
	if kontak != "" {
		pdf.CellFormat(0, 5, kontak, "", 1, "C", false, 0, "")
	}
	pdf.Line(15, 42, 195, 42)
	pdf.Line(15, 43, 195, 43)
	pdf.Ln(8)
}

func studentBukuIndukFormatScore(v float64) string {
	if v == 0 {
		return "-"
	}
	return fmt.Sprintf("%.0f", v)
}
