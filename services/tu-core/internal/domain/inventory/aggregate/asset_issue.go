package aggregate

import (
	"time"

	domainErrors "sekolah-platform/services/tu-core/internal/domain/inventory/errors"
	"sekolah-platform/services/tu-core/internal/domain/inventory/event"
	"sekolah-platform/services/tu-core/internal/domain/inventory/models/value_object"
)

type AssetIssueStatus string

const (
	StatusReported   AssetIssueStatus = "Reported"
	StatusVerified   AssetIssueStatus = "Verified"
	StatusAssigned   AssetIssueStatus = "Assigned"
	StatusInProgress AssetIssueStatus = "InProgress"
	StatusClosed     AssetIssueStatus = "Closed"
)

type AssetIssue struct {
	id              string
	inventoryItemID *string // nullable (untuk discovery, bisa kosong)
	reporterID      string
	roomID          string
	issueType       value_object.IssueType
	description     string
	photos          []string // URL
	status          AssetIssueStatus
	assignedTo      *string // technician ID
	verifiedBy      *string
	verifiedAt      *time.Time
	closedAt        *time.Time
	resolution      *string
	createdAt       time.Time
	updatedAt       time.Time

	events []event.DomainEvent
}

// ----- Constructor -----
func NewAssetIssue(
	id string,
	inventoryItemID *string,
	reporterID string,
	roomID string,
	issueType value_object.IssueType,
	description string,
	photos []string,
) (*AssetIssue, error) {
	if id == "" || reporterID == "" || roomID == "" {
		return nil, domainErrors.ErrInvalidReporter
	}
	if description == "" {
		return nil, domainErrors.ErrEmptyIssueDescription
	}
	if !issueType.IsValid() {
		return nil, domainErrors.ErrInvalidIssueStatus
	}

	now := time.Now()
	issue := &AssetIssue{
		id:              id,
		inventoryItemID: inventoryItemID,
		reporterID:      reporterID,
		roomID:          roomID,
		issueType:       issueType,
		description:     description,
		photos:          photos,
		status:          StatusReported,
		assignedTo:      nil,
		verifiedBy:      nil,
		verifiedAt:      nil,
		closedAt:        nil,
		resolution:      nil,
		createdAt:       now,
		updatedAt:       now,
		events:          []event.DomainEvent{},
	}

	issue.addEvent(event.AssetIssueCreated{
		BaseEvent: event.BaseEvent{
			OccurredAt:  now,
			AggregateID: id,
		},
		IssueType:   string(issueType),
		ReporterID:  reporterID,
		RoomID:      roomID,
		Description: description,
	})
	return issue, nil
}

// ----- Getter -----
func (a *AssetIssue) GetID() string                  { return a.id }
func (a *AssetIssue) GetStatus() AssetIssueStatus    { return a.status }
func (a *AssetIssue) GetEvents() []event.DomainEvent { return a.events }

// ----- Behaviour -----
func (a *AssetIssue) Verify(verifierID string) error {
	if a.status == StatusClosed {
		return domainErrors.ErrIssueAlreadyClosed
	}
	if a.status != StatusReported {
		return domainErrors.ErrInvalidIssueStatus
	}
	now := time.Now()
	a.status = StatusVerified
	a.verifiedBy = &verifierID
	a.verifiedAt = &now
	a.updatedAt = now
	a.addEvent(event.AssetIssueVerified{
		BaseEvent: event.BaseEvent{
			OccurredAt:  now,
			AggregateID: a.id,
		},
		VerifiedBy: verifierID,
		VerifiedAt: now,
	})
	return nil
}

func (a *AssetIssue) Assign(technicianID string) error {
	if a.status == StatusClosed {
		return domainErrors.ErrIssueAlreadyClosed
	}
	if a.status != StatusVerified && a.status != StatusAssigned {
		return domainErrors.ErrInvalidIssueStatus
	}
	a.status = StatusAssigned
	a.assignedTo = &technicianID
	a.updatedAt = time.Now()
	// tidak wajib event, bisa ditambahkan jika perlu
	return nil
}

func (a *AssetIssue) StartWork() error {
	if a.status == StatusClosed {
		return domainErrors.ErrIssueAlreadyClosed
	}
	if a.status != StatusAssigned {
		return domainErrors.ErrInvalidIssueStatus
	}
	a.status = StatusInProgress
	a.updatedAt = time.Now()
	return nil
}

func (a *AssetIssue) Close(resolution string) error {
	if a.status == StatusClosed {
		return domainErrors.ErrIssueAlreadyClosed
	}
	if a.status != StatusInProgress && a.status != StatusAssigned {
		return domainErrors.ErrInvalidIssueStatus
	}
	now := time.Now()
	a.status = StatusClosed
	a.closedAt = &now
	a.resolution = &resolution
	a.updatedAt = now
	a.addEvent(event.AssetIssueClosed{
		BaseEvent: event.BaseEvent{
			OccurredAt:  now,
			AggregateID: a.id,
		},
		ClosedAt:   now,
		Resolution: resolution,
	})
	return nil
}

func (a *AssetIssue) addEvent(e event.DomainEvent) {
	a.events = append(a.events, e)
}

func (a *AssetIssue) GetReporterID() string                { return a.reporterID }
func (a *AssetIssue) GetRoomID() string                    { return a.roomID }
func (a *AssetIssue) GetIssueType() value_object.IssueType { return a.issueType }
func (a *AssetIssue) GetDescription() string               { return a.description }

// func (a *AssetIssue) GetStatus() AssetIssueStatus          { return a.status }
func (a *AssetIssue) GetAssignedTo() *string    { return a.assignedTo }
func (a *AssetIssue) GetVerifiedBy() *string    { return a.verifiedBy }
func (a *AssetIssue) GetVerifiedAt() *time.Time { return a.verifiedAt }
func (a *AssetIssue) GetClosedAt() *time.Time   { return a.closedAt }
func (a *AssetIssue) GetResolution() *string    { return a.resolution }
func (a *AssetIssue) GetInventoryItemID() *string {
	return a.inventoryItemID
}
