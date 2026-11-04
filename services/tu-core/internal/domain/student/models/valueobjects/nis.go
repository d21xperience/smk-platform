package valueobjects

type NIS struct {
	value string
}

// NewNIS membuat NIS baru.
// DILONGGARKAN: Mengizinkan string kosong agar tidak crash saat migrasi data lama.
func NewNIS(value string) (NIS, error) {
	return NIS{value: value}, nil
}

func (n NIS) Value() string { return n.value }
func (n NIS) String() string { return n.value }
func (n NIS) Equals(other NIS) bool { return n.value == other.value }
