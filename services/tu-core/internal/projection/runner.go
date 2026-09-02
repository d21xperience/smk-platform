package projection

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// RunnerConfig konfigurasi untuk projection runner.
type RunnerConfig struct {
	MaxRetries    int
	RetryDelay    time.Duration
	EnableLogging bool
}

// DefaultRunnerConfig memberikan konfigurasi default.
func DefaultRunnerConfig() RunnerConfig {
	return RunnerConfig{
		MaxRetries:    3,
		RetryDelay:    1 * time.Second,
		EnableLogging: true,
	}
}

// Runner mengeksekusi projection handlers.
type Runner struct {
	manager *ProjectionManager
	config  RunnerConfig
}

// NewRunner membuat instance baru dari Runner.
func NewRunner(manager *ProjectionManager, config RunnerConfig) *Runner {
	return &Runner{
		manager: manager,
		config:  config,
	}
}

// Run memproses raw event message (dari NATS) dan mendelegasikan ke manager.
func (r *Runner) Run(ctx context.Context, rawMessage []byte) error {
	// Parse raw message ke DomainEvent
	var event DomainEvent
	if err := json.Unmarshal(rawMessage, &event); err != nil {
		return fmt.Errorf("%w: gagal parse event: %v", ErrInvalidEventPayload, err)
	}

	// Validasi operational context (Prinsip 5)
	if event.SchoolID == "" || event.AcademicPeriodID == "" {
		return fmt.Errorf("%w: schoolId dan academicPeriodId wajib diisi", ErrInvalidEventPayload)
	}

	if r.config.EnableLogging {
		log.Printf("[ProjectionRunner] Memproses event '%s' untuk aggregate '%s' (CorrelationID: %s)",
			event.EventType, event.AggregateID, event.CorrelationID)
	}

	// Eksekusi dengan retry
	var lastErr error
	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			if r.config.EnableLogging {
				log.Printf("[ProjectionRunner] Retry attempt %d/%d untuk event '%s'",
					attempt, r.config.MaxRetries, event.EventType)
			}
			time.Sleep(r.config.RetryDelay * time.Duration(attempt))
		}

		err := r.manager.HandleEvent(ctx, event)
		if err == nil {
			if r.config.EnableLogging {
				log.Printf("[ProjectionRunner] Sukses proses event '%s' (CorrelationID: %s)",
					event.EventType, event.CorrelationID)
			}
			return nil
		}

		lastErr = err
		log.Printf("[ProjectionRunner] Attempt %d gagal untuk event '%s': %v",
			attempt+1, event.EventType, err)
	}

	return fmt.Errorf("%w: semua retry gagal untuk event '%s': %v",
		ErrProjectionFailed, event.EventType, lastErr)
}
