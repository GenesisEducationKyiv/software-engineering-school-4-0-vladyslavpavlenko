package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// SubscriptionRepository is a models.Subscription repository.
type SubscriptionRepository struct {
	*GormDB
}

// NewSubscriptionRepository creates a new GormSubscriptionRepository.
func NewSubscriptionRepository(db *GormDB) *SubscriptionRepository {
	return &SubscriptionRepository{db}
}

// Create creates a new Subscription record.
func (s *SubscriptionRepository) AddSubscription(email string) error {
	subscription := models.Subscription{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := s.DB.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return gorm.ErrDuplicatedKey
		}
		return result.Error
	}
	return nil
}

// GetAll returns all the subscriptions.
func (s *SubscriptionRepository) GetSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := s.DB.Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
