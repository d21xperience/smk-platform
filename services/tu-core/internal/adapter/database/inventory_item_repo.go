package database

// import (
// 	"context"
// 	"time"

// 	"sekolah-platform/services/tu-core/internal/adapter/database/models"
// 	"sekolah-platform/services/tu-core/internal/domain/inventory/aggregate"
// 	"sekolah-platform/services/tu-core/internal/domain/inventory/errors"
// 	valueobject "sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
// 	"sekolah-platform/services/tu-core/internal/domain/inventory/repository"

// 	"gorm.io/gorm"
// )

// type inventoryItemRepository struct {
// 	db *gorm.DB
// }

// func NewInventoryItemRepository(db *gorm.DB) repository.InventoryItemRepository {
// 	return &inventoryItemRepository{db: db}
// }

// func toDomain(model *models.InventoryItemModel) (*aggregate.InventoryItem, error) {
// 	identity, _ := valueobject.NewAssetIdentity(model.ID, model.QRCode, model.Name, model.SerialNumber, time.Unix(model.PurchaseDate, 0))
// 	location, _ := valueobject.NewAssetLocation(model.BuildingID, model.FloorID, model.RoomID, model.LocationDetail)
// 	ownership := valueobject.NewAssetOwnership(model.OwnerID, valueobject.OwnerType(model.OwnerType))
// 	spec := valueobject.AssetSpecification{
// 		Brand:      model.Brand,
// 		Model:      model.Model,
// 		Year:       model.Year,
// 		Color:      model.Color,
// 		Size:       model.Size,
// 		Additional: map[string]string{},
// 	}
// 	procurement := aggregate.ProcurementInformation{
// 		Vendor:     model.Vendor,
// 		Price:      model.Price,
// 		InvoiceNo:  model.InvoiceNo,
// 		ReceivedAt: time.Unix(model.ReceivedAt, 0),
// 	}
// 	// Hydrate aggregate (tanpa events)
// 	return aggregate.HydrateInventoryItem(identity, valueobject.AssetCondition(model.Condition), location, ownership, spec, procurement), nil
// }

// func toModel(item *aggregate.InventoryItem) *models.InventoryItemModel {
// 	return &models.InventoryItemModel{
// 		ID:             item.GetID(),
// 		QRCode:         item.GetIdentity().QRCode,
// 		Name:           item.GetIdentity().Name,
// 		SerialNumber:   item.GetIdentity().SerialNumber,
// 		PurchaseDate:   item.GetIdentity().PurchaseDate.Unix(),
// 		BuildingID:     item.GetLocation().BuildingID,
// 		FloorID:        item.GetLocation().FloorID,
// 		RoomID:         item.GetLocation().RoomID,
// 		LocationDetail: item.GetLocation().Detail,
// 		Condition:      string(item.GetCondition()),
// 		OwnerID:        item.GetOwnership().OwnerID,
// 		OwnerType:      string(item.GetOwnership().OwnerType),
// 		Brand:          item.GetSpecification().Brand,
// 		Model:          item.GetSpecification().Model,
// 		Year:           item.GetSpecification().Year,
// 		Color:          item.GetSpecification().Color,
// 		Size:           item.GetSpecification().Size,
// 		Vendor:         item.GetProcurement().Vendor,
// 		Price:          item.GetProcurement().Price,
// 		InvoiceNo:      item.GetProcurement().InvoiceNo,
// 		ReceivedAt:     item.GetProcurement().ReceivedAt.Unix(),
// 		CreatedAt:      time.Now().Unix(),
// 		UpdatedAt:      time.Now().Unix(),
// 	}
// }

// func (r *inventoryItemRepository) Save(ctx context.Context, item *aggregate.InventoryItem) error {
// 	model := toModel(item)
// 	var existing models.InventoryItemModel
// 	err := r.db.WithContext(ctx).Where("id = ?", model.ID).First(&existing).Error
// 	if errors.Is(err, gorm.ErrRecordNotFound) {
// 		model.CreatedAt = time.Now().Unix()
// 		return r.db.WithContext(ctx).Create(model).Error
// 	}
// 	model.CreatedAt = existing.CreatedAt
// 	model.UpdatedAt = time.Now().Unix()
// 	return r.db.WithContext(ctx).Save(model).Error
// }

// // Implementasi FindByID, FindByQRCode, FindByRoom, FindBySerialNumber...
