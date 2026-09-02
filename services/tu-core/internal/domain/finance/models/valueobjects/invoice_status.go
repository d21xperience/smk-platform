package valueobjects

import "fmt"

// InvoiceStatus adalah value object untuk status tagihan.
type InvoiceStatus struct {
	code  string
	label string
}

var (
	InvoiceStatusPending   = InvoiceStatus{code: "PENDING", label: "Belum Bayar"}
	InvoiceStatusPartial   = InvoiceStatus{code: "PARTIAL", label: "Bayar Sebagian"}
	InvoiceStatusSettled   = InvoiceStatus{code: "SETTLED", label: "Lunas"}
	InvoiceStatusOverdue   = InvoiceStatus{code: "OVERDUE", label: "Jatuh Tempo"}
	InvoiceStatusCancelled = InvoiceStatus{code: "CANCELLED", label: "Dibatalkan"}
)

// NewInvoiceStatus membuat InvoiceStatus dari code.
func NewInvoiceStatus(code string) (InvoiceStatus, error) {
	switch code {
	case "PENDING":
		return InvoiceStatusPending, nil
	case "PARTIAL":
		return InvoiceStatusPartial, nil
	case "SETTLED":
		return InvoiceStatusSettled, nil
	case "OVERDUE":
		return InvoiceStatusOverdue, nil
	case "CANCELLED":
		return InvoiceStatusCancelled, nil
	default:
		return InvoiceStatus{}, fmt.Errorf("status tagihan tidak valid: %s", code)
	}
}

// Code mengembalikan code status.
func (s InvoiceStatus) Code() string { return s.code }

// Label mengembalikan label status.
func (s InvoiceStatus) Label() string { return s.label }

// IsPayable mengecek apakah invoice masih bisa dibayar.
func (s InvoiceStatus) IsPayable() bool {
	return s.code == "PENDING" || s.code == "PARTIAL" || s.code == "OVERDUE"
}

// IsSettled mengecek apakah invoice sudah lunas.
func (s InvoiceStatus) IsSettled() bool { return s.code == "SETTLED" }

// IsCancelled mengecek apakah invoice dibatalkan.
func (s InvoiceStatus) IsCancelled() bool { return s.code == "CANCELLED" }
