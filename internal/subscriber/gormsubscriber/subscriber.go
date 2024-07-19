package gormsubscriber

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"
	"gorm.io/gorm"
)

var (
	ErrDuplicateSubscription   = errors.New("subscription already exists")
	ErrNonExistentSubscription = errors.New("subscription does not exist")
	ErrInternal                = errors.New("internal error")
)

type Subscriber struct {
	db     *gorm.DB
	logger *logger.Logger
}

// NewSubscriber creates a new Subscriber.
func NewSubscriber(db *gorm.DB, l *logger.Logger) (*Subscriber, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	if l == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &Subscriber{
		db:     db,
		logger: l,
	}, nil
}

// AddSubscription creates a new models.Subscription record.
func (s *Subscriber) AddSubscription(email string) error {
	orchestrator, err := NewSagaOrchestrator(email, s.db)
	if err != nil {
		s.logger.Error("failed to create orchestrator", zap.Error(err))
		return ErrInternal
	}

	err = orchestrator.Run(s)
	if err != nil {
		if errors.Is(err, ErrDuplicateSubscription) {
			s.logger.Info(ErrDuplicateSubscription.Error(),
				zap.String("email", email),
			)
			return ErrDuplicateSubscription
		}

		s.logger.Error("error adding subscription",
			zap.String("email", email),
			zap.Error(err),
		)

		return ErrInternal
	}

	s.logger.Info("new subscription", zap.String("email", email))

	return nil
}

// DeleteSubscription deletes a models.Subscription record.
func (s *Subscriber) DeleteSubscription(email string) error {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).Where("email = ?", email).Delete(&models.Subscription{})
	if result.Error != nil {
		s.logger.Error("failed to delete subscription",
			zap.String("email", email),
			zap.Error(result.Error))
		return ErrInternal
	}

	if result.RowsAffected == 0 {
		return ErrNonExistentSubscription
	}

	s.logger.Info("subscription deleted", zap.String("email", email))

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
