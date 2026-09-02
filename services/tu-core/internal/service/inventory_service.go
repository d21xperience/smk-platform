package service

// import (
// 	"context"

// 	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
// 	"sekolah-platform/services/tu-core/internal/domain/inventory/commands"
// )

// type InventoryService struct {
// 	createAssetHandler  *commands.CreateAssetHandler
// 	moveAssetHandler    *commands.MoveAssetHandler
// 	borrowAssetHandler  *commands.BorrowAssetHandler
// 	reportDamageHandler *commands.ReportDamageHandler
// 	// ... lainnya
// 	getAssetHandler *query.GetAssetHandler // query di folder internal/domain/inventory/queries
// }

// func NewInventoryService(
// 	create *commands.CreateAssetHandler,
// 	move *commands.MoveAssetHandler,
// 	borrow *commands.BorrowAssetHandler,
// 	report *commands.ReportDamageHandler,
// 	// ...
// 	get *query.GetAssetHandler,
// ) *InventoryService {
// 	return &InventoryService{
// 		createAssetHandler:  create,
// 		moveAssetHandler:    move,
// 		borrowAssetHandler:  borrow,
// 		reportDamageHandler: report,
// 		// ...
// 		getAssetHandler: get,
// 	}
// }

// func (s *InventoryService) CreateAsset(ctx context.Context, cmd commands.CreateAssetCommand) (*aggregate.InventoryItem, error) {
// 	return s.createAssetHandler.Handle(ctx, cmd)
// }

// func (s *InventoryService) MoveAsset(ctx context.Context, id, newBuilding, newFloor, newRoom, detail string) (*aggregate.InventoryItem, error) {
// 	return s.moveAssetHandler.Handle(ctx, id, newBuilding, newFloor, newRoom, detail)
// }

// // ... method lainnya
