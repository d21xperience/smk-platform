package commands

import "errors"

// EnrollStudentCommand adalah command untuk mendaftarkan siswa ke kelas.
type EnrollStudentCommand struct {
	StudentID    string
	EnrollmentID string
	ClassID      string
}

// Validate memvalidasi command.
func (c EnrollStudentCommand) Validate() error {
	if c.StudentID == "" {
		return errors.New("studentId wajib diisi")
	}
	if c.EnrollmentID == "" {
		return errors.New("enrollmentId wajib diisi")
	}
	if c.ClassID == "" {
		return errors.New("classId wajib diisi")
	}
	return nil
}
