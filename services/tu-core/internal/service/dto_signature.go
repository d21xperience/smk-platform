package service

import (
	"time"
)

type SignCorrespondenceRequest struct {
	SignerRole  string `json:"signer_role"`
	ImageBase64 string `json:"image_base64"` // Format: "data:image/png;base64,...."
}

type SignatureResponse struct {
	ID               string    `json:"id"`
	CorrespondenceID string    `json:"correspondence_id"`
	SignerUserID     string    `json:"signer_user_id"`
	SignerRole       string    `json:"signer_role"`
	ImageURL         string    `json:"image_url"`
	VerificationCode string    `json:"verification_code"`
	SignedAt         time.Time `json:"signed_at"`
}
