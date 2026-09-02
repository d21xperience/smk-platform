package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/finance/events"
	"sekolah-platform/services/tu-core/internal/domain/finance/valueobjects"

	"github.com/google/uuid"
)

type Invoice struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        string
	ComponentType    valueobjects.ComponentType
	Month            int
	Year             int
	Amount           float64
	DueDate          time.Time
	Status           valueobjects.InvoiceStatus
	PaidAmount       float64
	PaidAt           time.Time
	Notes            string
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
	domainEvents     []types.DomainEvent
}

type InvoiceData struct {
	ID            string
	StudentID     string
	ComponentType valueobjects.ComponentType
	Month         int
	Year          int
	Amount        float64
	DueDate       time.Time
	Notes         string
}

func NewInvoice(data InvoiceData, ctx *context.OperationalContext) (*Invoice, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	// Validasi
	if data.StudentID == "" {
		return nil, errors.New("student_id wajib diisi")
	}
	if data.Amount <= 0 {
		return nil, errors.New("amount harus lebih dari 0")
	}
	if data.DueDate.IsZero() {
		return nil, errors.New("due_date wajib diisi")
	}

	// Validasi khusus untuk SPP (periodic)
	if data.ComponentType == valueobjects.ComponentSPP {
		if data.Month < 1 || data.Month > 12 {
			return nil, errors.New("month harus 1-12 untuk komponen SPP")
		}
		if data.Year == 0 {
			return nil, errors.New("year wajib diisi untuk komponen SPP")
		}
	}

	i := &Invoice{
		ID:               data.ID,
		SchoolID:         ctx.SchoolID,
		AcademicPeriodID: ctx.AcademicPeriodID,
		StudentID:        data.StudentID,
		ComponentType:    data.ComponentType,
		Month:            data.Month,
		Year:             data.Year,
		Amount:           data.Amount,
		DueDate:          data.DueDate,
		Status:           valueobjects.StatusUnpaid,
		PaidAmount:       0,
		CreatedBy:        ctx.UserID,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		Notes:            data.Notes,
		domainEvents:     []types.DomainEvent{},
	}

	i.recordEvent(events.NewInvoiceCreated(events.InvoiceCreatedPayload{
		InvoiceID:     i.ID,
		StudentID:     i.StudentID,
		ComponentType: i.ComponentType.String(),
		Amount:        i.Amount,
		Month:         i.Month,
		Year:          i.Year,
	}, ctx))

	return i, nil
}

func (i *Invoice) ApplyPayment(amount float64, paymentID string, ctx *context.OperationalContext) error {
	if i.Status.IsCancelled() {
		return errors.New("tidak dapat melakukan pembayaran pada invoice yang dibatalkan")
	}
	if amount <= 0 {
		return errors.New("amount pembayaran harus lebih dari 0")
	}

	remaining := i.Amount - i.PaidAmount
	if amount > remaining {
		return errors.New("amount pembayaran melebihi sisa tagihan")
	}

	i.PaidAmount += amount
	i.UpdatedAt = time.Now().UTC()

	// Update status berdasarkan paid_amount
	if i.PaidAmount >= i.Amount {
		i.Status = valueobjects.StatusPaid
		i.PaidAt = time.Now().UTC()
	} else {
		i.Status = valueobjects.StatusPartial
	}

	i.recordEvent(events.NewPaymentApplied(events.PaymentAppliedPayload{
		InvoiceID: i.ID,
		PaymentID: paymentID,
		Amount:    amount,
		Status:    i.Status.String(),
	}, ctx))

	return nil
}

func (i *Invoice) Cancel(reason string, ctx *context.OperationalContext) error {
	if i.Status.IsPaid() {
		return errors.New("tidak dapat membatalkan invoice yang sudah lunas")
	}

	i.Status = valueobjects.StatusCancelled
	i.Notes = reason
	i.UpdatedAt = time.Now().UTC()

	i.recordEvent(events.NewInvoiceCancelled(events.InvoiceCancelledPayload{
		InvoiceID: i.ID,
		Reason:    reason,
	}, ctx))

	return nil
}

func (i *Invoice) Outstanding() float64 {
	return i.Amount - i.PaidAmount
}

func (i *Invoice) recordEvent(event types.DomainEvent) {
	i.domainEvents = append(i.domainEvents, event)
}

func (i *Invoice) UncommittedEvents() []types.DomainEvent {
	return i.domainEvents
}

func (i *Invoice) ClearEvents() {
	i.domainEvents = []types.DomainEvent{}
}
