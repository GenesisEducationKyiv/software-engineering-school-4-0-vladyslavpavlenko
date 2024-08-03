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

// New creates `events` table to implement a transactional outbox.
// `events` table stores all the events ever occurred.
func New(db dbConnection) (*Outbox, error) {
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

	sData, err := data.Serialize()
	if err != nil {
		return errors.Wrap(err, "failed to serialize data")
	}

	event.Data = sData

	return o.db.AddEvent(event)
}
