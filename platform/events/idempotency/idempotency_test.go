package idempotency

import (
	"testing"
	"time"
)

func TestMarkFirstTime(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	marked, err := c.Mark("event-1", time.Hour)
	if err != nil {
		t.Fatalf("Mark failed: %v", err)
	}
	if !marked {
		t.Error("Expected mark to succeed first time")
	}
}

func TestMarkDuplicate(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	_, _ = c.Mark("event-1", time.Hour)
	marked, err := c.Mark("event-1", time.Hour)
	if err != nil {
		t.Fatalf("Mark failed: %v", err)
	}
	if marked {
		t.Error("Expected mark to fail for duplicate")
	}
}

func TestIsProcessed(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	processed, _ := c.IsProcessed("event-1")
	if processed {
		t.Error("Expected not processed initially")
	}

	_, _ = c.Mark("event-1", time.Hour)
	processed, _ = c.IsProcessed("event-1")
	if !processed {
		t.Error("Expected processed after mark")
	}
}

func TestTTLExpiration(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	// Mark dengan TTL sangat pendek
	_, _ = c.Mark("event-1", 10*time.Millisecond)

	// Tunggu sampai expired
	time.Sleep(20 * time.Millisecond)

	processed, _ := c.IsProcessed("event-1")
	if processed {
		t.Error("Expected not processed after TTL expiration")
	}

	// Mark lagi harus berhasil
	marked, _ := c.Mark("event-1", time.Hour)
	if !marked {
		t.Error("Expected mark to succeed after expiration")
	}
}

func TestClear(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	_, _ = c.Mark("event-1", time.Hour)
	_ = c.Clear("event-1")

	processed, _ := c.IsProcessed("event-1")
	if processed {
		t.Error("Expected not processed after clear")
	}
}

func TestClearAll(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	_, _ = c.Mark("event-1", time.Hour)
	_, _ = c.Mark("event-2", time.Hour)
	_ = c.ClearAll()

	if c.Size() != 0 {
		t.Errorf("Expected size 0 after ClearAll, got %d", c.Size())
	}
}

func TestGenerateKey(t *testing.T) {
	key := GenerateKey("StudentCreated", "evt-001")
	expected := "StudentCreated:evt-001"
	if key != expected {
		t.Errorf("Expected key '%s', got '%s'", expected, key)
	}
}

func TestConcurrentMarks(t *testing.T) {
	c := NewMemoryChecker()
	defer c.ClearAll()

	const goroutines = 100
	results := make(chan bool, goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			marked, _ := c.Mark("same-event", time.Hour)
			results <- marked
		}()
	}

	successCount := 0
	for i := 0; i < goroutines; i++ {
		if <-results {
			successCount++
		}
	}

	// Hanya satu yang seharusnya berhasil
	if successCount != 1 {
		t.Errorf("Expected exactly 1 successful mark, got %d", successCount)
	}
}
