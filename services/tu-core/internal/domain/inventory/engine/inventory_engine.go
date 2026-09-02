package engine

import (
	"context"
	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
	domainErrors "sekolah-platform/services/tu-core/internal/domain/inventory/errors"
	"sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
	"sekolah-platform/services/tu-core/internal/domain/inventory/repository"
	"time"
)

type InventoryEngine struct {
	itemRepo  repository.InventoryItemRepository
	issueRepo repository.AssetIssueRepository
}

func NewInventoryEngine(itemRepo repository.InventoryItemRepository, issueRepo repository.AssetIssueRepository) *InventoryEngine {
	return &InventoryEngine{
		itemRepo:  itemRepo,
		issueRepo: issueRepo,
	}
}

// ----- Command: Create Asset (digunakan oleh TU) -----
type CreateAssetCommand struct {
	ID           string
	QRCode       string
	Name         string
	SerialNumber string
	PurchaseDate time.Time
	BuildingID   string
	FloorID      string
	RoomID       string
	Detail       string
	OwnerID      string
	OwnerType    value_object.OwnerType
	Brand        string
	Model        string
	Year         int
	Color        string
	Size         string
	Vendor       string
	Price        float64
	InvoiceNo    string
	ReceivedAt   time.Time
}

func (e *InventoryEngine) CreateAsset(ctx context.Context, cmd CreateAssetCommand) (*aggregate.InventoryItem, error) {
	// 1. Build value objects
	identity, err := value_object.NewAssetIdentity(cmd.ID, cmd.QRCode, cmd.Name, cmd.SerialNumber, cmd.PurchaseDate)
	if err != nil {
		return nil, err
	}
	location, err := value_object.NewAssetLocation(cmd.BuildingID, cmd.FloorID, cmd.RoomID, cmd.Detail)
	if err != nil {
		return nil, err
	}
	ownership := value_object.NewAssetOwnership(cmd.OwnerID, cmd.OwnerType)
	spec := value_object.AssetSpecification{
		Brand:      cmd.Brand,
		Model:      cmd.Model,
		Year:       cmd.Year,
		Color:      cmd.Color,
		Size:       cmd.Size,
		Additional: map[string]string{},
	}
	procurement := aggregate.ProcurementInformation{
		Vendor:     cmd.Vendor,
		Price:      cmd.Price,
		InvoiceNo:  cmd.InvoiceNo,
		ReceivedAt: cmd.ReceivedAt,
	}

	// 2. Create aggregate
	item, err := aggregate.NewInventoryItem(identity, location, ownership, spec, procurement)
	if err != nil {
		return nil, err
	}

	// 3. Save
	if err := e.itemRepo.Save(ctx, item); err != nil {
		return nil, err
	}

	// 4. Events akan di-publish oleh application layer (dari item.GetEvents())
	return item, nil
}

// ----- Command: Report Issue (oleh Guru) -----
type ReportIssueCommand struct {
	ID              string
	InventoryItemID *string // bisa nil jika discovery
	ReporterID      string
	RoomID          string
	IssueType       value_object.IssueType
	Description     string
	Photos          []string
}

func (e *InventoryEngine) ReportIssue(ctx context.Context, cmd ReportIssueCommand) (*aggregate.AssetIssue, error) {
	issue, err := aggregate.NewAssetIssue(
		cmd.ID,
		cmd.InventoryItemID,
		cmd.ReporterID,
		cmd.RoomID,
		cmd.IssueType,
		cmd.Description,
		cmd.Photos,
	)
	if err != nil {
		return nil, err
	}

	// Jika issue terkait asset yang sudah ada, catat maintenance request
	if cmd.InventoryItemID != nil {
		item, err := e.itemRepo.FindByID(ctx, *cmd.InventoryItemID)
		if err != nil {
			return nil, domainErrors.ErrAssetNotFound
		}
		// request maintenance (hanya generate event, tidak ubah status)
		if err := item.RequestMaintenance(issue.GetID()); err != nil {
			return nil, err
		}
		// simpan perubahan pada item (events tambahan)
		if err := e.itemRepo.Save(ctx, item); err != nil {
			return nil, err
		}
	}

	// Simpan issue
	if err := e.issueRepo.Save(ctx, issue); err != nil {
		return nil, err
	}
	return issue, nil
}

// ----- Command: Verify Issue (oleh TU/Sarpras) -----
func (e *InventoryEngine) VerifyIssue(ctx context.Context, issueID, verifierID string) error {
	issue, err := e.issueRepo.FindByID(ctx, issueID)
	if err != nil {
		return domainErrors.ErrAssetNotFound
	}
	if err := issue.Verify(verifierID); err != nil {
		return err
	}
	return e.issueRepo.Save(ctx, issue)
}

// ----- Command: Close Issue (setelah maintenance selesai) -----
func (e *InventoryEngine) CloseIssue(ctx context.Context, issueID, resolution string, newCondition *value_object.AssetCondition) error {
	issue, err := e.issueRepo.FindByID(ctx, issueID)
	if err != nil {
		return domainErrors.ErrAssetNotFound
	}
	if err := issue.Close(resolution); err != nil {
		return err
	}

	// Jika issue terkait inventory, update kondisi aset setelah perbaikan
	if issue.GetInventoryItemID() != nil && newCondition != nil {
		item, err := e.itemRepo.FindByID(ctx, *issue.GetInventoryItemID())
		if err != nil {
			return domainErrors.ErrAssetNotFound
		}
		if err := item.CompleteMaintenance(*newCondition); err != nil {
			return err
		}
		if err := e.itemRepo.Save(ctx, item); err != nil {
			return err
		}
	}

	return e.issueRepo.Save(ctx, issue)
}
