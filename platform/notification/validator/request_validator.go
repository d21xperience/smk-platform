package validator

import (
	"fmt"
	"regexp"
	"strings"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/errors"
)

// RequestValidator memvalidasi NotificationRequest sebelum dispatch.
type RequestValidator struct {
	emailRegex *regexp.Regexp
	phoneRegex *regexp.Regexp
}

// NewRequestValidator membuat validator baru.
func NewRequestValidator() *RequestValidator {
	return &RequestValidator{
		emailRegex: regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		phoneRegex: regexp.MustCompile(`^[0-9]{10,15}$`), // Format: 628xxx atau 08xxx
	}
}

// Validate memvalidasi NotificationRequest.
func (v *RequestValidator) Validate(req contract.NotificationRequest) error {
	// Validasi channel
	if req.Channel == "" {
		return errors.NewNotificationError("Validate", "", req.Recipient, req.TemplateCode, req.CorrelationID,
			fmt.Errorf("channel tidak boleh kosong"))
	}

	// Validasi recipient
	if req.Recipient == "" {
		return errors.NewNotificationError("Validate", req.Channel, "", req.TemplateCode, req.CorrelationID,
			errors.ErrInvalidRecipient)
	}

	// Validasi format recipient berdasarkan channel
	switch req.Channel {
	case "email":
		if !v.emailRegex.MatchString(req.Recipient) {
			return errors.NewNotificationError("Validate", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
				fmt.Errorf("format email tidak valid"))
		}
	case "whatsapp":
		// Normalize phone number (hapus +, spasi, dash)
		phone := strings.ReplaceAll(req.Recipient, "+", "")
		phone = strings.ReplaceAll(phone, " ", "")
		phone = strings.ReplaceAll(phone, "-", "")
		if !v.phoneRegex.MatchString(phone) {
			return errors.NewNotificationError("Validate", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
				fmt.Errorf("format nomor telepon tidak valid"))
		}
	case "firebase":
		// FCM token biasanya panjang (100+ chars)
		if len(req.Recipient) < 50 {
			return errors.NewNotificationError("Validate", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
				fmt.Errorf("FCM token tidak valid (terlalu pendek)"))
		}
	}

	// Validasi template code
	if req.TemplateCode == "" {
		return errors.NewNotificationError("Validate", req.Channel, req.Recipient, "", req.CorrelationID,
			errors.ErrInvalidTemplateCode)
	}

	// Validasi operational context (Prinsip 5)
	if req.SchoolID == "" {
		return errors.NewNotificationError("Validate", req.Channel, req.Recipient, req.TemplateCode, req.CorrelationID,
			fmt.Errorf("schoolId wajib diisi (Prinsip 5)"))
	}

	return nil
}
