package types

// EventVersion merepresentasikan versi event.
type EventVersion int

// Versi standar.
const (
	Version1 EventVersion = 1
	Version2 EventVersion = 2
	Version3 EventVersion = 3
)

// EventVersioningStrategy adalah strategi versioning event.
type EventVersioningStrategy string

const (
	// StrategyBackwardCompatible menambahkan field baru tanpa menghapus field lama.
	StrategyBackwardCompatible EventVersioningStrategy = "backward_compatible"
	// StrategyBreakingChange mengubah struktur event secara signifikan.
	StrategyBreakingChange EventVersioningStrategy = "breaking_change"
)

// VersionInfo berisi informasi versi event.
type VersionInfo struct {
	Name     string                  `json:"name"`
	Version  EventVersion            `json:"version"`
	Strategy EventVersioningStrategy `json:"strategy"`
}
