package valueobjects

import "errors"

type InvoiceStatus string

const (
	StatusUnpaid    InvoiceStatus = "UNPAID"
	StatusPartial   InvoiceStatus = "PARTIAL"
	StatusPaid      InvoiceStatus = "PAID"
	StatusCancelled InvoiceStatus = "CANCELLED"
)

func NewInvoiceStatus(s string) (InvoiceStatus, error) {
	switch InvoiceStatus(s) {
	case StatusUnpaid, StatusPartial, StatusPaid, StatusCancelled:
		return InvoiceStatus(s), nil
	default:
		return "", errors.New("invalid invoice status")
	}
}

func (s InvoiceStatus) String() string {
	return string(s)
}

func (s InvoiceStatus) IsPaid() bool {
	return s == StatusPaid
}

func (s InvoiceStatus) IsCancelled() bool {
	return s == StatusCancelled
}
