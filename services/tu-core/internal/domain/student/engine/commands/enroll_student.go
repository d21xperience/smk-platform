package commands

import "errors"

// EnrollStudentCommand adalah command untuk mendaftarkan siswa ke kelas.
type EnrollStudentCommand struct {
	StudentID    string
	EnrollmentID string
	ClassID      string
	FromClassID  string // opsional, jika kosong diambil dari enrollment aktif
	ToClassID    string
	PeriodID     string
	Reason       string // alasan pindah, misal: "pindah kelas", "pindah jurusan"
	MovedBy      string // UserID
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

// MoveStudentCommand adalah command untuk memindahkan siswa ke kelas lain.
type MoveStudentCommand struct {
	StudentID string
	ToClassID string
	Reason    string
}

func (c MoveStudentCommand) Validate() error {
	if c.StudentID == "" {
		return errors.New("student id is required")
	}
	if c.ToClassID == "" {
		return errors.New("target class id is required")
	}
	return nil
}
