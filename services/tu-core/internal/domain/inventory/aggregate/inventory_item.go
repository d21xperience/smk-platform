package aggregate

import (
	"time"

	domainErrors "sekolah-platform/services/tu-core/internal/domain/inventory/errors"
	"sekolah-platform/services/tu-core/internal/domain/inventory/event"
	"sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
)

type InventoryItem struct {
	identity        value_object.AssetIdentity
	condition       value_object.AssetCondition
	location        value_object.AssetLocation
	ownership       value_object.AssetOwnership
	specification   value_object.AssetSpecification
	procurementInfo ProcurementInformation // kita definisikan nanti (sederhana)
	// histories (value object slices)
	maintenanceHistory []MaintenanceHistory
	loanHistory        []LoanHistory
	mutationHistory    []MutationHistory
	disposalHistory    []DisposalHistory

	events []event.DomainEvent
}

// ProcurementInformation (sederhana) - bisa ditempatkan di value_object/procurement.go
// Untuk Sprint INV-A kita buat inline dulu
type ProcurementInformation struct {
	Vendor     string
	Price      float64
	InvoiceNo  string
	ReceivedAt time.Time
}

type MaintenanceHistory struct {
	Date         time.Time
	Description  string
	TechnicianID string
}

type LoanHistory struct {
	BorrowerID string
	BorrowedAt time.Time
	ReturnedAt *time.Time
}

type MutationHistory struct {
	FromLocation string
	ToLocation   string
	MutatedAt    time.Time
	Reason       string
}

type DisposalHistory struct {
	DisposedAt time.Time
	Reason     string
	ApprovedBy string
}

// ----- Constructor -----
func NewInventoryItem(
	identity value_object.AssetIdentity,
	location value_object.AssetLocation,
	ownership value_object.AssetOwnership,
	spec value_object.AssetSpecification,
	procurement ProcurementInformation,
) (*InventoryItem, error) {
	if identity.ID == "" {
		return nil, domainErrors.ErrInvalidAssetName // placeholder
	}
	item := &InventoryItem{
		identity:           identity,
		condition:          value_object.ConditionGood, // default
		location:           location,
		ownership:          ownership,
		specification:      spec,
		procurementInfo:    procurement,
		maintenanceHistory: []MaintenanceHistory{},
		loanHistory:        []LoanHistory{},
		mutationHistory:    []MutationHistory{},
		disposalHistory:    []DisposalHistory{},
		events:             []event.DomainEvent{},
	}
	// Generate event
	item.addEvent(event.InventoryCreated{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: identity.ID,
		},
		Name:         identity.Name,
		QRCode:       identity.QRCode,
		SerialNumber: identity.SerialNumber,
		Location:     location.String(),
		Condition:    string(item.condition),
	})
	return item, nil
}

// ----- Getter / Accessor -----
func (i *InventoryItem) GetID() string                             { return i.identity.ID }
func (i *InventoryItem) GetCondition() value_object.AssetCondition { return i.condition }
func (i *InventoryItem) GetLocation() value_object.AssetLocation   { return i.location }
func (i *InventoryItem) GetEvents() []event.DomainEvent            { return i.events }

// ----- Domain Methods (Behaviour) -----

// MoveTo mengubah lokasi aset, menghasilkan event InventoryMoved
func (i *InventoryItem) MoveTo(newLocation value_object.AssetLocation) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrAssetAlreadyDisposed
	}
	if i.location.Equals(newLocation) {
		return nil // tidak ada perubahan
	}
	oldLocStr := i.location.String()
	i.location = newLocation
	i.addEvent(event.InventoryMoved{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: i.identity.ID,
		},
		OldLocation: oldLocStr,
		NewLocation: newLocation.String(),
	})
	return nil
}

// UpdateCondition mengubah kondisi aset, dengan validasi transisi
func (i *InventoryItem) UpdateCondition(newCond value_object.AssetCondition, reason string) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrAssetAlreadyDisposed
	}
	if !i.condition.CanTransitionTo(newCond) {
		return domainErrors.ErrInvalidConditionChange
	}
	// jika berubah menjadi Damaged, publish event Damaged
	if newCond == value_object.ConditionDamaged {
		i.addEvent(event.InventoryDamaged{
			BaseEvent: event.BaseEvent{
				OccurredAt:  time.Now(),
				AggregateID: i.identity.ID,
			},
			Condition: string(newCond),
			Reason:    reason,
		})
	}
	i.condition = newCond
	return nil
}

// Borrow mencatat peminjaman
func (i *InventoryItem) Borrow(borrowerID string) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrCannotBorrowDisposed
	}
	if i.condition == value_object.ConditionDamaged {
		return domainErrors.ErrCannotBorrowDamaged
	}
	// tambah ke loanHistory
	i.loanHistory = append(i.loanHistory, LoanHistory{
		BorrowerID: borrowerID,
		BorrowedAt: time.Now(),
		ReturnedAt: nil,
	})
	i.addEvent(event.InventoryBorrowed{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: i.identity.ID,
		},
		BorrowerID: borrowerID,
		BorrowedAt: time.Now(),
	})
	return nil
}

// Return mengembalikan aset
func (i *InventoryItem) Return(conditionAfterUse string) error {
	if len(i.loanHistory) == 0 || i.loanHistory[len(i.loanHistory)-1].ReturnedAt != nil {
		return domainErrors.ErrInvalidConditionChange // tidak ada pinjaman aktif
	}
	now := time.Now()
	i.loanHistory[len(i.loanHistory)-1].ReturnedAt = &now
	i.addEvent(event.InventoryReturned{
		BaseEvent: event.BaseEvent{
			OccurredAt:  now,
			AggregateID: i.identity.ID,
		},
		ReturnedAt: now,
		Condition:  conditionAfterUse,
	})
	return nil
}

// Dispose membuang aset
func (i *InventoryItem) Dispose(reason, approvedBy string) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrAssetAlreadyDisposed
	}
	i.condition = value_object.ConditionDisposed
	i.disposalHistory = append(i.disposalHistory, DisposalHistory{
		DisposedAt: time.Now(),
		Reason:     reason,
		ApprovedBy: approvedBy,
	})
	i.addEvent(event.InventoryDisposed{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: i.identity.ID,
		},
		DisposedAt: time.Now(),
		Reason:     reason,
	})
	return nil
}

// RequestMaintenance memicu event maintenance request (tanpa mengubah kondisi langsung)
func (i *InventoryItem) RequestMaintenance(issueID string) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrAssetAlreadyDisposed
	}
	i.addEvent(event.InventoryMaintenanceRequested{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: i.identity.ID,
		},
		IssueID: issueID,
	})
	return nil
}

// CompleteMaintenance mengubah kondisi setelah maintenance selesai
func (i *InventoryItem) CompleteMaintenance(newCondition value_object.AssetCondition) error {
	if i.condition == value_object.ConditionDisposed {
		return domainErrors.ErrAssetAlreadyDisposed
	}
	if newCondition != value_object.ConditionGood && newCondition != value_object.ConditionFair {
		return domainErrors.ErrInvalidConditionChange
	}
	i.condition = newCondition
	i.addEvent(event.InventoryMaintenanceCompleted{
		BaseEvent: event.BaseEvent{
			OccurredAt:  time.Now(),
			AggregateID: i.identity.ID,
		},
		CompletedAt: time.Now(),
	})
	return nil
}

// ----- Event Sourcing helper -----
func (i *InventoryItem) addEvent(e event.DomainEvent) {
	i.events = append(i.events, e)
}
