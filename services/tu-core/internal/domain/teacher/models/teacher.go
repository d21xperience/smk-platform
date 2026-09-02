// Package models menyediakan domain models untuk Teacher aggregate.
// SYMMETRIC dengan frontend (jika ada).
package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/teacher/events"
	"sekolah-platform/services/tu-core/internal/domain/teacher/models/valueobjects"

	"github.com/google/uuid"
)

// Teacher adalah Aggregate Root untuk domain Teacher.
type Teacher struct {
	teacherID      string
	schoolID       string
	nip            valueobjects.NIP
	nuptk          valueobjects.NUPTK
	fullName       FullName
	birthDate      time.Time
	gender         string // MALE, FEMALE
	phone          string
	email          string
	address        string
	subject        string // Mata pelajaran
	status         valueobjects.TeacherStatus
	joinDate       time.Time
	certifications []Certification
	events         []types.DomainEvent
	createdAt      time.Time
	updatedAt      time.Time
}

// FullName adalah value object untuk nama guru.
type FullName struct {
	FirstName  string
	MiddleName string
	LastName   string
	Title      string // Gelar depan (S.Pd., M.Pd., dll)
}

// DisplayName mengembalikan nama lengkap dengan gelar.
func (n FullName) DisplayName() string {
	parts := []string{}
	if n.Title != "" {
		parts = append(parts, n.Title)
	}
	parts = append(parts, n.FirstName)
	if n.MiddleName != "" {
		parts = append(parts, n.MiddleName)
	}
	if n.LastName != "" {
		parts = append(parts, n.LastName)
	}
	result := ""
	for i, p := range parts {
		if i == 1 && n.Title != "" {
			result += " " + p
		} else if i > 0 {
			result += " " + p
		} else {
			result = p
		}
	}
	return result
}

// Certification adalah value object untuk sertifikasi guru.
type Certification struct {
	ID         string
	Name       string
	Issuer     string
	IssuedDate time.Time
	ExpiryDate *time.Time
}

// TeacherData adalah data untuk membuat Teacher.
type TeacherData struct {
	TeacherID  string
	SchoolID   string
	NIP        string
	NUPTK      string
	FirstName  string
	MiddleName string
	LastName   string
	Title      string
	BirthDate  time.Time
	Gender     string
	Phone      string
	Email      string
	Address    string
	Subject    string
	JoinDate   time.Time
}

// NewTeacher membuat Teacher baru.
func NewTeacher(data TeacherData, ctx *context.OperationalContext) (*Teacher, error) {
	if data.TeacherID == "" {
		data.TeacherID = uuid.New().String()
	}

	nip, err := valueobjects.NewNIP(data.NIP)
	if err != nil {
		return nil, err
	}

	nuptk, err := valueobjects.NewNUPTK(data.NUPTK)
	if err != nil {
		return nil, err
	}

	if data.FirstName == "" {
		return nil, errors.New("nama depan wajib diisi")
	}

	teacher := &Teacher{
		teacherID: data.TeacherID,
		schoolID:  ctx.SchoolID,
		nip:       nip,
		nuptk:     nuptk,
		fullName: FullName{
			FirstName:  data.FirstName,
			MiddleName: data.MiddleName,
			LastName:   data.LastName,
			Title:      data.Title,
		},
		birthDate:      data.BirthDate,
		gender:         data.Gender,
		phone:          data.Phone,
		email:          data.Email,
		address:        data.Address,
		subject:        data.Subject,
		status:         valueobjects.TeacherStatusActive,
		joinDate:       data.JoinDate,
		certifications: []Certification{},
		events:         []types.DomainEvent{},
		createdAt:      time.Now().UTC(),
		updatedAt:      time.Now().UTC(),
	}

	teacher.recordEvent(events.NewTeacherCreated(events.TeacherCreatedPayload{
		TeacherID: teacher.teacherID,
		NIP:       teacher.nip.Value(),
		NUPTK:     teacher.nuptk.Value(),
		FullName:  teacher.fullName.DisplayName(),
		Subject:   teacher.subject,
		Gender:    teacher.gender,
	}, ctx))
	return teacher, nil
}

// Getters

func (t *Teacher) TeacherID() string                  { return t.teacherID }
func (t *Teacher) SchoolID() string                   { return t.schoolID }
func (t *Teacher) NIP() valueobjects.NIP              { return t.nip }
func (t *Teacher) NUPTK() valueobjects.NUPTK          { return t.nuptk }
func (t *Teacher) FullName() FullName                 { return t.fullName }
func (t *Teacher) BirthDate() time.Time               { return t.birthDate }
func (t *Teacher) Gender() string                     { return t.gender }
func (t *Teacher) Phone() string                      { return t.phone }
func (t *Teacher) Email() string                      { return t.email }
func (t *Teacher) Address() string                    { return t.address }
func (t *Teacher) Subject() string                    { return t.subject }
func (t *Teacher) Status() valueobjects.TeacherStatus { return t.status }
func (t *Teacher) JoinDate() time.Time                { return t.joinDate }
func (t *Teacher) Certifications() []Certification    { return t.certifications }
func (t *Teacher) CreatedAt() time.Time               { return t.createdAt }
func (t *Teacher) UpdatedAt() time.Time               { return t.updatedAt }

// Behaviors

// UpdateProfile mengupdate profil guru.
func (t *Teacher) UpdateProfile(phone, email, address string, ctx *context.OperationalContext) {
	if phone != "" {
		t.phone = phone
	}
	if email != "" {
		t.email = email
	}
	if address != "" {
		t.address = address
	}
	t.updatedAt = time.Now().UTC()
	t.recordEvent(events.NewTeacherProfileUpdated(events.TeacherProfileUpdatedPayload{
		TeacherID: t.teacherID,
		UpdatedAt: t.updatedAt,
	}, ctx))
}

// Deactivate menonaktifkan guru.
func (t *Teacher) Deactivate(reason string, ctx *context.OperationalContext) error {
	if !t.status.IsActive() {
		return errors.New("guru sudah tidak aktif")
	}
	t.status = valueobjects.TeacherStatusInactive
	t.updatedAt = time.Now().UTC()
	t.recordEvent(events.NewTeacherDeactivated(events.TeacherDeactivatedPayload{
		TeacherID:     t.teacherID,
		Reason:        reason,
		DeactivatedAt: time.Now().UTC(),
	}, ctx))
	return nil
}

// AddCertification menambahkan sertifikasi.
func (t *Teacher) AddCertification(cert Certification, ctx *context.OperationalContext) error {
	if cert.ID == "" {
		cert.ID = uuid.New().String()
	}
	if cert.Name == "" {
		return errors.New("nama sertifikasi wajib diisi")
	}
	t.certifications = append(t.certifications, cert)
	t.updatedAt = time.Now().UTC()
	t.recordEvent(events.NewTeacherCertificationAdded(events.TeacherCertificationAddedPayload{
		TeacherID:         t.teacherID,
		CertificationID:   cert.ID,
		CertificationName: cert.Name,
		Issuer:            cert.Issuer,
	}, ctx))
	return nil
}

// Event management

func (t *Teacher) recordEvent(event types.DomainEvent) {
	t.events = append(t.events, event)
}

// UncommittedEvents mengembalikan events yang belum di-commit.
func (t *Teacher) UncommittedEvents() []types.DomainEvent {
	return t.events
}

// ClearEvents membersihkan events.
func (t *Teacher) ClearEvents() {
	t.events = []types.DomainEvent{}
}
