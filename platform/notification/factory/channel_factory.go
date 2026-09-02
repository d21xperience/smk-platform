package factory

import (
	"time"

	"sekolah-platform/platform/notification/channels/email"
	"sekolah-platform/platform/notification/channels/firebase"
	"sekolah-platform/platform/notification/channels/whatsapp"
	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/templates"

	"firebase.google.com/go/messaging"
)

// ChannelFactory membuat channel notifikasi dari konfigurasi.
type ChannelFactory struct {
	templater *templates.TemplateEngine
}

// NewChannelFactory membuat factory baru.
func NewChannelFactory(templater *templates.TemplateEngine) *ChannelFactory {
	return &ChannelFactory{templater: templater}
}

// CreateEmailChannel membuat EmailChannel dari config map.
func (f *ChannelFactory) CreateEmailChannel(config map[string]string) contract.Notifier {
	cfg := email.Config{
		Host:      config["host"],
		Port:      config["port"],
		Username:  config["username"],
		Password:  config["password"],
		FromName:  config["from_name"],
		FromEmail: config["from_email"],
	}
	return email.NewEmailChannel(cfg, f.templater)
}

// CreateWhatsAppChannel membuat WhatsAppChannel dari config map.
func (f *ChannelFactory) CreateWhatsAppChannel(config map[string]string) contract.Notifier {
	timeout, _ := time.ParseDuration(config["timeout"])
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	cfg := whatsapp.Config{
		GatewayURL: config["gateway_url"],
		APIKey:     config["api_key"],
		Timeout:    timeout,
	}
	return whatsapp.NewWhatsAppChannel(cfg, f.templater)
}

// CreateFirebaseChannel membuat FirebaseChannel dari FCM client.
func (f *ChannelFactory) CreateFirebaseChannel(client *messaging.Client) contract.Notifier {
	return firebase.NewFirebaseChannel(client)
}
