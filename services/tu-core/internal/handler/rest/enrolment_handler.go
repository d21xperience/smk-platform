package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

// ============================================================
// ENROLLMENT HANDLERS (BARU)
// ============================================================

func (r *Router) handleClassMembers(w http.ResponseWriter, req *http.Request) {
	// Parse path: /api/v1/classes/{id}/members atau /api/v1/classes/{id}/members/{studentId}
	path := req.URL.Path
	// Ambil ID class: setelah "/api/v1/classes/"
	rest := path[len("/api/v1/classes/"):]
	parts := splitPath(rest)
	if len(parts) < 1 || parts[0] == "" {
		http.Error(w, `{"error": "Class ID is required"}`, http.StatusBadRequest)
		return
	}
	classID := parts[0]

	// Cek apakah ada subpath "members"
	if len(parts) < 2 || parts[1] != "members" {
		http.Error(w, `{"error": "Invalid endpoint"}`, http.StatusNotFound)
		return
	}

	ctx, err := r.enrichContext(req)
	if err != nil {
		http.Error(w, `{"error": "Missing operational context headers"}`, http.StatusBadRequest)
		return
	}

	// Jika ada studentId: DELETE /classes/:id/members/:studentId
	if len(parts) >= 3 && parts[2] != "" {
		if req.Method == http.MethodDelete {
			r.removeMember(w, req.WithContext(ctx), classID, parts[2])
			return
		}
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	switch req.Method {
	case http.MethodGet:
		r.listClassMembers(w, req.WithContext(ctx), classID)
	case http.MethodPost:
		r.addMember(w, req.WithContext(ctx), classID)
	default:
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func (r *Router) listClassMembers(w http.ResponseWriter, req *http.Request, classID string) {
	members, err := r.enrollmentService.GetClassMembers(req.Context(), classID)
	if err != nil {
		log.Printf("❌ [REST] Error listing class members: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data":  members,
		"total": len(members),
	})
}

func (r *Router) addMember(w http.ResponseWriter, req *http.Request, classID string) {
	var reqBody service.AddMemberRequest
	if err := json.NewDecoder(req.Body).Decode(&reqBody); err != nil {
		http.Error(w, `{"error": "Invalid request body"}`, http.StatusBadRequest)
		return
	}

	enrollment, err := r.enrollmentService.AddMember(req.Context(), classID, reqBody)
	if err != nil {
		if err == repository.ErrCapacityFull {
			http.Error(w, `{"error": "Class capacity full"}`, http.StatusBadRequest)
			return
		}
		log.Printf("❌ [REST] Error adding member: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrollment)
}

func (r *Router) removeMember(w http.ResponseWriter, req *http.Request, classID, studentID string) {
	err := r.enrollmentService.RemoveMember(req.Context(), classID, studentID)
	if err != nil {
		if err == repository.ErrNotFound {
			http.Error(w, `{"error": "Active enrollment not found"}`, http.StatusNotFound)
			return
		}
		log.Printf("❌ [REST] Error removing member: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Member removed successfully"})
}

func (r *Router) handleStudentHistory(w http.ResponseWriter, req *http.Request) {
	// Parse path: /api/v1/students/{id}/history
	path := req.URL.Path
	rest := path[len("/api/v1/students/"):]
	parts := splitPath(rest)
	if len(parts) < 2 || parts[0] == "" || parts[1] != "history" {
		http.Error(w, `{"error": "Invalid endpoint"}`, http.StatusNotFound)
		return
	}
	studentID := parts[0]

	if req.Method != http.MethodGet {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	history, err := r.enrollmentService.GetStudentHistory(req.Context(), studentID)
	if err != nil {
		log.Printf("❌ [REST] Error getting student history: %v", err)
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"data": history,
	})
}
