package valueobjects

import "fmt"

// AttendanceType adalah value object untuk tipe absensi.
type AttendanceType struct {
	code  string
	label string
}

var (
	AttendanceTypeDaily   = AttendanceType{code: "DAILY", label: "Harian"}
	AttendanceTypeSubject = AttendanceType{code: "SUBJECT", label: "Per Pelajaran"}
	AttendanceTypeEvent   = AttendanceType{code: "EVENT", label: "Kegiatan"}
)

// NewAttendanceType membuat AttendanceType dari code.
func NewAttendanceType(code string) (AttendanceType, error) {
	switch code {
	case "DAILY":
		return AttendanceTypeDaily, nil
	case "SUBJECT":
		return AttendanceTypeSubject, nil
	case "EVENT":
		return AttendanceTypeEvent, nil
	default:
		return AttendanceType{}, fmt.Errorf("tipe absensi tidak valid: %s", code)
	}
}

// Code mengembalikan code tipe.
func (t AttendanceType) Code() string { return t.code }

// Label mengembalikan label tipe.
func (t AttendanceType) Label() string { return t.label }
