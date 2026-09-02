// Package commands menyediakan command objects untuk StudentEngine.
package commands

import (
	"errors"
	"time"
)

// RegisterStudentCommand adalah command untuk mendaftarkan siswa baru.
type RegisterStudentCommand struct {
	StudentID          string
	NISN               string
	NIS                string
	FirstName          string
	MiddleName         string
	LastName           string
	BirthDate          time.Time
	Gender             string
	AddressStreet      string
	AddressRtRw        string
	AddressVillage     string
	AddressDistrict    string
	AddressCity        string
	AddressPostalCode  string
	ContactPhone       string
	ContactEmail       string
	GuardianName       string
	GuardianRelation   string
	GuardianPhone      string
	GuardianOccupation string
}

// Validate memvalidasi command.
func (c RegisterStudentCommand) Validate() error {
	if c.NISN == "" {
		return errors.New("NISN wajib diisi")
	}
	if c.NIS == "" {
		return errors.New("NIS wajib diisi")
	}
	if c.FirstName == "" {
		return errors.New("nama depan wajib diisi")
	}
	if c.BirthDate.IsZero() {
		return errors.New("tanggal lahir wajib diisi")
	}
	if c.Gender == "" {
		return errors.New("jenis kelamin wajib diisi")
	}
	if c.AddressStreet == "" || c.AddressVillage == "" || c.AddressCity == "" {
		return errors.New("alamat lengkap wajib diisi")
	}
	if c.ContactPhone == "" && c.ContactEmail == "" {
		return errors.New("minimal phone atau email harus diisi")
	}
	if c.GuardianName == "" {
		return errors.New("nama wali wajib diisi")
	}
	if c.GuardianRelation == "" {
		return errors.New("relasi wali wajib diisi")
	}
	return nil
}
