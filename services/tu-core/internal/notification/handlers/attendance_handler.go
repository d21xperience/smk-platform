package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/dispatcher"
)

// AttendanceRepository adalah interface untuk mengambil data kehadiran.
type AttendanceRepository interface {
	GetStudentContactInfo(ctx context.Context, studentID string) (*StudentContactInfo, error)
	GetClassName(ctx context.Context, classID string) (string, error)
}

// AttendanceHandler menangani notifikasi terkait kehadiran/absensi.
type AttendanceHandler struct {
	dispatcher *dispatcher.NotificationDispatcher
	repo       AttendanceRepository
}

// NewAttendanceHandler membuat instance baru dari AttendanceHandler.
func NewAttendanceHandler(dispatcher *dispatcher.NotificationDispatcher, repo AttendanceRepository) *AttendanceHandler {
	return &AttendanceHandler{
		dispatcher: dispatcher,
		repo:       repo,
	}
}

// attendanceSubmittedEvent adalah struktur event AttendanceSubmitted.
type attendanceSubmittedEvent struct {
	SchoolID         string             `json:"schoolId"`
	AcademicPeriodID string             `json:"academicPeriodId"`
	ClassID          string             `json:"classId"`
	Date             string             `json:"date"`
	TeacherID        string             `json:"teacherId"`
	Records          []attendanceRecord `json:"records"`
	CorrelationID    string             `json:"correlationId"`
}

// attendanceRecord adalah data kehadiran per siswa.
type attendanceRecord struct {
	StudentID string `json:"studentId"`
	Status    string `json:"status"` // PRESENT, ABSENT, LATE, SICK, PERMISSION
}

// Handle memproses event AttendanceSubmitted dan mengirim notifikasi untuk siswa yang absen.
func (h *AttendanceHandler) Handle(ctx context.Context, msg []byte) error {
	var event attendanceSubmittedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event AttendanceSubmitted: %w", err)
	}

	log.Printf("[AttendanceHandler] Memproses absensi kelas %s tanggal %s", event.ClassID, event.Date)

	// Fetch nama kelas
	className, err := h.repo.GetClassName(ctx, event.ClassID)
	if err != nil {
		log.Printf("[AttendanceHandler] Gagal fetch class name: %v", err)
		className = "Kelas"
	}

	// Proses setiap record
	for _, record := range event.Records {
		// Hanya kirim notifikasi untuk siswa yang ABSENT (tidak hadir)
		if record.Status != "ABSENT" {
			continue
		}

		// Fetch contact info
		contactInfo, err := h.repo.GetStudentContactInfo(ctx, record.StudentID)
		if err != nil {
			log.Printf("[AttendanceHandler] Gagal fetch contact info untuk student %s: %v", record.StudentID, err)
			continue
		}

		templateData := map[string]interface{}{
			"studentName": contactInfo.StudentName,
			"parentName":  contactInfo.ParentName,
			"className":   className,
			"date":        event.Date,
		}

		// Kirim notifikasi WhatsApp
		if contactInfo.ParentPhone != "" {
			reqWA := contract.NotificationRequest{
				Channel:          "whatsapp",
				Recipient:        contactInfo.ParentPhone,
				TemplateCode:     "attendance_absent_wa",
				Data:             templateData,
				SchoolID:         event.SchoolID,
				AcademicPeriodID: event.AcademicPeriodID,
				CorrelationID:    event.CorrelationID,
			}
			if err := h.dispatcher.Dispatch(ctx, reqWA); err != nil {
				log.Printf("[AttendanceHandler] Gagal kirim WA untuk student %s: %v", record.StudentID, err)
			}
		}

		// Kirim notifikasi Email
		if contactInfo.ParentEmail != "" {
			reqEmail := contract.NotificationRequest{
				Channel:          "email",
				Recipient:        contactInfo.ParentEmail,
				TemplateCode:     "attendance_absent_email",
				Data:             templateData,
				SchoolID:         event.SchoolID,
				AcademicPeriodID: event.AcademicPeriodID,
				CorrelationID:    event.CorrelationID,
			}
			if err := h.dispatcher.Dispatch(ctx, reqEmail); err != nil {
				log.Printf("[AttendanceHandler] Gagal kirim Email untuk student %s: %v", record.StudentID, err)
			}
		}
	}

	log.Printf("[AttendanceHandler] Notifikasi absensi kelas %s berhasil diproses", event.ClassID)
	return nil
}
