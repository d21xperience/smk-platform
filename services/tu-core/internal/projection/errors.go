package projection

import "errors"

var (
	// ErrProjectionNotFound error ketika projection handler tidak ditemukan.
	ErrProjectionNotFound = errors.New("projection handler tidak ditemukan")

	// ErrProjectionAlreadyRegistered error ketika projection sudah terdaftar.
	ErrProjectionAlreadyRegistered = errors.New("projection sudah terdaftar")

	// ErrInvalidEventPayload error ketika payload event tidak valid.
	ErrInvalidEventPayload = errors.New("payload event tidak valid")

	// ErrProjectionFailed error ketika eksekusi projection gagal.
	ErrProjectionFailed = errors.New("eksekusi projection gagal")
)
