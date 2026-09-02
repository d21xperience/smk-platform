package engine

import (
	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/teacher/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/teacher/models"
)

// Result adalah hasil dari engine operation.
type Result struct {
	Teacher *models.Teacher
	Events  []types.DomainEvent
	Error   error
}

func Success(teacher *models.Teacher) Result {
	return Result{
		Teacher: teacher,
		Events:  teacher.UncommittedEvents(),
		Error:   nil,
	}
}

func Failure(err error) Result {
	return Result{Teacher: nil, Events: nil, Error: err}
}

// TeacherEngine adalah business engine untuk domain Teacher.
type TeacherEngine struct{}

func NewTeacherEngine() *TeacherEngine {
	return &TeacherEngine{}
}

// RegisterTeacher mendaftarkan guru baru.
func (e *TeacherEngine) RegisterTeacher(cmd commands.RegisterTeacherCommand, ctx *context.OperationalContext) Result {
	if err := cmd.Validate(); err != nil {
		return Failure(err)
	}

	data := models.TeacherData{
		TeacherID:  cmd.TeacherID,
		SchoolID:   ctx.SchoolID,
		NIP:        cmd.NIP,
		NUPTK:      cmd.NUPTK,
		FirstName:  cmd.FirstName,
		MiddleName: cmd.MiddleName,
		LastName:   cmd.LastName,
		Title:      cmd.Title,
		BirthDate:  cmd.BirthDate,
		Gender:     cmd.Gender,
		Phone:      cmd.Phone,
		Email:      cmd.Email,
		Address:    cmd.Address,
		Subject:    cmd.Subject,
		JoinDate:   cmd.JoinDate,
	}

	teacher, err := models.NewTeacher(data, ctx)
	if err != nil {
		return Failure(err)
	}

	return Success(teacher)
}

// DeactivateTeacher menonaktifkan guru.
func (e *TeacherEngine) DeactivateTeacher(reason string, ctx *context.OperationalContext, currentTeacher *models.Teacher) Result {
	if currentTeacher == nil {
		return Failure(errTeacherNotFound)
	}
	if err := currentTeacher.Deactivate(reason, ctx); err != nil {
		return Failure(err)
	}
	return Success(currentTeacher)
}

var errTeacherNotFound = errString("guru tidak ditemukan")

type errString string

func (e errString) Error() string { return string(e) }
