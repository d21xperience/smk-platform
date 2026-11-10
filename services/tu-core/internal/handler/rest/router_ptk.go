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

func (r *Router) registerPtkRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/teachers", r.wrapWithContext(r.handleTeachers))
	// RUTE YANG HILANG (Dikembalikan)
	mux.HandleFunc("/api/v1/teachers/dashboard", r.wrapWithContext(r.handleTeacherDashboard))
	mux.HandleFunc("/api/v1/teachers/buku-induk", r.wrapWithContext(r.handleBukuIndukList))
	mux.HandleFunc("/api/v1/teachers/buku-induk/", r.wrapWithContext(r.handleBukuIndukDetail))
	
	mux.HandleFunc("/api/v1/teachers/workload", r.wrapWithContext(r.handleTeacherWorkload))
	mux.HandleFunc("/api/v1/teachers/photo/", r.wrapWithContext(r.handleTeacherPhoto))
	mux.HandleFunc("/api/v1/teachers/history/upload", r.wrapWithContext(r.handleHistoryDocumentUpload))
	mux.HandleFunc("/api/v1/teachers/history/", r.wrapWithContext(r.handleTeacherHistory))

	mux.HandleFunc("/api/v1/teachers/leaves", r.wrapWithContext(r.handleTeacherLeaves))
	mux.HandleFunc("/api/v1/teachers/leaves/", r.wrapWithContext(r.handleTeacherLeaveByID))

	mux.HandleFunc("/api/v1/teachers/archives/upload", r.wrapWithContext(r.handleArchiveUpload))
	mux.HandleFunc("/api/v1/teachers/archives", r.wrapWithContext(r.handleTeacherArchives))
	mux.HandleFunc("/api/v1/teachers/archives/", r.wrapWithContext(r.handleTeacherArchiveByID))

	mux.HandleFunc("/api/v1/kepegawaian/sk-pembagian-tugas/export/", r.wrapWithContext(r.handleSKExportExcel))
	mux.HandleFunc("/api/v1/kepegawaian/sk-pembagian-tugas", r.wrapWithContext(r.handleSKPembagianTugas))
	mux.HandleFunc("/api/v1/kepegawaian/sk-pembagian-tugas/", r.wrapWithContext(r.handleSKPembagianTugasByID))

	mux.HandleFunc("/api/v1/kepegawaian/statistik", r.wrapWithContext(r.handleStatistikDashboard))

	mux.HandleFunc("/api/v1/teachers/", r.wrapWithContext(r.handleTeacherByID))
}

// --- Handler yang dikembalikan ---

func (r *Router) handleTeacherDashboard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	stats, err := r.ptkDashboardService.GetDashboardStats(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

func (r *Router) handleBukuIndukList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	list, err := r.bukuIndukService.ListBukuInduk(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, list)
}

func (r *Router) handleBukuIndukDetail(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/buku-induk/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	detail, err := r.bukuIndukService.GetBukuIndukDetail(req.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Data tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

// --- Handler Lama (Tidak Diubah) ---

func (r *Router) handleArchiveUpload(w http.ResponseWriter, req *http.Request) {
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

	objectName := fmt.Sprintf("ptk_archive_%d_%s", time.Now().Unix(), handler.Filename)
	url, err := r.minioClient.UploadFile(req.Context(), "documents", objectName, file, handler)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"url": url})
}

func (r *Router) handleTeacherArchives(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		list, err := r.ptkArchiveService.ListArchives(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, list)

	case http.MethodPost:
		var body repository.Archive
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		if err := r.ptkArchiveService.CreateArchive(req.Context(), body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Dokumen berhasil disimpan"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleTeacherArchiveByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/teachers/archives/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	switch req.Method {
	case http.MethodPut:
		var body repository.Archive
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		body.ID = id
		if err := r.ptkArchiveService.UpdateArchive(req.Context(), body); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Dokumen berhasil diupdate"})

	case http.MethodDelete:
		if err := r.ptkArchiveService.DeleteArchive(req.Context(), id); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"message": "Dokumen dihapus"})

	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

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
	objectName := fmt.Sprintf("ptk_history_%d_%s", time.Now().Unix(), handler.Filename)
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

func (r *Router) handleSKPembagianTugas(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	switch req.Method {
	case http.MethodGet:
		list, err := r.skPembagianTugasService.ListSK(req.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, list)
	case http.MethodPost:
		var body repository.SKCreateRequest
		if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid body"})
			return
		}
		skID, nomorSK, err := r.skPembagianTugasService.CreateDraft(req.Context(), body, opCtx.SchoolID, opCtx.UserID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusCreated, map[string]string{"message": "Draft SK berhasil dibuat", "sk_id": skID, "nomor_sk": nomorSK})
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}
}

func (r *Router) handleSKPembagianTugasByID(w http.ResponseWriter, req *http.Request) {
	opCtx := middleware.GetOperationalContext(req.Context())
	if opCtx == nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Missing operational context"})
		return
	}
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/sk-pembagian-tugas/")
	parts := splitPath(path)
	if len(parts) < 2 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid URL"})
		return
	}
	id := parts[0]
	action := parts[1]
	if req.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	var body struct {
		Reason string `json:"reason"`
		Notes  string `json:"notes"`
	}
	if req.Body != http.NoBody {
		_ = json.NewDecoder(req.Body).Decode(&body)
	}
	var err error
	switch action {
	case "submit":
		err = r.skPembagianTugasService.SubmitSK(req.Context(), id)
	case "approve":
		err = r.skPembagianTugasService.ApproveAndIssueSK(req.Context(), id, opCtx.SchoolID)
	case "reject":
		err = r.skPembagianTugasService.RejectSK(req.Context(), id, body.Reason)
	case "revision":
		err = r.skPembagianTugasService.RequestRevision(req.Context(), id, body.Notes)
	default:
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "Unknown action"})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Aksi " + action + " berhasil"})
}

func (r *Router) handleSKExportExcel(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	id := strings.TrimPrefix(req.URL.Path, "/api/v1/kepegawaian/sk-pembagian-tugas/export/")
	if id == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}
	url, err := r.skPembagianTugasService.ExportExcel(req.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "Excel berhasil dibuat", "url": url})
}

func (r *Router) handleStatistikDashboard(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	stats, err := r.statistikService.GetDashboardStats(req.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, stats)
}
