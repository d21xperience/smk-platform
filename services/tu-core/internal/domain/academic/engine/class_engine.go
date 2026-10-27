package engine

import (
	"errors"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/academic/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/academic/models"
)

type ClassEngine struct{}

func NewClassEngine() *ClassEngine {
	return &ClassEngine{}
}

// ClassResult adalah hasil operasi engine.
type ClassResult struct {
	Class  *models.Class
	Events []types.DomainEvent
	Error  error
}

func successClass(c *models.Class) ClassResult {
	return ClassResult{
		Class:  c,
		Events: c.UncommittedEvents(),
		Error:  nil,
	}
}

func failureClass(err error) ClassResult {
	return ClassResult{Error: err}
}

// CreateClass membuat class baru.
func (e *ClassEngine) CreateClass(cmd commands.CreateClassCommand, ctx *context.OperationalContext) ClassResult {
	if err := cmd.Validate(); err != nil {
		return failureClass(err)
	}

	data := models.ClassData{
		ID:         cmd.ClassID,
		SchoolID:   ctx.SchoolID,
		PeriodID:   ctx.AcademicPeriodID,
		Name:       cmd.Name,
		Grade:      cmd.Grade,
		Level:      cmd.Level,
		Major:      cmd.Major,
		HomeroomID: cmd.HomeroomID,
		Capacity:   cmd.Capacity,
	}
	class, err := models.NewClass(data, ctx)
	if err != nil {
		return failureClass(err)
	}
	return successClass(class)
}

// UpdateClass memperbarui class.
func (e *ClassEngine) UpdateClass(cmd commands.UpdateClassCommand, current *models.Class, ctx *context.OperationalContext) ClassResult {
	if current == nil {
		return failureClass(errClassNotFound)
	}
	if err := cmd.Validate(); err != nil {
		return failureClass(err)
	}
	if err := current.UpdateClass(
		cmd.Name,
		cmd.Grade,
		cmd.Level,
		cmd.Major,
		cmd.HomeroomID,
		cmd.Capacity,
		cmd.Status,
		ctx,
	); err != nil {
		return failureClass(err)
	}
	return successClass(current)
}

// SetHomeroom mengatur wali kelas.
func (e *ClassEngine) SetHomeroom(classID, teacherID string, current *models.Class, ctx *context.OperationalContext) ClassResult {
	if current == nil {
		return failureClass(errClassNotFound)
	}
	if classID == "" || teacherID == "" {
		return failureClass(errors.New("classID dan teacherID wajib diisi"))
	}
	if err := current.SetHomeroom(teacherID, ctx); err != nil {
		return failureClass(err)
	}
	return successClass(current)
}

// CreateClass membuat kelas baru.
func (e *AcademicEngine) CreateClass(data models.ClassData, ctx *context.OperationalContext) Result {
	c, err := models.NewClass(data, ctx)
	if err != nil {
		return Failure(err)
	}
	return SuccessWithClass(c)
}

// SetClassHomeroom mengatur wali kelas.
func (e *AcademicEngine) SetClassHomeroom(c *models.Class, teacherID string, ctx *context.OperationalContext) Result {
	if c == nil {
		return Failure(errClassNotFound)
	}
	if err := c.SetHomeroom(teacherID, ctx); err != nil {
		return Failure(err)
	}
	return SuccessWithClass(c)
}
