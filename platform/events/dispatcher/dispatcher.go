// Package dispatcher menyediakan event dispatcher interface dan implementasi.
package dispatcher

import "sekolah-platform/platform/events/types"

// Dispatcher adalah interface untuk mendistribusikan domain events.
type Dispatcher interface {
	// Dispatch mendistribusikan event ke semua listener yang terdaftar.
	Dispatch(event types.DomainEvent) error

	// DispatchAsync mendistribusikan event secara asynchronous.
	DispatchAsync(event types.DomainEvent) error

	// DispatchAll mendistribusikan beberapa events sekaligus.
	DispatchAll(events []types.DomainEvent) error
}
