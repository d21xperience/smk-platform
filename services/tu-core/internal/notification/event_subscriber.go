package notification

import (
	"context"
	"fmt"
	"log"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/dispatcher"
	"sekolah-platform/services/tu-core/internal/notification/handlers"
)

// EventSubscriber mendengarkan Domain Events dari NATS dan memicu notifikasi.
// Sesuai Prinsip 11: Subscriber hanya orchestrator, business logic ada di handler.
type EventSubscriber struct {
	dispatcher        *dispatcher.NotificationDispatcher
	invoiceHandler    *handlers.InvoiceHandler
	documentHandler   *handlers.DocumentHandler
	attendanceHandler *handlers.AttendanceHandler
}

// NewEventSubscriber membuat instance baru dari EventSubscriber.
func NewEventSubscriber(
	dispatcher *dispatcher.NotificationDispatcher,
	invoiceHandler *handlers.InvoiceHandler,
	documentHandler *handlers.DocumentHandler,
	attendanceHandler *handlers.AttendanceHandler,
) *EventSubscriber {
	return &EventSubscriber{
		dispatcher:        dispatcher,
		invoiceHandler:    invoiceHandler,
		documentHandler:   documentHandler,
		attendanceHandler: attendanceHandler,
	}
}

// HandleInvoiceCreated menangani event InvoiceCreated.
func (s *EventSubscriber) HandleInvoiceCreated(ctx context.Context, msg []byte) error {
	log.Printf("[EventSubscriber] Menerima event InvoiceCreated")
	return s.invoiceHandler.Handle(ctx, msg)
}

// HandleDocumentIssued menangani event DocumentIssued.
func (s *EventSubscriber) HandleDocumentIssued(ctx context.Context, msg []byte) error {
	log.Printf("[EventSubscriber] Menerima event DocumentIssued")
	return s.documentHandler.Handle(ctx, msg)
}

// HandleAttendanceSubmitted menangani event AttendanceSubmitted.
func (s *EventSubscriber) HandleAttendanceSubmitted(ctx context.Context, msg []byte) error {
	log.Printf("[EventSubscriber] Menerima event AttendanceSubmitted")
	return s.attendanceHandler.Handle(ctx, msg)
}

// HandleStudentEnrolled menangani event StudentEnrolled (opsional: welcome notification).
func (s *EventSubscriber) HandleStudentEnrolled(ctx context.Context, msg []byte) error {
	log.Printf("[EventSubscriber] Menerima event StudentEnrolled")
	// Bisa ditambahkan handler khusus jika diperlukan
	return nil
}

// HandlePaymentReceived menangani event PaymentReceived.
func (s *EventSubscriber) HandlePaymentReceived(ctx context.Context, msg []byte) error {
	log.Printf("[EventSubscriber] Menerima event PaymentReceived")
	return s.invoiceHandler.HandlePaymentReceived(ctx, msg)
}

// GetDispatcher mengembalikan dispatcher untuk penggunaan langsung jika diperlukan.
func (s *EventSubscriber) GetDispatcher() *dispatcher.NotificationDispatcher {
	return s.dispatcher
}

// DispatchNotification adalah helper untuk mengirim notifikasi langsung (untuk kasus khusus).
func (s *EventSubscriber) DispatchNotification(ctx context.Context, req contract.NotificationRequest) error {
	if err := s.dispatcher.Dispatch(ctx, req); err != nil {
		return fmt.Errorf("gagal dispatch notifikasi: %w", err)
	}
	return nil
}
