package service

import "errors"

var (
	// ErrInvalidInput error ketika input tidak valid.
	ErrInvalidInput = errors.New("input tidak valid")

	// ErrNotFound error ketika data tidak ditemukan.
	ErrNotFound = errors.New("data tidak ditemukan")

	// ErrDuplicateEntry error ketika data sudah ada.
	ErrDuplicateEntry = errors.New("data sudah ada (duplikat)")

	// ErrMissingOperationalContext error ketika SchoolID/PeriodID kosong.
	ErrMissingOperationalContext = errors.New("schoolId dan academicPeriodId wajib diisi")

	// ErrUnauthorized error ketika user tidak memiliki akses.
	ErrUnauthorized = errors.New("akses ditolak")

	// ErrBusinessRuleViolation error ketika aturan bisnis dilanggar.
	ErrBusinessRuleViolation = errors.New("aturan bisnis dilanggar")

	// ErrEventPublishFailed error ketika gagal publish event.
	ErrEventPublishFailed = errors.New("gagal publish event")
)

// ServiceError adalah error wrapper dengan konteks tambahan.
type ServiceError struct {
	Op      string // Operasi yang gagal
	Entity  string // Entitas yang terlibat
	ID      string // ID entitas
	Message string
	Err     error
}

func (e *ServiceError) Error() string {
	if e.ID != "" {
		return "service error [" + e.Op + " " + e.Entity + ":" + e.ID + "]: " + e.Message
	}
	return "service error [" + e.Op + " " + e.Entity + "]: " + e.Message
}

func (e *ServiceError) Unwrap() error {
	return e.Err
}

// NewServiceError membuat error baru dengan konteks.
func NewServiceError(op, entity, id, message string, err error) *ServiceError {
	return &ServiceError{
		Op:      op,
		Entity:  entity,
		ID:      id,
		Message: message,
		Err:     err,
	}
}
