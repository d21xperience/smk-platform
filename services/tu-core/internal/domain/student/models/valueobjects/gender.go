package valueobjects

import "fmt"

type Gender struct {
	code  string
	label string
}

var (
	GenderMale   = Gender{code: "MALE", label: "Laki-laki"}
	GenderFemale = Gender{code: "FEMALE", label: "Perempuan"}
)

// NewGender membuat Gender dari code.
// Mendukung format DAPODIK ("L", "P") dan format Domain ("MALE", "FEMALE").
func NewGender(code string) (Gender, error) {
	switch code {
	case "MALE", "L", "Laki-laki", "1":
		return GenderMale, nil
	case "FEMALE", "P", "Perempuan", "2":
		return GenderFemale, nil
	default:
		return Gender{}, fmt.Errorf("gender tidak valid: %s", code)
	}
}

func (g Gender) Code() string { return g.code }
func (g Gender) Label() string { return g.label }
func (g Gender) Equals(other Gender) bool { return g.code == other.code }
