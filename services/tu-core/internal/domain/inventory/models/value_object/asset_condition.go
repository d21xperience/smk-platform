package value_object

type AssetCondition string

const (
	ConditionGood        AssetCondition = "Good"
	ConditionFair        AssetCondition = "Fair"
	ConditionDamaged     AssetCondition = "Damaged"
	ConditionUnderRepair AssetCondition = "UnderRepair"
	ConditionDisposed    AssetCondition = "Disposed"
)

func (c AssetCondition) IsValid() bool {
	switch c {
	case ConditionGood, ConditionFair, ConditionDamaged, ConditionUnderRepair, ConditionDisposed:
		return true
	}
	return false
}

// CanTransitionTo mengecek apakah transisi kondisi diperbolehkan
func (c AssetCondition) CanTransitionTo(newCondition AssetCondition) bool {
	if c == ConditionDisposed {
		return false // disposed tidak bisa berubah
	}
	if c == ConditionUnderRepair && newCondition == ConditionGood {
		return true
	}
	if c == ConditionDamaged && (newCondition == ConditionUnderRepair || newCondition == ConditionDisposed) {
		return true
	}
	if c == ConditionGood && newCondition == ConditionDamaged {
		return true
	}
	if c == ConditionGood && newCondition == ConditionFair {
		return true
	}
	if c == ConditionFair && newCondition == ConditionGood {
		return true
	}
	if c == ConditionFair && newCondition == ConditionDamaged {
		return true
	}
	return false
}
