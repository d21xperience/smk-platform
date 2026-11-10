package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
)

type BorrowingService struct {
	repo      repository.BorrowingRepository
	itemRepo  repository.InventoryItemRepository
}

func NewBorrowingService(repo repository.BorrowingRepository, itemRepo repository.InventoryItemRepository) *BorrowingService {
	return &BorrowingService{repo: repo, itemRepo: itemRepo}
}

func (s *BorrowingService) CreateBorrowing(ctx context.Context, req *repository.Borrowing, userID string, userName string) (*repository.Borrowing, error) {
	// 1. Validate dates
	if req.ExpectedReturnDate.Before(req.BorrowDate) {
		return nil, errors.New("tanggal kembali harus setelah tanggal pinjam")
	}
	
	// 2. Get asset details (to check status & value)
	asset, err := s.itemRepo.FindByID(ctx, req.InventoryItemID)
	if err != nil {
		return nil, errors.New("aset tidak ditemukan")
	}
	
	if asset.StatusAset == "DIPINJAM" {
		return nil, errors.New("aset sedang dipinjam")
	}
	if asset.StatusAset == "DIHAPUS" || asset.StatusAset == "DIUSULKAN_HAPUS" {
		return nil, errors.New("aset tidak dapat dipinjam karena status: " + asset.StatusAset)
	}
	if asset.Condition == "Rusak Berat" {
		return nil, errors.New("aset dalam kondisi rusak berat dan tidak dapat dipinjam")
	}
	
	// 3. Check conflict
	conflict, _ := s.repo.CountConflictingBorrowings(ctx, req.InventoryItemID, req.BorrowDate, req.ExpectedReturnDate)
	if conflict > 0 {
		return nil, errors.New("jadwal peminjaman bentrok dengan peminjaman lain pada periode yang sama")
	}
	
	// 4. Determine approval level
	durationDays := int(req.ExpectedReturnDate.Sub(req.BorrowDate).Hours()/24) + 1
	level := 1
	if durationDays > 1 { level = 2 }
	if durationDays > 3 { level = 3 }
	if durationDays > 7 { level = 4 }
	// +1 untuk luar sekolah
	if req.PurposeType == "DINAS" {
		level++
	}
	// +1 untuk aset mahal (di atas 10 juta)
	if asset.NilaiBuku > 10000000 {
		level++
	}
	if level > 5 { level = 5 }
	
	// 5. Generate borrowing number
	reqNo, err := s.repo.GenerateBorrowingNumber(ctx)
	if err != nil { return nil, err }
	
	// 6. Populate
	req.BorrowingNumber = reqNo
	req.BorrowerID = userID
	req.BorrowerName = userName
	req.DurationDays = durationDays
	req.ApprovalLevel = level
	req.Status = "PENDING_APPROVAL"
	req.CurrentApprovalStep = 0
	req.SemesterID = "20261" // TODO: Get from operational context
	
	// 7. Save
	if err := s.repo.Create(ctx, req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *BorrowingService) GetAll(ctx context.Context, semesterID string) ([]*repository.Borrowing, error) {
	return s.repo.GetAll(ctx, semesterID)
}

func (s *BorrowingService) GetByID(ctx context.Context, id string) (*repository.Borrowing, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BorrowingService) Approve(ctx context.Context, b *repository.Borrowing, approval *repository.BorrowingApproval) error {
	expectedStep := b.CurrentApprovalStep + 1
	if approval.ApprovalStep != expectedStep {
		return fmt.Errorf("step approval tidak sesuai. Expected: %d, Got: %d", expectedStep, approval.ApprovalStep)
	}
	return s.repo.Approve(ctx, b, approval)
}

func (s *BorrowingService) Pickup(ctx context.Context, b *repository.Borrowing, userID string, condition string, photos string, notes string) error {
	if b.Status != "APPROVED" {
		return errors.New("hanya peminjaman yang sudah APPROVED yang dapat diambil")
	}
	
	bastNo, err := s.repo.GenerateBASTNumber(ctx)
	if err != nil { return err }
	
	now := time.Now()
	b.Status = "PICKED_UP"
	b.ActualPickupDate = &now
	b.PickupBy = userID
	b.PickupCondition = condition
	b.PickupPhotos = photos
	b.PickupNotes = notes
	b.BastNumber = bastNo
	
	return s.repo.Pickup(ctx, b)
}

func (s *BorrowingService) Return(ctx context.Context, b *repository.Borrowing, userID string, condition string, photos string, notes string) error {
	if b.Status != "PICKED_UP" {
		return errors.New("hanya peminjaman dengan status PICKED_UP yang dapat dikembalikan")
	}
	
	bapNo, err := s.repo.GenerateBAPNumber(ctx)
	if err != nil { return err }
	
	now := time.Now()
	actualReturnDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	expectedReturnDate := time.Date(b.ExpectedReturnDate.Year(), b.ExpectedReturnDate.Month(), b.ExpectedReturnDate.Day(), 0, 0, 0, 0, time.UTC)
	
	daysLate := 0
	fineAmount := 0.0
	finePerDay := 10000.0 // TODO: Get from config
	
	if actualReturnDate.After(expectedReturnDate) {
		daysLate = int(actualReturnDate.Sub(expectedReturnDate).Hours() / 24)
		fineAmount = float64(daysLate) * finePerDay
	}
	
	b.Status = "RETURNED"
	b.ActualReturnDate = &actualReturnDate
	b.ReturnBy = userID
	b.ReturnCondition = condition
	b.ReturnPhotos = photos
	b.ReturnNotes = notes
	b.DaysLate = daysLate
	b.FineAmount = fineAmount
	b.FinePaid = fineAmount == 0
	b.BapNumber = bapNo
	
	var fine *repository.BorrowingFine
	if fineAmount > 0 {
		fine = &repository.BorrowingFine{
			BorrowingID: b.ID,
			BorrowerID:  b.BorrowerID,
			FineAmount:  fineAmount,
			DaysLate:    daysLate,
			FinePerDay:  finePerDay,
			FineDate:    actualReturnDate,
		}
	}
	
	return s.repo.Return(ctx, b, fine)
}
