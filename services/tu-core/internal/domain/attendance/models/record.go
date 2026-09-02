package models

import (
	"errors"
	"time"

	"sekolah-platform/services/tu-core/internal/domain/attendance/valueobjects"

	"github.com/google/uuid"
)

type Record struct {
	ID        string
	SessionID string
	StudentID string
	Status    valueobjects.AttendanceStatus
	Note      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RecordData struct {
	ID        string
	SessionID string
	StudentID string
	Status    valueobjects.AttendanceStatus
	Note      string
}

func NewRecord(data RecordData) (*Record, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.SessionID == "" {
		return nil, errors.New("session_id wajib diisi")
	}

	if data.StudentID == "" {
		return nil, errors.New("student_id wajib diisi")
	}

	r := &Record{
		ID:        data.ID,
		SessionID: data.SessionID,
		StudentID: data.StudentID,
		Status:    data.Status,
		Note:      data.Note,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	return r, nil
}

func (r *Record) UpdateStatus(status valueobjects.AttendanceStatus, note string) {
	r.Status = status
	r.Note = note
	r.UpdatedAt = time.Now().UTC()
}
