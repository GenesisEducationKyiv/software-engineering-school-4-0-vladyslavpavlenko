package gormrepo

import (
	"errors"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormSubscriptionRepository struct {
	db *gorm.DB
}

// NewGormSubscriptionRepository creates a new GormSubscriptionRepository.
func NewGormSubscriptionRepository(db *gorm.DB) dbrepo.SubscriptionRepository {
	return &GormSubscriptionRepository{db: db}
}

// Create creates a new Subscription record.
func (repo *GormSubscriptionRepository) Create(userID, baseID, targetID uint) (*dbrepo.Subscription, error) {
	subscription := dbrepo.Subscription{
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
func (repo *GormSubscriptionRepository) GetSubscriptions() ([]dbrepo.Subscription, error) {
	var subscriptions []dbrepo.Subscription
	result := repo.db.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
