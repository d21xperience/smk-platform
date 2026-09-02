package valueobjects

import "errors"

type CorrespondenceStatus string

const (
	StatusDraft     CorrespondenceStatus = "DRAFT"
	StatusProcessed CorrespondenceStatus = "PROCESSED"
	StatusArchived  CorrespondenceStatus = "ARCHIVED"
)

func NewCorrespondenceStatus(s string) (CorrespondenceStatus, error) {
	switch CorrespondenceStatus(s) {
	case StatusDraft, StatusProcessed, StatusArchived:
		return CorrespondenceStatus(s), nil
	default:
		return "", errors.New("invalid correspondence status")
	}
}

func (s CorrespondenceStatus) String() string {
	return string(s)
}

func (s CorrespondenceStatus) CanArchive() bool {
	return s == StatusDraft || s == StatusProcessed
}
