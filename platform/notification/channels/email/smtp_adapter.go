package email

import (
	"fmt"
	"net/smtp"
)

// SMTPAdapter mengimplementasikan pengiriman email via SMTP.
type SMTPAdapter struct {
	config Config
}

// NewSMTPAdapter membuat adapter baru.
func NewSMTPAdapter(config Config) *SMTPAdapter {
	return &SMTPAdapter{config: config}
}

// Send mengirimkan email via SMTP.
func (a *SMTPAdapter) Send(to, subject, htmlBody string) error {
	// Susun MIME Message
	msg := fmt.Appendf(nil, "To: %s\r\n"+
		"Subject: %s\r\n"+
		"MIME-Version: 1.0\r\n"+
		"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, htmlBody)

	// Autentikasi
	auth := smtp.PlainAuth("", a.config.Username, a.config.Password, a.config.Host)
	addr := fmt.Sprintf("%s:%s", a.config.Host, a.config.Port)

	// Kirim
	return smtp.SendMail(addr, auth, a.config.FromEmail, []string{to}, msg)
}
