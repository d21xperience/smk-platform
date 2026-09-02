package projections

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// FinanceProjectionHandler mendengarkan event pembayaran dan mengupdate read model keuangan.
type FinanceProjectionHandler struct {
	db *pgxpool.Pool
}

// NewFinanceProjectionHandler membuat instance baru.
func NewFinanceProjectionHandler(db *pgxpool.Pool) *FinanceProjectionHandler {
	return &FinanceProjectionHandler{db: db}
}

// paymentReceivedEvent merepresentasikan payload event pembayaran.
type paymentReceivedEvent struct {
	SchoolID         string  `json:"schoolId"`
	AcademicPeriodID string  `json:"academicPeriodId"`
	StudentID        string  `json:"studentId"`
	ComponentType    string  `json:"componentType"` // e.g., "SPP", "UANG_GEDUNG"
	Month            int     `json:"month"`         // 1-12
	Year             int     `json:"year"`          // e.g., 2026
	Amount           float64 `json:"amount"`
	CorrelationID    string  `json:"correlationId"`
}

// HandlePaymentReceived memproses event pembayaran dan mengagregasi ke tabel monthly.
func (h *FinanceProjectionHandler) HandlePaymentReceived(ctx context.Context, msg []byte) error {
	var event paymentReceivedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event PaymentReceived: %w", err)
	}

	// UPSERT query untuk agregasi bulanan
	// Kita menambahkan amount ke total_paid dan menambah payment_count
	query := `
		INSERT INTO reporting_finance_monthly
			(school_id, academic_period_id, month, year, component_type, total_paid, payment_count, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, 1, NOW())
		ON CONFLICT (school_id, academic_period_id, month, year, component_type) DO UPDATE SET
			total_paid = reporting_finance_monthly.total_paid + EXCLUDED.total_paid,
			payment_count = reporting_finance_monthly.payment_count + 1,
			total_outstanding = GREATEST(0, reporting_finance_monthly.total_outstanding - EXCLUDED.total_paid),
			updated_at = NOW()
	`

	_, err := h.db.Exec(ctx, query,
		event.SchoolID,
		event.AcademicPeriodID,
		event.Month,
		event.Year,
		event.ComponentType,
		event.Amount,
	)

	if err != nil {
		return fmt.Errorf("gagal upsert reporting_finance_monthly (CorrelationID: %s): %w", event.CorrelationID, err)
	}

	log.Printf("[Projection] Finance aggregated for %s %d/%d (CorrelationID: %s)",
		event.ComponentType, event.Month, event.Year, event.CorrelationID)

	return nil
}

// Helper untuk konversi string ke int jika month/year datang sebagai string dari event
func (h *FinanceProjectionHandler) parseMonthYear(monthStr, yearStr string) (int, int, error) {
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month: %w", err)
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year: %w", err)
	}
	return month, year, nil
}
