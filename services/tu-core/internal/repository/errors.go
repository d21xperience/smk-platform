package repository

import "errors"

var (
	// ErrNotFound error ketika data tidak ditemukan.
	ErrNotFound = errors.New("data tidak ditemukan")

	// ErrDuplicateEntry error ketika terjadi duplikasi data.
	ErrDuplicateEntry = errors.New("data sudah ada (duplikat)")

	// ErrInvalidInput error ketika input tidak valid.
	ErrInvalidInput = errors.New("input tidak valid")

	// ErrDatabaseError error ketika terjadi kesalahan database.
	ErrDatabaseError = errors.New("kesalahan database")

	// ErrOptimisticLock error ketika terjadi konflik versi data.
	ErrOptimisticLock = errors.New("konflik versi data (optimistic lock)")

	// ErrMissingOperationalContext error ketika SchoolID/AcademicPeriodID kosong.
	ErrMissingOperationalContext = errors.New("schoolId dan academicPeriodId wajib diisi")
)

// IsNotFound memeriksa apakah error adalah ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsDuplicate memeriksa apakah error adalah ErrDuplicateEntry.
func IsDuplicate(err error) bool {
	return errors.Is(err, ErrDuplicateEntry)
}
