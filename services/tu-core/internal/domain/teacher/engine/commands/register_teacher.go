package commands

import (
	"errors"
	"time"
)

// RegisterTeacherCommand adalah command untuk mendaftarkan guru baru.
type RegisterTeacherCommand struct {
	TeacherID  string
	NIP        string
	NUPTK      string
	FirstName  string
	MiddleName string
	LastName   string
	Title      string
	BirthDate  time.Time
	Gender     string
	Phone      string
	Email      string
	Address    string
	Subject    string
	JoinDate   time.Time
}

// Validate memvalidasi command.
func (c RegisterTeacherCommand) Validate() error {
	if c.NIP == "" {
		return errors.New("NIP wajib diisi")
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
	if c.Subject == "" {
		return errors.New("mata pelajaran wajib diisi")
	}
	if c.JoinDate.IsZero() {
		return errors.New("tanggal bergabung wajib diisi")
	}
	return nil
}
