package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/registration/v1"
)

type RegistrationGRPCHandler struct {
	pb.UnimplementedRegistrationServiceServer
	service *service.RegistrationService
}

func NewRegistrationGRPCHandler(svc *service.RegistrationService) *RegistrationGRPCHandler {
	return &RegistrationGRPCHandler{service: svc}
}

func (h *RegistrationGRPCHandler) CreateRegistration(ctx context.Context, req *pb.CreateRegistrationRequest) (*pb.CreateRegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	serviceReq := service.CreateRegistrationRequest{
		CalonNISN:         req.GetCalonNisn(),
		CalonNama:         req.GetCalonNama(),
		CalonJenisKelamin: req.GetCalonJenisKelamin(),
		CalonTempatLahir:  req.GetCalonTempatLahir(),
		CalonTanggalLahir: req.GetCalonTanggalLahir().AsTime(),
		CalonAlamat:       req.GetCalonAlamat(),
		CalonTelepon:      req.GetCalonTelepon(),
		CalonEmail:        req.GetCalonEmail(),
		NamaAyah:          req.GetNamaAyah(),
		NamaIbu:           req.GetNamaIbu(),
		TeleponOrtu:       req.GetTeleponOrtu(),
		PekerjaanAyah:     req.GetPekerjaanAyah(),
		PekerjaanIbu:      req.GetPekerjaanIbu(),
		AsalSekolah:       req.GetAsalSekolah(),
		JurusanDipilih:    req.GetJurusanDipilih(),
		AlasanMemilih:     req.GetAlasanMemilih(),
	}

	reg, err := h.service.CreateRegistration(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateRegistrationResponse{
		Registration: mapRegistrationResponseToProto(reg),
	}, nil
}

func (h *RegistrationGRPCHandler) SubmitRegistration(ctx context.Context, req *pb.SubmitRegistrationRequest) (*pb.SubmitRegistrationResponse, error) {
	reg, err := h.service.SubmitRegistration(ctx, req.GetRegistrationId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.SubmitRegistrationResponse{Registration: mapRegistrationResponseToProto(reg)}, nil
}

func (h *RegistrationGRPCHandler) VerifyRegistration(ctx context.Context, req *pb.VerifyRegistrationRequest) (*pb.VerifyRegistrationResponse, error) {
	reg, err := h.service.VerifyRegistration(ctx, req.GetRegistrationId(), req.GetCatatan())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.VerifyRegistrationResponse{Registration: mapRegistrationResponseToProto(reg)}, nil
}

func (h *RegistrationGRPCHandler) ApproveRegistration(ctx context.Context, req *pb.ApproveRegistrationRequest) (*pb.ApproveRegistrationResponse, error) {
	reg, err := h.service.ApproveRegistration(ctx, req.GetRegistrationId(), req.GetCatatan())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.ApproveRegistrationResponse{Registration: mapRegistrationResponseToProto(reg)}, nil
}

func (h *RegistrationGRPCHandler) RejectRegistration(ctx context.Context, req *pb.RejectRegistrationRequest) (*pb.RejectRegistrationResponse, error) {
	reg, err := h.service.RejectRegistration(ctx, req.GetRegistrationId(), req.GetAlasan())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.RejectRegistrationResponse{Registration: mapRegistrationResponseToProto(reg)}, nil
}

func (h *RegistrationGRPCHandler) GetRegistration(ctx context.Context, req *pb.GetRegistrationRequest) (*pb.GetRegistrationResponse, error) {
	reg, err := h.service.GetRegistration(ctx, req.GetRegistrationId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.GetRegistrationResponse{Registration: mapRegistrationResponseToProto(reg)}, nil
}

func (h *RegistrationGRPCHandler) ListRegistrations(ctx context.Context, req *pb.ListRegistrationsRequest) (*pb.ListRegistrationsResponse, error) {
	registrations, err := h.service.ListRegistrations(ctx, req.GetStatus(), req.GetSearch())
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbRegistrations []*pb.Registration
	for _, r := range registrations {
		pbRegistrations = append(pbRegistrations, mapRegistrationResponseToProto(r))
	}

	return &pb.ListRegistrationsResponse{Registrations: pbRegistrations, Total: int32(len(pbRegistrations))}, nil
}

func mapRegistrationResponseToProto(r *service.RegistrationResponse) *pb.Registration {
	if r == nil {
		return nil
	}

	return &pb.Registration{
		Id:                r.ID,
		SchoolId:          r.SchoolID,
		AcademicPeriodId:  r.AcademicPeriodID,
		CalonNisn:         r.CalonNISN,
		CalonNama:         r.CalonNama,
		CalonJenisKelamin: r.CalonJenisKelamin,
		CalonTempatLahir:  r.CalonTempatLahir,
		CalonTanggalLahir: timestamppb.New(r.CalonTanggalLahir),
		CalonAlamat:       r.CalonAlamat,
		CalonTelepon:      r.CalonTelepon,
		CalonEmail:        r.CalonEmail,
		NamaAyah:          r.NamaAyah,
		NamaIbu:           r.NamaIbu,
		TeleponOrtu:       r.TeleponOrtu,
		PekerjaanAyah:     r.PekerjaanAyah,
		PekerjaanIbu:      r.PekerjaanIbu,
		AsalSekolah:       r.AsalSekolah,
		JurusanDipilih:    r.JurusanDipilih,
		AlasanMemilih:     r.AlasanMemilih,
		Status:            r.Status,
		CatatanVerifikasi: r.CatatanVerifikasi,
		VerifiedBy:        r.VerifiedBy,
		VerifiedAt:        timestamppb.New(r.VerifiedAt),
		CreatedAt:         timestamppb.New(r.CreatedAt),
		UpdatedAt:         timestamppb.New(r.UpdatedAt),
	}
}
