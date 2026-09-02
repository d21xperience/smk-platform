package email

import (
	"context"
	"fmt"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/templates"
)

// Config konfigurasi SMTP.
type Config struct {
	Host      string
	Port      string
	Username  string
	Password  string
	FromName  string
	FromEmail string
}

// EmailChannel mengimplementasikan contract.Notifier untuk Email.
type EmailChannel struct {
	adapter   *SMTPAdapter
	templater *templates.TemplateEngine
}

// NewEmailChannel membuat channel baru.
func NewEmailChannel(cfg Config, templater *templates.TemplateEngine) *EmailChannel {
	return &EmailChannel{
		adapter:   NewSMTPAdapter(cfg),
		templater: templater,
	}
}

// Send mengirimkan email.
func (c *EmailChannel) Send(ctx context.Context, req contract.NotificationRequest) error {
	// Render template (sekarang return subject dan body)
	subject, body, err := c.templater.Render(req.TemplateCode, req.Data)
	if err != nil {
		return fmt.Errorf("render template gagal: %w", err)
	}

	// Jika subject kosong, gunakan default
	if subject == "" {
		subject = "Pemberitahuan Sekolah Digital Platform"
	}

	// Kirim via adapter
	return c.adapter.Send(req.Recipient, subject, body)
}
