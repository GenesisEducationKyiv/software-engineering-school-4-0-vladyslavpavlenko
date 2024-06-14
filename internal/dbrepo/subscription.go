package dbrepo

import (
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
