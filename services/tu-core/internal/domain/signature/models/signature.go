package models

import (
	"time"
)

type DigitalSignature struct {
	ID               string
	CorrespondenceID string
	SignerUserID     string
	SignerRole       string
	ImageURL         string
	VerificationCode string
	SignedAt         time.Time
}
