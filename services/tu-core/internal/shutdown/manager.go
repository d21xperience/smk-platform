package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// ShutdownFunc adalah fungsi yang dipanggil saat shutdown.
type ShutdownFunc func(ctx context.Context) error

// Manager mengelola graceful shutdown dengan timeout dan urutan yang benar.
type Manager struct {
	handlers []ShutdownFunc
	timeout  time.Duration
}

func NewManager(timeout time.Duration) *Manager {
	return &Manager{
		handlers: make([]ShutdownFunc, 0),
		timeout:  timeout,
	}
}

// Register menambahkan handler shutdown. Handler dipanggil dalam urutan LIFO (Last In, First Out).
func (m *Manager) Register(handler ShutdownFunc) {
	m.handlers = append(m.handlers, handler)
}

// WaitForShutdown memblokir hingga menerima sinyal SIGINT/SIGTERM, lalu menjalankan semua handler.
func (m *Manager) WaitForShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	log.Printf("[Shutdown] Menerima sinyal %v, memulai graceful shutdown...", sig)

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, len(m.handlers))

	// Panggil handler dalam urutan terbalik (LIFO)
	for i := len(m.handlers) - 1; i >= 0; i-- {
		wg.Add(1)
		go func(handler ShutdownFunc) {
			defer wg.Done()
			if err := handler(ctx); err != nil {
				errChan <- err
			}
		}(m.handlers[i])
	}

	wg.Wait()
	close(errChan)

	// Log semua error
	for err := range errChan {
		log.Printf("[Shutdown] ERROR: %v", err)
	}

	log.Println("[Shutdown] Graceful shutdown selesai.")
}
