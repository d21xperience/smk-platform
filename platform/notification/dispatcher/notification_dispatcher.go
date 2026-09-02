package dispatcher

import (
	"context"
	"log"

	"sekolah-platform/platform/notification/contract"
	notifErrors "sekolah-platform/platform/notification/errors"
	logNotif "sekolah-platform/platform/notification/log"
	"sekolah-platform/platform/notification/validator"
)

// NotificationDispatcher merutekan permintaan notifikasi ke channel yang tepat.
type NotificationDispatcher struct {
	channels  map[string]contract.Notifier
	validator *validator.RequestValidator
	logRepo   logNotif.NotificationLogRepository
}

// NewNotificationDispatcher membuat dispatcher baru.
func NewNotificationDispatcher(logRepo logNotif.NotificationLogRepository) *NotificationDispatcher {
	return &NotificationDispatcher{
		channels:  make(map[string]contract.Notifier),
		validator: validator.NewRequestValidator(),
		logRepo:   logRepo,
	}
}

// Register menambahkan channel ke dispatcher.
func (d *NotificationDispatcher) Register(channel string, notifier contract.Notifier) {
	d.channels[channel] = notifier
}

// Dispatch mengirim notifikasi ke channel yang sesuai.
func (d *NotificationDispatcher) Dispatch(ctx context.Context, req contract.NotificationRequest) error {
	// Validasi request
	if err := d.validator.Validate(req); err != nil {
		return err
	}

	// Cek channel terdaftar
	notifier, exists := d.channels[req.Channel]
	if !exists {
		return notifErrors.NewNotificationError("Dispatch", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
			notifErrors.ErrChannelNotFound)
	}

	// Buat log entry
	logEntry := &logNotif.NotificationLog{
		SchoolID:         req.SchoolID,
		AcademicPeriodID: req.AcademicPeriodID,
		CorrelationID:    req.CorrelationID,
		Channel:          req.Channel,
		Recipient:        req.Recipient,
		TemplateCode:     req.TemplateCode,
		Status:           "PENDING",
	}

	if d.logRepo != nil {
		if err := d.logRepo.Save(logEntry); err != nil {
			log.Printf("[Dispatcher] Gagal save log: %v", err)
		}
	}

	// Kirim notifikasi
	err := notifier.Send(ctx, req)

	// Update log status
	if d.logRepo != nil && logEntry.ID != "" {
		status := "SENT"
		errorMsg := ""
		if err != nil {
			status = "FAILED"
			errorMsg = err.Error()
		}
		if updateErr := d.logRepo.UpdateStatus(logEntry.ID, status, errorMsg); updateErr != nil {
			log.Printf("[Dispatcher] Gagal update log status: %v", updateErr)
		}
	}

	if err != nil {
		return notifErrors.NewNotificationError("Dispatch", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID, err)
	}

	return nil
}

// GetRegisteredChannels mengembalikan daftar channel yang terdaftar.
func (d *NotificationDispatcher) GetRegisteredChannels() []string {
	channels := make([]string, 0, len(d.channels))
	for ch := range d.channels {
		channels = append(channels, ch)
	}
	return channels
}
