package valueobjects

import (
	"errors"
	"regexp"
)

// NUPTK adalah value object untuk Nomor Unik Pendidik dan Tenaga Kependidikan.
type NUPTK struct {
	value string
}

var nuptkPattern = regexp.MustCompile(`^\d{16}$`)

// NewNUPTK membuat NUPTK baru dengan validasi.
// NUPTK bisa kosong (optional untuk guru honorer).
func NewNUPTK(value string) (NUPTK, error) {
	if value == "" {
		return NUPTK{}, nil // Optional
	}
	if !nuptkPattern.MatchString(value) {
		return NUPTK{}, errors.New("NUPTK harus 16 digit angka")
	}
	return NUPTK{value: value}, nil
}

// Value mengembalikan nilai NUPTK.
func (n NUPTK) Value() string { return n.value }

// IsEmpty mengecek apakah NUPTK kosong.
func (n NUPTK) IsEmpty() bool { return n.value == "" }

// String mengimplementasikan Stringer.
func (n NUPTK) String() string { return n.value }
