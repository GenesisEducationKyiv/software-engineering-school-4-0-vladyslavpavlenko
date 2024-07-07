package notifier

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	outboxpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"
)

const batchSize = 100

// dbConnection defines an interface for the database connection.
type dbConnection interface {
	GetSubscriptions(limit, offset int) ([]models.Subscription, error)
}

// fetcher defines an interface for the fetching data rates.
type fetcher interface {
	Fetch(ctx context.Context, base, target string) (string, error)
}

// outbox defines an interface for writing events to the outbox.
type outbox interface {
	AddEvent(data outboxpkg.Data) error
}

type Notifier struct {
	DB      dbConnection
	Fetcher fetcher
	Outbox  outbox
}

// NewNotifier creates a new Notifier.
func NewNotifier(db dbConnection, f fetcher, o outbox) *Notifier {
	return &Notifier{
		DB:      db,
		Fetcher: f,
		Outbox:  o,
	}
}

// Start handles producing events for currency rate update emails.
func (n *Notifier) Start() error {
	rate, err := n.Fetcher.Fetch(context.Background(), "USD", "UAH")
	if err != nil {
		return fmt.Errorf("failed to retrieve rate: %w", err)
	}

	floatRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return fmt.Errorf("failed to parse rate: %w", err)
	}

	var offset int
	errChan := make(chan error, 1)
	for {
		subscriptions, err := n.DB.GetSubscriptions(batchSize, offset)
		if err != nil {
			return err
		}
		if len(subscriptions) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, sub := range subscriptions {
			wg.Add(1)
			go func(sub models.Subscription) {
				defer wg.Done()

				data := outboxpkg.Data{
					Email: sub.Email,
					Rate:  floatRate,
				}
				if localErr := n.Outbox.AddEvent(data); localErr != nil {
					select {
					case errChan <- localErr:
					default:
					}
				}
			}(sub)
		}
		wg.Wait()

		select {
		case err := <-errChan:
			return err
		default:
		}

		offset += batchSize
	}
	return nil
}
