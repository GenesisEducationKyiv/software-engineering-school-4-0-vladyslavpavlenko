package gormrepo

import (
	"errors"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormSubscriptionRepository struct {
	db *gorm.DB
}

// NewGormSubscriptionRepository creates a new GormSubscriptionRepository.
func NewGormSubscriptionRepository(conn *GormDB) models.SubscriptionRepository {
	return &GormSubscriptionRepository{db: conn.DB}
}

// Create creates a new Subscription record.
func (repo *GormSubscriptionRepository) Create(userID, baseID, targetID uint) (*models.Subscription, error) {
	subscription := models.Subscription{
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
func (repo *GormSubscriptionRepository) GetSubscriptions() ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := repo.db.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
