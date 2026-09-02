package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

type InvoiceResult struct {
	Invoice *models.Invoice
	Events  []types.DomainEvent
	Error   error
}

func InvoiceSuccess(invoice *models.Invoice) InvoiceResult {
	return InvoiceResult{
		Invoice: invoice,
		Events:  invoice.UncommittedEvents(),
		Error:   nil,
	}
}

func InvoiceFailure(err error) InvoiceResult {
	return InvoiceResult{
		Invoice: nil,
		Events:  nil,
		Error:   err,
	}
}

type FinanceEngine struct{}

func NewFinanceEngine() *FinanceEngine {
	return &FinanceEngine{}
}

func (e *FinanceEngine) CreateInvoice(data models.InvoiceData, ctx *context.OperationalContext) InvoiceResult {
	i, err := models.NewInvoice(data, ctx)
	if err != nil {
		return InvoiceFailure(err)
	}
	return InvoiceSuccess(i)
}

func (e *FinanceEngine) ApplyPayment(i *models.Invoice, amount float64, paymentID string, ctx *context.OperationalContext) InvoiceResult {
	if i == nil {
		return InvoiceFailure(errInvoiceNotFound)
	}
	if err := i.ApplyPayment(amount, paymentID, ctx); err != nil {
		return InvoiceFailure(err)
	}
	return InvoiceSuccess(i)
}

func (e *FinanceEngine) CancelInvoice(i *models.Invoice, reason string, ctx *context.OperationalContext) InvoiceResult {
	if i == nil {
		return InvoiceFailure(errInvoiceNotFound)
	}
	if err := i.Cancel(reason, ctx); err != nil {
		return InvoiceFailure(err)
	}
	return InvoiceSuccess(i)
}

var errInvoiceNotFound = errString("invoice tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
