package commands

import "errors"

// CreateClassCommand untuk membuat rombel baru
type CreateClassCommand struct {
	ClassID    string
	PeriodID   string `json:"periodId"`
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
	Status     string
	CreatedBy  string `json:"createdBy"`
}

func (c CreateClassCommand) Validate() error {
	if c.ClassID == "" {
		return errors.New("class id is required")
	}
	if c.Name == "" {
		return errors.New("class name is required")
	}
	if c.Grade < 10 || c.Grade > 12 {
		return errors.New("grade must be between 10 and 12")
	}
	if c.Capacity <= 0 {
		return errors.New("capacity must be positive")
	}
	return nil
}

// UpdateClassCommand untuk update rombel
type UpdateClassCommand struct {
	ClassID    string `json:"classId"`
	Name       string `json:"name"`
	Grade      int    `json:"grade"`
	Level      string `json:"level"`
	Major      string `json:"major"`
	HomeroomID string `json:"homeroomId"`
	Capacity   int    `json:"capacity"`
	Status     string `json:"status"` // "active" | "inactive"
	UpdatedBy  string `json:"updatedBy"`
}

func (c UpdateClassCommand) Validate() error {
	if c.ClassID == "" {
		return errors.New("class id is required")
	}
	if c.Name == "" {
		return errors.New("class name is required")
	}
	if c.Grade < 10 || c.Grade > 12 {
		return errors.New("grade must be between 10 and 12")
	}
	if c.Capacity <= 0 {
		return errors.New("capacity must be positive")
	}
	if c.Status != "" && c.Status != "active" && c.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}
	return nil
}

// SetHomeroomCommand untuk mengatur wali kelas
type SetHomeroomCommand struct {
	ClassID    string `json:"classId"`
	HomeroomID string `json:"homeroomId"`
	UpdatedBy  string `json:"updatedBy"`
}
