package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type InventoryDocument struct {
	ID              string     `json:"id"`
	DocumentType    string     `json:"document_type"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	DocumentNumber  string     `json:"document_number"`
	DocumentDate    *time.Time `json:"document_date"`
	BucketName      string     `json:"bucket_name"`
	ObjectKey       string     `json:"object_key"`
	OriginalFilename string    `json:"original_filename"`
	MimeType        string     `json:"mime_type"`
	FileSize        int64      `json:"file_size"`
	Checksum        string     `json:"checksum"`
	Visibility      string     `json:"visibility"`
	Status          string     `json:"status"`
	RejectionReason string     `json:"rejection_reason"`
	UploadedBy      string     `json:"uploaded_by"`
	VerifiedBy      string     `json:"verified_by"`
	VerifiedAt      *time.Time `json:"verified_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type DocumentRepository interface {
	CreateDocument(ctx context.Context, doc *InventoryDocument) error
	LinkDocument(ctx context.Context, docID, entityType, entityID, relationType, userID string) error
	GetDocumentsByEntity(ctx context.Context, entityType, entityID string) ([]*InventoryDocument, error)
	GetDocumentByID(ctx context.Context, id string) (*InventoryDocument, error)
	UpdateStatus(ctx context.Context, id string, status string, userID string, reason string) error
}

type pgxDocumentRepository struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) DocumentRepository {
	return &pgxDocumentRepository{db: db}
}

func (r *pgxDocumentRepository) CreateDocument(ctx context.Context, doc *InventoryDocument) error {
	query := `
		INSERT INTO sarpras_inventory_documents 
		(document_type, title, description, document_number, document_date, bucket_name, object_key, original_filename, mime_type, file_size, visibility, status, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id, created_at
	`
	return r.db.QueryRow(ctx, query, doc.DocumentType, doc.Title, doc.Description, doc.DocumentNumber, doc.DocumentDate, doc.BucketName, doc.ObjectKey, doc.OriginalFilename, doc.MimeType, doc.FileSize, doc.Visibility, doc.Status, doc.UploadedBy).Scan(&doc.ID, &doc.CreatedAt)
}

func (r *pgxDocumentRepository) LinkDocument(ctx context.Context, docID, entityType, entityID, relationType, userID string) error {
	query := `INSERT INTO sarpras_inventory_document_links (document_id, entity_type, entity_id, relation_type, created_by) VALUES ($1, $2, $3, $4, $5)`
	_, err := r.db.Exec(ctx, query, docID, entityType, entityID, relationType, userID)
	return err
}

func (r *pgxDocumentRepository) GetDocumentsByEntity(ctx context.Context, entityType, entityID string) ([]*InventoryDocument, error) {
	query := `
		SELECT d.id::text, d.document_type, d.title, COALESCE(d.description, ''), COALESCE(d.document_number, ''), d.document_date,
		       d.bucket_name, d.object_key, d.original_filename, d.mime_type, d.file_size, COALESCE(d.checksum, ''),
		       d.visibility, d.status, COALESCE(d.rejection_reason, ''), d.uploaded_by, COALESCE(d.verified_by, ''), d.verified_at, d.created_at
		FROM sarpras_inventory_documents d
		JOIN sarpras_inventory_document_links l ON d.id = l.document_id
		WHERE l.entity_type = $1 AND l.entity_id = $2
		ORDER BY d.created_at DESC
	`
	rows, err := r.db.Query(ctx, query, entityType, entityID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*InventoryDocument
	for rows.Next() {
		var d InventoryDocument
		err := rows.Scan(&d.ID, &d.DocumentType, &d.Title, &d.Description, &d.DocumentNumber, &d.DocumentDate,
			&d.BucketName, &d.ObjectKey, &d.OriginalFilename, &d.MimeType, &d.FileSize, &d.Checksum,
			&d.Visibility, &d.Status, &d.RejectionReason, &d.UploadedBy, &d.VerifiedBy, &d.VerifiedAt, &d.CreatedAt)
		if err != nil { return nil, err }
		list = append(list, &d)
	}
	return list, nil
}

func (r *pgxDocumentRepository) GetDocumentByID(ctx context.Context, id string) (*InventoryDocument, error) {
	query := `
		SELECT id::text, document_type, title, COALESCE(description, ''), COALESCE(document_number, ''), document_date,
		       bucket_name, object_key, original_filename, mime_type, file_size, COALESCE(checksum, ''),
		       visibility, status, COALESCE(rejection_reason, ''), uploaded_by, COALESCE(verified_by, ''), verified_at, created_at
		FROM sarpras_inventory_documents WHERE id = $1
	`
	var d InventoryDocument
	err := r.db.QueryRow(ctx, query, id).Scan(&d.ID, &d.DocumentType, &d.Title, &d.Description, &d.DocumentNumber, &d.DocumentDate,
		&d.BucketName, &d.ObjectKey, &d.OriginalFilename, &d.MimeType, &d.FileSize, &d.Checksum,
		&d.Visibility, &d.Status, &d.RejectionReason, &d.UploadedBy, &d.VerifiedBy, &d.VerifiedAt, &d.CreatedAt)
	return &d, err
}

func (r *pgxDocumentRepository) UpdateStatus(ctx context.Context, id string, status string, userID string, reason string) error {
	query := `UPDATE sarpras_inventory_documents SET status = $1, verified_by = $2, verified_at = NOW(), rejection_reason = $3, updated_at = NOW() WHERE id = $4`
	_, err := r.db.Exec(ctx, query, status, userID, reason, id)
	return err
}
