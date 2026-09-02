package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sekolah-platform/services/tu-core/internal/domain/mutation/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/mutation/v1"
)

type MutationGRPCHandler struct {
	pb.UnimplementedMutationServiceServer
	service *service.MutationService
}

func NewMutationGRPCHandler(svc *service.MutationService) *MutationGRPCHandler {
	return &MutationGRPCHandler{service: svc}
}

func (h *MutationGRPCHandler) CreateMutation(ctx context.Context, req *pb.CreateMutationRequest) (*pb.CreateMutationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	mType, _ := valueobjects.NewMutationType(req.GetMutationType())

	serviceReq := service.CreateMutationRequest{
		MutationType:  mType,
		CalonNama:     req.GetCalonNama(),
		CalonNISN:     req.GetCalonNisn(),
		AsalSekolah:   req.GetAsalSekolah(),
		TujuanSekolah: req.GetTujuanSekolah(),
		Alasan:        req.GetAlasan(),
		DokumenURL:    req.GetDokumenUrl(),
	}

	if req.GetStudentId() != "" {
		serviceReq.StudentID = &req.StudentId
	}

	m, err := h.service.CreateMutation(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateMutationResponse{
		Mutation: mapMutationResponseToProto(m),
	}, nil
}

func (h *MutationGRPCHandler) ApproveMutation(ctx context.Context, req *pb.ApproveMutationRequest) (*pb.ApproveMutationResponse, error) {
	m, err := h.service.ApproveMutation(ctx, req.GetMutationId(), req.GetCatatan())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.ApproveMutationResponse{Mutation: mapMutationResponseToProto(m)}, nil
}

func (h *MutationGRPCHandler) RejectMutation(ctx context.Context, req *pb.RejectMutationRequest) (*pb.RejectMutationResponse, error) {
	m, err := h.service.RejectMutation(ctx, req.GetMutationId(), req.GetCatatan())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.RejectMutationResponse{Mutation: mapMutationResponseToProto(m)}, nil
}

func (h *MutationGRPCHandler) ListMutations(ctx context.Context, req *pb.ListMutationsRequest) (*pb.ListMutationsResponse, error) {
	mutations, err := h.service.ListMutations(ctx, req.GetStatus(), req.GetType())
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbMutations []*pb.Mutation
	for _, m := range mutations {
		pbMutations = append(pbMutations, mapMutationResponseToProto(m))
	}

	return &pb.ListMutationsResponse{Mutations: pbMutations}, nil
}

func mapMutationResponseToProto(m *service.MutationResponse) *pb.Mutation {
	if m == nil {
		return nil
	}

	var studentID string
	if m.StudentID != nil {
		studentID = *m.StudentID
	}

	return &pb.Mutation{
		Id:                m.ID,
		SchoolId:          m.SchoolID,
		AcademicPeriodId:  m.AcademicPeriodID,
		StudentId:         studentID,
		MutationType:      m.MutationType,
		CalonNama:         m.CalonNama,
		CalonNisn:         m.CalonNISN,
		AsalSekolah:       m.AsalSekolah,
		TujuanSekolah:     m.TujuanSekolah,
		Alasan:            m.Alasan,
		DokumenUrl:        m.DokumenURL,
		Status:            m.Status,
		CatatanVerifikasi: m.CatatanVerifikasi,
		CreatedAt:         m.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
