package value_object

type OwnerType string

const (
	OwnerSchool     OwnerType = "School"
	OwnerTeacher    OwnerType = "Teacher"
	OwnerDepartment OwnerType = "Department"
)

type AssetOwnership struct {
	OwnerID   string
	OwnerType OwnerType
}

func NewAssetOwnership(ownerID string, ownerType OwnerType) AssetOwnership {
	return AssetOwnership{
		OwnerID:   ownerID,
		OwnerType: ownerType,
	}
}
