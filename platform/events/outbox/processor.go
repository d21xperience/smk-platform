package outbox

import (
	"context"
	"log"
	"sync"
	"time"

	"sekolah-platform/platform/events/dlq"
	"sekolah-platform/platform/events/publisher"
	"sekolah-platform/platform/events/retry"
)

// ProcessorConfig adalah konfigurasi untuk Processor.
type ProcessorConfig struct {
	// PollInterval adalah interval polling outbox.
	PollInterval time.Duration

	// BatchSize adalah jumlah entry yang di-process per batch.
	BatchSize int

	// RetryPolicy adalah retry policy untuk publish yang gagal.
	RetryPolicy retry.Policy

	// DLQ adalah dead letter queue untuk event yang gagal setelah retry.
	DLQ dlq.DLQ

	// Logger untuk logging (optional).
	Logger Logger
}

// Logger adalah interface untuk logging.
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// defaultLogger adalah logger default yang menggunakan standard log.
type defaultLogger struct{}

func (l *defaultLogger) Info(msg string, args ...interface{})  { log.Printf("[INFO] "+msg, args...) }
func (l *defaultLogger) Error(msg string, args ...interface{}) { log.Printf("[ERROR] "+msg, args...) }

// DefaultProcessorConfig adalah konfigurasi default.
var DefaultProcessorConfig = ProcessorConfig{
	PollInterval: 1 * time.Second,
	BatchSize:    100,
	RetryPolicy:  retry.DefaultPolicy,
	Logger:       &defaultLogger{},
}

// Processor adalah background worker yang memproses outbox entries.
type Processor struct {
	config    ProcessorConfig
	outbox    Outbox
	publisher publisher.Publisher
	stopCh    chan struct{}
	wg        sync.WaitGroup
	running   bool
	mu        sync.Mutex
}

// NewProcessor membuat Processor baru.
func NewProcessor(outbox Outbox, pub publisher.Publisher, config ProcessorConfig) *Processor {
	if config.Logger == nil {
		config.Logger = &defaultLogger{}
	}
	if config.PollInterval == 0 {
		config.PollInterval = DefaultProcessorConfig.PollInterval
	}
	if config.BatchSize == 0 {
		config.BatchSize = DefaultProcessorConfig.BatchSize
	}

	return &Processor{
		config:    config,
		outbox:    outbox,
		publisher: pub,
		stopCh:    make(chan struct{}),
	}
}

// Start memulai processor sebagai background goroutine.
func (p *Processor) Start(ctx context.Context) {
	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return
	}
	p.running = true
	p.mu.Unlock()

	p.wg.Add(1)
	go p.run(ctx)
}

// Stop menghentikan processor dan menunggu semua goroutine selesai.
func (p *Processor) Stop() {
	p.mu.Lock()
	if !p.running {
		p.mu.Unlock()
		return
	}
	p.running = false
	p.mu.Unlock()

	close(p.stopCh)
	p.wg.Wait()
}

// run adalah main loop processor.
func (p *Processor) run(ctx context.Context) {
	defer p.wg.Done()

	ticker := time.NewTicker(p.config.PollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stopCh:
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

// processBatch memproses satu batch outbox entries.
func (p *Processor) processBatch(ctx context.Context) {
	entries, err := p.outbox.FetchBatch(ctx, p.config.BatchSize)
	if err != nil {
		p.config.Logger.Error("Failed to fetch outbox batch: %v", err)
		return
	}

	if len(entries) == 0 {
		return
	}

	for _, entry := range entries {
		p.processEntry(ctx, entry)
	}
}

// processEntry memproses satu outbox entry.
func (p *Processor) processEntry(ctx context.Context, entry OutboxEntry) {
	// Check apakah sudah waktunya retry
	if !entry.NextRetry.IsZero() && timeNow().Before(entry.NextRetry) {
		return
	}

	// Publish event
	err := p.publisher.Publish(ctx, entry.Event)
	if err == nil {
		// Sukses, hapus dari outbox
		if removeErr := p.outbox.Remove(ctx, entry.ID); removeErr != nil {
			p.config.Logger.Error("Failed to remove outbox entry %s: %v", entry.ID, removeErr)
		}
		return
	}

	// Gagal, update attempt
	entry.Attempts++
	p.config.Logger.Error("Failed to publish event %s (attempt %d): %v",
		entry.Event.GetEventID(), entry.Attempts, err)

	// Check apakah retry habis
	if entry.Attempts > p.config.RetryPolicy.MaxRetries {
		// Kirim ke DLQ
		if p.config.DLQ != nil {
			deadLetter := dlq.DeadLetter{
				OriginalEvent: entry.Event,
				ConsumerName:  "outbox-processor",
				FailedAt:      timeNow(),
				FailureCount:  entry.Attempts,
				LastError:     err.Error(),
			}
			if dlqErr := p.config.DLQ.Store(deadLetter); dlqErr != nil {
				p.config.Logger.Error("Failed to store to DLQ: %v", dlqErr)
			}
			// Hapus dari outbox (sudah masuk DLQ)
			_ = p.outbox.Remove(ctx, entry.ID)
		} else {
			// Tidak ada DLQ, biarkan di outbox untuk retry manual
			p.config.Logger.Error("No DLQ configured, entry %s remains in outbox", entry.ID)
		}
		return
	}

	// Hitung delay untuk retry berikutnya
	delay := p.config.RetryPolicy.DelayForAttempt(entry.Attempts - 1)
	nextRetry := timeNow().Add(delay)

	// Update attempt di outbox
	if updateErr := p.outbox.UpdateAttempt(ctx, entry.ID, entry.Attempts, err.Error(), nextRetry); updateErr != nil {
		p.config.Logger.Error("Failed to update outbox entry %s: %v", entry.ID, updateErr)
	}
}

// IsRunning mengecek apakah processor sedang berjalan.
func (p *Processor) IsRunning() bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.running
}
