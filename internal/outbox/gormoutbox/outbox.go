package gormoutbox

import (
	"fmt"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Outbox is the repository for this package.
type Outbox struct {
	db *gorm.DB
}

// NewOutbox creates a new `events` table to serve as an outbox.
func NewOutbox(db *gorm.DB) (*Outbox, error) {
	err := db.AutoMigrate(&outbox.Event{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate events")
	}
	return &Outbox{db: db}, nil
}

// AddEvent creates a new Event record.
func (o *Outbox) AddEvent(data outbox.Data) error {
	event := &outbox.Event{
		Published: false,
		CreatedAt: time.Now(),
	}

	err := event.SerializeData(data)
	if err != nil {
		return fmt.Errorf("failed to serialize data: %w", err)
	}

	return o.db.Create(event).Error
}

// GetUnpublishedEvents retrieves all events that haven't been published yet.
func (o *Outbox) GetUnpublishedEvents() ([]outbox.Event, error) {
	var events []outbox.Event
	err := o.db.Where("published = ?", false).Find(&events).Error
	return events, err
}

// MarkEventAsPublished marks an event as published.
func (o *Outbox) MarkEventAsPublished(eventID uint) error {
	return o.db.Model(&outbox.Event{}).Where("id = ?", eventID).Update("published", true).Error
}

// Cleanup deletes all the Event records that have already been published or are outdated.
func (o *Outbox) Cleanup() {
	o.db.Where("published = ? AND created_at <= ?", true, time.Now().AddDate(0, 0, -1)).Delete(&outbox.Event{})
}
