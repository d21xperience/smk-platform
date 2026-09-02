package queries

import (
	"context"
	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
	"sekolah-platform/services/tu-core/internal/domain/inventory/repository"
)

type GetAssetHandler struct {
	repo repository.InventoryItemRepository
}

func NewGetAssetHandler(repo repository.InventoryItemRepository) *GetAssetHandler {
	return &GetAssetHandler{repo: repo}
}

func (h *GetAssetHandler) Handle(ctx context.Context, id string) (*aggregate.InventoryItem, error) {
	return h.repo.FindByID(ctx, id)
}
