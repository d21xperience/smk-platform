package valueobjects

import "fmt"

// PaymentMethod adalah value object untuk metode pembayaran.
type PaymentMethod struct {
	code  string
	label string
}

var (
	PaymentMethodCash       = PaymentMethod{code: "CASH", label: "Tunai"}
	PaymentMethodTransfer   = PaymentMethod{code: "TRANSFER", label: "Transfer Bank"}
	PaymentMethodEWallet    = PaymentMethod{code: "E_WALLET", label: "E-Wallet"}
	PaymentMethodQRIS       = PaymentMethod{code: "QRIS", label: "QRIS"}
	PaymentMethodInstallment = PaymentMethod{code: "INSTALLMENT", label: "Cicilan"}
)

// NewPaymentMethod membuat PaymentMethod dari code.
func NewPaymentMethod(code string) (PaymentMethod, error) {
	switch code {
	case "CASH":
		return PaymentMethodCash, nil
	case "TRANSFER":
		return PaymentMethodTransfer, nil
	case "E_WALLET":
		return PaymentMethodEWallet, nil
	case "QRIS":
		return PaymentMethodQRIS, nil
	case "INSTALLMENT":
		return PaymentMethodInstallment, nil
	default:
		return PaymentMethod{}, fmt.Errorf("metode pembayaran tidak valid: %s", code)
	}
}

// Code mengembalikan code metode.
func (m PaymentMethod) Code() string { return m.code }

// Label mengembalikan label metode.
func (m PaymentMethod) Label() string { return m.label }
