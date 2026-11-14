package rest

import (
	"encoding/json"
	"net/http"

	"sekolah-platform/services/tu-core/internal/domain/finance/models"
)

// registerKeuanganTTDRoutes — endpoint PIN TTD digital.
//
// User ID diambil dari header `X-User-ID` (sesuai konvensi semua endpoint
// keuangan lain yang belum pakai auth middleware terpusat).
// Untuk `unlock`, admin menentukan user_id target via body JSON.
func (r *Router) registerKeuanganTTDRoutes(mux *http.ServeMux) {
	mux.HandleFunc(financeTransaksiBasePath+"/ttd/pin/setup", r.handleTTDPinSetup)
	mux.HandleFunc(financeTransaksiBasePath+"/ttd/pin/verify", r.handleTTDPinVerify)
	mux.HandleFunc(financeTransaksiBasePath+"/ttd/pin/change", r.handleTTDPinChange)
	mux.HandleFunc(financeTransaksiBasePath+"/ttd/pin/unlock", r.handleTTDPinUnlock)
	mux.HandleFunc(financeTransaksiBasePath+"/ttd/pin/status", r.handleTTDPinStatus)

	// FASE 5B-2: E-Sign penutupan + verify publik
	mux.HandleFunc(financeTransaksiBasePath+"/penutupan/{id}/sign", r.handlePenutupanSign)
	mux.HandleFunc(financeTransaksiBasePath+"/verify/{hash}", r.handleTTDVerify)
}

// resolveUserID — ambil user_id dari header X-User-ID
func resolveUserID(req *http.Request) string {
	return req.Header.Get("X-User-ID")
}

// =====================================================================
// POST /ttd/pin/setup
// =====================================================================
func (r *Router) handleTTDPinSetup(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := resolveUserID(req)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "header X-User-ID wajib diisi"})
		return
	}

	var payload models.SetupPinRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.UserID = userID

	item, err := r.userPinTTDService.SetupPin(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user_id":  item.UserID,
			"is_setup": true,
			"message":  "PIN TTD berhasil disetup",
		},
	})
}

// =====================================================================
// POST /ttd/pin/verify
// =====================================================================
func (r *Router) handleTTDPinVerify(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := resolveUserID(req)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "header X-User-ID wajib diisi"})
		return
	}

	var payload models.VerifyPinRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.UserID = userID

	resp, err := r.userPinTTDService.VerifyPin(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": resp})
}

// =====================================================================
// PUT /ttd/pin/change
// =====================================================================
func (r *Router) handleTTDPinChange(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := resolveUserID(req)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "header X-User-ID wajib diisi"})
		return
	}

	var payload models.ChangePinRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.UserID = userID

	if err := r.userPinTTDService.ChangePin(req.Context(), payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "PIN berhasil diganti",
	})
}

// =====================================================================
// POST /ttd/pin/unlock
//   Body: {"user_id": "<target user>"}  — admin action
// =====================================================================
func (r *Router) handleTTDPinUnlock(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var payload models.UnlockPinRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}

	// Fallback: kalau body kosong, coba header X-User-ID
	if payload.UserID == "" {
		payload.UserID = resolveUserID(req)
	}
	if payload.UserID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "user_id wajib diisi (body atau header X-User-ID)"})
		return
	}

	if err := r.userPinTTDService.Unlock(req.Context(), payload.UserID); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "PIN berhasil di-unlock",
	})
}

// =====================================================================
// GET /ttd/pin/status
// =====================================================================
func (r *Router) handleTTDPinStatus(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	userID := resolveUserID(req)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "header X-User-ID wajib diisi"})
		return
	}

	status, err := r.userPinTTDService.GetStatus(req.Context(), userID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"success": true, "data": status})
}

// =====================================================================
// POST /penutupan/{id}/sign — TTD penutupan oleh Kepala Sekolah
//
// Flow:
//   1. Verifikasi PIN user via userPinTTDService.VerifyPin
//   2. Kalau valid → panggil penutupanService.SignPenutupan
//   3. Return dokumen yang sudah di-sign
//
// Body: {"pin": "123456"}
// =====================================================================
func (r *Router) handlePenutupanSign(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	schoolID := resolveSchoolID(req)
	if schoolID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "school_id wajib diisi"})
		return
	}
	signerID := resolveUserID(req)
	if signerID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "header X-User-ID wajib diisi"})
		return
	}
	id := req.PathValue("id")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id wajib diisi"})
		return
	}

	var payload models.VerifyPinRequest
	if err := json.NewDecoder(req.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "format JSON tidak valid: " + err.Error()})
		return
	}
	payload.UserID = signerID

	// 1. Verify PIN
	verifyResp, err := r.userPinTTDService.VerifyPin(req.Context(), payload)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if !verifyResp.Valid {
		writeJSON(w, http.StatusUnauthorized, map[string]interface{}{
			"error":  "PIN tidak valid",
			"detail": verifyResp,
		})
		return
	}

	// 2. Sign
	signed, err := r.penutupanService.SignPenutupan(req.Context(), schoolID, id, signerID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Penutupan berhasil ditandatangani",
		"data":    signed,
	})
}

// =====================================================================
// GET /verify/{hash} — verifikasi publik dokumen via hash TTD
//
// Endpoint publik: TIDAK butuh auth. Berguna untuk QR code di dokumen.
// =====================================================================
func (r *Router) handleTTDVerify(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	hash := req.PathValue("hash")
	if hash == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "hash wajib diisi"})
		return
	}

	doc, err := r.penutupanService.VerifyByTTDHash(req.Context(), hash)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]interface{}{
			"success": false,
			"valid":   false,
			"error":   "dokumen dengan hash tersebut tidak ditemukan",
		})
		return
	}

	// Response sengaja disederhanakan — tidak expose data sensitif seperti
	// catatan internal, hanya data inti untuk verifikasi publik.
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"valid":   true,
		"data": map[string]interface{}{
			"doc_type":         "PENUTUPAN_BUKU_KAS",
			"doc_id":           doc.ID,
			"periode_bulan":    doc.PeriodeBulan,
			"periode_tahun":    doc.PeriodeTahun,
			"status":           doc.Status,
			"saldo_akhir_buku": doc.SaldoAkhirBuku,
			"signed_by":        doc.TTDKepsekNama,
			"signed_at":        doc.TTDKepsekAt,
			"hash":             doc.TTDKepsekHash,
		},
	})
}

