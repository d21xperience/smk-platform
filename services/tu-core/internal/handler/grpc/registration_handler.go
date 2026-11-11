package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
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

	// Parsing Tanggal Lahir ke *time.Time
	var birthDate *time.Time
	if req.GetCalonTanggalLahir() != nil {
		t := req.GetCalonTanggalLahir().AsTime()
		birthDate = &t
	}

	// Field disesuaikan dengan CreateRegistrationRequest di service yang baru
	// (Menghapus field yang tidak ada di tabel registrations DAPODIK)
	serviceReq := service.CreateRegistrationRequest{
		CalonNISN:         req.GetCalonNisn(),
		CalonNama:         req.GetCalonNama(),
		CalonJenisKelamin: req.GetCalonJenisKelamin(),
		CalonTempatLahir:  req.GetCalonTempatLahir(),
		CalonTanggalLahir: birthDate,
		AsalSekolah:       req.GetAsalSekolah(),
		JurusanDipilih:    req.GetJurusanDipilih(),
		// Field keuangan dan berkas di-set default karena tidak ada di proto lama
		BiayaStandar:     0,
		BiayaPendaftaran: 0,
		BerkasIjazah:     false,
		BerkasKK:         false,
		BerkasAkta:       false,
	}

	reg, err := h.service.CreateRegistration(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateRegistrationResponse{
		Registration: mapRegistrationResponseToProto(reg),
	}, nil
}

// --- STUB FUNCTIONS (Untuk mempertahankan kompilasi gRPC) ---

func (h *RegistrationGRPCHandler) SubmitRegistration(ctx context.Context, req *pb.SubmitRegistrationRequest) (*pb.SubmitRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitRegistration belum diimplementasikan untuk skema DAPODIK")
}

func (h *RegistrationGRPCHandler) VerifyRegistration(ctx context.Context, req *pb.VerifyRegistrationRequest) (*pb.VerifyRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyRegistration belum diimplementasikan untuk skema DAPODIK")
}

func (h *RegistrationGRPCHandler) ApproveRegistration(ctx context.Context, req *pb.ApproveRegistrationRequest) (*pb.ApproveRegistrationResponse, error) {
	// Kita bisa memetakan "Approve" ke "ActivateRegistration" di service baru
	err := h.service.ActivateRegistration(ctx, req.GetRegistrationId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	
	// Ambil data terbaru untuk response (opsional, karena service hanya return error)
	// Kita bisa mengembalikan response kosong atau mengambil ulang data
	return &pb.ApproveRegistrationResponse{}, nil
}

func (h *RegistrationGRPCHandler) RejectRegistration(ctx context.Context, req *pb.RejectRegistrationRequest) (*pb.RejectRegistrationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RejectRegistration belum diimplementasikan untuk skema DAPODIK")
}

func (h *RegistrationGRPCHandler) GetRegistration(ctx context.Context, req *pb.GetRegistrationRequest) (*pb.GetRegistrationResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}
	
	reg, err := h.service.GetRegistrationByID(ctx, req.GetRegistrationId())
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
		// Konversi dari []repository.Registration ke []*pb.Registration
		pbRegistrations = append(pbRegistrations, mapRegistrationResponseToProto(&r))
	}

	return &pb.ListRegistrationsResponse{Registrations: pbRegistrations, Total: int32(len(pbRegistrations))}, nil
}

// Mapping dari repository.Registration ke proto
func mapRegistrationResponseToProto(r *repository.Registration) *pb.Registration {
	if r == nil {
		return nil
	}

	// Handle *time.Time untuk timestamppb
	var pbBirthDate *timestamppb.Timestamp
	if r.CalonTanggalLahir != nil {
		pbBirthDate = timestamppb.New(*r.CalonTanggalLahir)
	}

	return &pb.Registration{
		Id:                r.ID,
		SchoolId:          r.SchoolID,
		AcademicPeriodId:  r.AcademicPeriodID,
		CalonNisn:         r.CalonNISN,
		CalonNama:         r.CalonNama,
		CalonJenisKelamin: r.CalonJenisKelamin,
		CalonTempatLahir:  r.CalonTempatLahir,
		CalonTanggalLahir: pbBirthDate,
		// Field yang sudah dihapus tidak diisi
		CalonAlamat:       "",
		CalonTelepon:      "",
		CalonEmail:        "",
		NamaAyah:          "",
		NamaIbu:           "",
		TeleponOrtu:       "",
		PekerjaanAyah:     "",
		PekerjaanIbu:      "",
		AsalSekolah:       r.AsalSekolah,
		JurusanDipilih:    r.JurusanDipilih,
		AlasanMemilih:     "",
		Status:            r.Status,
		CatatanVerifikasi: "",
		VerifiedBy:        "",
		VerifiedAt:        nil,
		CreatedAt:         timestamppb.New(r.CreatedAt),
		UpdatedAt:         timestamppb.New(r.UpdatedAt),
	}
}
