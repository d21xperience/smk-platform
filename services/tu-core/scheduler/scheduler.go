package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

// Job adalah interface untuk semua background job.
type Job interface {
	Run(ctx context.Context) error
}

// SDPScheduler membungkus cron scheduler.
type SDPScheduler struct {
	cron *cron.Cron
}

// NewSDPScheduler membuat instance scheduler baru.
func NewSDPScheduler() *SDPScheduler {
	return &SDPScheduler{
		cron: cron.New(cron.WithChain(
			cron.Recover(cron.DefaultLogger),
		)),
	}
}

// RegisterJob mendaftarkan job dengan cron expression (misal: "0 2 * * *").
func (s *SDPScheduler) RegisterJob(spec string, job Job) {
	s.cron.AddFunc(spec, func() {
		if err := job.Run(context.Background()); err != nil {
			log.Printf("[Scheduler] Job error: %v", err)
		}
	})
}

// Start memulai scheduler.
func (s *SDPScheduler) Start() {
	log.Println("✅ Scheduler dimulai")
	s.cron.Start()
}

// Stop menghentikan scheduler secara graceful.
func (s *SDPScheduler) Stop() {
	log.Println("🛑 Scheduler dihentikan")
	ctx := s.cron.Stop()
	<-ctx.Done()
}
