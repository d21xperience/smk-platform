package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/dispatcher"
)

// InvoiceRepository adalah interface untuk mengambil data invoice dan student.
// Didefinisikan di sini untuk menghindari circular dependency.
type InvoiceRepository interface {
	GetStudentContactInfo(ctx context.Context, studentID string) (*StudentContactInfo, error)
}

// StudentContactInfo berisi informasi kontak siswa dan orang tua.
type StudentContactInfo struct {
	StudentID   string
	StudentName string
	ParentName  string
	ParentPhone string
	ParentEmail string
}

// InvoiceHandler menangani notifikasi terkait invoice/tagihan.
type InvoiceHandler struct {
	dispatcher *dispatcher.NotificationDispatcher
	repo       InvoiceRepository
}

// NewInvoiceHandler membuat instance baru dari InvoiceHandler.
func NewInvoiceHandler(dispatcher *dispatcher.NotificationDispatcher, repo InvoiceRepository) *InvoiceHandler {
	return &InvoiceHandler{
		dispatcher: dispatcher,
		repo:       repo,
	}
}

// invoiceCreatedEvent adalah struktur event InvoiceCreated.
type invoiceCreatedEvent struct {
	InvoiceID        string  `json:"invoiceId"`
	SchoolID         string  `json:"schoolId"`
	AcademicPeriodID string  `json:"academicPeriodId"`
	StudentID        string  `json:"studentId"`
	ComponentType    string  `json:"componentType"` // SPP, UNIFORM, dll
	Amount           float64 `json:"amount"`
	DueDate          string  `json:"dueDate"` // Format YYYY-MM-DD
	CorrelationID    string  `json:"correlationId"`
}

// Handle memproses event InvoiceCreated dan mengirim notifikasi.
func (h *InvoiceHandler) Handle(ctx context.Context, msg []byte) error {
	// 1. Parse event
	var event invoiceCreatedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event InvoiceCreated: %w", err)
	}

	log.Printf("[InvoiceHandler] Memproses invoice %s untuk student %s", event.InvoiceID, event.StudentID)

	// 2. Fetch data kontak orang tua dari repository
	contactInfo, err := h.repo.GetStudentContactInfo(ctx, event.StudentID)
	if err != nil {
		return fmt.Errorf("gagal fetch contact info untuk student %s: %w", event.StudentID, err)
	}

	// 3. Siapkan data untuk template
	templateData := map[string]interface{}{
		"studentName":   contactInfo.StudentName,
		"parentName":    contactInfo.ParentName,
		"componentType": event.ComponentType,
		"amount":        event.Amount,
		"dueDate":       event.DueDate,
		"invoiceId":     event.InvoiceID,
	}

	// 4. Kirim notifikasi WhatsApp
	if contactInfo.ParentPhone != "" {
		reqWA := contract.NotificationRequest{
			Channel:          "whatsapp",
			Recipient:        contactInfo.ParentPhone,
			TemplateCode:     "invoice_created_wa",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqWA); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim WA untuk invoice %s: %v", event.InvoiceID, err)
			// Jangan return error, lanjutkan ke channel lain
		}
	}

	// 5. Kirim notifikasi Email
	if contactInfo.ParentEmail != "" {
		reqEmail := contract.NotificationRequest{
			Channel:          "email",
			Recipient:        contactInfo.ParentEmail,
			TemplateCode:     "invoice_created_email",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqEmail); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim Email untuk invoice %s: %v", event.InvoiceID, err)
		}
	}

	log.Printf("[InvoiceHandler] Notifikasi invoice %s berhasil diproses", event.InvoiceID)
	return nil
}

// // paymentReceivedEvent adalah struktur event PaymentReceived.
// type paymentReceivedEvent struct {
// 	InvoiceID        string  `json:"invoiceId"`
// 	SchoolID         string  `json:"schoolId"`
// 	AcademicPeriodID string  `json:"academicPeriodId"`
// 	StudentID        string  `json:"studentId"`
// 	Amount           float64 `json:"amount"`
// 	PaymentDate      string  `json:"paymentDate"`
// 	PaymentMethod    string  `json:"paymentMethod"`
// 	CorrelationID    string  `json:"correlationId"`
// }

// // HandlePaymentReceived memproses event PaymentReceived.
// func (h *InvoiceHandler) HandlePaymentReceived(ctx context.Context, msg []byte) error {
// 	var event paymentReceivedEvent
// 	if err := json.Unmarshal(msg, &event); err != nil {
// 		return fmt.Errorf("gagal parse event PaymentReceived: %w", err)
// 	}

// 	log.Printf("[InvoiceHandler] Memproses pembayaran invoice %s", event.InvoiceID)

// 	// Fetch contact info
// 	contactInfo, err := h.repo.GetStudentContactInfo(ctx, event.StudentID)
// 	if err != nil {
// 		return fmt.Errorf("gagal fetch contact info: %w", err)
// 	}

// 	templateData := map[string]interface{}{
// 		"studentName":   contactInfo.StudentName,
// 		"parentName":    contactInfo.ParentName,
// 		"amount":        event.Amount,
// 		"paymentDate":   event.PaymentDate,
// 		"paymentMethod": event.PaymentMethod,
// 		"invoiceId":     event.InvoiceID,
// 	}

// 	// Kirim notifikasi konfirmasi pembayaran
// 	if contactInfo.ParentPhone != "" {
// 		reqWA := contract.NotificationRequest{
// 			Channel:          "whatsapp",
// 			Recipient:        contactInfo.ParentPhone,
// 			TemplateCode:     "payment_received_wa",
// 			Data:             templateData,
// 			SchoolID:         event.SchoolID,
// 			AcademicPeriodID: event.AcademicPeriodID,
// 			CorrelationID:    event.CorrelationID,
// 		}
// 		h.dispatcher.Dispatch(ctx, reqWA)
// 	}

// 	if contactInfo.ParentEmail != "" {
// 		reqEmail := contract.NotificationRequest{
// 			Channel:          "email",
// 			Recipient:        contactInfo.ParentEmail,
// 			TemplateCode:     "payment_received_email",
// 			Data:             templateData,
// 			SchoolID:         event.SchoolID,
// 			AcademicPeriodID: event.AcademicPeriodID,
// 			CorrelationID:    event.CorrelationID,
// 		}
// 		h.dispatcher.Dispatch(ctx, reqEmail)
// 	}

// 	return nil
// }

// paymentReceivedEvent adalah struktur event saat pembayaran berhasil.
type paymentReceivedEvent struct {
	InvoiceID        string  `json:"invoiceId"`
	SchoolID         string  `json:"schoolId"`
	SchoolName       string  `json:"schoolName"`
	AcademicPeriodID string  `json:"academicPeriodId"`
	StudentID        string  `json:"studentId"`
	ComponentType    string  `json:"componentType"` // e.g., "SPP"
	Month            string  `json:"month"`         // e.g., "Agustus"
	Year             string  `json:"year"`          // e.g., "2026"
	Amount           float64 `json:"amount"`
	PaymentDate      string  `json:"paymentDate"`
	CorrelationID    string  `json:"correlationId"`
}

// HandlePaymentReceived memproses event PaymentReceived.
func (h *InvoiceHandler) HandlePaymentReceived(ctx context.Context, msg []byte) error {
	var event paymentReceivedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event PaymentReceived: %w", err)
	}

	log.Printf("[InvoiceHandler] Memproses notifikasi pembayaran invoice %s", event.InvoiceID)

	contactInfo, err := h.repo.GetStudentContactInfo(ctx, event.StudentID)
	if err != nil {
		return fmt.Errorf("gagal fetch contact info: %w", err)
	}

	templateData := map[string]interface{}{
		"studentName":   contactInfo.StudentName,
		"parentName":    contactInfo.ParentName,
		"schoolName":    event.SchoolName,
		"componentType": event.ComponentType,
		"month":         event.Month,
		"year":          event.Year,
		"amount":        event.Amount,
		"paymentDate":   event.PaymentDate,
		"invoiceId":     event.InvoiceID,
	}

	// 1. Kirim WhatsApp
	if contactInfo.ParentPhone != "" {
		reqWA := contract.NotificationRequest{
			Channel:          "whatsapp",
			Recipient:        contactInfo.ParentPhone,
			TemplateCode:     "finance_payment_received_wa",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqWA); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim WA pembayaran: %v", err)
		}
	}

	// 2. Kirim Email
	if contactInfo.ParentEmail != "" {
		reqEmail := contract.NotificationRequest{
			Channel:          "email",
			Recipient:        contactInfo.ParentEmail,
			TemplateCode:     "finance_payment_received_email",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqEmail); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim Email pembayaran: %v", err)
		}
	}

	return nil
}

// invoiceOverdueEvent adalah struktur event saat tagihan melewati jatuh tempo.
type invoiceOverdueEvent struct {
	InvoiceID        string  `json:"invoiceId"`
	SchoolID         string  `json:"schoolId"`
	SchoolName       string  `json:"schoolName"`
	AcademicPeriodID string  `json:"academicPeriodId"`
	StudentID        string  `json:"studentId"`
	ComponentType    string  `json:"componentType"`
	Month            string  `json:"month"`
	Year             string  `json:"year"`
	Amount           float64 `json:"amount"`
	DueDate          string  `json:"dueDate"`
	CorrelationID    string  `json:"correlationId"`
}

// HandleInvoiceOverdue memproses event InvoiceOverdue (Pengingat Bulanan).
func (h *InvoiceHandler) HandleInvoiceOverdue(ctx context.Context, msg []byte) error {
	var event invoiceOverdueEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event InvoiceOverdue: %w", err)
	}

	log.Printf("[InvoiceHandler] Memproses notifikasi overdue invoice %s", event.InvoiceID)

	contactInfo, err := h.repo.GetStudentContactInfo(ctx, event.StudentID)
	if err != nil {
		return fmt.Errorf("gagal fetch contact info: %w", err)
	}

	templateData := map[string]interface{}{
		"studentName":   contactInfo.StudentName,
		"parentName":    contactInfo.ParentName,
		"schoolName":    event.SchoolName,
		"componentType": event.ComponentType,
		"month":         event.Month,
		"year":          event.Year,
		"amount":        event.Amount,
		"dueDate":       event.DueDate,
		"invoiceId":     event.InvoiceID,
	}

	// 1. Kirim WhatsApp Pengingat
	if contactInfo.ParentPhone != "" {
		reqWA := contract.NotificationRequest{
			Channel:          "whatsapp",
			Recipient:        contactInfo.ParentPhone,
			TemplateCode:     "finance_payment_reminder_wa",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqWA); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim WA reminder: %v", err)
		}
	}

	// 2. Kirim Email Pengingat
	if contactInfo.ParentEmail != "" {
		reqEmail := contract.NotificationRequest{
			Channel:          "email",
			Recipient:        contactInfo.ParentEmail,
			TemplateCode:     "finance_payment_reminder_email",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqEmail); err != nil {
			log.Printf("[InvoiceHandler] Gagal kirim Email reminder: %v", err)
		}
	}

	return nil
}
