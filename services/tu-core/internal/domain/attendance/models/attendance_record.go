package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"sekolah-platform/platform/context"
	"sekolah-platform/services/tu-core/internal/domain/attendance/models/valueobjects"
)

// AttendanceRecord adalah Entity dalam AttendanceSession Aggregate.
// Satu record = satu siswa dalam satu sesi.
type AttendanceRecord struct {
	recordID   string
	sessionID  string
	schoolID   string
	studentID  string
	status     valueobjects.AttendanceStatus
	notes      string
	corrected  bool
	correctedAt *time.Time
	correctedReason string
	createdAt  time.Time
	updatedAt  time.Time
}

// AttendanceRecordData adalah data untuk membuat AttendanceRecord.
type AttendanceRecordData struct {
	RecordID  string
	StudentID string
	Status    string
	Notes     string
}

// NewAttendanceRecord membuat AttendanceRecord baru.
func NewAttendanceRecord(data AttendanceRecordData, sessionID, schoolID string, ctx *context.OperationalContext) (*AttendanceRecord, error) {
	if data.RecordID == "" {
		data.RecordID = uuid.New().String()
	}

	if data.StudentID == "" {
		return nil, errors.New("studentId wajib diisi")
	}

	status, err := valueobjects.NewAttendanceStatus(data.Status)
	if err != nil {
		return nil, err
	}

	return &AttendanceRecord{
		recordID:  data.RecordID,
		sessionID: sessionID,
		schoolID:  schoolID,
		studentID: data.StudentID,
		status:    status,
		notes:     data.Notes,
		createdAt: time.Now().UTC(),
		updatedAt: time.Now().UTC(),
	}, nil
}

// Getters

func (r *AttendanceRecord) RecordID() string                    { return r.recordID }
func (r *AttendanceRecord) SessionID() string                   { return r.sessionID }
func (r *AttendanceRecord) SchoolID() string                    { return r.schoolID }
func (r *AttendanceRecord) StudentID() string                   { return r.studentID }
func (r *AttendanceRecord) Status() valueobjects.AttendanceStatus { return r.status }
func (r *AttendanceRecord) Notes() string                       { return r.notes }
func (r *AttendanceRecord) IsCorrected() bool                   { return r.corrected }
func (r *AttendanceRecord) CreatedAt() time.Time                { return r.createdAt }
func (r *AttendanceRecord) UpdatedAt() time.Time                { return r.updatedAt }

// Correct mengoreksi status absensi.
func (r *AttendanceRecord) Correct(newStatus, reason string) error {
	if reason == "" {
		return errors.New("alasan koreksi wajib diisi")
	}

	status, err := valueobjects.NewAttendanceStatus(newStatus)
	if err != nil {
		return err
	}

	r.status = status
	r.corrected = true
	now := time.Now().UTC()
	r.correctedAt = &now
	r.correctedReason = reason
	r.updatedAt = now

	return nil
}
