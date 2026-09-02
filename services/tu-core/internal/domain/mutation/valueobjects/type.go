package valueobjects

import "errors"

type MutationType string

const (
	TypeMasuk  MutationType = "MASUK"
	TypeKeluar MutationType = "KELUAR"
)

func NewMutationType(t string) (MutationType, error) {
	switch MutationType(t) {
	case TypeMasuk, TypeKeluar:
		return MutationType(t), nil
	default:
		return "", errors.New("invalid mutation type: must be MASUK or KELUAR")
	}
}

func (t MutationType) String() string {
	return string(t)
}
