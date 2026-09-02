package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// StepStatus adalah status step.
type StepStatus string

const (
	StepStatusPending   StepStatus = "PENDING"
	StepStatusCompleted StepStatus = "COMPLETED"
	StepStatusRejected  StepStatus = "REJECTED"
	StepStatusSkipped   StepStatus = "SKIPPED"
)

// WorkflowStep adalah Entity dalam WorkflowInstance Aggregate.
type WorkflowStep struct {
	stepID      string
	instanceID  string
	name        string
	role        string
	order       int
	status      StepStatus
	action      string // "approve", "reject", dll
	comment     string
	completedBy string
	completedAt *time.Time
}

// NewWorkflowStep membuat WorkflowStep baru dari definition.
func NewWorkflowStep(def WorkflowStepDefinition, instanceID string) *WorkflowStep {
	return &WorkflowStep{
		stepID:     uuid.New().String(),
		instanceID: instanceID,
		name:       def.Name,
		role:       def.Role,
		order:      def.Order,
		status:     StepStatusPending,
	}
}

// Getters

func (s *WorkflowStep) StepID() string          { return s.stepID }
func (s *WorkflowStep) InstanceID() string      { return s.instanceID }
func (s *WorkflowStep) Name() string            { return s.name }
func (s *WorkflowStep) Role() string            { return s.role }
func (s *WorkflowStep) Order() int              { return s.order }
func (s *WorkflowStep) Status() StepStatus      { return s.status }
func (s *WorkflowStep) Action() string          { return s.action }
func (s *WorkflowStep) Comment() string         { return s.comment }
func (s *WorkflowStep) CompletedBy() string     { return s.completedBy }
func (s *WorkflowStep) CompletedAt() *time.Time { return s.completedAt }

// Complete menyelesaikan step.
func (s *WorkflowStep) Complete(action, comment, userID string) error {
	if s.status != StepStatusPending {
		return errors.New("step sudah selesai")
	}
	if action == "" {
		return errors.New("action wajib diisi")
	}
	if userID == "" {
		return errors.New("userId wajib diisi")
	}

	s.status = StepStatusCompleted
	s.action = action
	s.comment = comment
	s.completedBy = userID
	now := time.Now().UTC()
	s.completedAt = &now

	return nil
}

// Reject menolak step.
func (s *WorkflowStep) Reject(reason, userID string) error {
	if s.status != StepStatusPending {
		return errors.New("step sudah selesai")
	}
	if reason == "" {
		return errors.New("alasan penolakan wajib diisi")
	}

	s.status = StepStatusRejected
	s.action = "reject"
	s.comment = reason
	s.completedBy = userID
	now := time.Now().UTC()
	s.completedAt = &now

	return nil
}

// IsPending mengecek apakah step masih pending.
func (s *WorkflowStep) IsPending() bool { return s.status == StepStatusPending }

// IsCompleted mengecek apakah step sudah completed.
func (s *WorkflowStep) IsCompleted() bool { return s.status == StepStatusCompleted }
