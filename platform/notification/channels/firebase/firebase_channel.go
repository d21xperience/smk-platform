package firebase

import (
	"context"
	"fmt"

	"sekolah-platform/platform/notification/contract"

	"firebase.google.com/go/messaging"
)

// FirebaseChannel mengimplementasikan contract.Notifier untuk FCM.
type FirebaseChannel struct {
	adapter *FCMAdapter
}

// NewFirebaseChannel membuat channel baru.
func NewFirebaseChannel(client *messaging.Client) *FirebaseChannel {
	return &FirebaseChannel{
		adapter: NewFCMAdapter(client),
	}
}

// Send mengirimkan push notification.
func (c *FirebaseChannel) Send(ctx context.Context, req contract.NotificationRequest) error {
	// Konversi map[string]interface{} ke map[string]string (syarat FCM)
	data := make(map[string]string)
	for key, value := range req.Data {
		data[key] = fmt.Sprintf("%v", value)
	}

	// Tambahkan metadata tracing
	data["correlation_id"] = req.CorrelationID
	data["school_id"] = req.SchoolID

	// Kirim via adapter
	return c.adapter.Send(ctx, req.Recipient, data)
}
