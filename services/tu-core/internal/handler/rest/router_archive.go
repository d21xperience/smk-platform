package rest

import (
	"net/http"
)

func (r *Router) registerArchiveRoutes(mux *http.ServeMux) {
	handler := NewArchiveHandler(r.archiveService)

	mux.HandleFunc("/api/v1/archives", r.wrapWithContext(handler.HandleCollection))
	mux.HandleFunc("/api/v1/archives/", r.wrapWithContext(handler.HandleByID))
}
