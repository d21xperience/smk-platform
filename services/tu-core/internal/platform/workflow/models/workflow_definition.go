// Package models menyediakan domain models untuk Workflow Engine (Platform Runtime).
package models

import (
	"errors"
	"time"
)

// WorkflowDefinition adalah template workflow.
type WorkflowDefinition struct {
	definitionID string
	name         string
	description  string
	steps        []WorkflowStepDefinition
	timeout      time.Duration
	isActive     bool
	createdAt    time.Time
}

// WorkflowStepDefinition adalah definisi step dalam workflow.
type WorkflowStepDefinition struct {
	StepID     string
	Name       string
	Role       string // Role yang bisa complete step ini
	Order      int
	IsRequired bool
	CanReject  bool
}

// WorkflowDefinitionData adalah data untuk membuat WorkflowDefinition.
type WorkflowDefinitionData struct {
	DefinitionID string
	Name         string
	Description  string
	Steps        []WorkflowStepDefinition
	Timeout      time.Duration
}

// NewWorkflowDefinition membuat WorkflowDefinition baru.
func NewWorkflowDefinition(data WorkflowDefinitionData) (*WorkflowDefinition, error) {
	if data.Name == "" {
		return nil, errors.New("nama workflow wajib diisi")
	}
	if len(data.Steps) == 0 {
		return nil, errors.New("minimal satu step wajib ada")
	}

	// Validate step order
	for i, step := range data.Steps {
		if step.Order != i+1 {
			return nil, errors.New("urutan step tidak valid")
		}
	}

	if data.Timeout == 0 {
		data.Timeout = 7 * 24 * time.Hour // Default 7 hari
	}

	return &WorkflowDefinition{
		definitionID: data.DefinitionID,
		name:         data.Name,
		description:  data.Description,
		steps:        data.Steps,
		timeout:      data.Timeout,
		isActive:     true,
		createdAt:    time.Now().UTC(),
	}, nil
}

// Getters

func (d *WorkflowDefinition) DefinitionID() string            { return d.definitionID }
func (d *WorkflowDefinition) Name() string                    { return d.name }
func (d *WorkflowDefinition) Description() string             { return d.description }
func (d *WorkflowDefinition) Steps() []WorkflowStepDefinition { return d.steps }
func (d *WorkflowDefinition) Timeout() time.Duration          { return d.timeout }
func (d *WorkflowDefinition) IsActive() bool                  { return d.isActive }
func (d *WorkflowDefinition) CreatedAt() time.Time            { return d.createdAt }

// TotalSteps mengembalikan jumlah step.
func (d *WorkflowDefinition) TotalSteps() int {
	return len(d.steps)
}

// GetStepByOrder mengambil step berdasarkan urutan.
func (d *WorkflowDefinition) GetStepByOrder(order int) *WorkflowStepDefinition {
	for _, step := range d.steps {
		if step.Order == order {
			return &step
		}
	}
	return nil
}
