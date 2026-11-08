package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerSarprasRoutes(mux *http.ServeMux) {
	// Logistik
	mux.HandleFunc("/api/v1/logistics", r.wrapWithContext(r.handleLogistics))
	mux.HandleFunc("/api/v1/logistics/", r.wrapWithContext(r.handleLogisticByID))

	// KIR
	mux.HandleFunc("/api/v1/sarpras/kir/rooms", r.wrapWithContext(r.handleKIRRooms))
	mux.HandleFunc("/api/v1/sarpras/kir/items", r.wrapWithContext(r.handleKIRItems))
	mux.HandleFunc("/api/v1/sarpras/kir/items/", r.wrapWithContext(r.handleKIRItemByID))
	mux.HandleFunc("/api/v1/sarpras/kir/print-qr", r.wrapWithContext(r.handlePrintQR))

	// Master Data
	mux.HandleFunc("/api/v1/sarpras/master/categories", r.wrapWithContext(r.handleMasterCategories))
	mux.HandleFunc("/api/v1/sarpras/master/categories/", r.wrapWithContext(r.handleMasterCategoryByID))
	mux.HandleFunc("/api/v1/sarpras/master/catalogs", r.wrapWithContext(r.handleMasterCatalogs))
	mux.HandleFunc("/api/v1/sarpras/master/catalogs/", r.wrapWithContext(r.handleMasterCatalogByID))
	
	// Endpoint Upload Foto Katalog
	mux.HandleFunc("/api/v1/sarpras/master/catalogs/upload", r.wrapWithContext(r.handleCatalogImageUpload))
}

func (r *Router) handleLogistics(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		data, err := r.logisticService.ListLogistics(req.Context())
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]repository.Logistic, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
	case http.MethodPost:
		var body repository.Logistic
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		if err := r.logisticService.CreateLogistic(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Item logistik berhasil ditambahkan"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleLogisticByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/logistics/")
	parts := splitPath(path)
	if len(parts) < 1 { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	id := parts[0]

	if len(parts) == 2 {
		action := parts[1]
		if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }

		var body struct { Jumlah int `json:"jumlah"`; Keterangan string `json:"keterangan"` }
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }

		switch action {
		case "restock":
			if err := r.logisticService.RestockLogistic(req.Context(), id, body.Jumlah, body.Keterangan); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, map[string]string{"message": "Restock berhasil"})
		case "use":
			if err := r.logisticService.UseLogistic(req.Context(), id, body.Jumlah, body.Keterangan); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, map[string]string{"message": "Penggunaan berhasil dicatat"})
		default:
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid action"})
		}
		return
	}
	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method or path not allowed"})
}

func (r *Router) handleKIRRooms(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	data, err := r.kirService.GetRooms(req.Context())
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
	if data == nil { data = make([]service.KIRRoom, 0) }
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

func (r *Router) handleKIRItems(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		roomID := req.URL.Query().Get("room_id")
		if roomID != "" {
			data, err := r.kirService.GetItemsByRoom(req.Context(), roomID)
			if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			if data == nil { data = make([]repository.InventoryItem, 0) }
			writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
		} else {
			data, err := r.kirService.GetAllItems(req.Context())
			if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			if data == nil { data = make([]repository.InventoryItem, 0) }
			writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
		}

	case http.MethodPost:
		var body struct {
			CatalogID    string `json:"catalog_id"`
			Condition    string `json:"condition"`
			RoomID       string `json:"room_id"`
			PurchaseDate string `json:"purchase_date"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		
		pDate, _ := time.Parse("2006-01-02", body.PurchaseDate)
		item := repository.InventoryItem{
			CatalogID: body.CatalogID,
			Condition: body.Condition,
			RoomID:    body.RoomID,
			PurchaseDate: pDate,
		}

		if err := r.kirService.CreateItem(req.Context(), &item); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Aset berhasil ditambahkan dengan ID: " + item.ID})
	
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleKIRItemByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/kir/items/")
	if id == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }

	switch req.Method {
	case http.MethodPut:
		var body repository.InventoryItem
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		body.ID = id
		if err := r.kirService.UpdateItem(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Aset berhasil diupdate"})
	
	case http.MethodDelete:
		if err := r.kirService.DeleteItem(req.Context(), id); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Aset berhasil dihapus"})
	
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handlePrintQR(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	var body struct { ItemIDs []string `json:"item_ids"` }
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
	if len(body.ItemIDs) == 0 { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "No items selected"}); return }

	pdfBuf, err := r.kirService.GenerateStickerPDF(req.Context(), body.ItemIDs)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"Stiker_QR_Aset.pdf\"")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuf.Bytes())
}

// ==========================================
// MASTER DATA HANDLERS
// ==========================================

func (r *Router) handleMasterCategories(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		data, err := r.kirService.GetMasterCategories(req.Context())
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]repository.MasterCategory, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
	case http.MethodPost:
		var body repository.MasterCategory
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		if err := r.kirService.CreateCategory(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Kategori berhasil ditambahkan"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleMasterCategoryByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/master/categories/")
	if id == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	switch req.Method {
	case http.MethodPut:
		var body repository.MasterCategory
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		body.ID = id
		if err := r.kirService.UpdateCategory(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Kategori diupdate"})
	case http.MethodDelete:
		if err := r.kirService.DeleteCategory(req.Context(), id); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Kategori dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleMasterCatalogs(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		data, err := r.kirService.GetMasterCatalogs(req.Context())
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]repository.MasterCatalog, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
	case http.MethodPost:
		var body repository.MasterCatalog
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		if err := r.kirService.CreateCatalog(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Katalog berhasil ditambahkan"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleMasterCatalogByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/master/catalogs/")
	if id == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }
	switch req.Method {
	case http.MethodPut:
		var body repository.MasterCatalog
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		body.ID = id
		if err := r.kirService.UpdateCatalog(req.Context(), &body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Katalog diupdate"})
	case http.MethodDelete:
		if err := r.kirService.DeleteCatalog(req.Context(), id); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		writeJSON(w, http.StatusOK, map[string]string{"message": "Katalog dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleCatalogImageUpload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	
	err := req.ParseMultipartForm(10 << 20) // Max 10MB
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"})
		return
	}
	
	file, handler, err := req.FormFile("image")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File 'image' tidak ditemukan"})
		return
	}
	defer file.Close()

	objectName := fmt.Sprintf("catalog_%d_%s", time.Now().Unix(), handler.Filename)
	// Menggunakan bucket 'documents' sesuai pola PTK
	url, err := r.minioClient.UploadFile(req.Context(), "documents", objectName, file, handler)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}
