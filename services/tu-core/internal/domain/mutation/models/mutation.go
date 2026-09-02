package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/mutation/events"
	"sekolah-platform/services/tu-core/internal/domain/mutation/valueobjects"

	"github.com/google/uuid"
)

type Mutation struct {
	ID                string
	SchoolID          string
	AcademicPeriodID  string
	StudentID         *string
	MutationType      valueobjects.MutationType
	CalonNama         string
	CalonNISN         string
	AsalSekolah       string
	TujuanSekolah     string
	Alasan            string
	DokumenURL        string
	Status            valueobjects.MutationStatus
	CatatanVerifikasi string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	domainEvents      []types.DomainEvent
}

type MutationData struct {
	ID               string
	SchoolID         string
	AcademicPeriodID string
	StudentID        *string
	MutationType     valueobjects.MutationType
	CalonNama        string
	CalonNISN        string
	AsalSekolah      string
	TujuanSekolah    string
	Alasan           string
	DokumenURL       string
}

func NewMutation(data MutationData, ctx *context.OperationalContext) (*Mutation, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}
	if data.Alasan == "" {
		return nil, errors.New("alasan mutasi wajib diisi")
	}
	if data.MutationType == valueobjects.TypeMasuk {
		if data.CalonNama == "" || data.CalonNISN == "" {
			return nil, errors.New("data calon siswa wajib diisi untuk mutasi masuk")
		}
	} else if data.MutationType == valueobjects.TypeKeluar {
		if data.StudentID == nil || *data.StudentID == "" {
			return nil, errors.New("student_id wajib diisi untuk mutasi keluar")
		}
		if data.TujuanSekolah == "" {
			return nil, errors.New("tujuan sekolah wajib diisi untuk mutasi keluar")
		}
	}

	m := &Mutation{
		ID:               data.ID,
		SchoolID:         ctx.SchoolID,
		AcademicPeriodID: ctx.AcademicPeriodID,
		StudentID:        data.StudentID,
		MutationType:     data.MutationType,
		CalonNama:        data.CalonNama,
		CalonNISN:        data.CalonNISN,
		AsalSekolah:      data.AsalSekolah,
		TujuanSekolah:    data.TujuanSekolah,
		Alasan:           data.Alasan,
		DokumenURL:       data.DokumenURL,
		Status:           valueobjects.StatusPending,
		CreatedAt:        time.Now().UTC(),
		UpdatedAt:        time.Now().UTC(),
		domainEvents:     []types.DomainEvent{},
	}

	m.recordEvent(events.NewMutationCreated(events.MutationCreatedPayload{
		MutationID:    m.ID,
		StudentID:     m.StudentID,
		MutationType:  m.MutationType.String(),
		CalonNama:     m.CalonNama,
		CalonNISN:     m.CalonNISN,
		AsalSekolah:   m.AsalSekolah,
		TujuanSekolah: m.TujuanSekolah,
		Alasan:        m.Alasan,
	}, ctx))

	return m, nil
}

func (m *Mutation) Approve(catatan string, ctx *context.OperationalContext) error {
	if !m.Status.IsPending() {
		return errors.New("hanya mutasi dengan status PENDING yang dapat disetujui")
	}
	m.Status = valueobjects.StatusDisetujui
	m.CatatanVerifikasi = catatan
	m.UpdatedAt = time.Now().UTC()

	m.recordEvent(events.NewMutationApproved(events.MutationApprovedPayload{
		MutationID:    m.ID,
		StudentID:     m.StudentID,
		MutationType:  m.MutationType.String(),
		CalonNama:     m.CalonNama,
		CalonNISN:     m.CalonNISN,
		AsalSekolah:   m.AsalSekolah,
		TujuanSekolah: m.TujuanSekolah,
		Catatan:       catatan,
	}, ctx))

	return nil
}

func (m *Mutation) Reject(catatan string, ctx *context.OperationalContext) error {
	if !m.Status.IsPending() {
		return errors.New("hanya mutasi dengan status PENDING yang dapat ditolak")
	}
	if catatan == "" {
		return errors.New("catatan penolakan wajib diisi")
	}
	m.Status = valueobjects.StatusDitolak
	m.CatatanVerifikasi = catatan
	m.UpdatedAt = time.Now().UTC()

	m.recordEvent(events.NewMutationRejected(events.MutationRejectedPayload{
		MutationID: m.ID,
		StudentID:  m.StudentID,
		Catatan:    catatan,
	}, ctx))

	return nil
}

func (m *Mutation) recordEvent(event types.DomainEvent) {
	m.domainEvents = append(m.domainEvents, event)
}

func (m *Mutation) UncommittedEvents() []types.DomainEvent {
	return m.domainEvents
}

func (m *Mutation) ClearEvents() {
	m.domainEvents = []types.DomainEvent{}
}
