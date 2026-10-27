package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
	"strconv"
)

// enrichContext mengekstrak header dan menyuntikkan OperationalContext ke dalam context
// func (r *Router) enrichContext(req *http.Request) (context.Context, error) {
// 	schoolID := req.Header.Get("x-school-id")
// 	academicYear := req.Header.Get("x-academic-year")
// 	semesterStr := req.Header.Get("x-semester")
// 	academicPeriodID := req.Header.Get("x-academic-period")
// 	userID := req.Header.Get("x-user-id")
// 	role := req.Header.Get("x-user-role")

// 	if schoolID == "" || academicYear == "" || userID == "" || role == "" {
// 		return nil, http.ErrAbortHandler
// 	}

// 	semester := 1
// 	if semesterStr != "" {
// 		if s, err := strconv.Atoi(semesterStr); err == nil {
// 			semester = s
// 		}
// 	}

// 	opCtx := &platformctx.OperationalContext{
// 		SchoolID:         schoolID,
// 		AcademicYear:     academicYear,
// 		Semester:         semester,
// 		AcademicPeriodID: academicPeriodID,
// 		UserID:           userID,
// 		Role:             role,
// 	}

// 	ctx := middleware.WithOperationalContext(req.Context(), opCtx)
// 	return ctx, nil
// }

// ============================================================
// CLASS HANDLERS (BARU)
// ============================================================

func (r *Router) handleClasses(w http.ResponseWriter, req *http.Request) {
	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodPost:
		r.createClass(w, req.WithContext(ctx))
	case http.MethodGet:
		r.listClasses(w, req.WithContext(ctx))
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) createClass(w http.ResponseWriter, req *http.Request) {
	var reqBody service.CreateClassRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	class, err := r.classService.CreateClass(req.Context(), reqBody)
	if err != nil {
		log.Printf("❌ [REST] Error creating class: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(class)
}

func (r *Router) listClasses(w http.ResponseWriter, req *http.Request) {
	limitStr := req.URL.Query().Get("limit")
	offsetStr := req.URL.Query().Get("offset")
	status := req.URL.Query().Get("status")
	keyword := req.URL.Query().Get("keyword")

	limit, _ := strconv.Atoi(limitStr)
	if limit <= 0 {
		limit = 20
	}
	offset, _ := strconv.Atoi(offsetStr)
	if offset < 0 {
		offset = 0
	}

	filter := repository.ClassFilter{
		Status:  status,
		Keyword: keyword,
		Limit:   limit,
		Offset:  offset,
	}

	classes, err := r.classService.ListClasses(req.Context(), filter)
	if err != nil {
		log.Printf("❌ [REST] Error listing classes: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": classes,
		// "total": total,
	})
}

func (r *Router) handleClassByID(w http.ResponseWriter, req *http.Request) {
	// Parse ID dari path: /api/v1/classes/{id} atau /api/v1/classes/{id}/...
	path := req.URL.Path
	// sederhana: ambil setelah "/api/v1/classes/"
	id := path[len("/api/v1/classes/"):]
	// jika masih ada slash, ambil bagian pertama
	for i, ch := range id {
		if ch == '/' {
			id = id[:i]
			break
		}
	}
	if id == "" {
		http.Error(w, `{"error": "Class ID is required"}`, http.StatusBadRequest)
		return
	}

	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		r.getClass(w, req.WithContext(ctx), id)
	case http.MethodPut:
		r.updateClass(w, req.WithContext(ctx), id)
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) getClass(w http.ResponseWriter, req *http.Request, classID string) {
	class, err := r.classService.GetClass(req.Context(), classID)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, `{"error": "Class not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("❌ [REST] Error getting class: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}

func (r *Router) updateClass(w http.ResponseWriter, req *http.Request, classID string) {
	var reqBody service.UpdateClassRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	class, err := r.classService.UpdateClass(req.Context(), classID, reqBody)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, `{"error": "Class not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("❌ [REST] Error updating class: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(class)
}
