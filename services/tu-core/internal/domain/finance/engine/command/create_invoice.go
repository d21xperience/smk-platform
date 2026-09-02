package commands

import (
	"errors"
	"time"
)

// CreateInvoiceCommand adalah command untuk membuat tagihan.
type CreateInvoiceCommand struct {
	InvoiceID   string
	StudentID   string
	InvoiceType string
	Description string
	Amount      int64
	DueDate     time.Time
}

// Validate memvalidasi command.
func (c CreateInvoiceCommand) Validate() error {
	if c.StudentID == "" {
		return errors.New("studentId wajib diisi")
	}
	if c.InvoiceType == "" {
		return errors.New("invoiceType wajib diisi")
	}
	if c.Amount <= 0 {
		return errors.New("jumlah tagihan harus > 0")
	}
	if c.DueDate.IsZero() {
		return errors.New("tanggal jatuh tempo wajib diisi")
	}
	if c.DueDate.Before(time.Now()) {
		return errors.New("tanggal jatuh tempo tidak boleh di masa lalu")
	}
	return nil
}
