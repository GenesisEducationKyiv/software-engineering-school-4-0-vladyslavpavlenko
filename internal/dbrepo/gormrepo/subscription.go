package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

var ErrDuplicateSubscription = errors.New("subscription already exists")

// SubscriptionRepository is a models.Subscription repository.
type SubscriptionRepository struct {
	*GormDB
}

// NewSubscriptionRepository creates a new GormSubscriptionRepository.
func NewSubscriptionRepository(db *GormDB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

// AddSubscription creates a new Subscription record.
func (s *SubscriptionRepository) AddSubscription(email string) error {
	subscription := models.Subscription{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := s.DB.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return ErrDuplicateSubscription
		}
		return result.Error
	}
	return nil
}

// GetSubscriptions returns a paginated list of subscriptions. Limit specify the number of records to be retrieved
// Limit conditions can be canceled by using `Limit(-1)`. Offset specify the number of records to skip before starting
// to return the records. Offset conditions can be canceled by using `Offset(-1)`.
func (s *SubscriptionRepository) GetSubscriptions(limit, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := s.DB.Limit(limit).Offset(offset).Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
