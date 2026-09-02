package engine

import (
	"time"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/academic/models"
)

// Result adalah hasil dari engine operation.
type Result struct {
	AcademicYear *models.AcademicYear
	Class        *models.Class
	Semester     *models.Semester
	Events       []types.DomainEvent
	Error        error
}

func SuccessWithYear(ay *models.AcademicYear) Result {
	return Result{AcademicYear: ay, Events: ay.UncommittedEvents()}
}

func SuccessWithClass(c *models.Class) Result {
	return Result{Class: c, Events: c.UncommittedEvents()}
}

func Failure(err error) Result {
	return Result{Error: err}
}

// AcademicEngine adalah business engine untuk domain Academic.
type AcademicEngine struct{}

func NewAcademicEngine() *AcademicEngine {
	return &AcademicEngine{}
}

// CreateAcademicYear membuat tahun ajaran baru.
func (e *AcademicEngine) CreateAcademicYear(startYear, endYear int, startDate, endDate time.Time, isActive bool, ctx *context.OperationalContext) Result {
	data := models.AcademicYearData{
		SchoolID:  ctx.SchoolID,
		StartYear: startYear,
		EndYear:   endYear,
		StartDate: startDate,
		EndDate:   endDate,
		IsActive:  isActive,
	}

	ay, err := models.NewAcademicYear(data, ctx)
	if err != nil {
		return Failure(err)
	}

	return SuccessWithYear(ay)
}

// AddSemester menambahkan semester ke tahun ajaran.
func (e *AcademicEngine) AddSemester(ay *models.AcademicYear, number int, startDate, endDate time.Time, ctx *context.OperationalContext) Result {
	if ay == nil {
		return Failure(errAcademicYearNotFound)
	}

	semester, err := ay.AddSemester(number, startDate, endDate, ctx)
	if err != nil {
		return Failure(err)
	}

	return Result{
		AcademicYear: ay,
		Semester:     semester,
		Events:       ay.UncommittedEvents(),
	}
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

var (
	errAcademicYearNotFound = errString("tahun ajaran tidak ditemukan")
	errClassNotFound        = errString("kelas tidak ditemukan")
)

type errString string

func (e errString) Error() string { return string(e) }
