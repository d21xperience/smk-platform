package models

import "gorm.io/gorm"

type InventoryItemModel struct {
	ID             string `gorm:"primaryKey"`
	QRCode         string `gorm:"uniqueIndex"`
	Name           string
	SerialNumber   string `gorm:"uniqueIndex"`
	PurchaseDate   int64
	BuildingID     string
	FloorID        string
	RoomID         string
	LocationDetail string
	Condition      string
	OwnerID        string
	OwnerType      string
	Brand          string
	Model          string
	Year           int
	Color          string
	Size           string
	Vendor         string
	Price          float64
	InvoiceNo      string
	ReceivedAt     int64
	CreatedAt      int64
	UpdatedAt      int64
	DeletedAt      gorm.DeletedAt `gorm:"index"`
}

func (InventoryItemModel) TableName() string { return "inventory_items" }

// Model untuk AssetIssue
type AssetIssueModel struct {
	ID              string `gorm:"primaryKey"`
	InventoryItemID *string
	ReporterID      string
	RoomID          string
	IssueType       string
	Description     string
	Photos          string // JSON array
	Status          string
	AssignedTo      *string
	VerifiedBy      *string
	VerifiedAt      *int64
	ClosedAt        *int64
	Resolution      *string
	CreatedAt       int64
	UpdatedAt       int64
}

func (AssetIssueModel) TableName() string { return "asset_issues" }
