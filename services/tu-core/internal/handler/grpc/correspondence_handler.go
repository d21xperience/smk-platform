package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/correspondence/v1"
)

type CorrespondenceGRPCHandler struct {
	pb.UnimplementedCorrespondenceServiceServer
	service *service.CorrespondenceService
}

func NewCorrespondenceGRPCHandler(svc *service.CorrespondenceService) *CorrespondenceGRPCHandler {
	return &CorrespondenceGRPCHandler{service: svc}
}

func (h *CorrespondenceGRPCHandler) CreateCorrespondence(ctx context.Context, req *pb.CreateCorrespondenceRequest) (*pb.CreateCorrespondenceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	serviceReq := service.CreateCorrespondenceRequest{
		Type:          req.GetType(),
		Number:        req.GetNumber(),
		Date:          req.GetDate().AsTime(),
		Subject:       req.GetSubject(),
		From:          req.GetFrom(),
		To:            req.GetTo(),
		Description:   req.GetDescription(),
		AttachmentURL: req.GetAttachmentUrl(),
	}

	c, err := h.service.CreateCorrespondence(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateCorrespondenceResponse{
		Correspondence: mapCorrespondenceResponseToProto(c),
	}, nil
}

func (h *CorrespondenceGRPCHandler) UpdateCorrespondence(ctx context.Context, req *pb.UpdateCorrespondenceRequest) (*pb.UpdateCorrespondenceResponse, error) {
	c, err := h.service.UpdateCorrespondence(ctx, req.GetId(), req.GetSubject(), req.GetDescription(), req.GetAttachmentUrl())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.UpdateCorrespondenceResponse{Correspondence: mapCorrespondenceResponseToProto(c)}, nil
}

func (h *CorrespondenceGRPCHandler) ArchiveCorrespondence(ctx context.Context, req *pb.ArchiveCorrespondenceRequest) (*pb.ArchiveCorrespondenceResponse, error) {
	c, err := h.service.ArchiveCorrespondence(ctx, req.GetId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.ArchiveCorrespondenceResponse{Correspondence: mapCorrespondenceResponseToProto(c)}, nil
}

func (h *CorrespondenceGRPCHandler) GetCorrespondence(ctx context.Context, req *pb.GetCorrespondenceRequest) (*pb.GetCorrespondenceResponse, error) {
	c, err := h.service.GetCorrespondence(ctx, req.GetId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.GetCorrespondenceResponse{Correspondence: mapCorrespondenceResponseToProto(c)}, nil
}

func (h *CorrespondenceGRPCHandler) ListCorrespondences(ctx context.Context, req *pb.ListCorrespondencesRequest) (*pb.ListCorrespondencesResponse, error) {
	correspondences, err := h.service.ListCorrespondences(ctx, req.GetType(), req.GetStatus(), req.GetSearch())
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbCorrespondences []*pb.Correspondence
	for _, c := range correspondences {
		pbCorrespondences = append(pbCorrespondences, mapCorrespondenceResponseToProto(c))
	}

	return &pb.ListCorrespondencesResponse{
		Correspondences: pbCorrespondences,
		Total:           int32(len(pbCorrespondences)),
	}, nil
}

func mapCorrespondenceResponseToProto(c *service.CorrespondenceResponse) *pb.Correspondence {
	if c == nil {
		return nil
	}
	return &pb.Correspondence{
		Id:            c.ID,
		Type:          c.Type,
		Number:        c.Number,
		Date:          timestamppb.New(c.Date),
		Subject:       c.Subject,
		From:          c.From,
		To:            c.To,
		Description:   c.Description,
		AttachmentUrl: c.AttachmentURL,
		Status:        c.Status,
		CreatedAt:     timestamppb.New(c.CreatedAt),
	}
}
