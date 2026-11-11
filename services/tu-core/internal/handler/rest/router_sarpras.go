package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerSarprasRoutes(mux *http.ServeMux) {
	// Logistik
	mux.HandleFunc("/api/v1/logistics", r.wrapWithContext(r.handleLogistics))
	mux.HandleFunc("/api/v1/logistics/", r.wrapWithContext(r.handleLogisticByID))

	// KIR & Master Data
	mux.HandleFunc("/api/v1/sarpras/kir/rooms", r.wrapWithContext(r.handleKIRRooms))
	mux.HandleFunc("/api/v1/sarpras/kir/items", r.wrapWithContext(r.handleKIRItems))
	mux.HandleFunc("/api/v1/sarpras/kir/items/", r.wrapWithContext(r.handleKIRItemByID))
	mux.HandleFunc("/api/v1/sarpras/kir/print-qr", r.wrapWithContext(r.handlePrintQR))
	mux.HandleFunc("/api/v1/sarpras/condition-book", r.wrapWithContext(r.handleConditionBook))
	mux.HandleFunc("/api/v1/sarpras/kir/report/", r.wrapWithContext(r.handleKIRReport))

	// Master Data
	mux.HandleFunc("/api/v1/sarpras/master/categories", r.wrapWithContext(r.handleMasterCategories))
	mux.HandleFunc("/api/v1/sarpras/master/categories/", r.wrapWithContext(r.handleMasterCategoryByID))
	mux.HandleFunc("/api/v1/sarpras/master/catalogs", r.wrapWithContext(r.handleMasterCatalogs))
	mux.HandleFunc("/api/v1/sarpras/master/catalogs/", r.wrapWithContext(r.handleMasterCatalogByID))
	mux.HandleFunc("/api/v1/sarpras/master/catalogs/upload", r.wrapWithContext(r.handleCatalogImageUpload))

	// Mutasi & Laporan Buku Induk
	mux.HandleFunc("/api/v1/sarpras/inventory/report", r.wrapWithContext(r.handleBukuIndukReport))
	
	// Handler dinamis untuk mutate & history (menggunakan prefix)
	mux.HandleFunc("/api/v1/sarpras/inventory/", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/mutate") {
			r.handleInventoryMutate(w, req)
		} else if strings.HasSuffix(req.URL.Path, "/history") {
			r.handleInventoryHistory(w, req)
		} else {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "Endpoint not found"})
		}
	}))
}

// ... (fungsi handleLogistics, handleLogisticByID, handleKIR* dll tetap sama seperti sebelumnya) ...
// Saya sertakan kembali di bawah ini agar file utuh.

func (r *Router) handleLogistics(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		opCtx := middleware.GetOperationalContext(req.Context())
		semesterID := "20261"
		if opCtx != nil { semesterID = opCtx.AcademicPeriodID }
		data, err := r.logisticService.ListLogistics(req.Context(), semesterID)
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		if data == nil { data = make([]repository.Logistic, 0) }
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
	case http.MethodPost:
		var body repository.Logistic
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
		opCtxPost := middleware.GetOperationalContext(req.Context())
		semesterIDPost := "20261"
		if opCtxPost != nil { semesterIDPost = opCtxPost.AcademicPeriodID }
		if err := r.logisticService.CreateLogistic(req.Context(), &body, semesterIDPost); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
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
		opCtx := middleware.GetOperationalContext(req.Context())
		semesterID := "20261"
		if opCtx != nil { semesterID = opCtx.AcademicPeriodID }
		
		item := repository.InventoryItem{ CatalogID: body.CatalogID, Condition: body.Condition, RoomID: body.RoomID, PurchaseDate: pDate }
		if err := r.kirService.CreateItem(req.Context(), &item, semesterID); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
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
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	err := req.ParseMultipartForm(10 << 20)
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"}); return }
	file, handler, err := req.FormFile("image")
	if err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File 'image' tidak ditemukan"}); return }
	defer file.Close()

	objectName := fmt.Sprintf("catalog_%d_%s", time.Now().Unix(), handler.Filename)
	url, err := r.minioClient.UploadFile(req.Context(), "documents", objectName, file, handler)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

// ==========================================
// HANDLER MUTASI & LAPORAN (Buku Induk)
// ==========================================

func (r *Router) handleInventoryMutate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/inventory/")
	id = strings.TrimSuffix(id, "/mutate")
	if id == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }

	var body struct {
		Action       string `json:"action"`
		NewRoomID    string `json:"new_room_id"`
		NewCondition string `json:"new_condition"`
		StatusAset   string `json:"status_aset"`
		Reason       string `json:"reason"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }

	opCtx := middleware.GetOperationalContext(req.Context())
	userID := "system"
	if opCtx != nil { userID = opCtx.UserID }

	if err := r.kirService.MutateItem(req.Context(), id, body.Action, body.NewRoomID, body.NewCondition, body.StatusAset, body.Reason, userID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Mutasi berhasil dicatat"})
}

func (r *Router) handleInventoryHistory(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/inventory/")
	id = strings.TrimSuffix(id, "/history")
	if id == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"}); return }

	data, err := r.kirService.GetItemHistory(req.Context(), id)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
	if data == nil { data = make([]repository.InventoryItemHistory, 0) }
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}

func (r *Router) handleBukuIndukReport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	
	pdfBuf, err := r.kirService.GenerateBukuIndukPDF(req.Context())
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"Buku_Induk_Inventaris.pdf\"")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuf.Bytes())
}

func (r *Router) handleKIRReport(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	
	roomID := strings.TrimPrefix(req.URL.Path, "/api/v1/sarpras/kir/report/")
	if roomID == "" { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Room ID required"}); return }

	pdfBuf, err := r.kirService.GenerateKIRReportPDF(req.Context(), roomID)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\"KIR_Report.pdf\"")
	w.WriteHeader(http.StatusOK)
	w.Write(pdfBuf.Bytes())
}

func (r *Router) handleConditionBook(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet { writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"}); return }
	
	isPDF := req.URL.Query().Get("format") == "pdf"
	
	if isPDF {
		pdfBuf, err := r.kirService.GenerateConditionBookPDF(req.Context())
		if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"Buku_Keadaan_Barang.pdf\"")
		w.WriteHeader(http.StatusOK)
		w.Write(pdfBuf.Bytes())
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	semesterID := "20261"
	if opCtx != nil { semesterID = opCtx.AcademicPeriodID }

	data, err := r.kirService.GetConditionBook(req.Context(), semesterID)
	if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
	if data == nil { data = make([]*repository.ConditionBookEntry, 0) }
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})
}
