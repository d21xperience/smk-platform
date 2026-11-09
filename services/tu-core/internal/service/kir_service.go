package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jung-kurt/gofpdf"
	qrcode "github.com/skip2/go-qrcode"
)

type KIRRoom struct {
	ID              string `json:"id"`
	Nama            string `json:"nama"`
	PenanggungJawab string `json:"penanggungJawab"`
}

type KIRService struct {
	db       *pgxpool.Pool
	itemRepo repository.InventoryItemRepository
}

func NewKIRService(db *pgxpool.Pool, itemRepo repository.InventoryItemRepository) *KIRService {
	return &KIRService{db: db, itemRepo: itemRepo}
}

func (s *KIRService) GetRooms(ctx context.Context) ([]KIRRoom, error) {
	query := `SELECT k.rombongan_belajar_id::text, k.nm_kelas, COALESCE(tp.nama, 'Belum Ditunjuk') as pjr FROM tabel_kelas k LEFT JOIN tabel_ptk tp ON k.ptk_id = tp.ptk_id ORDER BY k.nm_kelas ASC`
	rows, err := s.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []KIRRoom
	for rows.Next() {
		var r KIRRoom
		if err := rows.Scan(&r.ID, &r.Nama, &r.PenanggungJawab); err != nil {
			return nil, err
		}
		rooms = append(rooms, r)
	}
	return rooms, nil
}

func (s *KIRService) GetItemsByRoom(ctx context.Context, roomID string) ([]repository.InventoryItem, error) {
	return s.itemRepo.FindByRoomID(ctx, roomID)
}

func (s *KIRService) GetAllItems(ctx context.Context) ([]repository.InventoryItem, error) {
	return s.itemRepo.FindAllItems(ctx)
}

func (s *KIRService) CreateItem(ctx context.Context, item *repository.InventoryItem) error {
	if item.CatalogID == "" {
		return errors.New("Katalog Barang wajib dipilih")
	}

	newID, err := s.itemRepo.GenerateID(ctx, item.CatalogID, item.PurchaseDate.Year(), item.SumberPerolehan)
	if err != nil {
		return err
	}
	item.ID = newID

	return s.itemRepo.Create(ctx, item)
}

func (s *KIRService) UpdateItem(ctx context.Context, item *repository.InventoryItem) error {
	return s.itemRepo.Update(ctx, item)
}

func (s *KIRService) DeleteItem(ctx context.Context, id string) error {
	return s.itemRepo.Delete(ctx, id)
}

func (s *KIRService) GetMasterCategories(ctx context.Context) ([]repository.MasterCategory, error) {
	return s.itemRepo.GetMasterCategories(ctx)
}

func (s *KIRService) CreateCategory(ctx context.Context, cat *repository.MasterCategory) error {
	return s.itemRepo.CreateCategory(ctx, cat)
}
func (s *KIRService) UpdateCategory(ctx context.Context, cat *repository.MasterCategory) error {
	return s.itemRepo.UpdateCategory(ctx, cat)
}
func (s *KIRService) DeleteCategory(ctx context.Context, id string) error {
	return s.itemRepo.DeleteCategory(ctx, id)
}

func (s *KIRService) GetMasterCatalogs(ctx context.Context) ([]repository.MasterCatalog, error) {
	return s.itemRepo.GetMasterCatalogs(ctx)
}
func (s *KIRService) CreateCatalog(ctx context.Context, cat *repository.MasterCatalog) error {
	return s.itemRepo.CreateCatalog(ctx, cat)
}
func (s *KIRService) UpdateCatalog(ctx context.Context, cat *repository.MasterCatalog) error {
	return s.itemRepo.UpdateCatalog(ctx, cat)
}
func (s *KIRService) DeleteCatalog(ctx context.Context, id string) error {
	return s.itemRepo.DeleteCatalog(ctx, id)
}

func (s *KIRService) GenerateStickerPDF(ctx context.Context, itemIDs []string) (*bytes.Buffer, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 10)

	x, y := 10.0, 10.0
	colWidth, rowHeight := 60.0, 50.0
	maxCols := 3
	col := 0

	for _, id := range itemIDs {
		item, err := s.itemRepo.FindByID(ctx, id)
		if err != nil {
			continue
		}

		qr, err := qrcode.New(item.ID, qrcode.Medium)
		if err != nil {
			continue
		}
		var qrBuf bytes.Buffer
		if err := qr.Write(128, &qrBuf); err != nil {
			continue
		}

		if col >= maxCols {
			col = 0
			y += rowHeight
			x = 10.0
			if y > 240 {
				pdf.AddPage()
				y = 10.0
			}
		}

		pdf.Rect(x, y, colWidth, rowHeight, "D")

		opts := gofpdf.ImageOptions{ImageType: "PNG", ReadDpi: true}
		pdf.RegisterImageOptionsReader("qr_"+item.ID, opts, &qrBuf)
		pdf.ImageOptions("qr_"+item.ID, x+2, y+2, 25, 25, false, opts, 0, "")

		pdf.SetXY(x+30, y+5)
		pdf.CellFormat(28, 5, item.Name, "", 1, "L", false, 0, "")
		pdf.SetXY(x+30, y+11)
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(28, 5, "Merk: "+item.Brand, "", 1, "L", false, 0, "")
		pdf.SetXY(x+30, y+16)
		pdf.CellFormat(28, 5, "Kategori: "+item.CategoryCode, "", 1, "L", false, 0, "")

		pdf.SetFont("Arial", "B", 8)
		pdf.SetXY(x+2, y+28)
		pdf.MultiCell(colWidth-4, 4, item.ID, "", "C", false)

		x += colWidth
		col++
	}

	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return nil, err
	}
	return &buf, nil
}

// ==========================================
// FITUR MUTASI & LAPORAN BUKU INDUK
// ==========================================

func (s *KIRService) MutateItem(ctx context.Context, itemID string, action string, newRoomID string, newCondition string, statusAset string, reason string, userID string) error {
	item, err := s.itemRepo.FindByID(ctx, itemID)
	if err != nil {
		return err
	}

	oldData := fmt.Sprintf(`{"room_id":"%s","condition":"%s","status":"%s"}`, item.RoomID, item.Condition, item.StatusAset)

	// Update Data
	if newRoomID != "" {
		item.RoomID = newRoomID
	}
	if newCondition != "" {
		item.Condition = newCondition
	}
	if statusAset != "" {
		item.StatusAset = statusAset
		if statusAset == "DIHAPUS" {
			now := time.Now()
			item.TanggalPenghapusan = &now
			item.AlasanPenghapusan = reason
		}
	}

	newData := fmt.Sprintf(`{"room_id":"%s","condition":"%s","status":"%s"}`, item.RoomID, item.Condition, item.StatusAset)

	return s.itemRepo.MutateItem(ctx, item, action, oldData, newData, userID)
}

func (s *KIRService) GetItemHistory(ctx context.Context, itemID string) ([]repository.InventoryItemHistory, error) {
	return s.itemRepo.GetHistory(ctx, itemID)
}

func (s *KIRService) GenerateBukuIndukPDF(ctx context.Context) (*bytes.Buffer, error) {
	items, err := s.itemRepo.FindAllItems(ctx)
	if err != nil {
		return nil, err
	}

	pdf := gofpdf.New("L", "mm", "A4", "") // Landscape
	pdf.AddPage()

	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(0, 10, "BUKU INDUK INVENTARIS SEKOLAH", "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(0, 5, "Dicetak pada: "+time.Now().Format("02 January 2006"), "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// Header Tabel
	pdf.SetFont("Arial", "B", 8)
	headers := []struct {
		w float64
		t string
	}{
		{8, "No"}, {25, "No. Registrasi"}, {35, "Kode Aset"}, {50, "Nama Barang"},
		{25, "Merk"}, {15, "Tahun"}, {25, "Harga"}, {20, "Kondisi"}, {25, "Status"},
	}
	for _, h := range headers {
		pdf.CellFormat(h.w, 7, h.t, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Body Tabel
	pdf.SetFont("Arial", "", 7)
	for i, item := range items {
		pdf.CellFormat(8, 6, fmt.Sprint(i+1), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, item.NomorRegistrasi, "1", 0, "L", false, 0, "")
		pdf.CellFormat(35, 6, item.ID, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 6, item.Name, "1", 0, "L", false, 0, "")
		pdf.CellFormat(25, 6, item.Brand, "1", 0, "L", false, 0, "")
		pdf.CellFormat(15, 6, fmt.Sprint(item.PurchaseDate.Year()), "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, fmt.Sprintf("Rp %.0f", item.NilaiBuku), "1", 0, "R", false, 0, "")
		pdf.CellFormat(20, 6, item.Condition, "1", 0, "C", false, 0, "")
		pdf.CellFormat(25, 6, item.StatusAset, "1", 1, "C", false, 0, "")
	}

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	return &buf, err
}
