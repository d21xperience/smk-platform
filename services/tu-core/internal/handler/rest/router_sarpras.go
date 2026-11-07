package rest

import (
	"encoding/json"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/repository"
)

func (r *Router) registerSarprasRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/logistics", r.wrapWithContext(r.handleLogistics))
	mux.HandleFunc("/api/v1/logistics/", r.wrapWithContext(r.handleLogisticByID))
}

func (r *Router) handleLogistics(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		data, err := r.logisticService.ListLogistics(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		if data == nil {
			data = make([]repository.Logistic, 0)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": data})

	case http.MethodPost:
		var body repository.Logistic
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.logisticService.CreateLogistic(req.Context(), &body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Item logistik berhasil ditambahkan"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleLogisticByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/logistics/")
	parts := splitPath(path)
	if len(parts) < 1 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]

	if len(parts) == 2 {
		action := parts[1]
		if req.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			return
		}

		var body struct {
			Jumlah     int    `json:"jumlah"`
			Keterangan string `json:"keterangan"`
		}
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}

		if action == "restock" {
			if err := r.logisticService.RestockLogistic(req.Context(), id, body.Jumlah, body.Keterangan); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Restock berhasil"})
			return
		} else if action == "use" {
			if err := r.logisticService.UseLogistic(req.Context(), id, body.Jumlah, body.Keterangan); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Penggunaan barang berhasil dicatat"})
			return
		}
	}

	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method or path not allowed"})
}
