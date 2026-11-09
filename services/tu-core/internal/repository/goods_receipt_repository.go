package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodsReceipt struct {
	ID                   string     `json:"id"`
	BastNumber           string     `json:"bast_number"`
	ReceiptDate          time.Time  `json:"receipt_date"`
	SourceType           string     `json:"source_type"`
	SupplierName         string     `json:"supplier_name"`
	InvoiceNumber        string     `json:"invoice_number"`
	FundingSource        string     `json:"funding_source"`
	ProcurementRequestID *string    `json:"procurement_request_id"`
	ReceivedBy           string     `json:"received_by"`
	VerifiedBy           *string    `json:"verified_by"`
	Status               string     `json:"status"`
	SemesterID           string     `json:"semester_id"`
	RejectionReason      *string    `json:"rejection_reason"`
	DocumentURLs         []string   `json:"document_urls"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	Items                []GoodsReceiptItem `json:"items,omitempty"`
}

type GoodsReceiptItem struct {
	ID             string  `json:"id"`
	GoodsReceiptID string  `json:"goods_receipt_id"`
	CatalogID      string  `json:"catalog_id"`
	ItemName       string  `json:"item_name"`
	Specification  string  `json:"specification"`
	Quantity       int     `json:"quantity"`
	Unit           string  `json:"unit"`
	UnitPrice      float64 `json:"unit_price"`
	TotalPrice     float64 `json:"total_price"`
	Condition      string  `json:"condition"`
	IsInventory    bool    `json:"is_inventory"`
}

type GoodsReceiptRepository interface {
	Create(ctx context.Context, receipt *GoodsReceipt) error
	GetAll(ctx context.Context, semesterID string) ([]*GoodsReceipt, error)
	GetByID(ctx context.Context, id string) (*GoodsReceipt, error)
	Verify(ctx context.Context, receiptID string, verifierID string) error
	Reject(ctx context.Context, receiptID string, verifierID string, reason string) error
	GenerateBastNumber(ctx context.Context, semesterID string) (string, error)
}

type pgxGoodsReceiptRepository struct {
	db *pgxpool.Pool
}

func NewGoodsReceiptRepository(db *pgxpool.Pool) GoodsReceiptRepository {
	return &pgxGoodsReceiptRepository{db: db}
}

func (r *pgxGoodsReceiptRepository) GenerateBastNumber(ctx context.Context, semesterID string) (string, error) {
	year := time.Now().Year()
	var count int
	query := `SELECT COUNT(*) FROM goods_receipts WHERE EXTRACT(YEAR FROM receipt_date) = $1 AND semester_id = $2`
	if err := r.db.QueryRow(ctx, query, year, semesterID).Scan(&count); err != nil {
		return "", err
	}
	var schoolCode string
	r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schoolCode)
	return fmt.Sprintf("BAST/%s/%d/%03d", schoolCode, year, count+1), nil
}

func (r *pgxGoodsReceiptRepository) Create(ctx context.Context, receipt *GoodsReceipt) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	docs, _ := json.Marshal(receipt.DocumentURLs)
	query := `
		INSERT INTO goods_receipts (bast_number, receipt_date, source_type, supplier_name, invoice_number, funding_source, procurement_request_id, received_by, status, semester_id, document_urls)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	err = tx.QueryRow(ctx, query,
		receipt.BastNumber, receipt.ReceiptDate, receipt.SourceType, receipt.SupplierName,
		receipt.InvoiceNumber, receipt.FundingSource, receipt.ProcurementRequestID,
		receipt.ReceivedBy, receipt.Status, receipt.SemesterID, docs,
	).Scan(&receipt.ID)
	if err != nil { return err }

	for i := range receipt.Items {
		item := &receipt.Items[i]
		item.GoodsReceiptID = receipt.ID
		item.TotalPrice = float64(item.Quantity) * item.UnitPrice
		itemQuery := `
			INSERT INTO goods_receipt_items (goods_receipt_id, catalog_id, item_name, specification, quantity, unit, unit_price, total_price, condition, is_inventory)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`
		_, err = tx.Exec(ctx, itemQuery, item.GoodsReceiptID, item.CatalogID, item.ItemName, item.Specification, item.Quantity, item.Unit, item.UnitPrice, item.TotalPrice, item.Condition, item.IsInventory)
		if err != nil { return err }
	}

	return tx.Commit(ctx)
}

func (r *pgxGoodsReceiptRepository) GetAll(ctx context.Context, semesterID string) ([]*GoodsReceipt, error) {
	query := `
		SELECT id::text, bast_number, receipt_date, source_type, supplier_name, COALESCE(invoice_number, ''), 
		       COALESCE(funding_source, ''), procurement_request_id::text, received_by, verified_by, status, semester_id, 
		       rejection_reason, document_urls, created_at, updated_at
		FROM goods_receipts WHERE semester_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*GoodsReceipt
	for rows.Next() {
		var gr GoodsReceipt
		var docs []byte
		err := rows.Scan(&gr.ID, &gr.BastNumber, &gr.ReceiptDate, &gr.SourceType, &gr.SupplierName, &gr.InvoiceNumber,
			&gr.FundingSource, &gr.ProcurementRequestID, &gr.ReceivedBy, &gr.VerifiedBy, &gr.Status, &gr.SemesterID,
			&gr.RejectionReason, &docs, &gr.CreatedAt, &gr.UpdatedAt)
		if err != nil { return nil, err }
		if docs != nil { json.Unmarshal(docs, &gr.DocumentURLs) }
		list = append(list, &gr)
	}
	return list, nil
}

func (r *pgxGoodsReceiptRepository) GetByID(ctx context.Context, id string) (*GoodsReceipt, error) {
	query := `
		SELECT id::text, bast_number, receipt_date, source_type, supplier_name, COALESCE(invoice_number, ''), 
		       COALESCE(funding_source, ''), procurement_request_id::text, received_by, verified_by, status, semester_id, 
		       rejection_reason, document_urls, created_at, updated_at
		FROM goods_receipts WHERE id = $1
	`
	var gr GoodsReceipt
	var docs []byte
	err := r.db.QueryRow(ctx, query, id).Scan(&gr.ID, &gr.BastNumber, &gr.ReceiptDate, &gr.SourceType, &gr.SupplierName, &gr.InvoiceNumber,
		&gr.FundingSource, &gr.ProcurementRequestID, &gr.ReceivedBy, &gr.VerifiedBy, &gr.Status, &gr.SemesterID,
		&gr.RejectionReason, &docs, &gr.CreatedAt, &gr.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, errors.New("penerimaan tidak ditemukan") }
		return nil, err
	}
	if docs != nil { json.Unmarshal(docs, &gr.DocumentURLs) }

	// Load Items
	itemQuery := `SELECT id::text, goods_receipt_id::text, COALESCE(catalog_id::text, ''), item_name, COALESCE(specification, ''), quantity, unit, COALESCE(unit_price, 0), COALESCE(total_price, 0), condition, is_inventory FROM goods_receipt_items WHERE goods_receipt_id = $1`
	rows, _ := r.db.Query(ctx, itemQuery, id)
	defer rows.Close()
	for rows.Next() {
		var item GoodsReceiptItem
		rows.Scan(&item.ID, &item.GoodsReceiptID, &item.CatalogID, &item.ItemName, &item.Specification, &item.Quantity, &item.Unit, &item.UnitPrice, &item.TotalPrice, &item.Condition, &item.IsInventory)
		gr.Items = append(gr.Items, item)
	}
	return &gr, nil
}

func (r *pgxGoodsReceiptRepository) Verify(ctx context.Context, receiptID string, verifierID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	var sourceType string
	var reqID *string
	var status string
	var receiptDate time.Time
	err = tx.QueryRow(ctx, `SELECT source_type, procurement_request_id, status, receipt_date FROM goods_receipts WHERE id = $1 FOR UPDATE`, receiptID).Scan(&sourceType, &reqID, &status, &receiptDate)
	if err != nil { return err }
	if status != "DRAFT" { return errors.New("hanya status DRAFT yang dapat diverifikasi") }

	_, err = tx.Exec(ctx, `UPDATE goods_receipts SET status = 'VERIFIED', verified_by = $1, updated_at = NOW() WHERE id = $2`, verifierID, receiptID)
	if err != nil { return err }

	if sourceType == "PROCUREMENT" && reqID != nil {
		_, err = tx.Exec(ctx, `UPDATE procurement_requests SET status = 'PROCURED', updated_at = NOW() WHERE id = $1`, *reqID)
		if err != nil { return err }
	}

	rows, err := tx.Query(ctx, `SELECT catalog_id::text, quantity, is_inventory, unit_price, condition, item_name FROM goods_receipt_items WHERE goods_receipt_id = $1`, receiptID)
	if err != nil { return err }
	defer rows.Close()

	var schCode string
	tx.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&schCode)

	romanMonths := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	monthRoman := romanMonths[int(receiptDate.Month())]

	for rows.Next() {
		var catalogID, condition, itemName string
		var quantity int
		var isInventory bool
		var unitPrice float64
		rows.Scan(&catalogID, &quantity, &isInventory, &unitPrice, &condition, &itemName)

		if isInventory {
			for i := 1; i <= quantity; i++ {
				var catCode string
				tx.QueryRow(ctx, `SELECT code FROM master_categories WHERE id = (SELECT category_id FROM master_catalogs WHERE id = $1)`, catalogID).Scan(&catCode)
				
				// Generate ID dengan Prefix Sumber
				prefix := ""
				switch sourceType {
				case "PROCUREMENT", "PEMBELIAN": prefix = "P"
				case "HIBAH": prefix = "H"
				case "DONASI": prefix = "D"
				case "MUTASI": prefix = "M"
				}
				
				var countAset int
				tx.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_items WHERE sumber_perolehan = $1 AND EXTRACT(YEAR FROM purchase_date) = $2`, sourceType, receiptDate.Year()).Scan(&countAset)
				
				var assetCode string
				if prefix != "" {
					assetCode = fmt.Sprintf("%s/%s/%s/%d/%03d", schCode, prefix, catCode, receiptDate.Year(), countAset+i)
				} else {
					assetCode = fmt.Sprintf("%s/%s/%d/%03d", schCode, catCode, receiptDate.Year(), countAset+i)
				}

				var countReg int
				tx.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_items WHERE nomor_registrasi IS NOT NULL AND EXTRACT(YEAR FROM tanggal_pencatatan) = $1`, receiptDate.Year()).Scan(&countReg)
				regNumber := fmt.Sprintf("%03d/INV/%s/%s/%d", countReg+i, schCode, monthRoman, receiptDate.Year())

				_, err = tx.Exec(ctx, `INSERT INTO inventory_items (id, nomor_registrasi, catalog_id, condition, purchase_date, goods_receipt_id, sumber_perolehan, nilai_buku, status_aset, tanggal_pencatatan) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'AKTIF', $9)`, 
					assetCode, regNumber, catalogID, condition, receiptDate, receiptID, sourceType, unitPrice, receiptDate)
				if err != nil { return err }

				_, err = tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, description, changed_by) VALUES ($1, 'CREATED', 'Dicatat otomatis dari BAST', $2)`, assetCode, verifierID)
				if err != nil { return err }
			}
		} else {
			var logisticID string
			err = tx.QueryRow(ctx, `SELECT id::text FROM sarpras_logistics WHERE nama = $1 LIMIT 1`, itemName).Scan(&logisticID)
			if err != nil {
				err = tx.QueryRow(ctx, `INSERT INTO sarpras_logistics (kode, nama, kategori, satuan, stok) VALUES ($1, $2, 'Penerimaan', 'Unit', $3) RETURNING id::text`, "LOG-"+time.Now().Format("0601021504"), itemName, quantity).Scan(&logisticID)
				if err != nil { return err }
			} else {
				_, err = tx.Exec(ctx, `UPDATE sarpras_logistics SET stok = stok + $1, updated_at = NOW() WHERE id = $2`, quantity, logisticID)
				if err != nil { return err }
			}
			_, err = tx.Exec(ctx, `INSERT INTO sarpras_logistic_mutations (logistic_id, tipe, jumlah, keterangan) VALUES ($1, 'IN', $2, 'Penerimaan BAST')`, logisticID, quantity)
			if err != nil { return err }
		}
	}

	docRows, err := tx.Query(ctx, `SELECT document_id::text FROM sarpras_inventory_document_links WHERE entity_type = 'GOODS_RECEIPT' AND entity_id = $1`, receiptID)
	if err == nil {
		var docIDs []string
		for docRows.Next() { var did string; docRows.Scan(&did); docIDs = append(docIDs, did) }
		docRows.Close()

		assetRows, err := tx.Query(ctx, `SELECT id FROM inventory_items WHERE goods_receipt_id = $1`, receiptID)
		if err == nil {
			var assetIDs []string
			for assetRows.Next() { var aid string; assetRows.Scan(&aid); assetIDs = append(assetIDs, aid) }
			assetRows.Close()

			for _, did := range docIDs {
				for _, aid := range assetIDs {
					tx.Exec(ctx, `INSERT INTO sarpras_inventory_document_links (document_id, entity_type, entity_id, relation_type, created_by) VALUES ($1, 'INVENTORY_ITEM', $2, 'SUPPORTING', $3)`, did, aid, verifierID)
				}
			}
		}
	}

	return tx.Commit(ctx)
}

func (r *pgxGoodsReceiptRepository) Reject(ctx context.Context, receiptID string, verifierID string, reason string) error {
	_, err := r.db.Exec(ctx, `UPDATE goods_receipts SET status = 'REJECTED', verified_by = $1, rejection_reason = $2, updated_at = NOW() WHERE id = $3`, verifierID, reason, receiptID)
	return err
}
