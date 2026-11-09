package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
)

type GoodsReceiptService struct {
	repo repository.GoodsReceiptRepository
}

func NewGoodsReceiptService(repo repository.GoodsReceiptRepository) *GoodsReceiptService {
	return &GoodsReceiptService{repo: repo}
}

func (s *GoodsReceiptService) CreateReceipt(ctx context.Context, gr *repository.GoodsReceipt) (*repository.GoodsReceipt, error) {
	if gr.SupplierName == "" || len(gr.Items) == 0 {
		return nil, errors.New("supplier dan detail barang wajib diisi")
	}

	bastNo, err := s.repo.GenerateBastNumber(ctx, gr.SemesterID)
	if err != nil {
		return nil, err
	}
	gr.BastNumber = bastNo
	gr.Status = "DRAFT"

	if err := s.repo.Create(ctx, gr); err != nil {
		return nil, err
	}
	return gr, nil
}

func (s *GoodsReceiptService) GetAllReceipts(ctx context.Context, semesterID string) ([]*repository.GoodsReceipt, error) {
	return s.repo.GetAll(ctx, semesterID)
}

func (s *GoodsReceiptService) VerifyReceipt(ctx context.Context, id string, verifierID string) error {
	return s.repo.Verify(ctx, id, verifierID)
}

func (s *GoodsReceiptService) RejectReceipt(ctx context.Context, id string, verifierID string, reason string) error {
	if reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}
	return s.repo.Reject(ctx, id, verifierID, reason)
}

func (s *GoodsReceiptService) GeneratePDF(ctx context.Context, id string) (*bytes.Buffer, error) {
	gr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "BERITA ACARA SERAH TERIMA (BAST)", "", 1, "C", false, 0, "")
	pdf.Ln(2)

	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(40, 6, "No. BAST", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, ": "+gr.BastNumber, "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 6, "Tanggal Penerimaan", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, ": "+gr.ReceiptDate.Format("02 Januari 2006"), "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 6, "Sumber", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, ": "+gr.SourceType, "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 6, "Supplier/Donor", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, ": "+gr.SupplierName, "", 1, "L", false, 0, "")
	pdf.CellFormat(40, 6, "No. Faktur", "", 0, "L", false, 0, "")
	pdf.CellFormat(0, 6, ": "+gr.InvoiceNumber, "", 1, "L", false, 0, "")
	pdf.Ln(5)

	// Table Header
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(10, 7, "No", "1", 0, "C", false, 0, "")
	pdf.CellFormat(70, 7, "Nama Barang", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 7, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Harga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 7, "Total", "1", 1, "C", false, 0, "")

	// Table Body
	pdf.SetFont("Arial", "", 9)
	for i, item := range gr.Items {
		pdf.CellFormat(10, 6, fmt.Sprint(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(70, 6, item.ItemName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 6, fmt.Sprintf("%d %s", item.Quantity, item.Unit), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("Rp %.2f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 6, fmt.Sprintf("Rp %.2f", item.TotalPrice), "1", 1, "R", false, 0, "")
	}

	// Tanda Tangan
	pdf.Ln(20)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(90, 6, "Yang Menyerahkan,", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 6, "Yang Menerima,", "", 1, "C", false, 0, "")
	pdf.Ln(20)
	pdf.CellFormat(90, 6, "(....................)", "", 0, "C", false, 0, "")
	pdf.CellFormat(0, 6, gr.ReceivedBy, "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}
