package valueobjects

import "fmt"

// AttendanceStatus adalah value object untuk status kehadiran.
type AttendanceStatus struct {
	code  string
	label string
}

var (
	AttendanceStatusPresent = AttendanceStatus{code: "PRESENT", label: "Hadir"}
	AttendanceStatusAbsent  = AttendanceStatus{code: "ABSENT", label: "Alpha"}
	AttendanceStatusSick    = AttendanceStatus{code: "SICK", label: "Sakit"}
	AttendanceStatusLeave   = AttendanceStatus{code: "LEAVE", label: "Izin"}
)

// NewAttendanceStatus membuat AttendanceStatus dari code.
func NewAttendanceStatus(code string) (AttendanceStatus, error) {
	switch code {
	case "PRESENT":
		return AttendanceStatusPresent, nil
	case "ABSENT":
		return AttendanceStatusAbsent, nil
	case "SICK":
		return AttendanceStatusSick, nil
	case "LEAVE":
		return AttendanceStatusLeave, nil
	default:
		return AttendanceStatus{}, fmt.Errorf("status kehadiran tidak valid: %s", code)
	}
}

// Code mengembalikan code status.
func (s AttendanceStatus) Code() string { return s.code }

// Label mengembalikan label status.
func (s AttendanceStatus) Label() string { return s.label }

// IsPresent mengecek apakah status hadir.
func (s AttendanceStatus) IsPresent() bool { return s.code == "PRESENT" }

// Equals mengecek kesamaan.
func (s AttendanceStatus) Equals(other AttendanceStatus) bool {
	return s.code == other.code
}
