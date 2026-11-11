package rest

import (
	"net/http"
	"strings"
)

func (r *Router) registerCorrespondenceRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/correspondences", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			r.correspondenceHandler.Create(w, req)
		case http.MethodGet:
			r.correspondenceHandler.List(w, req)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
		}
	}))

	mux.HandleFunc("/api/v1/correspondences/", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasSuffix(req.URL.Path, "/archive") {
			if req.Method == http.MethodPut {
				r.correspondenceHandler.Archive(w, req)
				return
			}
		} else if strings.HasSuffix(req.URL.Path, "/pdf") {
			if req.Method == http.MethodGet {
				r.correspondenceHandler.DownloadPDF(w, req)
				return
			}
		} else if strings.HasSuffix(req.URL.Path, "/dispositions") {
			id := strings.Split(strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/"), "/")[0]
			switch req.Method {
			case http.MethodPost:
				r.dispositionHandler.Create(w, req, id)
			case http.MethodGet:
				r.dispositionHandler.ListByCorrespondence(w, req, id)
			default:
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			}
			return
		} else if strings.HasSuffix(req.URL.Path, "/sign") {
			if req.Method == http.MethodPost {
				id := strings.Split(strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/"), "/")[0]
				r.signatureHandler.Sign(w, req, id)
				return
			}
		} else if strings.HasSuffix(req.URL.Path, "/expeditions") {
			id := strings.Split(strings.TrimPrefix(req.URL.Path, "/api/v1/correspondences/"), "/")[0]
			switch req.Method {
			case http.MethodPost:
				r.expeditionHandler.Create(w, req, id)
			case http.MethodGet:
				r.expeditionHandler.ListByCorrespondence(w, req, id)
			default:
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			}
			return
		} else {
			switch req.Method {
			case http.MethodGet:
				r.correspondenceHandler.GetByID(w, req)
			case http.MethodPut:
				r.correspondenceHandler.Update(w, req)
			case http.MethodDelete:
				r.correspondenceHandler.Delete(w, req)
			default:
				writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
			}
		}
	}))

	mux.HandleFunc("/api/v1/dispositions/inbox", r.wrapWithContext(r.dispositionHandler.GetMyInbox))
	mux.HandleFunc("/api/v1/dispositions/", r.wrapWithContext(func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodPut && strings.HasSuffix(req.URL.Path, "/status") {
			path := strings.TrimPrefix(req.URL.Path, "/api/v1/dispositions/")
			id := strings.Split(path, "/")[0]
			r.dispositionHandler.UpdateStatus(w, req, id)
			return
		}
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "Method not allowed"})
	}))

	mux.HandleFunc("/api/v1/signatures/verify/", r.signatureHandler.Verify)
}
