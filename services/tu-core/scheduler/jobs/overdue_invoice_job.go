package jobs

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

// OverdueInvoiceJob memeriksa dan memperbarui status tagihan yang lewat jatuh tempo.
type OverdueInvoiceJob struct {
	db *pgxpool.Pool // UBAH DARI *pgx.Conn MENJADI *pgxpool.Pool
	js nats.JetStreamContext
}

// NewOverdueInvoiceJob menerima *pgxpool.Pool
func NewOverdueInvoiceJob(db *pgxpool.Pool, js nats.JetStreamContext) *OverdueInvoiceJob {
	return &OverdueInvoiceJob{
		db: db,
		js: js,
	}
}

func (j *OverdueInvoiceJob) Run(ctx context.Context) error {
	log.Println("[Scheduler] Menjalankan OverdueInvoiceJob...")

	// ... (Logika bisnis query update invoice Anda di sini) ...
	// Pastikan semua query menggunakan j.db (yang sekarang bertipe *pgxpool.Pool)

	log.Println("[Scheduler] OverdueInvoiceJob selesai.")
	return nil
}
