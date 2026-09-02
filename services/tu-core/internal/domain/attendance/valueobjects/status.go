package valueobjects

import "errors"

type AttendanceStatus string

const (
	StatusPresent    AttendanceStatus = "PRESENT"
	StatusAbsent     AttendanceStatus = "ABSENT"
	StatusLate       AttendanceStatus = "LATE"
	StatusSick       AttendanceStatus = "SICK"
	StatusPermission AttendanceStatus = "PERMISSION"
)

type SessionStatus string

const (
	SessionDraft     SessionStatus = "DRAFT"
	SessionSubmitted SessionStatus = "SUBMITTED"
	SessionClosed    SessionStatus = "CLOSED"
)

func NewAttendanceStatus(s string) (AttendanceStatus, error) {
	switch AttendanceStatus(s) {
	case StatusPresent, StatusAbsent, StatusLate, StatusSick, StatusPermission:
		return AttendanceStatus(s), nil
	default:
		return "", errors.New("invalid attendance status")
	}
}

func (s AttendanceStatus) String() string {
	return string(s)
}

func NewSessionStatus(s string) (SessionStatus, error) {
	switch SessionStatus(s) {
	case SessionDraft, SessionSubmitted, SessionClosed:
		return SessionStatus(s), nil
	default:
		return "", errors.New("invalid session status")
	}
}

func (s SessionStatus) String() string {
	return string(s)
}

func (s SessionStatus) IsDraft() bool {
	return s == SessionDraft
}

func (s SessionStatus) CanSubmit() bool {
	return s == SessionDraft
}

func (s SessionStatus) CanClose() bool {
	return s == SessionSubmitted
}
