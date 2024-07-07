package consumer

import (
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"
)

// ConsumedEvent represents an event consumed by the consumer.
type ConsumedEvent struct {
	ID         uint         `gorm:"not null;index"`
	Event      outbox.Event `gorm:"foreignKey:ID" json:"-"`
	Data       string
	ConsumedAt time.Time
	UpdatedAt  time.Time
}
