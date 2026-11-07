package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
	"sekolah-platform/services/tu-core/internal/repository"
	"sekolah-platform/services/tu-core/internal/service"
)

func (r *Router) registerPtkRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/teachers", r.wrapWithContext(r.handleTeachers))
	mux.HandleFunc("/api/v1/teachers/workload", r.wrapWithContext(r.handleTeacherWorkload))
	mux.HandleFunc("/api/v1/teachers/photo/", r.wrapWithContext(r.handleTeacherPhoto))
	mux.HandleFunc("/api/v1/teachers/history/upload", r.wrapWithContext(r.handleHistoryDocumentUpload))
	mux.HandleFunc("/api/v1/teachers/history/", r.wrapWithContext(r.handleTeacherHistory))
	
	// Rute Cuti & Kinerja
	mux.HandleFunc("/api/v1/teachers/leaves", r.wrapWithContext(r.handleTeacherLeaves))
	mux.HandleFunc("/api/v1/teachers/leaves/", r.wrapWithContext(r.handleTeacherLeaveByID))
	mux.HandleFunc("/api/v1/teachers/performances", r.wrapWithContext(r.handleTeacherPerformances))
	mux.HandleFunc("/api/v1/teachers/performances/", r.wrapWithContext(r.handleTeacherPerformanceByID))

	mux.HandleFunc("/api/v1/teachers/", r.wrapWithContext(r.handleTeacherByID))
}

func (r *Router) handleTeacherPerformances(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		list, err := r.ptkPerformanceService.ListPerformances(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var body repository.Performance
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.ptkPerformanceService.CreatePerformance(req.Context(), body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Penilaian kinerja berhasil ditambahkan"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleTeacherPerformanceByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/performances/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	if req.Method == http.MethodDelete {
		if err := r.ptkPerformanceService.DeletePerformance(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Penilaian kinerja dihapus"})
		return
	}

	writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
}

// --- HANDLER CUTI SEBELUMNYA ---
func (r *Router) handleTeacherLeaves(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		leaves, err := r.ptkLeaveService.ListLeaves(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, leaves)

	case http.MethodPost:
		var body repository.LeaveRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.ptkLeaveService.CreateLeave(req.Context(), body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Pengajuan cuti berhasil dibuat"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleTeacherLeaveByID(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/leaves/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	id := parts[0]

	switch req.Method {
	case http.MethodPut:
		if len(parts) > 1 && parts[1] == "status" {
			var body struct {
				Status string `json:"status"`
				Note   string `json:"note"`
			}
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
				writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
				return
			}
			if err := r.ptkLeaveService.UpdateStatus(req.Context(), id, body.Status, body.Note); err != nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{"message": "Status cuti diperbarui"})
		} else {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid path for PUT"})
		}

	case http.MethodDelete:
		if err := r.ptkLeaveService.DeleteLeave(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Pengajuan cuti dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

// -------------------------------------------------------------------------
// HANDLER SEBELUMNYA (Workload, Photo, History, Teachers, dll) - TETAP SAMA
// -------------------------------------------------------------------------
func (r *Router) handleTeacherWorkload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	workloads, err := r.teacherService.ListWorkloads(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, workloads)
}

func (r *Router) handleHistoryDocumentUpload(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"})
		return
	}
	file, handler, err := req.FormFile("document")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File 'document' tidak ditemukan"})
		return
	}
	defer file.Close()
	objectName := fmt.Sprintf("ptk_history_%s", handler.Filename)
	url, err := r.minioClient.UploadFile(req.Context(), "documents", objectName, file, handler)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (r *Router) handleTeachers(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		keyword := req.URL.Query().Get("keyword")
		teachers, total, err := r.teacherService.ListTeachers(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID, keyword)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"data": teachers, "total": total})
	case http.MethodPost:
		var body service.TeacherRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.teacherService.CreateTeacher(req.Context(), body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Guru berhasil ditambahkan"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleTeacherByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	pathParts := strings.Split(req.URL.Path, "/")
	id := pathParts[len(pathParts)-1]
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	switch req.Method {
	case http.MethodPut:
		var body service.TeacherRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.teacherService.UpdateTeacher(req.Context(), id, body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Guru berhasil diupdate"})
	case http.MethodDelete:
		if err := r.teacherService.DeleteTeacher(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Guru berhasil dihapus"})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleTeacherPhoto(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/photo/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	err := req.ParseMultipartForm(10 << 20)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Gagal parse form"})
		return
	}
	file, handler, err := req.FormFile("photo")
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "File 'photo' tidak ditemukan"})
		return
	}
	defer file.Close()
	objectName := fmt.Sprintf("ptk_%s_%s", id, handler.Filename)
	url, err := r.minioClient.UploadFile(req.Context(), "teachers", objectName, file, handler)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := r.teacherService.UpdatePhoto(req.Context(), id, url); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Foto berhasil diupload", "url": url})
}

func (r *Router) handleTeacherHistory(w http.ResponseWriter, req *http.Request) {
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/history/")
	parts := splitPath(path)
	if len(parts) < 2 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL format"})
		return
	}
	ptkID := parts[0]
	historyType := parts[1]
	var itemID string
	if len(parts) > 2 { itemID = parts[2] }

	switch historyType {
	case "educations":
		if req.Method == http.MethodGet {
			list, err := r.teacherHistoryService.ListEducations(req.Context(), ptkID)
			if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, list)
		} else if req.Method == http.MethodPost {
			var body repository.EducationHistory
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
			if err := r.teacherHistoryService.AddEducation(req.Context(), ptkID, body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusCreated, map[string]string{"message": "Pendidikan ditambahkan"})
		} else if req.Method == http.MethodDelete && itemID != "" {
			if err := r.teacherHistoryService.DeleteEducation(req.Context(), itemID); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, map[string]string{"message": "Pendidikan dihapus"})
		}
	case "certifications":
		if req.Method == http.MethodGet {
			list, err := r.teacherHistoryService.ListCertifications(req.Context(), ptkID)
			if err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, list)
		} else if req.Method == http.MethodPost {
			var body repository.CertificationHistory
			if err := json.NewDecoder(req.Body).Decode(&body); err != nil { writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"}); return }
			if err := r.teacherHistoryService.AddCertification(req.Context(), ptkID, body); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusCreated, map[string]string{"message": "Sertifikasi ditambahkan"})
		} else if req.Method == http.MethodDelete && itemID != "" {
			if err := r.teacherHistoryService.DeleteCertification(req.Context(), itemID); err != nil { writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()}); return }
			writeJSON(w, http.StatusOK, map[string]string{"message": "Sertifikasi dihapus"})
		}
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown history type"})
	}
}
