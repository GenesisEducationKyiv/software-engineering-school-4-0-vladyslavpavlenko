package gormstorage

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/producer"
)

// GetLastOffset retrieves the last offset for a given topic from the database.
func (c *Connection) GetLastOffset(topic string, partition int) (producer.Offset, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	var lastOffset producer.Offset
	err := c.db.WithContext(ctx).Where("topic = ? AND partition = ?", topic, partition).First(&lastOffset).Error
	if err != nil {
		return producer.Offset{}, err
	}
	return lastOffset, nil
}

// UpdateOffset updates the offset in the database to reflect the latest published
// event's ID.
func (c *Connection) UpdateOffset(offset *producer.Offset) error {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	return c.db.WithContext(ctx).Save(offset).Error
}
