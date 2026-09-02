package firebase

import (
	"context"
	"fmt"

	"firebase.google.com/go/messaging"
)

// FCMAdapter mengimplementasikan pengiriman push notification via Firebase Cloud Messaging.
type FCMAdapter struct {
	client *messaging.Client
}

// NewFCMAdapter membuat adapter baru.
func NewFCMAdapter(client *messaging.Client) *FCMAdapter {
	return &FCMAdapter{client: client}
}

// Send mengirimkan push notification.
func (a *FCMAdapter) Send(ctx context.Context, token string, data map[string]string) error {
	msg := &messaging.Message{
		Token: token,
		Data:  data,
	}

	_, err := a.client.Send(ctx, msg)
	if err != nil {
		return fmt.Errorf("gagal kirim FCM: %w", err)
	}

	return nil
}
