package models

import (
	"errors"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/finance/valueobjects"

	"github.com/google/uuid"
)

type Payment struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	InvoiceID        string
	Amount           float64
	PaymentMethod    valueobjects.PaymentMethod
	PaymentDate      time.Time
	ReferenceNumber  string
	Notes            string
	CreatedBy        string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type PaymentData struct {
	ID              string
	InvoiceID       string
	Amount          float64
	PaymentMethod   valueobjects.PaymentMethod
	PaymentDate     time.Time
	ReferenceNumber string
	Notes           string
}

func NewPayment(data PaymentData, schoolID, academicPeriodID, createdBy string) (*Payment, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.InvoiceID == "" {
		return nil, errors.New("invoice_id wajib diisi")
	}
	if data.Amount <= 0 {
		return nil, errors.New("amount harus lebih dari 0")
	}
	if data.PaymentDate.IsZero() {
		return nil, errors.New("payment_date wajib diisi")
	}

	p := &Payment{
		ID:               data.ID,
		SchoolID:         schoolID,
		AcademicPeriodID: academicPeriodID,
		InvoiceID:        data.InvoiceID,
		Amount:           data.Amount,
		PaymentMethod:    data.PaymentMethod,
		PaymentDate:      data.PaymentDate,
		ReferenceNumber:  data.ReferenceNumber,
		Notes:            data.Notes,
		CreatedBy:        createdBy,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
	}

	return p, nil
}
