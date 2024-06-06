package models

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"time"
)

// Subscription is a GORM subscription model.
type Subscription struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	UserID           uint      `gorm:"not null;index" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"-"`
	BaseCurrencyID   uint      `gorm:"not null;index" json:"base_currency_id"`
	BaseCurrency     Currency  `gorm:"foreignKey:BaseCurrencyID" json:"-"`
	TargetCurrencyID uint      `gorm:"not null;index" json:"target_currency_id"`
	TargetCurrency   Currency  `gorm:"foreignKey:TargetCurrencyID" json:"-"`
	CreatedAt        time.Time `json:"created_at"`
}

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
