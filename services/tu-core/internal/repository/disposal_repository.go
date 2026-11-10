package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DisposalRequest struct {
	ID                     string                `json:"id"`
	RequestNumber          string                `json:"request_number"`
	ProposedBy             string                `json:"proposed_by"`
	ProposalDate           time.Time             `json:"proposal_date"`
	ReasonCategory         string                `json:"reason_category"`
	ProposedMethod         string                `json:"proposed_method"`
	TotalAssets            int                   `json:"total_assets"`
	TotalValue             float64               `json:"total_value"`
	ApprovalLevel          int                   `json:"approval_level"`
	Status                 string                `json:"status"`
	CurrentApprovalStep    int                   `json:"current_approval_step"`
	RejectionReason        string                `json:"rejection_reason"`
	ExecutionDate          *time.Time            `json:"execution_date"`
	ExecutionMethod        string                `json:"execution_method"`
	ExecutionLocation      string                `json:"execution_location"`
	Executors              string                `json:"executors"`
	Witnesses              string                `json:"witnesses"`
	ExecutionDocumentation string                `json:"execution_documentation"`
	ExecutionNotes         string                `json:"execution_notes"`
	SaleProceeds           float64               `json:"sale_proceeds"`
	SaleEvidenceURL        string                `json:"sale_evidence_url"`
	BastNumber             string                `json:"bast_number"`
	BastPdfUrl             string                `json:"bast_pdf_url"`
	VerificationHash       string                `json:"verification_hash"`
	CreatedAt              time.Time             `json:"created_at"`
	UpdatedAt              time.Time             `json:"updated_at"`
	Items                  []DisposalItem        `json:"items,omitempty"`
	Verification           *DisposalVerification `json:"verification,omitempty"`
	Approvals              []DisposalApproval    `json:"approvals,omitempty"`
}

type DisposalItem struct {
	ID                 string `json:"id"`
	DisposalRequestID  string `json:"disposal_request_id"`
	InventoryItemID    string `json:"inventory_item_id"`
	SpecificReason     string `json:"specific_reason"`
	PhotoURLs          string `json:"photo_urls"`
	VerificationStatus string `json:"verification_status"`
	VerifiedCondition  string `json:"verified_condition"`
	VerificationNotes  string `json:"verification_notes"`

	// Join fields
	AssetCode        string  `json:"asset_code"`
	AssetName        string  `json:"asset_name"`
	MerkType         string  `json:"merk_type"`
	AcquisitionYear  int     `json:"acquisition_year"`
	AcquisitionPrice float64 `json:"acquisition_price"`
	LastCondition    string  `json:"last_condition"`
	RoomID           string  `json:"room_id"`
}

type DisposalVerification struct {
	ID                      string    `json:"id"`
	DisposalRequestID       string    `json:"disposal_request_id"`
	VerificationDate        time.Time `json:"verification_date"`
	VerificationTeam        string    `json:"verification_team"`
	VerificationDocumentURL string    `json:"verification_document_url"`
	TeamRecommendation      string    `json:"team_recommendation"`
	TeamNotes               string    `json:"team_notes"`
	VerifiedBy              string    `json:"verified_by"`
	CreatedAt               time.Time `json:"created_at"`
}

type DisposalApproval struct {
	ID                  string    `json:"id"`
	DisposalRequestID   string    `json:"disposal_request_id"`
	ApprovalStep        int       `json:"approval_step"`
	ApproverID          string    `json:"approver_id"`
	Decision            string    `json:"decision"`
	ApprovalDocumentURL string    `json:"approval_document_url"`
	Notes               string    `json:"notes"`
	ApprovalDate        time.Time `json:"approval_date"`
}

type DisposalRepository interface {
	CreateDisposal(ctx context.Context, req *DisposalRequest, items []DisposalItem, semesterID string) error
	GetAllDisposals(ctx context.Context, semesterID string) ([]*DisposalRequest, error)
	GetDisposalByID(ctx context.Context, id string) (*DisposalRequest, error)
	VerifyDisposal(ctx context.Context, req *DisposalRequest, verif *DisposalVerification) error
	ApproveDisposal(ctx context.Context, req *DisposalRequest, approval *DisposalApproval) error
	ExecuteDisposal(ctx context.Context, req *DisposalRequest) error
	GenerateRequestNumber(ctx context.Context) (string, error)
}

type pgxDisposalRepository struct {
	db *pgxpool.Pool
}

func NewDisposalRepository(db *pgxpool.Pool) DisposalRepository {
	return &pgxDisposalRepository{db: db}
}

func (r *pgxDisposalRepository) GenerateRequestNumber(ctx context.Context) (string, error) {
	var count int
	year := time.Now().Year()
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM disposal_requests WHERE EXTRACT(YEAR FROM proposal_date) = $1`, year).Scan(&count)
	if err != nil {
		return "", err
	}

	var schCode string
	r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schCode)

	return fmt.Sprintf("USUL-HPS/%s/%d/%03d", schCode, year, count+1), nil
}

func (r *pgxDisposalRepository) CreateDisposal(ctx context.Context, req *DisposalRequest, items []DisposalItem, semesterID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert Header
	query := `
		INSERT INTO disposal_requests
		(request_number, proposed_by, proposal_date, reason_category, proposed_method, total_assets, total_value, approval_level, status, current_approval_step)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'DRAFT', 0)
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, query, req.RequestNumber, req.ProposedBy, req.ProposalDate, req.ReasonCategory, req.ProposedMethod, req.TotalAssets, req.TotalValue, req.ApprovalLevel, semesterID).Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert Items & Update Status
	for _, item := range items {
		_, err = tx.Exec(ctx, `INSERT INTO disposal_request_items (disposal_request_id, inventory_item_id, specific_reason, photo_urls) VALUES ($1, $2, $3, $4)`, req.ID, item.InventoryItemID, item.SpecificReason, item.PhotoURLs)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'DIUSULKAN_HAPUS', disposal_request_id = $1 WHERE id = $2`, req.ID, item.InventoryItemID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'DIUSULKAN_HAPUS', $2, $3)`, item.InventoryItemID, "Diusulkan hapus via "+req.RequestNumber, req.ProposedBy)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *pgxDisposalRepository) GetAllDisposals(ctx context.Context, semesterID string) ([]*DisposalRequest, error) {
	query := `SELECT id::text, request_number, proposed_by, proposal_date, reason_category, proposed_method, total_assets, total_value, approval_level, status, current_approval_step, COALESCE(rejection_reason, ''), execution_date, COALESCE(execution_method, ''), COALESCE(bast_number, '') FROM disposal_requests WHERE ($1 = '' OR semester_id = $1) ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []*DisposalRequest
	for rows.Next() {
		var d DisposalRequest
		if err := rows.Scan(&d.ID, &d.RequestNumber, &d.ProposedBy, &d.ProposalDate, &d.ReasonCategory, &d.ProposedMethod, &d.TotalAssets, &d.TotalValue, &d.ApprovalLevel, &d.Status, &d.CurrentApprovalStep, &d.RejectionReason, &d.ExecutionDate, &d.ExecutionMethod, &d.BastNumber); err != nil {
			return nil, err
		}
		list = append(list, &d)
	}
	return list, nil
}

func (r *pgxDisposalRepository) GetDisposalByID(ctx context.Context, id string) (*DisposalRequest, error) {
	query := `SELECT id::text, request_number, proposed_by, proposal_date, reason_category, proposed_method, total_assets, total_value, approval_level, status, current_approval_step, COALESCE(rejection_reason, ''), execution_date, COALESCE(execution_method, ''), COALESCE(bast_number, ''), COALESCE(bast_pdf_url, '') FROM disposal_requests WHERE id = $1`
	var d DisposalRequest
	err := r.db.QueryRow(ctx, query, id).Scan(&d.ID, &d.RequestNumber, &d.ProposedBy, &d.ProposalDate, &d.ReasonCategory, &d.ProposedMethod, &d.TotalAssets, &d.TotalValue, &d.ApprovalLevel, &d.Status, &d.CurrentApprovalStep, &d.RejectionReason, &d.ExecutionDate, &d.ExecutionMethod, &d.BastNumber, &d.BastPdfUrl)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("usulan tidak ditemukan")
		}
		return nil, err
	}

	// Get Items
	itemQuery := `
		SELECT dri.id::text, dri.disposal_request_id::text, dri.inventory_item_id, dri.specific_reason, COALESCE(dri.photo_urls::text, '[]'), dri.verification_status, COALESCE(dri.verified_condition, ''), COALESCE(dri.verification_notes, ''),
		       i.id, c.name, COALESCE(c.brand, ''), EXTRACT(YEAR FROM i.purchase_date), COALESCE(i.nilai_buku, 0), COALESCE(i.condition, ''), COALESCE(i.room_id::text, '')
		FROM disposal_request_items dri
		JOIN inventory_items i ON dri.inventory_item_id = i.id
		JOIN master_catalogs c ON i.catalog_id = c.id
		WHERE dri.disposal_request_id = $1
	`
	rows, _ := r.db.Query(ctx, itemQuery, id)
	defer rows.Close()
	for rows.Next() {
		var i DisposalItem
		rows.Scan(&i.ID, &i.DisposalRequestID, &i.InventoryItemID, &i.SpecificReason, &i.PhotoURLs, &i.VerificationStatus, &i.VerifiedCondition, &i.VerificationNotes,
			&i.AssetCode, &i.AssetName, &i.MerkType, &i.AcquisitionYear, &i.AcquisitionPrice, &i.LastCondition, &i.RoomID)
		d.Items = append(d.Items, i)
	}

	return &d, nil
}

func (r *pgxDisposalRepository) VerifyDisposal(ctx context.Context, req *DisposalRequest, verif *DisposalVerification) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert Verifikasi
	query := `INSERT INTO disposal_verifications (disposal_request_id, verification_date, verification_team, verification_document_url, team_recommendation, team_notes, verified_by) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(ctx, query, req.ID, verif.VerificationDate, verif.VerificationTeam, verif.VerificationDocumentURL, verif.TeamRecommendation, verif.TeamNotes, verif.VerifiedBy)
	if err != nil {
		return err
	}

	// Update Status Request
	_, err = tx.Exec(ctx, `UPDATE disposal_requests SET status = $1, current_approval_step = 1, updated_at = NOW() WHERE id = $2`, "VERIFIED_KTU", req.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgxDisposalRepository) ApproveDisposal(ctx context.Context, req *DisposalRequest, approval *DisposalApproval) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `INSERT INTO disposal_approvals (disposal_request_id, approval_step, approver_id, decision, approval_document_url, notes) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(ctx, query, req.ID, approval.ApprovalStep, approval.ApproverID, approval.Decision, approval.ApprovalDocumentURL, approval.Notes)
	if err != nil {
		return err
	}

	newStatus := fmt.Sprintf("APPROVED_STEP_%d", approval.ApprovalStep)
	if approval.Decision == "REJECTED" {
		newStatus = "REJECTED"
		// Kembalikan status aset
		tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'AKTIF', disposal_request_id = NULL WHERE disposal_request_id = $1`, req.ID)
	} else if approval.ApprovalStep == req.ApprovalLevel {
		newStatus = "APPROVED"
	}

	_, err = tx.Exec(ctx, `UPDATE disposal_requests SET status = $1, current_approval_step = $2, updated_at = NOW() WHERE id = $3`, newStatus, approval.ApprovalStep, req.ID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *pgxDisposalRepository) ExecuteDisposal(ctx context.Context, req *DisposalRequest) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update Request
	query := `UPDATE disposal_requests SET status = 'EXECUTED', execution_date = $1, execution_method = $2, execution_location = $3, executors = $4, witnesses = $5, execution_documentation = $6, execution_notes = $7, sale_proceeds = $8, sale_evidence_url = $9, bast_number = $10, bast_pdf_url = $11, verification_hash = $12, updated_at = NOW() WHERE id = $13`
	_, err = tx.Exec(ctx, query, req.ExecutionDate, req.ExecutionMethod, req.ExecutionLocation, req.Executors, req.Witnesses, req.ExecutionDocumentation, req.ExecutionNotes, req.SaleProceeds, req.SaleEvidenceURL, req.BastNumber, req.BastPdfUrl, req.VerificationHash, req.ID)
	if err != nil {
		return err
	}

	// Update Items
	for _, item := range req.Items {
		tx.Exec(ctx, `UPDATE inventory_items SET status_aset = 'DIHAPUS', tanggal_penghapusan = $1, metode_penghapusan = $2, keterangan_penghapusan = $3 WHERE id = $4`, req.ExecutionDate, req.ExecutionMethod, "BAP "+req.BastNumber, item.InventoryItemID)
		tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'DIHAPUS', $2, $3)`, item.InventoryItemID, "Dihapus via "+req.BastNumber, req.ProposedBy)
	}

	return tx.Commit(ctx)
}
