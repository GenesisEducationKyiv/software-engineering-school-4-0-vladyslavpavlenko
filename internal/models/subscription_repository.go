package models

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type SubscriptionRepository struct {
	db *gorm.DB
}

// NewSubscriptionRepository creates a new SubscriptionRepository.
func NewSubscriptionRepository(db *gorm.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create creates a new Subscription record.
func (repo *SubscriptionRepository) Create(userID, baseID, targetID uint) (*Subscription, error) {
	subscription := Subscription{
		UserID:           userID,
		BaseCurrencyID:   baseID,
		TargetCurrencyID: targetID,
	}
	result := repo.db.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, result.Error
	}
	return &subscription, nil
}

// GetSubscriptions returns all the subscriptions.
func (repo *SubscriptionRepository) GetSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	result := repo.db.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
