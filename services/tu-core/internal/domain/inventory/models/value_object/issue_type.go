package value_object

type IssueType string

const (
	DamageReport       IssueType = "DamageReport"
	MaintenanceRequest IssueType = "MaintenanceRequest"
	MissingAsset       IssueType = "MissingAsset"
	InventoryDiscovery IssueType = "InventoryDiscovery"
)

func (i IssueType) IsValid() bool {
	switch i {
	case DamageReport, MaintenanceRequest, MissingAsset, InventoryDiscovery:
		return true
	}
	return false
}
