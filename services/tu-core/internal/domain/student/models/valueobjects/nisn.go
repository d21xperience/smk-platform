package valueobjects

import (
	"regexp"
)

type NISN struct {
	value string
}

var nisnPattern = regexp.MustCompile(`^\d{10}$`)

// NewNISN membuat NISN baru.
// DILONGGARKAN: DAPODIK mengizinkan NISN kosong (NULL) atau format lama.
func NewNISN(value string) (NISN, error) {
	if value == "" {
		return NISN{value: ""}, nil // Izinkan kosong
	}
	if !nisnPattern.MatchString(value) {
		// Jika tidak 10 digit, tetap izinkan (mungkin data legacy DAPODIK)
		return NISN{value: value}, nil 
	}
	return NISN{value: value}, nil
}

func (n NISN) Value() string { return n.value }
func (n NISN) String() string { return n.value }
func (n NISN) Equals(other NISN) bool { return n.value == other.value }
