package gormsubscriber

import (
	"context"

	"github.com/pkg/errors"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"
	"gorm.io/gorm"
)

var (
	ErrorDuplicateSubscription   = errors.New("subscription already exists")
	ErrorNonExistentSubscription = errors.New("subscription does not exist")
	ErrorSubscriptionFailed      = errors.New("subscription failed")
)

type Subscriber struct {
	db *gorm.DB
}

// NewSubscriber creates a new Subscriber.
func NewSubscriber(db *gorm.DB) *Subscriber {
	return &Subscriber{
		db: db,
	}
}

// AddSubscription creates a new models.Subscription record.
func (s *Subscriber) AddSubscription(email string) error {
	orchestrator, err := NewSagaOrchestrator(email, s.db)
	if err != nil {
		return errors.Wrap(err, "failed to create orchestrator")
	}

	err = orchestrator.Run(s)
	if err != nil {
		return err
	}
	if orchestrator.State.Status == StatusCompleted {
		return nil
	}
	return err
}

// DeleteSubscription deletes a models.Subscription record.
func (s *Subscriber) DeleteSubscription(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).Where("email = ?", email).Delete(&models.Subscription{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorNonExistentSubscription
	}

	return nil
}

// GetSubscriptions returns a paginated list of subscriptions. Limit specifies the number of records to be retrieved
// Limit conditions can be canceled by using `Limit(-1)`. Offset specify the number of records to skip before starting
// to return the records. Offset conditions can be canceled by using `Offset(-1)`.
func (s *Subscriber) GetSubscriptions(limit, offset int) ([]models.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	var subscriptions []models.Subscription
	result := s.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
