package models

import (
	"errors"
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/student/events"
	"sekolah-platform/services/tu-core/internal/domain/student/models/valueobjects"

	"github.com/google/uuid"
)

// ============================================================================
// VALUE OBJECTS (Local Structs untuk Data Transfer)
// ============================================================================

// Address adalah struct untuk data alamat.
type Address struct {
	Street     string
	RtRw       string
	Village    string
	District   string
	City       string
	PostalCode string
}

// ContactInfo adalah struct untuk data kontak.
type ContactInfo struct {
	Phone string
	Email string
}

// GuardianInfo adalah struct untuk data wali/orang tua.
type GuardianInfo struct {
	Name       string
	Relation   string // father, mother, guardian
	Phone      string
	Occupation string
}

// ============================================================================
// AGGREGATE ROOT
// ============================================================================

// Student adalah Aggregate Root untuk domain Student.
type Student struct {
	studentID    string
	schoolID     string
	nisn         valueobjects.NISN
	nis          valueobjects.NIS
	fullName     valueobjects.FullName
	birthDate    time.Time
	gender       valueobjects.Gender
	address      Address
	contactInfo  ContactInfo
	guardianInfo GuardianInfo
	status       valueobjects.StudentStatus
	enrollments  []*Enrollment
	events       []types.DomainEvent
	createdAt    time.Time
	updatedAt    time.Time
}

// StudentData adalah data input untuk membuat Student.
type StudentData struct {
	StudentID    string
	SchoolID     string
	NISN         string
	NIS          string
	FirstName    string
	MiddleName   string
	LastName     string
	BirthDate    time.Time
	Gender       string
	Address      Address
	ContactInfo  ContactInfo
	GuardianInfo GuardianInfo
}

// NewStudent membuat Student baru.
func NewStudent(data StudentData, ctx *context.OperationalContext) (*Student, error) {
	if data.StudentID == "" {
		data.StudentID = uuid.New().String()
	}

	nisn, err := valueobjects.NewNISN(data.NISN)
	if err != nil {
		return nil, err
	}

	nis, err := valueobjects.NewNIS(data.NIS)
	if err != nil {
		return nil, err
	}

	fullName, err := valueobjects.NewFullName(data.FirstName, data.MiddleName, data.LastName)
	if err != nil {
		return nil, err
	}

	gender, err := valueobjects.NewGender(data.Gender)
	if err != nil {
		return nil, err
	}

	student := &Student{
		studentID:    data.StudentID,
		schoolID:     ctx.SchoolID,
		nisn:         nisn,
		nis:          nis,
		fullName:     fullName,
		birthDate:    data.BirthDate,
		gender:       gender,
		address:      data.Address,
		contactInfo:  data.ContactInfo,
		guardianInfo: data.GuardianInfo,
		status:       valueobjects.StudentStatusActive,
		enrollments:  []*Enrollment{},
		events:       []types.DomainEvent{},
		createdAt:    time.Now().UTC(),
		updatedAt:    time.Now().UTC(),
	}

	// ✅ AMAN: Mengirim Payload Struct, BUKAN pointer ke struct Student
	student.recordEvent(events.NewStudentCreated(events.StudentCreatedPayload{
		StudentID: student.studentID,
		NISN:      student.nisn.Value(),
		NIS:       student.nis.Value(),
		FullName:  student.fullName.DisplayName(),
		Gender:    student.gender.Code(),
		BirthDate: student.birthDate,
	}, ctx))

	return student, nil
}

// === GETTERS ===
func (s *Student) StudentID() string                  { return s.studentID }
func (s *Student) SchoolID() string                   { return s.schoolID }
func (s *Student) NISN() valueobjects.NISN            { return s.nisn }
func (s *Student) NIS() valueobjects.NIS              { return s.nis }
func (s *Student) FullName() valueobjects.FullName    { return s.fullName }
func (s *Student) BirthDate() time.Time               { return s.birthDate }
func (s *Student) Gender() valueobjects.Gender        { return s.gender }
func (s *Student) Address() Address                   { return s.address }
func (s *Student) ContactInfo() ContactInfo           { return s.contactInfo }
func (s *Student) GuardianInfo() GuardianInfo         { return s.guardianInfo }
func (s *Student) Status() valueobjects.StudentStatus { return s.status }
func (s *Student) Enrollments() []*Enrollment         { return s.enrollments }
func (s *Student) CreatedAt() time.Time               { return s.createdAt }
func (s *Student) UpdatedAt() time.Time               { return s.updatedAt }

// === BEHAVIORS ===

func (s *Student) Enroll(enrollmentID, periodID, classID string, ctx *context.OperationalContext) error {
	if !s.status.CanBeEnrolled() {
		return errors.New("siswa dengan status " + s.status.Label() + " tidak bisa di-enroll")
	}

	for _, e := range s.enrollments {
		if e.PeriodID() == periodID && e.IsActive() {
			return errors.New("siswa sudah memiliki enrollment aktif di periode " + periodID)
		}
	}

	enrollment := NewEnrollment(enrollmentID, s.schoolID, periodID, classID)
	s.enrollments = append(s.enrollments, enrollment)
	s.updatedAt = time.Now().UTC()

	// ✅ AMAN: Mengirim Payload Struct
	s.recordEvent(events.NewStudentEnrolled(events.StudentEnrolledPayload{
		StudentID:    s.studentID,
		EnrollmentID: enrollment.EnrollmentID(),
		PeriodID:     enrollment.PeriodID(),
		ClassID:      enrollment.ClassID(),
	}, ctx))

	return nil
}

func (s *Student) Graduate(ctx *context.OperationalContext) error {
	if !s.status.IsActive() {
		return errors.New("siswa dengan status " + s.status.Label() + " tidak bisa diluluskan")
	}

	for _, e := range s.enrollments {
		if e.IsActive() {
			e.Complete()
		}
	}

	s.status = valueobjects.StudentStatusGraduated
	s.updatedAt = time.Now().UTC()

	// ✅ AMAN: Mengirim Payload Struct
	s.recordEvent(events.NewStudentGraduated(events.StudentGraduatedPayload{
		StudentID:      s.studentID,
		GraduationDate: time.Now().UTC(),
	}, ctx))

	return nil
}

func (s *Student) Transfer(targetSchool, reason string, ctx *context.OperationalContext) error {
	if !s.status.IsActive() {
		return errors.New("siswa dengan status " + s.status.Label() + " tidak bisa dipindahkan")
	}

	for _, e := range s.enrollments {
		if e.IsActive() {
			e.Complete()
		}
	}

	s.status = valueobjects.StudentStatusTransferred
	s.updatedAt = time.Now().UTC()

	// ✅ AMAN: Mengirim Payload Struct
	s.recordEvent(events.NewStudentTransferred(events.StudentTransferredPayload{
		StudentID:    s.studentID,
		TargetSchool: targetSchool,
		Reason:       reason,
		TransferDate: time.Now().UTC(),
	}, ctx))

	return nil
}

// === EVENT MANAGEMENT ===

func (s *Student) recordEvent(event types.DomainEvent) {
	s.events = append(s.events, event)
}

func (s *Student) UncommittedEvents() []types.DomainEvent {
	return s.events
}

func (s *Student) ClearEvents() {
	s.events = []types.DomainEvent{}
}
