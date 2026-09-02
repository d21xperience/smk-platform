package scheduler

import (
	"context"
	"log"

	"github.com/robfig/cron/v3"
)

type Job interface {
	Run(ctx context.Context) error
}

type Scheduler struct {
	cron   *cron.Cron
	ctx    context.Context
	cancel context.CancelFunc
}

func NewScheduler() *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())

	// Chain: Recover dari panic agar satu job yang error tidak mematikan seluruh scheduler
	c := cron.New(cron.WithChain(
		cron.Recover(cron.DefaultLogger),
	))

	return &Scheduler{
		cron:   c,
		ctx:    ctx,
		cancel: cancel,
	}
}

func (s *Scheduler) RegisterJob(spec string, job Job) {
	s.cron.AddFunc(spec, func() {
		// Gunakan context dari scheduler agar bisa di-cancel saat shutdown
		if err := job.Run(s.ctx); err != nil {
			log.Printf("[Scheduler] ERROR pada job: %v", err)
		}
	})
}

func (s *Scheduler) Start() {
	log.Println("[Scheduler] Cron scheduler dimulai.")
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	log.Println("[Scheduler] Menghentikan cron scheduler (graceful)...")
	s.cancel()           // Cancel semua context yang sedang berjalan di job
	ctx := s.cron.Stop() // Stop menerima trigger baru, tunggu job yang sedang berjalan selesai
	<-ctx.Done()
	log.Println("[Scheduler] Cron scheduler berhasil dihentikan.")
}
