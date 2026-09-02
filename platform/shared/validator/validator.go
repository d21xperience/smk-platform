// Package validator menyediakan helper untuk validasi input.
package validator

import (
	"fmt"
	"regexp"
	"strings"

	"sekolah-platform/platform/shared/errors"
)

// ValidationError adalah kumpulan error validasi per field.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors adalah kumpulan validation errors.
type ValidationErrors []ValidationError

// Error mengimplementasikan error interface.
func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return ""
	}
	msgs := make([]string, len(ve))
	for i, e := range ve {
		msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return strings.Join(msgs, "; ")
}

// ToAppError mengkonversi ke AppError.
func (ve ValidationErrors) ToAppError() *errors.AppError {
	details := make(map[string]string)
	for _, e := range ve {
		details[e.Field] = e.Message
	}
	return errors.New(errors.CodeValidation, "Validasi gagal", details)
}

// Validator adalah builder untuk validasi.
type Validator struct {
	errs ValidationErrors
}

// New membuat Validator baru.
func New() *Validator {
	return &Validator{}
}

// Required memvalidasi field tidak kosong.
func (v *Validator) Required(field, value string) *Validator {
	if strings.TrimSpace(value) == "" {
		v.errs = append(v.errs, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s wajib diisi", field),
		})
	}
	return v
}

// MinLength memvalidasi panjang minimum.
func (v *Validator) MinLength(field, value string, min int) *Validator {
	if len(value) < min {
		v.errs = append(v.errs, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s minimal %d karakter", field, min),
		})
	}
	return v
}

// MaxLength memvalidasi panjang maksimum.
func (v *Validator) MaxLength(field, value string, max int) *Validator {
	if len(value) > max {
		v.errs = append(v.errs, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s maksimal %d karakter", field, max),
		})
	}
	return v
}

// Pattern memvalidasi dengan regex.
func (v *Validator) Pattern(field, value, pattern, message string) *Validator {
	if value == "" {
		return v // Skip jika kosong (gunakan Required untuk cek kosong)
	}
	matched, err := regexp.MatchString(pattern, value)
	if err != nil || !matched {
		v.errs = append(v.errs, ValidationError{
			Field:   field,
			Message: message,
		})
	}
	return v
}

// NISN memvalidasi format NISN (10 digit angka).
func (v *Validator) NISN(field, value string) *Validator {
	return v.Pattern(field, value, `^\d{10}$`, "NISN harus 10 digit angka")
}

// Email memvalidasi format email.
func (v *Validator) Email(field, value string) *Validator {
	if value == "" {
		return v
	}
	return v.Pattern(field, value, `^[^\s@]+@[^\s@]+\.[^\s@]+$`, "Format email tidak valid")
}

// OneOf memvalidasi value ada dalam daftar.
func (v *Validator) OneOf(field, value string, options []string) *Validator {
	if value == "" {
		v.errs = append(v.errs, ValidationError{
			Field:   field,
			Message: fmt.Sprintf("%s wajib dipilih", field),
		})
		return v
	}
	for _, opt := range options {
		if value == opt {
			return v
		}
	}
	v.errs = append(v.errs, ValidationError{
		Field:   field,
		Message: fmt.Sprintf("%s harus salah satu dari: %s", field, strings.Join(options, ", ")),
	})
	return v
}

// Validate mengembalikan error jika ada validation errors.
func (v *Validator) Validate() error {
	if len(v.errs) == 0 {
		return nil
	}
	return v.errs
}

// Errors mengembalikan validation errors (bisa kosong).
func (v *Validator) Errors() ValidationErrors {
	return v.errs
}

// IsValid mengecek apakah validasi lolos.
func (v *Validator) IsValid() bool {
	return len(v.errs) == 0
}
