package service

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jung-kurt/gofpdf"
)

type SarprasReportService struct {
	db *pgxpool.Pool
}

func NewSarprasReportService(db *pgxpool.Pool) *SarprasReportService {
	return &SarprasReportService{db: db}
}

type DashboardStats struct {
	TotalAssets     int     `json:"total_assets"`
	TotalValue      float64 `json:"total_value"`
	ConditionBaik   int     `json:"condition_baik"`
	ConditionRR     int     `json:"condition_rr"`
	ConditionRB     int     `json:"condition_rb"`
	ActiveBorrowings int    `json:"active_borrowings"`
	PendingDisposals int    `json:"pending_disposals"`
	LowStockItems   int     `json:"low_stock_items"`
}

func (s *SarprasReportService) GetDashboardStats(ctx context.Context, semesterID string) (*DashboardStats, error) {
	stats := &DashboardStats{}

	// 1. Asset Stats (Aktif & Dipinjam)
	assetQuery := `
		SELECT 
			COUNT(*), 
			COALESCE(SUM(nilai_buku), 0),
			SUM(CASE WHEN condition = 'Bagus' THEN 1 ELSE 0 END),
			SUM(CASE WHEN condition = 'Rusak Ringan' THEN 1 ELSE 0 END),
			SUM(CASE WHEN condition = 'Rusak Berat' THEN 1 ELSE 0 END)
		FROM inventory_items 
		WHERE status_aset IN ('AKTIF', 'DIPINJAM')
	`
	err := s.db.QueryRow(ctx, assetQuery).Scan(
		&stats.TotalAssets, &stats.TotalValue, 
		&stats.ConditionBaik, &stats.ConditionRR, &stats.ConditionRB,
	)
	if err != nil { return nil, err }

	// 2. Borrowings Aktif (Status PICKED_UP)
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM borrowings WHERE status = 'PICKED_UP' AND semester_id = $1`, semesterID).Scan(&stats.ActiveBorrowings)

	// 3. Disposal Pending (Status DRAFT / VERIFIED_KTU)
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM disposal_requests WHERE status IN ('DRAFT', 'VERIFIED_KTU', 'APPROVED_STEP_1')`).Scan(&stats.PendingDisposals)

	// 4. Logistik Stok Kritis
	s.db.QueryRow(ctx, `SELECT COUNT(*) FROM sarpras_logistics WHERE stok <= stok_minimum`).Scan(&stats.LowStockItems)

	return stats, nil
}

func (s *SarprasReportService) GenerateComprehensivePDF(ctx context.Context, startDate, endDate time.Time) (*bytes.Buffer, error) {
	// 1. Ambil Data Ringkasan
	stats, err := s.GetDashboardStats(ctx, "")
	if err != nil { return nil, err }

	// 2. Ambil Data Aset Masuk (Penerimaan)
	type ReceiptData struct { Number, Supplier string; Date time.Time; Total float64 }
	var receipts []ReceiptData
	recRows, err := s.db.Query(ctx, `
		SELECT bast_number, supplier_name, receipt_date, total_estimated 
		FROM goods_receipts 
		WHERE status = 'VERIFIED' AND receipt_date BETWEEN $1 AND $2
		ORDER BY receipt_date ASC
	`, startDate, endDate)
	if err == nil {
		defer recRows.Close()
		for recRows.Next() {
			var r ReceiptData
			recRows.Scan(&r.Number, &r.Supplier, &r.Date, &r.Total)
			receipts = append(receipts, r)
		}
	}

	// 3. Ambil Data Peminjaman
	type BorrowData struct { Number, Borrower, Asset string; Date time.Time }
	var borrowings []BorrowData
	bwRows, err := s.db.Query(ctx, `
		SELECT b.borrowing_number, b.borrower_name, c.name, b.borrow_date
		FROM borrowings b
		JOIN inventory_items i ON b.inventory_item_id = i.id
		JOIN master_catalogs c ON i.catalog_id = c.id
		WHERE b.borrow_date BETWEEN $1 AND $2
		ORDER BY b.borrow_date ASC
	`, startDate, endDate)
	if err == nil {
		defer bwRows.Close()
		for bwRows.Next() {
			var b BorrowData
			bwRows.Scan(&b.Number, &b.Borrower, &b.Asset, &b.Date)
			borrowings = append(borrowings, b)
		}
	}

	// 4. Generate PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Kop Sekolah
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "LAPORAN KOMPREHENSIF INVENTARIS SARPRAS", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, "SMK Pasundan Jatinangor", "", 1, "C", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf("Periode: %s s.d. %s", startDate.Format("02 Jan 2006"), endDate.Format("02 Jan 2006")), "", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Ringkasan Eksekutif
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "RINGKASAN EKSEKUTIF", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(60, 6, "Total Aset Aktif", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf(": %d Unit", stats.TotalAssets), "", 1, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Total Nilai Aset", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf(": Rp %.2f", stats.TotalValue), "", 1, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Kondisi Aset", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf(": %d Baik | %d Rusak Ringan | %d Rusak Berat", stats.ConditionBaik, stats.ConditionRR, stats.ConditionRB), "", 1, "L", false, 0, "")
	pdf.CellFormat(60, 6, "Peminjaman Aktif", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, fmt.Sprintf(": %d Transaksi", stats.ActiveBorrowings), "", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Tabel Penerimaan
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "A. REKAP PENERIMAAN BARANG", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(40, 7, "No. BAST", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 7, "Supplier/Donor", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Tanggal", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Nilai (Rp)", "1", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 8)
	if len(receipts) == 0 {
		pdf.CellFormat(170, 6, "Tidak ada data penerimaan pada periode ini.", "1", 1, "C", false, 0, "")
	} else {
		var totalRec float64
		for _, r := range receipts {
			pdf.CellFormat(40, 6, r.Number, "1", 0, "L", false, 0, "")
			pdf.CellFormat(60, 6, r.Supplier, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 6, r.Date.Format("02/01/2006"), "1", 0, "C", false, 0, "")
			pdf.CellFormat(40, 6, fmt.Sprintf("%.2f", r.Total), "1", 1, "R", false, 0, "")
			totalRec += r.Total
		}
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(130, 7, "TOTAL", "1", 0, "R", false, 0, "")
		pdf.CellFormat(40, 7, fmt.Sprintf("%.2f", totalRec), "1", 1, "R", false, 0, "")
	}
	pdf.Ln(10)

	// Tabel Peminjaman
	pdf.SetFont("Arial", "B", 11)
	pdf.CellFormat(0, 8, "B. REKAP PEMINJAMAN ASET", "", 1, "L", false, 0, "")
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(40, 7, "No. Peminjaman", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 7, "Peminjam", "1", 0, "C", false, 0, "")
	pdf.CellFormat(60, 7, "Aset", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Tgl Pinjam", "1", 1, "C", false, 0, "")
	
	pdf.SetFont("Arial", "", 8)
	if len(borrowings) == 0 {
		pdf.CellFormat(170, 6, "Tidak ada data peminjaman pada periode ini.", "1", 1, "C", false, 0, "")
	} else {
		for _, b := range borrowings {
			pdf.CellFormat(40, 6, b.Number, "1", 0, "L", false, 0, "")
			pdf.CellFormat(40, 6, b.Borrower, "1", 0, "L", false, 0, "")
			pdf.CellFormat(60, 6, b.Asset, "1", 0, "L", false, 0, "")
			pdf.CellFormat(30, 6, b.Date.Format("02/01/2006"), "1", 1, "C", false, 0, "")
		}
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil { return nil, err }
	return &buf, nil
}
