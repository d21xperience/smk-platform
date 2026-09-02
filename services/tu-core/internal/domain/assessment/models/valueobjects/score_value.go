package valueobjects

import (
	"errors"
	"fmt"
)

// ScoreValue adalah value object untuk nilai (0-100).
type ScoreValue struct {
	value float64
}

// NewScoreValue membuat ScoreValue baru dengan validasi.
func NewScoreValue(value float64) (ScoreValue, error) {
	if value < 0 || value > 100 {
		return ScoreValue{}, fmt.Errorf("nilai harus antara 0-100, got: %.2f", value)
	}
	return ScoreValue{value: value}, nil
}

// Value mengembalikan nilai.
func (s ScoreValue) Value() float64 { return s.value }

// IsPassing mengecek apakah nilai lulus (>= 75).
func (s ScoreValue) IsPassing() bool { return s.value >= 75 }

// Grade mengembalikan grade (A, B, C, D, E).
func (s ScoreValue) Grade() string {
	switch {
	case s.value >= 90:
		return "A"
	case s.value >= 80:
		return "B"
	case s.value >= 70:
		return "C"
	case s.value >= 60:
		return "D"
	default:
		return "E"
	}
}

// Equals mengecek kesamaan.
func (s ScoreValue) Equals(other ScoreValue) bool {
	return s.value == other.value
}

// Placeholder untuk errors
var _ = errors.New
