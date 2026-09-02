package valueobjects

import "errors"

type MutationStatus string

const (
	StatusPending   MutationStatus = "PENDING"
	StatusDisetujui MutationStatus = "DISETUJUI"
	StatusDitolak   MutationStatus = "DITOLAK"
)

func NewMutationStatus(s string) (MutationStatus, error) {
	switch MutationStatus(s) {
	case StatusPending, StatusDisetujui, StatusDitolak:
		return MutationStatus(s), nil
	default:
		return "", errors.New("invalid mutation status")
	}
}

func (s MutationStatus) String() string {
	return string(s)
}

func (s MutationStatus) IsPending() bool {
	return s == StatusPending
}
