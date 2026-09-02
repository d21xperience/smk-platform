package valueobjects

import (
	"errors"
	"fmt"
)

// Amount adalah value object untuk jumlah uang (dalam Rupiah, tanpa desimal).
type Amount struct {
	value int64 // dalam Rupiah (e.g., 500000 = Rp 500.000)
}

// NewAmount membuat Amount baru dengan validasi.
func NewAmount(value int64) (Amount, error) {
	if value < 0 {
		return Amount{}, errors.New("jumlah tidak boleh negatif")
	}
	if value > 1_000_000_000 { // Max 1 Miliar
		return Amount{}, fmt.Errorf("jumlah melebihi batas maksimal: %d", value)
	}
	return Amount{value: value}, nil
}

// Value mengembalikan nilai dalam Rupiah.
func (a Amount) Value() int64 { return a.value }

// IsZero mengecek apakah jumlah nol.
func (a Amount) IsZero() bool { return a.value == 0 }

// Add menambahkan dua Amount.
func (a Amount) Add(other Amount) Amount {
	return Amount{value: a.value + other.value}
}

// Subtract mengurangkan Amount.
func (a Amount) Subtract(other Amount) (Amount, error) {
	if a.value < other.value {
		return Amount{}, errors.New("hasil pengurangan negatif")
	}
	return Amount{value: a.value - other.value}, nil
}

// GreaterThanOrEqual mengecek apakah a >= other.
func (a Amount) GreaterThanOrEqual(other Amount) bool {
	return a.value >= other.value
}

// Format mengembalikan format Rupiah (e.g., "Rp 500.000").
func (a Amount) Format() string {
	// Simplified formatting
	return fmt.Sprintf("Rp %d", a.value)
}

// Equals mengecek kesamaan.
func (a Amount) Equals(other Amount) bool { return a.value == other.value }
