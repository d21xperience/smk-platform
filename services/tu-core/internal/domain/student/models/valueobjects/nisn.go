// Package valueobjects menyediakan value objects untuk domain Student.
// SYMMETRIC dengan frontend (apps/siakad-tu/src/domain/student/models/value-objects/).
package valueobjects

import (
	"errors"
	"fmt"
	"regexp"
)

// NISN adalah value object untuk Nomor Induk Siswa Nasional.
type NISN struct {
	value string
}

var nisnPattern = regexp.MustCompile(`^\d{10}$`)

// NewNISN membuat NISN baru dengan validasi.
func NewNISN(value string) (NISN, error) {
	if value == "" {
		return NISN{}, errors.New("NISN tidak boleh kosong")
	}
	if !nisnPattern.MatchString(value) {
		return NISN{}, fmt.Errorf("NISN harus 10 digit angka, got: %s", value)
	}
	return NISN{value: value}, nil
}

// Value mengembalikan nilai NISN.
func (n NISN) Value() string { return n.value }

// String mengimplementasikan Stringer.
func (n NISN) String() string { return n.value }

// Equals mengecek kesamaan.
func (n NISN) Equals(other NISN) bool { return n.value == other.value }
