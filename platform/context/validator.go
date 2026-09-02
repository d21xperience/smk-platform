package context

import (
	"sekolah-platform/platform/shared/errors"
)

// ValidateAndWrap memvalidasi OperationalContext dan mengembalikan AppError jika invalid.
func ValidateAndWrap(ctx *OperationalContext) *errors.AppError {
	if ctx == nil {
		return errors.InvalidContext("OperationalContext tidak boleh nil")
	}

	if err := ctx.Validate(); err != nil {
		return errors.InvalidContext(err.Error())
	}

	return nil
}

// MustGet mengambil OperationalContext atau panic jika nil/invalid.
// Gunakan hanya di tempat yang sudah pasti context ada (setelah middleware).
func MustGet(ctx *OperationalContext) *OperationalContext {
	if ctx == nil {
		panic("OperationalContext is nil")
	}
	if err := ctx.Validate(); err != nil {
		panic("OperationalContext is invalid: " + err.Error())
	}
	return ctx
}
