// Package engine menyediakan business logic untuk domain Student.
// Engine adalah PURE GO, tidak boleh import GORM, Redis, atau infrastruktur lain.
// SYMMETRIC dengan frontend (apps/siakad-tu/src/engine/student/).
package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/student/models"
)

// Result adalah hasil dari engine operation.
// SYMMETRIC dengan Result di frontend.
type Result struct {
	Student *models.Student
	Events  []types.DomainEvent
	Error   error
}

// Success membuat Result sukses.
func Success(student *models.Student) Result {
	return Result{
		Student: student,
		Events:  student.UncommittedEvents(),
		Error:   nil,
	}
}

// Failure membuat Result gagal.
func Failure(err error) Result {
	return Result{
		Student: nil,
		Events:  nil,
		Error:   err,
	}
}

// StudentEngine adalah business engine untuk domain Student.
// SYMMETRIC dengan StudentEngine.js di frontend.
type StudentEngine struct{}

// NewStudentEngine membuat StudentEngine baru.
func NewStudentEngine() *StudentEngine {
	return &StudentEngine{}
}

// RegisterStudent mendaftarkan siswa baru.
func (e *StudentEngine) RegisterStudent(cmd commands.RegisterStudentCommand, ctx *context.OperationalContext) Result {
	// Validasi command
	if err := cmd.Validate(); err != nil {
		return Failure(err)
	}

	// Bangun Student aggregate
	data := models.StudentData{
		StudentID:  cmd.StudentID,
		SchoolID:   ctx.SchoolID,
		NISN:       cmd.NISN,
		NIS:        cmd.NIS,
		FirstName:  cmd.FirstName,
		MiddleName: cmd.MiddleName,
		LastName:   cmd.LastName,
		BirthDate:  cmd.BirthDate,
		Gender:     cmd.Gender,
		Address: models.Address{
			Street:     cmd.AddressStreet,
			RtRw:       cmd.AddressRtRw,
			Village:    cmd.AddressVillage,
			District:   cmd.AddressDistrict,
			City:       cmd.AddressCity,
			PostalCode: cmd.AddressPostalCode,
		},
		ContactInfo: models.ContactInfo{
			Phone: cmd.ContactPhone,
			Email: cmd.ContactEmail,
		},
		GuardianInfo: models.GuardianInfo{
			Name:       cmd.GuardianName,
			Relation:   cmd.GuardianRelation,
			Phone:      cmd.GuardianPhone,
			Occupation: cmd.GuardianOccupation,
		},
	}

	student, err := models.NewStudent(data, ctx)
	if err != nil {
		return Failure(err)
	}

	return Success(student)
}

// EnrollStudent mendaftarkan siswa ke kelas.
func (e *StudentEngine) EnrollStudent(cmd commands.EnrollStudentCommand, ctx *context.OperationalContext, currentStudent *models.Student) Result {
	if err := cmd.Validate(); err != nil {
		return Failure(err)
	}

	if currentStudent == nil {
		return Failure(errStudentNotFound)
	}

	// Enroll
	if err := currentStudent.Enroll(cmd.EnrollmentID, ctx.AcademicPeriodID, cmd.ClassID, ctx); err != nil {
		return Failure(err)
	}

	return Success(currentStudent)
}

// GraduateStudent meluluskan siswa.
func (e *StudentEngine) GraduateStudent(ctx *context.OperationalContext, currentStudent *models.Student) Result {
	if currentStudent == nil {
		return Failure(errStudentNotFound)
	}

	if err := currentStudent.Graduate(ctx); err != nil {
		return Failure(err)
	}

	return Success(currentStudent)
}

// TransferStudent memindahkan siswa.
func (e *StudentEngine) TransferStudent(targetSchool, reason string, ctx *context.OperationalContext, currentStudent *models.Student) Result {
	if currentStudent == nil {
		return Failure(errStudentNotFound)
	}

	if err := currentStudent.Transfer(targetSchool, reason, ctx); err != nil {
		return Failure(err)
	}

	return Success(currentStudent)
}

var errStudentNotFound = errString("siswa tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
