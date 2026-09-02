package valueobjects

import "fmt"

// StudentStatus adalah value object untuk status siswa.
type StudentStatus struct {
	code  string
	label string
}

var (
	StudentStatusActive      = StudentStatus{code: "ACTIVE", label: "Aktif"}
	StudentStatusGraduated   = StudentStatus{code: "GRADUATED", label: "Lulus"}
	StudentStatusTransferred = StudentStatus{code: "TRANSFERRED", label: "Pindah"}
	StudentStatusDropped     = StudentStatus{code: "DROPPED", label: "Keluar"}
	StudentStatusAlumni      = StudentStatus{code: "ALUMNI", label: "Alumni"}
)

// NewStudentStatus membuat StudentStatus dari code.
func NewStudentStatus(code string) (StudentStatus, error) {
	switch code {
	case "ACTIVE":
		return StudentStatusActive, nil
	case "GRADUATED":
		return StudentStatusGraduated, nil
	case "TRANSFERRED":
		return StudentStatusTransferred, nil
	case "DROPPED":
		return StudentStatusDropped, nil
	case "ALUMNI":
		return StudentStatusAlumni, nil
	default:
		return StudentStatus{}, fmt.Errorf("status tidak valid: %s", code)
	}
}

// Code mengembalikan code status.
func (s StudentStatus) Code() string { return s.code }

// Label mengembalikan label status.
func (s StudentStatus) Label() string { return s.label }

// IsActive mengecek apakah status aktif.
func (s StudentStatus) IsActive() bool { return s.code == "ACTIVE" }

// CanBeEnrolled mengecek apakah siswa bisa di-enroll.
func (s StudentStatus) CanBeEnrolled() bool { return s.code == "ACTIVE" }

// Equals mengecek kesamaan.
func (s StudentStatus) Equals(other StudentStatus) bool { return s.code == other.code }
