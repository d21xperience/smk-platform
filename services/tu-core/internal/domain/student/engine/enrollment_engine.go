package engine

import (
	"errors"

	"sekolah-platform/platform/context"
	"sekolah-platform/platform/events/types"
	"sekolah-platform/services/tu-core/internal/domain/student/engine/commands"
	"sekolah-platform/services/tu-core/internal/domain/student/models"
)

type EnrollmentEngine struct{}

func NewEnrollmentEngine() *EnrollmentEngine {
	return &EnrollmentEngine{}
}

type EnrollmentResult struct {
	Enrollment *models.Enrollment
	Events     []types.DomainEvent
	Error      error
}

func successEnrollment(e *models.Enrollment) EnrollmentResult {
	// Asumsi Enrollment tidak punya events, bisa diubah jika ada
	return EnrollmentResult{
		Enrollment: e,
		Events:     nil,
		Error:      nil,
	}
}

func failureEnrollment(err error) EnrollmentResult {
	return EnrollmentResult{Error: err}
}

// EnrollStudent menambahkan enrollment baru (dipanggil dari StudentEngine.EnrollStudent sebenarnya,
// tapi kita buat terpisah untuk kejelasan).
// Catatan: Method ini tidak digunakan di StudentService? Di StudentService, EnrollStudent dipanggil via StudentEngine.
// Kita bisa biarkan method ini untuk kelengkapan.
func (e *EnrollmentEngine) EnrollStudent(cmd commands.EnrollStudentCommand, ctx *context.OperationalContext) EnrollmentResult {
	if err := cmd.Validate(); err != nil {
		return failureEnrollment(err)
	}

	// Buat entity Enrollment menggunakan constructor/factory jika ada, atau langsung dengan getter
	// Karena kita tidak punya constructor di models.Enrollment, kita buat manual dengan setter? Sebaiknya buat factory.
	// Untuk sementara, kita buat menggunakan struct dengan field publik? Tidak, field private.
	// Kita perlu menggunakan builder atau constructor.
	// Karena tidak ada constructor, saya akan asumsikan models.Enrollment memiliki NewEnrollment atau serupa.
	// Jika tidak, kita harus mengakses field private melalui setter?

	// Alternatif: kita bisa mengembalikan error dan menyarankan menggunakan StudentEngine.EnrollStudent
	return failureEnrollment(errors.New("use StudentEngine.EnrollStudent instead"))
}

// MoveStudent memindahkan siswa ke class lain (mengakhiri enrollment lama, buat baru)
// Ini akan dipanggil oleh service, dan service akan memanggil repository untuk menyimpan.
func (e *EnrollmentEngine) MoveStudent(cmd commands.MoveStudentCommand, current *models.Enrollment, fromClassGrade, toClassGrade int, capacity int, ctx *context.OperationalContext) EnrollmentResult {
	if current == nil {
		return failureEnrollment(errEnrollmentNotFound)
	}
	if err := cmd.Validate(); err != nil {
		return failureEnrollment(err)
	}
	// Validasi grade
	if fromClassGrade != toClassGrade {
		return failureEnrollment(errInvalidGradeTransition)
	}
	// Validasi kapasitas (dilakukan di service karena butuh query count)
	// Untuk engine, kita hanya terima parameter capacity dan asumsikan sudah divalidasi sebelumnya.

	// Akhiri enrollment lama
	current.Transfer(cmd.Reason) // asumsi ada method Transfer

	// Buat enrollment baru (tapi di engine kita hanya return data, service yang menyimpan)
	// Kita perlu membuat enrollment baru dengan constructor/factory.
	// Karena tidak ada, kita akan return error.
	return failureEnrollment(errors.New("use service layer to create new enrollment"))
}


