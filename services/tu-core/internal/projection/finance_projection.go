package projection

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FinanceProjection menangani projection untuk domain Finance.
type FinanceProjection struct {
	db *pgxpool.Pool
}

// NewFinanceProjection membuat instance baru.
func NewFinanceProjection(db *pgxpool.Pool) *FinanceProjection {
	return &FinanceProjection{db: db}
}

// Name mengembalikan nama projection.
func (p *FinanceProjection) Name() string {
	return "finance_projection"
}

// SupportedEvents mengembalikan daftar event yang didukung.
func (p *FinanceProjection) SupportedEvents() []string {
	return []string{
		"InvoiceCreated",
		"PaymentReceived",
		"InvoiceOverdue",
	}
}

// HandleEvent memproses event finance.
func (p *FinanceProjection) HandleEvent(ctx context.Context, event DomainEvent) error {
	switch event.EventType {
	case "InvoiceCreated":
		return p.handleInvoiceCreated(ctx, event)
	case "PaymentReceived":
		return p.handlePaymentReceived(ctx, event)
	case "InvoiceOverdue":
		return p.handleInvoiceOverdue(ctx, event)
	default:
		return fmt.Errorf("event type '%s' tidak didukung", event.EventType)
	}
}

// handleInvoiceCreated mengupdate read model saat invoice dibuat.
func (p *FinanceProjection) handleInvoiceCreated(ctx context.Context, event DomainEvent) error {
	invoiceID := event.Payload["invoiceId"].(string)
	studentID := event.Payload["studentId"].(string)
	amount := event.Payload["amount"].(float64)
	status := "UNPAID"

	query := `
		INSERT INTO projection_invoice_summary
			(school_id, academic_period_id, invoice_id, student_id, amount, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
		ON CONFLICT (invoice_id) DO UPDATE SET
			amount = EXCLUDED.amount,
			status = EXCLUDED.status,
			updated_at = NOW()
	`

	_, err := p.db.Exec(ctx, query,
		event.SchoolID,
		event.AcademicPeriodID,
		invoiceID,
		studentID,
		amount,
		status,
	)
	if err != nil {
		return fmt.Errorf("gagal upsert invoice summary: %w", err)
	}

	log.Printf("[FinanceProjection] Invoice %s created (CorrelationID: %s)",
		invoiceID, event.CorrelationID)
	return nil
}

// handlePaymentReceived mengupdate read model saat pembayaran diterima.
func (p *FinanceProjection) handlePaymentReceived(ctx context.Context, event DomainEvent) error {
	invoiceID := event.Payload["invoiceId"].(string)
	amount := event.Payload["amount"].(float64)

	query := `
		UPDATE projection_invoice_summary
		SET status = 'PAID', paid_amount = $1, paid_at = NOW(), updated_at = NOW()
		WHERE invoice_id = $2
	`

	_, err := p.db.Exec(ctx, query, amount, invoiceID)
	if err != nil {
		return fmt.Errorf("gagal update invoice paid: %w", err)
	}

	log.Printf("[FinanceProjection] Invoice %s paid (CorrelationID: %s)",
		invoiceID, event.CorrelationID)
	return nil
}

// handleInvoiceOverdue mengupdate status invoice menjadi overdue.
func (p *FinanceProjection) handleInvoiceOverdue(ctx context.Context, event DomainEvent) error {
	invoiceID := event.Payload["invoiceId"].(string)

	query := `
		UPDATE projection_invoice_summary
		SET status = 'OVERDUE', updated_at = NOW()
		WHERE invoice_id = $1
	`

	_, err := p.db.Exec(ctx, query, invoiceID)
	if err != nil {
		return fmt.Errorf("gagal update invoice overdue: %w", err)
	}

	log.Printf("[FinanceProjection] Invoice %s marked as overdue (CorrelationID: %s)",
		invoiceID, event.CorrelationID)
	return nil
}
