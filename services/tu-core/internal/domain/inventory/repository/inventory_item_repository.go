package repository

import (
	"context"

	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
	"sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
)

type InventoryItemRepository interface {
	Save(ctx context.Context, item *aggregate.InventoryItem) error
	FindByID(ctx context.Context, id string) (*aggregate.InventoryItem, error)
	FindByQRCode(ctx context.Context, qrCode string) (*aggregate.InventoryItem, error)
	FindByRoom(ctx context.Context, location value_object.AssetLocation) ([]*aggregate.InventoryItem, error)
	FindBySerialNumber(ctx context.Context, serial string) (*aggregate.InventoryItem, error)
}
