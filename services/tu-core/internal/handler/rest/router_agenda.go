package rest

import "net/http"

func (r *Router) registerAgendaRoutes(mux *http.ServeMux) {
	handler := NewAgendaHandler(r.agendaService)
	mux.HandleFunc("/api/v1/agenda", r.wrapWithContext(handler.List))
	mux.HandleFunc("/api/v1/agenda/export", r.wrapWithContext(handler.Export))
}
