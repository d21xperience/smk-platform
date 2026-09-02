package valueobjects

import "fmt"

// Gender adalah value object untuk jenis kelamin.
type Gender struct {
	code  string
	label string
}

var (
	GenderMale   = Gender{code: "MALE", label: "Laki-laki"}
	GenderFemale = Gender{code: "FEMALE", label: "Perempuan"}
)

// NewGender membuat Gender dari code.
func NewGender(code string) (Gender, error) {
	switch code {
	case "MALE":
		return GenderMale, nil
	case "FEMALE":
		return GenderFemale, nil
	default:
		return Gender{}, fmt.Errorf("gender tidak valid: %s", code)
	}
}

// Code mengembalikan code gender.
func (g Gender) Code() string { return g.code }

// Label mengembalikan label gender.
func (g Gender) Label() string { return g.label }

// Equals mengecek kesamaan.
func (g Gender) Equals(other Gender) bool { return g.code == other.code }
