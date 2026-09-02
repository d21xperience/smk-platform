package whatsapp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/templates"
)

// Config konfigurasi WhatsApp Gateway.
type Config struct {
	GatewayURL string
	APIKey     string
	Timeout    time.Duration
}

// WhatsAppChannel mengimplementasikan contract.Notifier untuk WhatsApp.
type WhatsAppChannel struct {
	adapter   *WAGatewayAdapter
	templater *templates.TemplateEngine
}

// NewWhatsAppChannel membuat channel baru.
func NewWhatsAppChannel(cfg Config, templater *templates.TemplateEngine) *WhatsAppChannel {
	client := &http.Client{Timeout: cfg.Timeout}
	return &WhatsAppChannel{
		adapter:   NewWAGatewayAdapter(client, cfg.GatewayURL, cfg.APIKey),
		templater: templater,
	}
}

// Send mengirimkan pesan WhatsApp.
func (c *WhatsAppChannel) Send(ctx context.Context, req contract.NotificationRequest) error {
	// Render template (hanya gunakan body untuk WhatsApp)
	_, message, err := c.templater.Render(req.TemplateCode, req.Data)
	if err != nil {
		return fmt.Errorf("render template gagal: %w", err)
	}

	// Kirim via adapter
	return c.adapter.Send(ctx, req.Recipient, message)
}
