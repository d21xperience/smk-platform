package outbox

import (
	"fmt"
	"sync/atomic"
	"time"
)

var (
	entryIDCounter int64
	timeNow        = time.Now
)

func generateEntryID() string {
	id := atomic.AddInt64(&entryIDCounter, 1)
	return fmt.Sprintf("outbox-%d", id)
}

// SetTimeNow mengganti timeNow function (untuk testing).
func SetTimeNow(fn func() time.Time) {
	timeNow = fn
}

// ResetTimeNow mengembalikan timeNow ke time.Now.
func ResetTimeNow() {
	timeNow = time.Now
}
