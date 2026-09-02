package service

import (
	"context"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/engine"
	"sekolah-platform/services/tu-core/internal/domain/finance/models"
	"sekolah-platform/services/tu-core/internal/domain/finance/valueobjects"
	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/google/uuid"
)

type FinanceService struct {
	repo        repository.FinanceRepository
	studentRepo repository.StudentRepository
	finEngine   *engine.FinanceEngine
	publisher   EventPublisher
}

func NewFinanceService(
	repo repository.FinanceRepository,
	studentRepo repository.StudentRepository,
	finEngine *engine.FinanceEngine,
	publisher EventPublisher,
) *FinanceService {
	if publisher == nil {
		publisher = NewNoopEventPublisher()
	}
	return &FinanceService{
		repo:        repo,
		studentRepo: studentRepo,
		finEngine:   finEngine,
		publisher:   publisher,
	}
}

type CreateInvoiceServiceRequest struct {
	StudentID     string
	ComponentType valueobjects.ComponentType
	Month         int
	Year          int
	Amount        float64
	DueDate       time.Time
	Notes         string
}

func (s *FinanceService) CreateInvoice(ctx context.Context, req CreateInvoiceServiceRequest) (*InvoiceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	data := models.InvoiceData{
		StudentID:     req.StudentID,
		ComponentType: req.ComponentType,
		Month:         req.Month,
		Year:          req.Year,
		Amount:        req.Amount,
		DueDate:       req.DueDate,
		Notes:         req.Notes,
	}

	result := s.finEngine.CreateInvoice(data, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreateInvoice", "Invoice", "", "invalid invoice data", result.Error)
	}

	if err := s.repo.CreateInvoice(ctx, result.Invoice); err != nil {
		return nil, NewServiceError("CreateInvoice", "Invoice", result.Invoice.ID, "failed to save to database", err)
	}

	event := DomainEvent{
		EventType:     "InvoiceCreated",
		AggregateType: "Invoice",
		AggregateID:   result.Invoice.ID,
		Payload: map[string]interface{}{
			"invoiceId":     result.Invoice.ID,
			"studentId":     result.Invoice.StudentID,
			"componentType": result.Invoice.ComponentType.String(),
			"amount":        result.Invoice.Amount,
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return toInvoiceResponse(result.Invoice), nil
}

func (s *FinanceService) GetInvoice(ctx context.Context, invoiceID string) (*InvoiceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	i, err := s.repo.GetInvoiceByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, invoiceID)
	if err != nil {
		return nil, NewServiceError("GetInvoice", "Invoice", invoiceID, "failed to get invoice", err)
	}

	return toInvoiceResponse(i), nil
}

func (s *FinanceService) ListInvoices(ctx context.Context, studentID, componentType, status string, month, year int) ([]*InvoiceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.InvoiceFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        studentID,
		ComponentType:    componentType,
		Status:           status,
		Month:            month,
		Year:             year,
		Limit:            100,
		Offset:           0,
	}

	invoices, _, err := s.repo.ListInvoices(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListInvoices", "Invoice", "", "failed to list invoices", err)
	}

	var responses []*InvoiceResponse
	for _, i := range invoices {
		responses = append(responses, toInvoiceResponse(&i))
	}

	return responses, nil
}

func (s *FinanceService) GetStudentOutstanding(ctx context.Context, studentID string) (*StudentOutstandingResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.InvoiceFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		StudentID:        studentID,
		Status:           "UNPAID",
		Limit:            1000,
		Offset:           0,
	}

	invoices, _, err := s.repo.ListInvoices(ctx, filter)
	if err != nil {
		return nil, NewServiceError("GetStudentOutstanding", "Invoice", studentID, "failed to get invoices", err)
	}

	var totalOutstanding float64
	var invoiceResponses []InvoiceResponse
	for _, i := range invoices {
		totalOutstanding += i.Outstanding()
		invoiceResponses = append(invoiceResponses, *toInvoiceResponse(&i))
	}

	return &StudentOutstandingResponse{
		StudentID:          studentID,
		TotalOutstanding:   totalOutstanding,
		UnpaidInvoices:     len(invoices),
		UnpaidInvoicesList: invoiceResponses,
	}, nil
}

type CreatePaymentServiceRequest struct {
	InvoiceID       string
	Amount          float64
	PaymentMethod   valueobjects.PaymentMethod
	PaymentDate     time.Time
	ReferenceNumber string
	Notes           string
}

func (s *FinanceService) CreatePayment(ctx context.Context, req CreatePaymentServiceRequest) (*PaymentWithInvoiceResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	// 1. Ambil invoice
	invoice, err := s.repo.GetInvoiceByID(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, req.InvoiceID)
	if err != nil {
		return nil, NewServiceError("CreatePayment", "Invoice", req.InvoiceID, "failed to get invoice", err)
	}

	// 2. Buat payment object
	paymentData := models.PaymentData{
		InvoiceID:       req.InvoiceID,
		Amount:          req.Amount,
		PaymentMethod:   req.PaymentMethod,
		PaymentDate:     req.PaymentDate,
		ReferenceNumber: req.ReferenceNumber,
		Notes:           req.Notes,
	}

	payment, err := models.NewPayment(paymentData, opCtx.SchoolID, opCtx.AcademicPeriodID, opCtx.UserID)
	if err != nil {
		return nil, NewServiceError("CreatePayment", "Payment", "", "invalid payment data", err)
	}

	// 3. Apply payment ke invoice (domain logic)
	result := s.finEngine.ApplyPayment(invoice, req.Amount, payment.ID, opCtx)
	if result.Error != nil {
		return nil, NewServiceError("CreatePayment", "Invoice", req.InvoiceID, "failed to apply payment", result.Error)
	}

	// 4. Simpan payment
	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, NewServiceError("CreatePayment", "Payment", payment.ID, "failed to save payment", err)
	}

	// 5. Update invoice
	if err := s.repo.UpdateInvoice(ctx, result.Invoice); err != nil {
		return nil, NewServiceError("CreatePayment", "Invoice", req.InvoiceID, "failed to update invoice", err)
	}

	// 6. Publish event
	event := DomainEvent{
		EventType:     "PaymentApplied",
		AggregateType: "Invoice",
		AggregateID:   invoice.ID,
		Payload: map[string]interface{}{
			"invoiceId": invoice.ID,
			"paymentId": payment.ID,
			"amount":    req.Amount,
			"method":    req.PaymentMethod.String(),
		},
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		CorrelationID:    uuid.New().String(),
		Version:          "v1",
	}
	_ = s.publisher.Publish(ctx, event)

	return &PaymentWithInvoiceResponse{
		Payment:        toPaymentResponse(payment),
		UpdatedInvoice: toInvoiceResponse(result.Invoice),
	}, nil
}

func (s *FinanceService) ListPayments(ctx context.Context, invoiceID, paymentMethod string, startDate, endDate time.Time) ([]*PaymentResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	filter := repository.PaymentFilter{
		SchoolID:         opCtx.SchoolID,
		AcademicPeriodID: opCtx.AcademicPeriodID,
		InvoiceID:        invoiceID,
		PaymentMethod:    paymentMethod,
		StartDate:        startDate,
		EndDate:          endDate,
		Limit:            100,
		Offset:           0,
	}

	payments, _, err := s.repo.ListPayments(ctx, filter)
	if err != nil {
		return nil, NewServiceError("ListPayments", "Payment", "", "failed to list payments", err)
	}

	var responses []*PaymentResponse
	for _, p := range payments {
		responses = append(responses, toPaymentResponse(&p))
	}

	return responses, nil
}

func (s *FinanceService) GetPaymentSummary(ctx context.Context, month, year int, componentType string) (*PaymentSummaryResponse, error) {
	opCtx := middleware.GetOperationalContext(ctx)
	if opCtx == nil {
		return nil, ErrMissingOperationalContext
	}

	summary, err := s.repo.GetPaymentSummary(ctx, opCtx.SchoolID, opCtx.AcademicPeriodID, month, year, componentType)
	if err != nil {
		return nil, NewServiceError("GetPaymentSummary", "Invoice", "", "failed to get summary", err)
	}

	return &PaymentSummaryResponse{
		TotalInvoices:     summary.TotalInvoices,
		PaidInvoices:      summary.PaidInvoices,
		UnpaidInvoices:    summary.UnpaidInvoices,
		TotalAmount:       summary.TotalAmount,
		TotalPaid:         summary.TotalPaid,
		TotalOutstanding:  summary.TotalOutstanding,
		PaymentPercentage: summary.PaymentPercentage,
	}, nil
}

type PaymentWithInvoiceResponse struct {
	Payment        *PaymentResponse
	UpdatedInvoice *InvoiceResponse
}

type StudentOutstandingResponse struct {
	StudentID          string
	TotalOutstanding   float64
	UnpaidInvoices     int
	UnpaidInvoicesList []InvoiceResponse
}

type PaymentSummaryResponse struct {
	TotalInvoices     int
	PaidInvoices      int
	UnpaidInvoices    int
	TotalAmount       float64
	TotalPaid         float64
	TotalOutstanding  float64
	PaymentPercentage float64
}

func toInvoiceResponse(i *models.Invoice) *InvoiceResponse {
	if i == nil {
		return nil
	}

	var paidAt *time.Time
	if !i.PaidAt.IsZero() {
		paidAt = &i.PaidAt
	}

	return &InvoiceResponse{
		ID:               i.ID,
		SchoolID:         i.SchoolID,
		AcademicPeriodID: i.AcademicPeriodID,
		StudentID:        i.StudentID,
		ComponentType:    i.ComponentType.String(),
		Month:            i.Month,
		Year:             i.Year,
		Amount:           i.Amount,
		DueDate:          i.DueDate,
		Status:           i.Status.String(),
		PaidAmount:       i.PaidAmount,
		PaidAt:           paidAt,
		Notes:            i.Notes,
		CreatedBy:        i.CreatedBy,
		CreatedAt:        i.CreatedAt,
		UpdatedAt:        i.UpdatedAt,
	}
}

func toPaymentResponse(p *models.Payment) *PaymentResponse {
	if p == nil {
		return nil
	}
	return &PaymentResponse{
		ID:               p.ID,
		SchoolID:         p.SchoolID,
		AcademicPeriodID: p.AcademicPeriodID,
		InvoiceID:        p.InvoiceID,
		Amount:           p.Amount,
		PaymentMethod:    p.PaymentMethod.String(),
		PaymentDate:      p.PaymentDate,
		ReferenceNumber:  p.ReferenceNumber,
		Notes:            p.Notes,
		CreatedBy:        p.CreatedBy,
		CreatedAt:        p.CreatedAt,
	}
}
