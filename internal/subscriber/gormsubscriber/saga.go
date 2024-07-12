package gormsubscriber

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"
	"gorm.io/gorm"
)

const (
	StatusCompleted  = "completed"
	StatusInProgress = "in_progress"
	StatusFailed     = "failed"
)

// State represents the current state of the SAGA transaction.
type State struct {
	ID             string `gorm:"primary_key"`
	CurrentStep    int
	Email          string
	IsCompensating bool
	Status         string // StatusCompleted, StatusInProgress, StatusFailed
}

// Step represents a single step in the SAGA transaction.
type Step struct {
	Action       func(saga *State, s *Subscriber) error
	Compensation func(saga *State, s *Subscriber) error
}

// Orchestrator manages the execution of SAGA steps.
type Orchestrator struct {
	Steps []Step
	State State
	db    *gorm.DB
}

// NewSagaOrchestrator creates a new SAGA Orchestrator.
func NewSagaOrchestrator(email string, db *gorm.DB) (*Orchestrator, error) {
	err := db.AutoMigrate(&State{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate events")
	}

	return &Orchestrator{
		Steps: []Step{
			{
				Action:       validateSubscription,
				Compensation: nil,
			},
			{
				Action:       addSubscription,
				Compensation: deleteSubscription,
			},
		},
		State: State{
			ID:             uuid.New().String(),
			CurrentStep:    0,
			Email:          email,
			IsCompensating: false,
			Status:         StatusInProgress,
		},
		db: db,
	}, nil
}

// Run runs the SAGA Orchestrator.
func (o *Orchestrator) Run(s *Subscriber) error {
	for o.State.CurrentStep < len(o.Steps) {
		step := o.Steps[o.State.CurrentStep]
		var err error

		if o.State.IsCompensating {
			if step.Compensation != nil {
				err = step.Compensation(&o.State, s)
			}
		} else {
			err = step.Action(&o.State, s)
		}

		if err != nil {
			o.State.IsCompensating = true
			err = o.saveState()
			if err != nil {
				return err
			}
			continue
		}

		if o.State.IsCompensating {
			o.State.CurrentStep--
			if o.State.CurrentStep < 0 {
				o.State.Status = StatusFailed
				err = o.saveState()
				if err != nil {
					return err
				}
				return err
			}
		} else {
			o.State.CurrentStep++
		}
		err = o.saveState()
		if err != nil {
			return err
		}
	}

	if !o.State.IsCompensating {
		o.State.Status = StatusCompleted
		err := o.saveState()
		if err != nil {
			return err
		}
	}
	return nil
}

// saveState saves the SAGA transaction to the database.
func (o *Orchestrator) saveState() error {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	return o.db.WithContext(ctx).Save(&o.State).Error
}
