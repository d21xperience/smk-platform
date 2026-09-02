package outbox

import (
	"context"

	"sekolah-platform/platform/events/types"
)

// TransactionalOperation adalah operasi bisnis yang dijalankan dalam transaksi.
// Return value adalah event yang akan ditulis ke outbox.
type TransactionalOperation func(ctx context.Context) (types.DomainEvent, error)

// Writer menulis event ke outbox dalam transaksi yang sama dengan business data.
type Writer struct {
	outbox Outbox
}

// NewWriter membuat Writer baru.
func NewWriter(outbox Outbox) *Writer {
	return &Writer{outbox: outbox}
}

// WriteInTransaction menjalankan operasi bisnis dan menulis event ke outbox.
// Dalam implementasi nyata, operasi bisnis dan outbox save harus dalam SATU transaksi DB.
// Implementasi ini menyederhanakan dengan memanggil operasi lalu save ke outbox.
// Service layer yang sebenarnya harus wrap dengan DB transaction.
func (w *Writer) WriteInTransaction(ctx context.Context, op TransactionalOperation) error {
	event, err := op(ctx)
	if err != nil {
		return err
	}

	if event == nil {
		return nil // Tidak ada event yang dihasilkan
	}

	entry := OutboxEntry{
		ID:        generateEntryID(),
		Event:     event,
		CreatedAt: timeNow(),
		Attempts:  0,
	}

	return w.outbox.Save(ctx, entry)
}

// WriteEvent langsung menulis event ke outbox (tanpa operasi bisnis).
// Gunakan ini untuk testing atau saat event sudah dihasilkan di luar transaksi.
func (w *Writer) WriteEvent(ctx context.Context, event types.DomainEvent) error {
	if event == nil {
		return nil
	}

	entry := OutboxEntry{
		ID:        generateEntryID(),
		Event:     event,
		CreatedAt: timeNow(),
		Attempts:  0,
	}

	return w.outbox.Save(ctx, entry)
}
