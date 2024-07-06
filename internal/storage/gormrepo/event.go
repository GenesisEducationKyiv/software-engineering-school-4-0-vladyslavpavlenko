package gormrepo

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"
)

// AddEvent creates a new outbox.Event record.
func (c *Connection) AddEvent(event *outbox.Event) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := c.db.WithContext(ctx).Create(event)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// FetchUnpublishedEvents retrieves all events from the database that have not been
// published after the specified offset.
func (c *Connection) FetchUnpublishedEvents(lastOffset uint) ([]outbox.Event, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	var events []outbox.Event
	err := c.db.WithContext(ctx).Where("id > ?", lastOffset).Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}
