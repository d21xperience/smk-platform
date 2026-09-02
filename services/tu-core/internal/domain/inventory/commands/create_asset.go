package commands

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
	"sekolah-platform/services/tu-core/internal/domain/inventory/engine"
	valueobject "sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
	"sekolah-platform/services/tu-core/internal/domain/inventory/repository"
)

type CreateAssetCommand struct {
	ID             string
	QRCode         string
	Name           string
	SerialNumber   string
	PurchaseDate   time.Time
	BuildingID     string
	FloorID        string
	RoomID         string
	LocationDetail string
	OwnerID        string
	OwnerType      string
	Brand          string
	Model          string
	Year           int
	Color          string
	Size           string
	Vendor         string
	Price          float64
	InvoiceNo      string
	ReceivedAt     time.Time
}

type CreateAssetHandler struct {
	engine *engine.InventoryEngine
	repo   repository.InventoryItemRepository
}

func NewCreateAssetHandler(engine *engine.InventoryEngine, repo repository.InventoryItemRepository) *CreateAssetHandler {
	return &CreateAssetHandler{engine: engine, repo: repo}
}

func (h *CreateAssetHandler) Handle(ctx context.Context, cmd CreateAssetCommand) (*aggregate.InventoryItem, error) {
	// Gunakan engine method CreateAsset yang sudah dibuat
	engineCmd := engine.CreateAssetCommand{
		ID:           cmd.ID,
		QRCode:       cmd.QRCode,
		Name:         cmd.Name,
		SerialNumber: cmd.SerialNumber,
		PurchaseDate: cmd.PurchaseDate,
		BuildingID:   cmd.BuildingID,
		FloorID:      cmd.FloorID,
		RoomID:       cmd.RoomID,
		Detail:       cmd.LocationDetail,
		OwnerID:      cmd.OwnerID,
		OwnerType:    valueobject.OwnerType(cmd.OwnerType),
		Brand:        cmd.Brand,
		Model:        cmd.Model,
		Year:         cmd.Year,
		Color:        cmd.Color,
		Size:         cmd.Size,
		Vendor:       cmd.Vendor,
		Price:        cmd.Price,
		InvoiceNo:    cmd.InvoiceNo,
		ReceivedAt:   cmd.ReceivedAt,
	}
	return h.engine.CreateAsset(ctx, engineCmd)
}
