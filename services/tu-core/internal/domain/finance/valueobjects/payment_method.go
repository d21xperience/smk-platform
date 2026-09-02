package valueobjects

import "errors"

type PaymentMethod string

const (
	MethodCash     PaymentMethod = "CASH"
	MethodTransfer PaymentMethod = "TRANSFER"
	MethodVA       PaymentMethod = "VA"
	MethodQRIS     PaymentMethod = "QRIS"
)

func NewPaymentMethod(m string) (PaymentMethod, error) {
	switch PaymentMethod(m) {
	case MethodCash, MethodTransfer, MethodVA, MethodQRIS:
		return PaymentMethod(m), nil
	default:
		return "", errors.New("invalid payment method: must be CASH, TRANSFER, VA, or QRIS")
	}
}

func (m PaymentMethod) String() string {
	return string(m)
}
