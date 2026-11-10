package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ... (Struct MasterCategory, MasterCatalog, InventoryItem, History SAMA SEPERTI SEBELUMNYA) ...

type MasterCategory struct {
	ID   string `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type MasterCatalog struct {
	ID         string `json:"id"`
	CategoryID string `json:"category_id"`
	Name       string `json:"name"`
	Brand      string `json:"brand"`
	ImageURL   string `json:"image_url"`
}

type InventoryItem struct {
	ID                  string    `json:"id"`
	NomorRegistrasi     string    `json:"nomor_registrasi"`
	CatalogID           string    `json:"catalog_id"`
	Name                string    `json:"name"`
	Brand               string    `json:"brand"`
	CategoryName        string    `json:"category_name"`
	CategoryCode        string    `json:"category_code"`
	Condition           string    `json:"condition"`
	RoomID              string    `json:"room_id"`
	PurchaseDate        time.Time `json:"purchase_date"`
	SumberPerolehan     string    `json:"sumber_perolehan"`
	Keterangan          string    `json:"keterangan"`
	MasaManfaatTahun    int       `json:"masa_manfaat_tahun"`
	NilaiBuku           float64   `json:"nilai_buku"`
	StatusAset          string    `json:"status_aset"`
	TanggalPencatatan   *time.Time `json:"tanggal_pencatatan"`
	TanggalPenghapusan  *time.Time `json:"tanggal_penghapusan"`
	AlasanPenghapusan   string    `json:"alasan_penghapusan"`
}


type ConditionBookEntry struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Brand          string     `json:"brand"`
	CategoryName   string     `json:"category_name"`
	RoomName       string     `json:"room_name"`
	Condition      string     `json:"condition"`
	LastChangeDate *time.Time `json:"last_change_date"`
	LastChangeNote string     `json:"last_change_note"`
	Recommendation string     `json:"recommendation"`
}

type InventoryItemHistory struct {
	ID              string    `json:"id"`
	InventoryItemID string    `json:"inventory_item_id"`
	Action          string    `json:"action"`
	OldData         string    `json:"old_data"`
	NewData         string    `json:"new_data"`
	Description     string    `json:"description"`
	ChangedBy       string    `json:"changed_by"`
	ChangedAt       time.Time `json:"changed_at"`
}

type InventoryItemRepository interface {
	FindAllItems(ctx context.Context) ([]InventoryItem, error)
	FindByRoomID(ctx context.Context, roomID string) ([]InventoryItem, error)
	FindByID(ctx context.Context, id string) (*InventoryItem, error)
	Create(ctx context.Context, item *InventoryItem, semesterID string) error
	Update(ctx context.Context, item *InventoryItem) error
	Delete(ctx context.Context, id string) error
	MutateItem(ctx context.Context, item *InventoryItem, action string, oldData string, newData string, userID string) error
	GetHistory(ctx context.Context, itemID string) ([]InventoryItemHistory, error)
	GetConditionBook(ctx context.Context, semesterID string) ([]*ConditionBookEntry, error)

	
	GetMasterCategories(ctx context.Context) ([]MasterCategory, error)
	CreateCategory(ctx context.Context, cat *MasterCategory) error
	UpdateCategory(ctx context.Context, cat *MasterCategory) error
	DeleteCategory(ctx context.Context, id string) error

	GetMasterCatalogs(ctx context.Context) ([]MasterCatalog, error)
	CreateCatalog(ctx context.Context, cat *MasterCatalog) error
	UpdateCatalog(ctx context.Context, cat *MasterCatalog) error
	DeleteCatalog(ctx context.Context, id string) error

	GetSchoolCode(ctx context.Context) (string, error)
	CountItemsForAutoID(ctx context.Context, categoryID string, year int, sourceType string) (int, error)
	GenerateID(ctx context.Context, catalogID string, year int, sourceType string) (string, error)
}

type pgxInventoryItemRepository struct {
	db *pgxpool.Pool
}

func NewInventoryItemRepository(db *pgxpool.Pool) InventoryItemRepository {
	return &pgxInventoryItemRepository{db: db}
}

// ... (Fungsi FindAllItems, FindByRoomID, FindByID, scanItems, Create, Update, Delete, MutateItem, GetHistory, MasterData SAMA SEPERTI SEBELUMNYA) ...
// (Saya potong untuk menghemat tempat, Tuan bisa copy dari file sebelumnya, karena tidak ada perubahan di bagian ini)

func (r *pgxInventoryItemRepository) FindAllItems(ctx context.Context) ([]InventoryItem, error) {
	query := `
		SELECT i.id, COALESCE(i.nomor_registrasi, ''), i.catalog_id::text, c.name, c.brand, cat.name, cat.code, COALESCE(i.condition, ''), COALESCE(i.room_id::text, ''), i.purchase_date,
		       COALESCE(i.sumber_perolehan, ''), COALESCE(i.keterangan, ''), COALESCE(i.masa_manfaat_tahun, 0), COALESCE(i.nilai_buku, 0), COALESCE(i.status_aset, 'AKTIF'), i.tanggal_pencatatan, i.tanggal_penghapusan, COALESCE(i.alasan_penghapusan, '')
		FROM inventory_items i
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		ORDER BY i.id ASC
	`
	return r.scanItems(ctx, query)
}

func (r *pgxInventoryItemRepository) FindByRoomID(ctx context.Context, roomID string) ([]InventoryItem, error) {
	query := `
		SELECT i.id, COALESCE(i.nomor_registrasi, ''), i.catalog_id::text, c.name, c.brand, cat.name, cat.code, COALESCE(i.condition, ''), COALESCE(i.room_id::text, ''), i.purchase_date,
		       COALESCE(i.sumber_perolehan, ''), COALESCE(i.keterangan, ''), COALESCE(i.masa_manfaat_tahun, 0), COALESCE(i.nilai_buku, 0), COALESCE(i.status_aset, 'AKTIF'), i.tanggal_pencatatan, i.tanggal_penghapusan, COALESCE(i.alasan_penghapusan, '')
		FROM inventory_items i
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		WHERE i.room_id = $1 AND i.status_aset != 'DIHAPUS'
		ORDER BY i.id ASC
	`
	return r.scanItems(ctx, query, roomID)
}

func (r *pgxInventoryItemRepository) FindByID(ctx context.Context, id string) (*InventoryItem, error) {
	query := `
		SELECT i.id, COALESCE(i.nomor_registrasi, ''), i.catalog_id::text, c.name, c.brand, cat.name, cat.code, COALESCE(i.condition, ''), COALESCE(i.room_id::text, ''), i.purchase_date,
		       COALESCE(i.sumber_perolehan, ''), COALESCE(i.keterangan, ''), COALESCE(i.masa_manfaat_tahun, 0), COALESCE(i.nilai_buku, 0), COALESCE(i.status_aset, 'AKTIF'), i.tanggal_pencatatan, i.tanggal_penghapusan, COALESCE(i.alasan_penghapusan, '')
		FROM inventory_items i
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		WHERE i.id = $1
	`
	var i InventoryItem
	err := r.db.QueryRow(ctx, query, id).Scan(
		&i.ID, &i.NomorRegistrasi, &i.CatalogID, &i.Name, &i.Brand, &i.CategoryName, &i.CategoryCode, &i.Condition, &i.RoomID, &i.PurchaseDate,
		&i.SumberPerolehan, &i.Keterangan, &i.MasaManfaatTahun, &i.NilaiBuku, &i.StatusAset, &i.TanggalPencatatan, &i.TanggalPenghapusan, &i.AlasanPenghapusan,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { return nil, errors.New("aset tidak ditemukan") }
		return nil, err
	}
	return &i, nil
}

func (r *pgxInventoryItemRepository) scanItems(ctx context.Context, query string, args ...interface{}) ([]InventoryItem, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []InventoryItem
	for rows.Next() {
		var i InventoryItem
		if err := rows.Scan(
			&i.ID, &i.NomorRegistrasi, &i.CatalogID, &i.Name, &i.Brand, &i.CategoryName, &i.CategoryCode, &i.Condition, &i.RoomID, &i.PurchaseDate,
			&i.SumberPerolehan, &i.Keterangan, &i.MasaManfaatTahun, &i.NilaiBuku, &i.StatusAset, &i.TanggalPencatatan, &i.TanggalPenghapusan, &i.AlasanPenghapusan,
		); err != nil {
			return nil, err
		}
		list = append(list, i)
	}
	return list, nil
}

func (r *pgxInventoryItemRepository) Create(ctx context.Context, item *InventoryItem, semesterID string) error {
	// Generate Nomor Registrasi
	year := item.PurchaseDate.Year()
	var countReg int
	r.db.QueryRow(ctx, `SELECT COUNT(*) FROM inventory_items WHERE nomor_registrasi IS NOT NULL AND nomor_registrasi != '' AND EXTRACT(YEAR FROM purchase_date) = $1`, year).Scan(&countReg)
	
	romanMonths := []string{"", "I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X", "XI", "XII"}
	monthRoman := romanMonths[int(item.PurchaseDate.Month())]
	
	schoolCode, _ := r.GetSchoolCode(ctx)
	regNumber := fmt.Sprintf("%03d/INV/%s/%s/%d", countReg+1, schoolCode, monthRoman, year)

	query := `INSERT INTO inventory_items (id, nomor_registrasi, catalog_id, condition, room_id, purchase_date, sumber_perolehan, status_aset, tanggal_pencatatan) VALUES ($1, $2, $3, $4, $5, $6, $7, 'AKTIF', NOW())`
	_, err := r.db.Exec(ctx, query, item.ID, regNumber, item.CatalogID, item.Condition, item.RoomID, item.PurchaseDate, item.SumberPerolehan)
	return err
}

func (r *pgxInventoryItemRepository) Update(ctx context.Context, item *InventoryItem) error {
	query := `UPDATE inventory_items SET condition=$1, status_aset=$2 WHERE id=$3`
	_, err := r.db.Exec(ctx, query, item.Condition, item.StatusAset, item.ID)
	return err
}

func (r *pgxInventoryItemRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM inventory_items WHERE id = $1`, id)
	return err
}

func (r *pgxInventoryItemRepository) MutateItem(ctx context.Context, item *InventoryItem, action string, oldData string, newData string, userID string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil { return err }
	defer tx.Rollback(ctx)

	query := `UPDATE inventory_items SET room_id=$1, condition=$2, status_aset=$3, tanggal_penghapusan=$4, alasan_penghapusan=$5 WHERE id=$6`
	_, err = tx.Exec(ctx, query, item.RoomID, item.Condition, item.StatusAset, item.TanggalPenghapusan, item.AlasanPenghapusan, item.ID)
	if err != nil { return err }

	_, err = tx.Exec(ctx, `INSERT INTO inventory_item_history (inventory_item_id, action, old_data, new_data, description, changed_by) VALUES ($1, $2, $3, $4, $5, $6)`, item.ID, action, oldData, newData, action, userID)
	if err != nil { return err }

	return tx.Commit(ctx)
}

func (r *pgxInventoryItemRepository) GetHistory(ctx context.Context, itemID string) ([]InventoryItemHistory, error) {
	rows, err := r.db.Query(ctx, `SELECT id::text, inventory_item_id, action, COALESCE(old_data::text, ''), COALESCE(new_data::text, ''), COALESCE(description, ''), changed_by, changed_at FROM inventory_item_history WHERE inventory_item_id = $1 ORDER BY changed_at DESC`, itemID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []InventoryItemHistory
	for rows.Next() {
		var h InventoryItemHistory
		rows.Scan(&h.ID, &h.InventoryItemID, &h.Action, &h.OldData, &h.NewData, &h.Description, &h.ChangedBy, &h.ChangedAt)
		list = append(list, h)
	}
	return list, nil
}

// ==========================================
// MASTER DATA
// ==========================================

func (r *pgxInventoryItemRepository) GetMasterCategories(ctx context.Context) ([]MasterCategory, error) {
	rows, err := r.db.Query(ctx, `SELECT id::text, code, name FROM master_categories ORDER BY name ASC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []MasterCategory
	for rows.Next() {
		var c MasterCategory
		rows.Scan(&c.ID, &c.Code, &c.Name)
		list = append(list, c)
	}
	return list, nil
}

func (r *pgxInventoryItemRepository) CreateCategory(ctx context.Context, cat *MasterCategory) error {
	if cat.ID == "" { cat.ID = uuid.New().String() }
	_, err := r.db.Exec(ctx, `INSERT INTO master_categories (id, code, name) VALUES ($1, $2, $3)`, cat.ID, cat.Code, cat.Name)
	return err
}

func (r *pgxInventoryItemRepository) UpdateCategory(ctx context.Context, cat *MasterCategory) error {
	_, err := r.db.Exec(ctx, `UPDATE master_categories SET code=$1, name=$2 WHERE id=$3`, cat.Code, cat.Name, cat.ID)
	return err
}

func (r *pgxInventoryItemRepository) DeleteCategory(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM master_categories WHERE id = $1`, id)
	return err
}

func (r *pgxInventoryItemRepository) GetMasterCatalogs(ctx context.Context) ([]MasterCatalog, error) {
	rows, err := r.db.Query(ctx, `SELECT id::text, category_id::text, name, COALESCE(brand, ''), COALESCE(image_url, '') FROM master_catalogs ORDER BY name ASC`)
	if err != nil { return nil, err }
	defer rows.Close()
	var list []MasterCatalog
	for rows.Next() {
		var c MasterCatalog
		rows.Scan(&c.ID, &c.CategoryID, &c.Name, &c.Brand, &c.ImageURL)
		list = append(list, c)
	}
	return list, nil
}

func (r *pgxInventoryItemRepository) CreateCatalog(ctx context.Context, cat *MasterCatalog) error {
	if cat.ID == "" { cat.ID = uuid.New().String() }
	_, err := r.db.Exec(ctx, `INSERT INTO master_catalogs (id, category_id, name, brand, image_url) VALUES ($1, $2, $3, $4, $5)`, cat.ID, cat.CategoryID, cat.Name, cat.Brand, cat.ImageURL)
	return err
}

func (r *pgxInventoryItemRepository) UpdateCatalog(ctx context.Context, cat *MasterCatalog) error {
	_, err := r.db.Exec(ctx, `UPDATE master_catalogs SET category_id=$1, name=$2, brand=$3, image_url=$4 WHERE id=$5`, cat.CategoryID, cat.Name, cat.Brand, cat.ImageURL, cat.ID)
	return err
}

func (r *pgxInventoryItemRepository) DeleteCatalog(ctx context.Context, id string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM master_catalogs WHERE id = $1`, id)
	return err
}

func (r *pgxInventoryItemRepository) GetSchoolCode(ctx context.Context) (string, error) {
	var code string
	err := r.db.QueryRow(ctx, `SELECT COALESCE(kode_sekolah, 'SMKPasja') FROM tabel_sekolah LIMIT 1`).Scan(&code)
	return code, err
}

// Fungsi getPrefix untuk memetakan sumber ke prefix
func getPrefix(sourceType string) string {
	switch strings.ToUpper(sourceType) {
	case "PEMBELIAN", "PROCUREMENT":
		return "P"
	case "HIBAH":
		return "H"
	case "DONASI":
		return "D"
	case "MUTASI":
		return "M"
	default:
		return ""
	}
}

func (r *pgxInventoryItemRepository) CountItemsForAutoID(ctx context.Context, categoryID string, year int, sourceType string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM inventory_items i 
	          JOIN master_catalogs c ON i.catalog_id = c.id 
	          WHERE c.category_id = $1 AND EXTRACT(YEAR FROM i.purchase_date) = $2`
	
	args := []interface{}{categoryID, year}
	
	// Jika sourceType tidak kosong, filter berdasarkan sumber_perolehan
	if sourceType != "" {
		query += ` AND UPPER(i.sumber_perolehan) = UPPER($3)`
		args = append(args, sourceType)
	}
	
	err := r.db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *pgxInventoryItemRepository) GenerateID(ctx context.Context, catalogID string, year int, sourceType string) (string, error) {
	var categoryID, categoryCode string
	err := r.db.QueryRow(ctx, `SELECT category_id::text, (SELECT code FROM master_categories WHERE id = category_id) FROM master_catalogs WHERE id = $1`, catalogID).Scan(&categoryID, &categoryCode)
	if err != nil { return "", err }
	
	schoolCode, _ := r.GetSchoolCode(ctx)
	count, _ := r.CountItemsForAutoID(ctx, categoryID, year, sourceType)
	
	prefix := getPrefix(sourceType)
	if prefix != "" {
		return fmt.Sprintf("%s/%s/%s/%d/%03d", schoolCode, prefix, categoryCode, year, count+1), nil
	}
	
	// Fallback jika tidak ada prefix
	return fmt.Sprintf("%s/%s/%d/%03d", schoolCode, categoryCode, year, count+1), nil
}

func (r *pgxInventoryItemRepository) GetConditionBook(ctx context.Context, semesterID string) ([]*ConditionBookEntry, error) {
	query := `
		SELECT i.id, c.name, COALESCE(c.brand,''), cat.name, COALESCE(k.nm_kelas,''), COALESCE(i.condition,''), 
		       h.changed_at, COALESCE(h.description,'')
		FROM inventory_items i
		JOIN master_catalogs c ON i.catalog_id = c.id
		JOIN master_categories cat ON c.category_id = cat.id
		LEFT JOIN tabel_kelas k ON i.room_id = k.rombongan_belajar_id
		LEFT JOIN LATERAL (
			SELECT changed_at, description
			FROM inventory_item_history
			WHERE inventory_item_id = i.id AND action = 'UBAH_KONDISI'
			ORDER BY changed_at DESC LIMIT 1
		) h ON true
		WHERE i.status_aset IN ('AKTIF', 'DIPINJAM', 'DALAM_PERBAIKAN', 'RUSAK_BERAT')
		AND ($1 = '' OR i.semester_id = $1)
		ORDER BY i.id
	`
	rows, err := r.db.Query(ctx, query, semesterID)
	if err != nil { return nil, err }
	defer rows.Close()

	var list []*ConditionBookEntry
	for rows.Next() {
		var e ConditionBookEntry
		if err := rows.Scan(&e.ID, &e.Name, &e.Brand, &e.CategoryName, &e.RoomName, &e.Condition, &e.LastChangeDate, &e.LastChangeNote); err != nil {
			return nil, err
		}
		list = append(list, &e)
	}
	return list, nil
}
