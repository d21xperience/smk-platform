package rest

import (
	"net/http"
)

func (r *Router) registerLetterTypeRoutes(mux *http.ServeMux) {
	handler := NewLetterTypeHandler(r.letterTypeService)
	
	mux.HandleFunc("/api/v1/letter-types", r.wrapWithContext(handler.HandleCollection))
	mux.HandleFunc("/api/v1/letter-types/", r.wrapWithContext(handler.HandleByID))
}
