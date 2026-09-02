package grpc

// import (
// 	"context"
// 	"time"

// 	"sekolah-platform/services/tu-core/internal/domain/inventory/commands"
// 	"sekolah-platform/services/tu-core/internal/service"

// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// )

// type InventoryHandler struct {
// 	pb.UnimplementedInventoryServiceServer
// 	service *service.InventoryService
// }

// func NewInventoryHandler(svc *service.InventoryService) *InventoryHandler {
// 	return &InventoryHandler{service: svc}
// }

// func (h *InventoryHandler) CreateAsset(ctx context.Context, req *pb.CreateAssetRequest) (*pb.CreateAssetResponse, error) {
// 	cmd := commands.CreateAssetCommand{
// 		ID:             req.Id,
// 		QRCode:         req.QrCode,
// 		Name:           req.Name,
// 		SerialNumber:   req.SerialNumber,
// 		PurchaseDate:   time.Unix(req.PurchaseDate, 0),
// 		BuildingID:     req.BuildingId,
// 		FloorID:        req.FloorId,
// 		RoomID:         req.RoomId,
// 		LocationDetail: req.LocationDetail,
// 		OwnerID:        req.OwnerId,
// 		OwnerType:      req.OwnerType,
// 		Brand:          req.Brand,
// 		Model:          req.Model,
// 		Year:           int(req.Year),
// 		Color:          req.Color,
// 		Size:           req.Size,
// 		Vendor:         req.Vendor,
// 		Price:          req.Price,
// 		InvoiceNo:      req.InvoiceNo,
// 		ReceivedAt:     time.Unix(req.ReceivedAt, 0),
// 	}
// 	item, err := h.service.CreateAsset(ctx, cmd)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, err.Error())
// 	}
// 	return &pb.CreateAssetResponse{Id: item.GetID()}, nil
// }

// // ... implementasi method lain (MoveAsset, BorrowAsset, ReportDamage, dll.)
