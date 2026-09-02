package commands

import (
	"errors"
	"time"
)

// ProcessPaymentCommand adalah command untuk memproses pembayaran.
type ProcessPaymentCommand struct {
	PaymentID    string
	InvoiceID    string
	Amount       int64
	Method       string
	Reference    string
	ReceivedBy   string
	ReceivedDate time.Time
	Notes        string
}

// Validate memvalidasi command.
func (c ProcessPaymentCommand) Validate() error {
	if c.InvoiceID == "" {
		return errors.New("invoiceId wajib diisi")
	}
	if c.Amount <= 0 {
		return errors.New("jumlah pembayaran harus > 0")
	}
	if c.Method == "" {
		return errors.New("metode pembayaran wajib diisi")
	}
	if c.ReceivedBy == "" {
		return errors.New("penerima pembayaran wajib diisi")
	}
	return nil
}
