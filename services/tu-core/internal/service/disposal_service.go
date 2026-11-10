package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type DisposalService struct {
	repo repository.DisposalRepository
}

func NewDisposalService(repo repository.DisposalRepository) *DisposalService {
	return &DisposalService{repo: repo}
}

func (s *DisposalService) CreateRequest(ctx context.Context, req *repository.DisposalRequest, items []repository.DisposalItem, userID string, semesterID string) (*repository.DisposalRequest, error) {
	if len(items) == 0 {
		return nil, errors.New("minimal 1 aset harus diusulkan")
	}

	var totalValue float64
	for _, item := range items {
		if item.SpecificReason == "" {
			return nil, errors.New("alasan spesifik wajib diisi untuk setiap aset")
		}
		// In real scenario, fetch asset value here. For now, assume value is passed in items.
		totalValue += item.AcquisitionPrice
	}

	reqNo, err := s.repo.GenerateRequestNumber(ctx)
	if err != nil { return nil, err }

	req.RequestNumber = reqNo
	req.ProposedBy = userID
	req.TotalAssets = len(items)
	req.TotalValue = totalValue

	if totalValue < 10000000 {
		req.ApprovalLevel = 1
	} else if totalValue < 50000000 {
		req.ApprovalLevel = 2
	} else {
		req.ApprovalLevel = 3
	}

	err = s.repo.CreateDisposal(ctx, req, items, semesterID)
	if err != nil { return nil, err }
	return req, nil
}

func (s *DisposalService) GetAll(ctx context.Context, semesterID string) ([]*repository.DisposalRequest, error) {
	return s.repo.GetAllDisposals(ctx, semesterID)
}

func (s *DisposalService) GetByID(ctx context.Context, id string) (*repository.DisposalRequest, error) {
	return s.repo.GetDisposalByID(ctx, id)
}

func (s *DisposalService) Verify(ctx context.Context, req *repository.DisposalRequest, verif *repository.DisposalVerification) error {
	if verif.TeamRecommendation == "TOLAK" {
		req.Status = "REJECTED"
		req.RejectionReason = verif.TeamNotes
	}
	return s.repo.VerifyDisposal(ctx, req, verif)
}

func (s *DisposalService) Approve(ctx context.Context, req *repository.DisposalRequest, approval *repository.DisposalApproval) error {
	expectedStep := req.CurrentApprovalStep + 1
	if approval.ApprovalStep != expectedStep {
		return fmt.Errorf("step approval tidak sesuai. Expected: %d", expectedStep)
	}
	return s.repo.ApproveDisposal(ctx, req, approval)
}

func (s *DisposalService) Execute(ctx context.Context, req *repository.DisposalRequest) error {
	// Generate BAP number
	year := time.Now().Year()
	req.BastNumber = fmt.Sprintf("BAP/SMKPASJA/%d/%03d", year, time.Now().Unix()%1000)
	
	// Generate Hash for QR
	req.VerificationHash = fmt.Sprintf("%x", time.Now().UnixNano())

	// Generate PDF (Basic mock, full implementation can be added)
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "BERITA ACARA PENGHAPUSAN", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 6, "Nomor: "+req.BastNumber, "", 1, "C", false, 0, "")
	pdf.Ln(10)
	
	// Basic QR Code
	qr, _ := qrcode.New("https://school.app/verify/"+req.ID, qrcode.Medium)
	var qrBuf bytes.Buffer
	qr.Write(128, &qrBuf)
	opts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
	pdf.RegisterImageOptionsReader("qr", opts, &qrBuf)
	pdf.ImageOptions("qr", 150, 20, 30, 30, false, opts, 0, "")

	var buf bytes.Buffer
	pdf.Output(&buf)
	req.BastPdfUrl = "https://minio.local/bap_" + req.BastNumber + ".pdf" // Mock URL

	return s.repo.ExecuteDisposal(ctx, req)
}
