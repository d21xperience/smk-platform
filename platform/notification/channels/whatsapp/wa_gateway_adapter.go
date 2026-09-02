package whatsapp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// WAGatewayAdapter mengimplementasikan pengiriman WhatsApp via HTTP Gateway.
type WAGatewayAdapter struct {
	client     *http.Client
	gatewayURL string
	apiKey     string
}

// NewWAGatewayAdapter membuat adapter baru.
func NewWAGatewayAdapter(client *http.Client, gatewayURL, apiKey string) *WAGatewayAdapter {
	return &WAGatewayAdapter{
		client:     client,
		gatewayURL: gatewayURL,
		apiKey:     apiKey,
	}
}

// Send mengirimkan pesan WhatsApp.
func (a *WAGatewayAdapter) Send(ctx context.Context, phone, message string) error {
	payload := map[string]interface{}{
		"target":  phone,
		"message": message,
	}

	payloadBytes, _ := json.Marshal(payload)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", a.gatewayURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return fmt.Errorf("gagal buat request: %w", err)
	}

	httpReq.Header.Set("Authorization", a.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("gagal kirim (network): %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("gagal kirim (status %d)", resp.StatusCode)
	}

	return nil
}
