package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/registration/events"
	"sekolah-platform/services/tu-core/internal/domain/registration/valueobjects"

	"github.com/google/uuid"
)

type Registration struct {
	ID                string
	SchoolID          string
	AcademicPeriodID  string
	CalonNISN         string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir time.Time
	CalonAlamat       string
	CalonTelepon      string
	CalonEmail        string
	NamaAyah          string
	NamaIbu           string
	TeleponOrtu       string
	PekerjaanAyah     string
	PekerjaanIbu      string
	AsalSekolah       string
	JurusanDipilih    string
	AlasanMemilih     string
	Status            valueobjects.RegistrationStatus
	CatatanVerifikasi string
	VerifiedBy        string
	VerifiedAt        time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
	domainEvents      []types.DomainEvent
}

type RegistrationData struct {
	ID                string
	CalonNISN         string
	CalonNama         string
	CalonJenisKelamin string
	CalonTempatLahir  string
	CalonTanggalLahir time.Time
	CalonAlamat       string
	CalonTelepon      string
	CalonEmail        string
	NamaAyah          string
	NamaIbu           string
	TeleponOrtu       string
	PekerjaanAyah     string
	PekerjaanIbu      string
	AsalSekolah       string
	JurusanDipilih    string
	AlasanMemilih     string
}

func NewRegistration(data RegistrationData, ctx *context.OperationalContext) (*Registration, error) {
	if data.ID == "" {
		data.ID = uuid.New().String()
	}

	if data.CalonNama == "" {
		return nil, errors.New("nama calon siswa wajib diisi")
	}

	if data.CalonJenisKelamin == "" {
		return nil, errors.New("jenis kelamin wajib diisi")
	}

	if data.CalonTanggalLahir.IsZero() {
		return nil, errors.New("tanggal lahir wajib diisi")
	}

	r := &Registration{
		ID:                data.ID,
		SchoolID:          ctx.SchoolID,
		AcademicPeriodID:  ctx.AcademicPeriodID,
		CalonNISN:         data.CalonNISN,
		CalonNama:         data.CalonNama,
		CalonJenisKelamin: data.CalonJenisKelamin,
		CalonTempatLahir:  data.CalonTempatLahir,
		CalonTanggalLahir: data.CalonTanggalLahir,
		CalonAlamat:       data.CalonAlamat,
		CalonTelepon:      data.CalonTelepon,
		CalonEmail:        data.CalonEmail,
		NamaAyah:          data.NamaAyah,
		NamaIbu:           data.NamaIbu,
		TeleponOrtu:       data.TeleponOrtu,
		PekerjaanAyah:     data.PekerjaanAyah,
		PekerjaanIbu:      data.PekerjaanIbu,
		AsalSekolah:       data.AsalSekolah,
		JurusanDipilih:    data.JurusanDipilih,
		AlasanMemilih:     data.AlasanMemilih,
		Status:            valueobjects.StatusDraft,
		CreatedAt:         time.Now().UTC(),
		UpdatedAt:         time.Now().UTC(),
		domainEvents:      []types.DomainEvent{},
	}

	r.recordEvent(events.NewRegistrationCreated(events.RegistrationCreatedPayload{
		RegistrationID: r.ID,
		CalonNISN:      r.CalonNISN,
		CalonNama:      r.CalonNama,
	}, ctx))

	return r, nil
}

func (r *Registration) Submit(ctx *context.OperationalContext) error {
	if !r.Status.CanSubmit() {
		return errors.New("hanya pendaftaran dengan status DRAFT yang dapat disubmit")
	}

	r.Status = valueobjects.StatusSubmitted
	r.UpdatedAt = time.Now().UTC()

	r.recordEvent(events.NewRegistrationSubmitted(events.RegistrationSubmittedPayload{
		RegistrationID: r.ID,
		CalonNISN:      r.CalonNISN,
		CalonNama:      r.CalonNama,
	}, ctx))

	return nil
}

func (r *Registration) Verify(catatan, verifiedBy string, ctx *context.OperationalContext) error {
	if !r.Status.CanVerify() {
		return errors.New("hanya pendaftaran dengan status SUBMITTED yang dapat diverifikasi")
	}

	r.Status = valueobjects.StatusVerified
	r.CatatanVerifikasi = catatan
	r.VerifiedBy = verifiedBy
	r.VerifiedAt = time.Now().UTC()
	r.UpdatedAt = time.Now().UTC()

	r.recordEvent(events.NewRegistrationVerified(events.RegistrationVerifiedPayload{
		RegistrationID: r.ID,
		VerifiedBy:     verifiedBy,
		Catatan:        catatan,
	}, ctx))

	return nil
}

func (r *Registration) Approve(catatan, approvedBy string, ctx *context.OperationalContext) error {
	if !r.Status.CanApprove() {
		return errors.New("hanya pendaftaran dengan status VERIFIED yang dapat disetujui")
	}

	r.Status = valueobjects.StatusApproved
	r.CatatanVerifikasi = catatan
	r.VerifiedBy = approvedBy
	r.VerifiedAt = time.Now().UTC()
	r.UpdatedAt = time.Now().UTC()

	r.recordEvent(events.NewRegistrationApproved(events.RegistrationApprovedPayload{
		RegistrationID: r.ID,
		CalonNISN:      r.CalonNISN,
		CalonNama:      r.CalonNama,
		ApprovedBy:     approvedBy,
	}, ctx))

	return nil
}

func (r *Registration) Reject(alasan, rejectedBy string, ctx *context.OperationalContext) error {
	if !r.Status.CanReject() {
		return errors.New("hanya pendaftaran dengan status SUBMITTED atau VERIFIED yang dapat ditolak")
	}

	if alasan == "" {
		return errors.New("alasan penolakan wajib diisi")
	}

	r.Status = valueobjects.StatusRejected
	r.CatatanVerifikasi = alasan
	r.VerifiedBy = rejectedBy
	r.VerifiedAt = time.Now().UTC()
	r.UpdatedAt = time.Now().UTC()

	r.recordEvent(events.NewRegistrationRejected(events.RegistrationRejectedPayload{
		RegistrationID: r.ID,
		CalonNISN:      r.CalonNISN,
		Alasan:         alasan,
		RejectedBy:     rejectedBy,
	}, ctx))

	return nil
}

func (r *Registration) recordEvent(event types.DomainEvent) {
	r.domainEvents = append(r.domainEvents, event)
}

func (r *Registration) UncommittedEvents() []types.DomainEvent {
	return r.domainEvents
}

func (r *Registration) ClearEvents() {
	r.domainEvents = []types.DomainEvent{}
}
