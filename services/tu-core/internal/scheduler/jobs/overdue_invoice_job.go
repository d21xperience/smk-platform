package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type OverdueInvoiceJob struct {
	db *pgxpool.Pool // Atau *pgxpool.Pool
	js nats.JetStreamContext
}

func NewOverdueInvoiceJob(db *pgxpool.Pool, js nats.JetStreamContext) *OverdueInvoiceJob {
	return &OverdueInvoiceJob{db: db, js: js}
}

// overdueInvoiceData digunakan untuk payload event
type overdueInvoiceData struct {
	InvoiceID        string  `json:"invoiceId"`
	SchoolID         string  `json:"schoolId"`
	SchoolName       string  `json:"schoolName"` // Bisa di-hardcode atau ambil dari config
	AcademicPeriodID string  `json:"academicPeriodId"`
	StudentID        string  `json:"studentId"`
	ComponentType    string  `json:"componentType"`
	Month            string  `json:"month"`
	Year             string  `json:"year"`
	Amount           float64 `json:"amount"`
	DueDate          string  `json:"dueDate"`
	CorrelationID    string  `json:"correlationId"`
}

func (j *OverdueInvoiceJob) Run(ctx context.Context) error {
	log.Println("[Scheduler] Memulai OverdueInvoiceJob...")
	start := time.Now()

	// 1. Update status dan RETURNING data yang dibutuhkan untuk event
	// Asumsi: tabel invoices memiliki kolom month, year, component_type
	query := `
		UPDATE invoices
		SET status = 'OVERDUE', updated_at = NOW()
		WHERE status = 'UNPAID' AND due_date < CURRENT_DATE
		RETURNING id, school_id, academic_period_id, student_id, component_type, month, year, amount, due_date
	`
	rows, err := j.db.Query(ctx, query)
	if err != nil {
		return fmt.Errorf("query update overdue gagal: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		var data overdueInvoiceData
		if err := rows.Scan(&data.InvoiceID, &data.SchoolID, &data.AcademicPeriodID,
			&data.StudentID, &data.ComponentType, &data.Month, &data.Year, &data.Amount, &data.DueDate); err != nil {
			log.Printf("[Scheduler] Gagal scan row: %v", err)
			continue
		}

		count++
		data.SchoolName = "SDP Sekolah" // Sebaiknya ambil dari school config
		data.CorrelationID = uuid.New().String()

		// 2. Publish Event ke NATS (Prinsip 7: Event-Driven)
		payload, _ := json.Marshal(data)
		msg := nats.Msg{
			Subject: "sdp.event.invoice.overdue",
			Data:    payload,
			Header: nats.Header{
				"X-Correlation-ID": []string{data.CorrelationID},
				"X-Event-Type":     []string{"InvoiceOverdue"},
			},
		}

		if _, err := j.js.PublishMsg(&msg); err != nil {
			log.Printf("[Scheduler] Gagal publish event InvoiceOverdue untuk invoice %s: %v", data.InvoiceID, err)
			// Catatan: Karena status sudah di-update di DB, kegagalan publish akan ditangani oleh
			// mekanisme retry atau manual reconciliation, bukan rollback transaksi ini.
		} else {
			log.Printf("[Scheduler] Invoice %s (Student: %s, Bulan: %s) ditandai OVERDUE dan event dipublish",
				data.InvoiceID, data.StudentID, data.Month)
		}
	}

	log.Printf("[Scheduler] OverdueInvoiceJob selesai. %d invoice diperbarui. Durasi: %v", count, time.Since(start))
	return nil
}
