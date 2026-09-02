package event

import "time"

// ---------- base event ----------
type BaseEvent struct {
	OccurredAt  time.Time
	AggregateID string
}

func (b BaseEvent) GetOccurredAt() time.Time { return b.OccurredAt }
func (b BaseEvent) GetAggregateID() string   { return b.AggregateID }

// ---------- Inventory Events ----------
type InventoryCreated struct {
	BaseEvent
	Name         string
	QRCode       string
	SerialNumber string
	Location     string // string representation
	Condition    string
}

func (e InventoryCreated) GetEventName() string    { return "InventoryCreated" }
func (e InventoryCreated) GetPayload() interface{} { return e }

type InventoryUpdated struct {
	BaseEvent
	Changes map[string]interface{}
}

func (e InventoryUpdated) GetEventName() string    { return "InventoryUpdated" }
func (e InventoryUpdated) GetPayload() interface{} { return e }

type InventoryMoved struct {
	BaseEvent
	OldLocation string
	NewLocation string
}

func (e InventoryMoved) GetEventName() string    { return "InventoryMoved" }
func (e InventoryMoved) GetPayload() interface{} { return e }

type InventoryDamaged struct {
	BaseEvent
	Condition string
	Reason    string
}

func (e InventoryDamaged) GetEventName() string    { return "InventoryDamaged" }
func (e InventoryDamaged) GetPayload() interface{} { return e }

type InventoryBorrowed struct {
	BaseEvent
	BorrowerID string
	BorrowedAt time.Time
}

func (e InventoryBorrowed) GetEventName() string    { return "InventoryBorrowed" }
func (e InventoryBorrowed) GetPayload() interface{} { return e }

type InventoryReturned struct {
	BaseEvent
	ReturnedAt time.Time
	Condition  string
}

func (e InventoryReturned) GetEventName() string    { return "InventoryReturned" }
func (e InventoryReturned) GetPayload() interface{} { return e }

type InventoryMaintenanceRequested struct {
	BaseEvent
	IssueID string
}

func (e InventoryMaintenanceRequested) GetEventName() string    { return "InventoryMaintenanceRequested" }
func (e InventoryMaintenanceRequested) GetPayload() interface{} { return e }

type InventoryMaintenanceCompleted struct {
	BaseEvent
	CompletedAt time.Time
}

func (e InventoryMaintenanceCompleted) GetEventName() string    { return "InventoryMaintenanceCompleted" }
func (e InventoryMaintenanceCompleted) GetPayload() interface{} { return e }

type InventoryDisposed struct {
	BaseEvent
	DisposedAt time.Time
	Reason     string
}

func (e InventoryDisposed) GetEventName() string    { return "InventoryDisposed" }
func (e InventoryDisposed) GetPayload() interface{} { return e }

type InventoryDiscoveryReported struct {
	BaseEvent
	FoundAt     time.Time
	Description string
}

func (e InventoryDiscoveryReported) GetEventName() string    { return "InventoryDiscoveryReported" }
func (e InventoryDiscoveryReported) GetPayload() interface{} { return e }

// ---------- Asset Issue Events ----------
type AssetIssueCreated struct {
	BaseEvent
	IssueType   string
	ReporterID  string
	RoomID      string
	Description string
}

func (e AssetIssueCreated) GetEventName() string    { return "AssetIssueCreated" }
func (e AssetIssueCreated) GetPayload() interface{} { return e }

type AssetIssueVerified struct {
	BaseEvent
	VerifiedBy string
	VerifiedAt time.Time
}

func (e AssetIssueVerified) GetEventName() string    { return "AssetIssueVerified" }
func (e AssetIssueVerified) GetPayload() interface{} { return e }

type AssetIssueClosed struct {
	BaseEvent
	ClosedAt   time.Time
	Resolution string
}

func (e AssetIssueClosed) GetEventName() string    { return "AssetIssueClosed" }
func (e AssetIssueClosed) GetPayload() interface{} { return e }
