package gormrepo

import (
	"errors"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	DB *gorm.DB
}

// NewGormSubscriptionRepository creates a new GormSubscriptionRepository.
func NewGormSubscriptionRepository(conn *GormDB) models.SubscriptionRepository {
	return &SubscriptionRepository{DB: conn.DB}
}

// Create creates a new Subscription record.
func (s *SubscriptionRepository) Create(userID, baseID, targetID uint) error {
	subscription := models.Subscription{
		UserID:           userID,
		BaseCurrencyID:   baseID,
		TargetCurrencyID: targetID,
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

// GetSubscriptions returns all the subscriptions.
func (s *SubscriptionRepository) GetSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := s.DB.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
