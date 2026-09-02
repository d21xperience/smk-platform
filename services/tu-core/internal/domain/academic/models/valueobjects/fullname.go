package valueobjects

import (
	"errors"
	"strings"
)

// FullName adalah value object untuk nama lengkap.
type FullName struct {
	firstName  string
	middleName string
	lastName   string
}

// NewFullName membuat FullName baru dengan validasi.
func NewFullName(firstName, middleName, lastName string) (FullName, error) {
	if strings.TrimSpace(firstName) == "" {
		return FullName{}, errors.New("nama depan tidak boleh kosong")
	}
	return FullName{
		firstName:  strings.TrimSpace(firstName),
		middleName: strings.TrimSpace(middleName),
		lastName:   strings.TrimSpace(lastName),
	}, nil
}

// Getters

func (n FullName) FirstName() string  { return n.firstName }
func (n FullName) MiddleName() string { return n.middleName }
func (n FullName) LastName() string   { return n.lastName }

// DisplayName mengembalikan nama lengkap untuk ditampilkan.
func (n FullName) DisplayName() string {
	parts := []string{n.firstName}
	if n.middleName != "" {
		parts = append(parts, n.middleName)
	}
	if n.lastName != "" {
		parts = append(parts, n.lastName)
	}
	return strings.Join(parts, " ")
}

// Equals mengecek kesamaan.
func (n FullName) Equals(other FullName) bool {
	return n.firstName == other.firstName &&
		n.middleName == other.middleName &&
		n.lastName == other.lastName
}
