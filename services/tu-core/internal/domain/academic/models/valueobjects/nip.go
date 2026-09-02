package valueobjects

import "errors"

// NIS adalah value object untuk Nomor Induk Siswa.
type NIS struct {
	value string
}

// NewNIS membuat NIS baru dengan validasi.
func NewNIS(value string) (NIS, error) {
	if value == "" {
		return NIS{}, errors.New("NIS tidak boleh kosong")
	}
	return NIS{value: value}, nil
}

// Value mengembalikan nilai NIS.
func (n NIS) Value() string { return n.value }

// String mengimplementasikan Stringer.
func (n NIS) String() string { return n.value }

// Equals mengecek kesamaan.
func (n NIS) Equals(other NIS) bool { return n.value == other.value }
