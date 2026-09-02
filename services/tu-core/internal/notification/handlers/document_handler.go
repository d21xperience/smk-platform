package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"sekolah-platform/platform/notification/contract"
	"sekolah-platform/platform/notification/dispatcher"
)

// DocumentRepository adalah interface untuk mengambil data dokumen dan siswa.
type DocumentRepository interface {
	GetStudentContactInfo(ctx context.Context, studentID string) (*StudentContactInfo, error)
}

// DocumentHandler menangani notifikasi terkait dokumen (SK, Surat Keterangan, dll).
type DocumentHandler struct {
	dispatcher *dispatcher.NotificationDispatcher
	repo       DocumentRepository
}

// NewDocumentHandler membuat instance baru dari DocumentHandler.
func NewDocumentHandler(dispatcher *dispatcher.NotificationDispatcher, repo DocumentRepository) *DocumentHandler {
	return &DocumentHandler{
		dispatcher: dispatcher,
		repo:       repo,
	}
}

// documentIssuedEvent adalah struktur event DocumentIssued.
type documentIssuedEvent struct {
	DocumentID       string `json:"documentId"`
	SchoolID         string `json:"schoolId"`
	AcademicPeriodID string `json:"academicPeriodId"`
	StudentID        string `json:"studentId"`
	DocumentType     string `json:"documentType"` // SK, SURAT_KETERANGAN, RAPOR, dll
	DocumentNumber   string `json:"documentNumber"`
	IssueDate        string `json:"issueDate"`
	CorrelationID    string `json:"correlationId"`
}

// Handle memproses event DocumentIssued.
func (h *DocumentHandler) Handle(ctx context.Context, msg []byte) error {
	var event documentIssuedEvent
	if err := json.Unmarshal(msg, &event); err != nil {
		return fmt.Errorf("gagal parse event DocumentIssued: %w", err)
	}

	log.Printf("[DocumentHandler] Memproses dokumen %s untuk student %s", event.DocumentID, event.StudentID)

	// Fetch contact info
	contactInfo, err := h.repo.GetStudentContactInfo(ctx, event.StudentID)
	if err != nil {
		return fmt.Errorf("gagal fetch contact info: %w", err)
	}

	templateData := map[string]interface{}{
		"studentName":    contactInfo.StudentName,
		"parentName":     contactInfo.ParentName,
		"documentType":   event.DocumentType,
		"documentNumber": event.DocumentNumber,
		"issueDate":      event.IssueDate,
		"documentId":     event.DocumentID,
	}

	// Kirim notifikasi via WhatsApp
	if contactInfo.ParentPhone != "" {
		reqWA := contract.NotificationRequest{
			Channel:          "whatsapp",
			Recipient:        contactInfo.ParentPhone,
			TemplateCode:     "document_issued_wa",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqWA); err != nil {
			log.Printf("[DocumentHandler] Gagal kirim WA: %v", err)
		}
	}

	// Kirim notifikasi via Email (dengan attachment jika perlu)
	if contactInfo.ParentEmail != "" {
		reqEmail := contract.NotificationRequest{
			Channel:          "email",
			Recipient:        contactInfo.ParentEmail,
			TemplateCode:     "document_issued_email",
			Data:             templateData,
			SchoolID:         event.SchoolID,
			AcademicPeriodID: event.AcademicPeriodID,
			CorrelationID:    event.CorrelationID,
		}
		if err := h.dispatcher.Dispatch(ctx, reqEmail); err != nil {
			log.Printf("[DocumentHandler] Gagal kirim Email: %v", err)
		}
	}

	// Kirim push notification ke aplikasi mobile siswa (jika ada)
	// Asumsi: FCM token disimpan di tabel students atau user_devices
	// fcmToken := h.getStudentFCMToken(ctx, event.StudentID)
	// if fcmToken != "" {
	// 	reqFCM := contract.NotificationRequest{
	// 		Channel:          "firebase",
	// 		Recipient:        fcmToken,
	// 		TemplateCode:     "document_issued_fcm",
	// 		Data:             templateData,
	// 		SchoolID:         event.SchoolID,
	// 		AcademicPeriodID: event.AcademicPeriodID,
	// 		CorrelationID:    event.CorrelationID,
	// 	}
	// 	h.dispatcher.Dispatch(ctx, reqFCM)
	// }

	log.Printf("[DocumentHandler] Notifikasi dokumen %s berhasil diproses", event.DocumentID)
	return nil
}
