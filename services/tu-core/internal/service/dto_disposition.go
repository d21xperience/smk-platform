package service

import (
	"time"
)

type CreateDispositionRequest struct {
	CorrespondenceID string `json:"correspondence_id"`
	ToUserID         string `json:"to_user_id"`
	Note             string `json:"note"`
}

type UpdateDispositionStatusRequest struct {
	Status string `json:"status"` // READ, DONE
}

type DispositionResponse struct {
	ID               string    `json:"id"`
	CorrespondenceID string    `json:"correspondence_id"`
	FromUserID       string    `json:"from_user_id"`
	ToUserID         string    `json:"to_user_id"`
	Note             string    `json:"note"`
	Status           string    `json:"status"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
