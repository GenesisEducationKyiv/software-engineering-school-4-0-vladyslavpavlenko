package gormrepo

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email/consumer"
)

// AddConsumedEvent creates a new consumer.ConsumedEvent record.
func (c *Connection) AddConsumedEvent(event consumer.ConsumedEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := c.db.WithContext(ctx).Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
