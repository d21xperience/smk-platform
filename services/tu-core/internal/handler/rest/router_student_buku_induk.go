package rest

import (
	"fmt"
	"net/http"
	"strings"

	"sekolah-platform/services/tu-core/internal/middleware"
)

// registerStudentBukuIndukRoutes mendaftarkan route Buku Induk Siswa.
// Prefix /api/v1/kesiswaan/buku-induk-siswa dipakai (bukan /api/v1/students/buku-induk)
// untuk menghindari konflik dengan pattern existing /api/v1/students/{id}/history
// di router_student.go (Go 1.22 ServeMux panic saat startup jika ada ambiguity).
func (r *Router) registerStudentBukuIndukRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/kesiswaan/buku-induk-siswa", r.wrapWithContext(r.handleStudentBukuIndukList))
	mux.HandleFunc("/api/v1/kesiswaan/buku-induk-siswa/", r.wrapWithContext(r.handleStudentBukuIndukByID))
}

func (r *Router) handleStudentBukuIndukList(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	opCtx := middleware.GetOperationalContext(req.Context())
	list, err := r.studentBukuIndukService.ListBukuInduk(req.Context(), opCtx.SchoolID, opCtx.AcademicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"data": list})
}

func (r *Router) handleStudentBukuIndukByID(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		return
	}
	path := strings.TrimPrefix(req.URL.Path, "/api/v1/kesiswaan/buku-induk-siswa/")
	parts := splitPath(path)
	if len(parts) == 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "ID required"})
		return
	}

	opCtx := middleware.GetOperationalContext(req.Context())
	studentID := parts[0]
	schoolID := opCtx.SchoolID
	academicPeriodID := opCtx.AcademicPeriodID

	if len(parts) == 1 {
		r.handleStudentBukuIndukDetail(w, req, studentID, schoolID, academicPeriodID)
		return
	}
	if len(parts) == 2 && parts[1] == "pdf" {
		r.handleStudentBukuIndukPDF(w, req, studentID, schoolID, academicPeriodID)
		return
	}
	writeJSON(w, http.StatusNotFound, map[string]string{"error": "Endpoint not found"})
}

func (r *Router) handleStudentBukuIndukDetail(w http.ResponseWriter, req *http.Request, studentID, schoolID, academicPeriodID string) {
	detail, err := r.studentBukuIndukService.GetBukuIndukDetail(req.Context(), studentID, schoolID, academicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "Data tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, detail)
}

func (r *Router) handleStudentBukuIndukPDF(w http.ResponseWriter, req *http.Request, studentID, schoolID, academicPeriodID string) {
	buf, nis, err := r.studentBukuIndukService.GenerateBukuIndukPDF(req.Context(), studentID, schoolID, academicPeriodID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	filename := fmt.Sprintf("buku-induk-%s.pdf", nis)
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", buf.Len()))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
