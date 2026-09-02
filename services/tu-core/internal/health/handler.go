package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

// HealthHandler menyediakan endpoint health check.
type HealthHandler struct {
	db *pgxpool.Pool
	nc *nats.Conn
	js nats.JetStreamContext
}

func NewHealthHandler(db *pgxpool.Pool, nc *nats.Conn, js nats.JetStreamContext) *HealthHandler {
	return &HealthHandler{db: db, nc: nc, js: js}
}

// ServeHTTP menangani request health check.
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	status := "healthy"
	checks := make(map[string]string)

	// Check Database
	if err := h.db.Ping(ctx); err != nil {
		status = "unhealthy"
		checks["database"] = "unhealthy: " + err.Error()
	} else {
		checks["database"] = "healthy"
	}

	// Check NATS
	if !h.nc.IsConnected() {
		status = "unhealthy"
		checks["nats"] = "unhealthy: not connected"
	} else {
		checks["nats"] = "healthy"
	}

	// Check JetStream
	if _, err := h.js.AccountInfo(); err != nil {
		status = "unhealthy"
		checks["jetstream"] = "unhealthy: " + err.Error()
	} else {
		checks["jetstream"] = "healthy"
	}

	response := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().UTC(),
		"checks":    checks,
	}

	w.Header().Set("Content-Type", "application/json")
	if status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	json.NewEncoder(w).Encode(response)
}
