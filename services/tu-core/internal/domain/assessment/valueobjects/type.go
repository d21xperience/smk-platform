package valueobjects

import "errors"

type AssessmentType string

const (
	TypeDaily   AssessmentType = "DAILY"
	TypeMidterm AssessmentType = "MIDTERM"
	TypeFinal   AssessmentType = "FINAL"
)

func NewAssessmentType(t string) (AssessmentType, error) {
	switch AssessmentType(t) {
	case TypeDaily, TypeMidterm, TypeFinal:
		return AssessmentType(t), nil
	default:
		return "", errors.New("invalid assessment type: must be DAILY, MIDTERM, or FINAL")
	}
}

func (t AssessmentType) String() string {
	return string(t)
}

// CalculateGrade menghitung predikat berdasarkan persentase nilai
func CalculateGrade(percentage float64) string {
	switch {
	case percentage >= 90:
		return "A"
	case percentage >= 80:
		return "B"
	case percentage >= 70:
		return "C"
	case percentage >= 60:
		return "D"
	default:
		return "E"
	}
}
