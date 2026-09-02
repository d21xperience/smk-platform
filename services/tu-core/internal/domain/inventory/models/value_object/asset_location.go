package value_object

import (
	"fmt"
	domainErrors "sekolah-platform/services/tu-core/internal/domain/inventory/errors"
)

type AssetLocation struct {
	BuildingID string
	FloorID    string
	RoomID     string
	Detail     string // opsional: "dekat jendela", "plafon", dll.
}

func NewAssetLocation(building, floor, room, detail string) (AssetLocation, error) {
	if building == "" || floor == "" || room == "" {
		return AssetLocation{}, domainErrors.ErrInvalidLocation
	}
	return AssetLocation{
		BuildingID: building,
		FloorID:    floor,
		RoomID:     room,
		Detail:     detail,
	}, nil
}

func (l AssetLocation) Equals(other AssetLocation) bool {
	return l.BuildingID == other.BuildingID &&
		l.FloorID == other.FloorID &&
		l.RoomID == other.RoomID
}

func (l AssetLocation) String() string {
	if l.Detail != "" {
		return fmt.Sprintf("%s - %s - %s (%s)", l.BuildingID, l.FloorID, l.RoomID, l.Detail)
	}
	return fmt.Sprintf("%s - %s - %s", l.BuildingID, l.FloorID, l.RoomID)
}
