package search

import "errors"

var (
	// ErrMissingOperationalContext error ketika SchoolID atau AcademicPeriodID kosong.
	ErrMissingOperationalContext = errors.New("schoolId dan academicPeriodId wajib diisi")

	// ErrInvalidQuery error ketika query pencarian tidak valid.
	ErrInvalidQuery = errors.New("query pencarian tidak valid")

	// ErrSearchFailed error ketika pencarian gagal.
	ErrSearchFailed = errors.New("pencarian gagal")
)
