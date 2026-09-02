package valueobjects

import "fmt"

// AssessmentType adalah value object untuk tipe penilaian.
type AssessmentType struct {
	code   string
	label  string
	weight float64 // Bobot nilai (0-100%)
}

var (
	AssessmentTypeAssignment = AssessmentType{code: "ASSIGNMENT", label: "Tugas", weight: 20}
	AssessmentTypeDaily      = AssessmentType{code: "DAILY", label: "Harian", weight: 20}
	AssessmentTypeMidterm    = AssessmentType{code: "MIDTERM", label: "UTS", weight: 30}
	AssessmentTypeFinal      = AssessmentType{code: "FINAL", label: "UAS", weight: 30}
)

// NewAssessmentType membuat AssessmentType dari code.
func NewAssessmentType(code string) (AssessmentType, error) {
	switch code {
	case "ASSIGNMENT":
		return AssessmentTypeAssignment, nil
	case "DAILY":
		return AssessmentTypeDaily, nil
	case "MIDTERM":
		return AssessmentTypeMidterm, nil
	case "FINAL":
		return AssessmentTypeFinal, nil
	default:
		return AssessmentType{}, fmt.Errorf("tipe penilaian tidak valid: %s", code)
	}
}

// Code mengembalikan code tipe.
func (t AssessmentType) Code() string { return t.code }

// Label mengembalikan label tipe.
func (t AssessmentType) Label() string { return t.label }

// Weight mengembalikan bobot tipe.
func (t AssessmentType) Weight() float64 { return t.weight }
