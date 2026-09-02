package valueobjects

import "errors"

type RegistrationStatus string

const (
	StatusDraft     RegistrationStatus = "DRAFT"
	StatusSubmitted RegistrationStatus = "SUBMITTED"
	StatusVerified  RegistrationStatus = "VERIFIED"
	StatusApproved  RegistrationStatus = "APPROVED"
	StatusRejected  RegistrationStatus = "REJECTED"
)

func NewRegistrationStatus(s string) (RegistrationStatus, error) {
	switch RegistrationStatus(s) {
	case StatusDraft, StatusSubmitted, StatusVerified, StatusApproved, StatusRejected:
		return RegistrationStatus(s), nil
	default:
		return "", errors.New("invalid registration status")
	}
}

func (s RegistrationStatus) String() string {
	return string(s)
}

func (s RegistrationStatus) IsDraft() bool {
	return s == StatusDraft
}

func (s RegistrationStatus) IsSubmitted() bool {
	return s == StatusSubmitted
}

func (s RegistrationStatus) IsVerified() bool {
	return s == StatusVerified
}

func (s RegistrationStatus) CanSubmit() bool {
	return s == StatusDraft
}

func (s RegistrationStatus) CanVerify() bool {
	return s == StatusSubmitted
}

func (s RegistrationStatus) CanApprove() bool {
	return s == StatusVerified
}

func (s RegistrationStatus) CanReject() bool {
	return s == StatusSubmitted || s == StatusVerified
}
