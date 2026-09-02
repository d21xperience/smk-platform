package valueobjects

import "errors"

type CorrespondenceType string

const (
	TypeIncoming CorrespondenceType = "INCOMING"
	TypeOutgoing CorrespondenceType = "OUTGOING"
)

func NewCorrespondenceType(t string) (CorrespondenceType, error) {
	switch CorrespondenceType(t) {
	case TypeIncoming, TypeOutgoing:
		return CorrespondenceType(t), nil
	default:
		return "", errors.New("invalid correspondence type: must be INCOMING or OUTGOING")
	}
}

func (t CorrespondenceType) String() string {
	return string(t)
}
