package service

import (
	"github.com/google/uuid"
)

// generateID menghasilkan UUID v4 baru sebagai string
func generateID() string {
	return uuid.New().String()
}
