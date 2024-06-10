package models

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Create creates a new User record.
func (s *Subscription) Create(userID uint, baseID uint, targetID uint) (*Subscription, error) {
	subscription := Subscription{
		UserID:           userID,
		BaseCurrencyID:   baseID,
		TargetCurrencyID: targetID,
	}
	result := db.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, result.Error
	}
	return &subscription, nil
}
