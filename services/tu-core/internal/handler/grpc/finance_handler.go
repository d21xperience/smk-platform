package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"sekolah-platform/services/tu-core/internal/domain/finance/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/service"
	pb "sekolah-platform/services/tu-core/proto/finance/v1"
)

type FinanceGRPCHandler struct {
	pb.UnimplementedFinanceServiceServer
	service *service.FinanceService
}

func NewFinanceGRPCHandler(svc *service.FinanceService) *FinanceGRPCHandler {
	return &FinanceGRPCHandler{service: svc}
}

func (h *FinanceGRPCHandler) CreateInvoice(ctx context.Context, req *pb.CreateInvoiceRequest) (*pb.CreateInvoiceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, status.Errorf(codes.InvalidArgument, "missing operational context")
	}

	componentType, _ := valueobjects.NewComponentType(req.GetComponentType())

	serviceReq := service.CreateInvoiceServiceRequest{
		StudentID:     req.GetStudentId(),
		ComponentType: componentType,
		Month:         int(req.GetMonth()),
		Year:          int(req.GetYear()),
		Amount:        req.GetAmount(),
		DueDate:       req.GetDueDate().AsTime(),
		Notes:         req.GetNotes(),
	}

	invoice, err := h.service.CreateInvoice(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreateInvoiceResponse{
		Invoice: mapInvoiceResponseToProto(invoice),
	}, nil
}

func (h *FinanceGRPCHandler) GetInvoice(ctx context.Context, req *pb.GetInvoiceRequest) (*pb.GetInvoiceResponse, error) {
	invoice, err := h.service.GetInvoice(ctx, req.GetInvoiceId())
	if err != nil {
		return nil, mapServiceError(err)
	}
	return &pb.GetInvoiceResponse{Invoice: mapInvoiceResponseToProto(invoice)}, nil
}

func (h *FinanceGRPCHandler) ListInvoices(ctx context.Context, req *pb.ListInvoicesRequest) (*pb.ListInvoicesResponse, error) {
	invoices, err := h.service.ListInvoices(ctx, req.GetStudentId(), req.GetComponentType(), req.GetStatus(), int(req.GetMonth()), int(req.GetYear()))
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbInvoices []*pb.Invoice
	for _, i := range invoices {
		pbInvoices = append(pbInvoices, mapInvoiceResponseToProto(i))
	}

	return &pb.ListInvoicesResponse{
		Invoices: pbInvoices,
		Total:    int32(len(pbInvoices)),
	}, nil
}

func (h *FinanceGRPCHandler) GetStudentOutstanding(ctx context.Context, req *pb.GetStudentOutstandingRequest) (*pb.GetStudentOutstandingResponse, error) {
	result, err := h.service.GetStudentOutstanding(ctx, req.GetStudentId())
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbInvoices []*pb.Invoice
	for _, i := range result.UnpaidInvoicesList {
		pbInvoices = append(pbInvoices, mapInvoiceResponseToProto(&i))
	}

	return &pb.GetStudentOutstandingResponse{
		StudentId:          result.StudentID,
		TotalOutstanding:   result.TotalOutstanding,
		UnpaidInvoices:     int32(result.UnpaidInvoices),
		UnpaidInvoicesList: pbInvoices,
	}, nil
}

func (h *FinanceGRPCHandler) CreatePayment(ctx context.Context, req *pb.CreatePaymentRequest) (*pb.CreatePaymentResponse, error) {
	paymentMethod, _ := valueobjects.NewPaymentMethod(req.GetPaymentMethod())

	serviceReq := service.CreatePaymentServiceRequest{
		InvoiceID:       req.GetInvoiceId(),
		Amount:          req.GetAmount(),
		PaymentMethod:   paymentMethod,
		PaymentDate:     req.GetPaymentDate().AsTime(),
		ReferenceNumber: req.GetReferenceNumber(),
		Notes:           req.GetNotes(),
	}

	result, err := h.service.CreatePayment(ctx, serviceReq)
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.CreatePaymentResponse{
		Payment:        mapPaymentResponseToProto(result.Payment),
		UpdatedInvoice: mapInvoiceResponseToProto(result.UpdatedInvoice),
	}, nil
}

func (h *FinanceGRPCHandler) ListPayments(ctx context.Context, req *pb.ListPaymentsRequest) (*pb.ListPaymentsResponse, error) {
	var startDate, endDate time.Time
	if req.GetStartDate() != nil {
		startDate = req.GetStartDate().AsTime()
	}
	if req.GetEndDate() != nil {
		endDate = req.GetEndDate().AsTime()
	}

	payments, err := h.service.ListPayments(ctx, req.GetInvoiceId(), req.GetPaymentMethod(), startDate, endDate)
	if err != nil {
		return nil, mapServiceError(err)
	}

	var pbPayments []*pb.Payment
	for _, p := range payments {
		pbPayments = append(pbPayments, mapPaymentResponseToProto(p))
	}

	return &pb.ListPaymentsResponse{
		Payments: pbPayments,
		Total:    int32(len(pbPayments)),
	}, nil
}

func (h *FinanceGRPCHandler) GetPaymentSummary(ctx context.Context, req *pb.GetPaymentSummaryRequest) (*pb.GetPaymentSummaryResponse, error) {
	summary, err := h.service.GetPaymentSummary(ctx, int(req.GetMonth()), int(req.GetYear()), req.GetComponentType())
	if err != nil {
		return nil, mapServiceError(err)
	}

	return &pb.GetPaymentSummaryResponse{
		TotalInvoices:     int32(summary.TotalInvoices),
		PaidInvoices:      int32(summary.PaidInvoices),
		UnpaidInvoices:    int32(summary.UnpaidInvoices),
		TotalAmount:       summary.TotalAmount,
		TotalPaid:         summary.TotalPaid,
		TotalOutstanding:  summary.TotalOutstanding,
		PaymentPercentage: summary.PaymentPercentage,
	}, nil
}

// Mapping helpers
// Mapping helpers
func mapInvoiceResponseToProto(i *service.InvoiceResponse) *pb.Invoice {
	if i == nil {
		return nil
	}

	var paidAt *timestamppb.Timestamp
	// ✅ PERBAIKAN: Cek nil pointer dan zero value sebelum convert
	if i.PaidAt != nil && !i.PaidAt.IsZero() {
		paidAt = timestamppb.New(*i.PaidAt)
	}

	return &pb.Invoice{
		Id:               i.ID,
		SchoolId:         i.SchoolID,
		AcademicPeriodId: i.AcademicPeriodID,
		StudentId:        i.StudentID,
		ComponentType:    i.ComponentType,
		Month:            int32(i.Month),
		Year:             int32(i.Year),
		Amount:           i.Amount,
		DueDate:          timestamppb.New(i.DueDate),
		Status:           i.Status,
		PaidAmount:       i.PaidAmount,
		PaidAt:           paidAt,
		Notes:            i.Notes,
		CreatedBy:        i.CreatedBy,
		CreatedAt:        timestamppb.New(i.CreatedAt),
		UpdatedAt:        timestamppb.New(i.UpdatedAt),
	}
}

func mapPaymentResponseToProto(p *service.PaymentResponse) *pb.Payment {
	if p == nil {
		return nil
	}
	return &pb.Payment{
		Id:               p.ID,
		SchoolId:         p.SchoolID,
		AcademicPeriodId: p.AcademicPeriodID,
		InvoiceId:        p.InvoiceID,
		Amount:           p.Amount,
		PaymentMethod:    p.PaymentMethod,
		PaymentDate:      timestamppb.New(p.PaymentDate),
		ReferenceNumber:  p.ReferenceNumber,
		Notes:            p.Notes,
		CreatedBy:        p.CreatedBy,
		CreatedAt:        timestamppb.New(p.CreatedAt),
	}
}
