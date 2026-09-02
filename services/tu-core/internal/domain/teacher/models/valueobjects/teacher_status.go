package valueobjects
import "fmt"

// TeacherStatus adalah value object untuk status guru.
type TeacherStatus struct {
	code  string
	label string
}

var (
	TeacherStatusActive   = TeacherStatus{code: "ACTIVE", label: "Aktif"}
	TeacherStatusInactive = TeacherStatus{code: "INACTIVE", label: "Tidak Aktif"}
	TeacherStatusResigned = TeacherStatus{code: "RESIGNED", label: "Resign"}
	TeacherStatusRetired  = TeacherStatus{code: "RETIRED", label: "Pensiun"}
)

// NewTeacherStatus membuat TeacherStatus dari code.
func NewTeacherStatus(code string) (TeacherStatus, error) {
	switch code {
	case "ACTIVE":
		return TeacherStatusActive, nil
	case "INACTIVE":
		return TeacherStatusInactive, nil
	case "RESIGNED":
		return TeacherStatusResigned, nil
	case "RETIRED":
		return TeacherStatusRetired, nil
	default:
		return TeacherStatus{}, fmt.Errorf("status guru tidak valid: %s", code)
	}
}

// Code mengembalikan code status.
func (s TeacherStatus) Code() string { return s.code }

// Label mengembalikan label status.
func (s TeacherStatus) Label() string { return s.label }

// IsActive mengecek apakah status aktif.
func (s TeacherStatus) IsActive() bool { return s.code == "ACTIVE" }
