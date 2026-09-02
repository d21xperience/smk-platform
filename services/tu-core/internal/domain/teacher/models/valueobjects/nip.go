package valueobjects

import (
	"errors"
	"regexp"
)

// NIP adalah value object untuk Nomor Induk Pegawai.
type NIP struct {
	value string
}

var nipPattern = regexp.MustCompile(`^\d{18}$`)

// NewNIP membuat NIP baru dengan validasi.
func NewNIP(value string) (NIP, error) {
	if value == "" {
		return NIP{}, errors.New("NIP tidak boleh kosong")
	}
	if !nipPattern.MatchString(value) {
		return NIP{}, errors.New("NIP harus 18 digit angka")
	}
	return NIP{value: value}, nil
}

// Value mengembalikan nilai NIP.
func (n NIP) Value() string { return n.value }

// String mengimplementasikan Stringer.
func (n NIP) String() string { return n.value }

// Equals mengecek kesamaan.
func (n NIP) Equals(other NIP) bool { return n.value == other.value }
