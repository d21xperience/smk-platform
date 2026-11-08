package models

import (
	"time"
)

type Disposition struct {
	ID               string
	CorrespondenceID string
	SchoolID         string
	FromUserID       string
	ToUserID         string
	Note             string
	Status           string // PENDING, READ, DONE
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
