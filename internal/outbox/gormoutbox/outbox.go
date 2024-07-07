package gormoutbox

import (
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/pkg/errors"
)

// dbConnection defines an interface for the database connection.
type dbConnection interface {
	Migrate(models ...any) error
	AddEvent(event *outbox.Event) error
}

// Outbox defines an interface for the transactional outbox.
type Outbox struct {
	db dbConnection
}

// NewOutbox creates `events` table to implement a transactional outbox.
// `events` table stores all the events ever occurred.
func NewOutbox(db dbConnection) (*Outbox, error) {
	err := db.Migrate(&outbox.Event{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate events")
	}
	return &Outbox{db: db}, nil
}

// AddEvent creates a new Event record.
func (o *Outbox) AddEvent(data outbox.Data) error {
	event := &outbox.Event{
		CreatedAt: time.Now(),
	}

	event.Data = data.Serialize()

	return o.db.AddEvent(event)
}
