package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

const batchSize = 100

type (
	// Sender defines an interface for sending emails.
	Sender interface {
		Send(emailConfig email.Config, params email.Params) error
	}

	// Fetcher interface defines an interface for fetching rates.
	Fetcher interface {
		Fetch(ctx context.Context, base, target string) (string, error)
	}

	// Subscriber interface defines methods to access models.Subscription data.
	Subscriber interface {
		AddSubscription(string) error
		GetSubscriptions(limit, offset int) ([]models.Subscription, error)
	}
)

// NotifySubscribers handles sending currency update emails to all the subscribers in batches.
func (m *Repository) NotifySubscribers() error {
	var offset int
	for {
		subscriptions, err := m.Services.Subscriber.GetSubscriptions(batchSize, offset)
		if err != nil {
			return err
		}
		if len(subscriptions) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, subscription := range subscriptions {
			wg.Add(1)
			go func(sub models.Subscription) {
				defer wg.Done()
				if err = m.sendEmail(sub); err != nil {
					log.Println(err)
				}
			}(subscription)
		}
		wg.Wait()

		offset += batchSize
	}

	return nil
}

// sendEmail is a controller function to prepare and send emails
func (m *Repository) sendEmail(subscription models.Subscription) error {
	price, err := m.Services.Fetcher.Fetch(context.Background(), "USD", "UAH")
	if err != nil {
		return fmt.Errorf("failed to retrieve rate: %w", err)
	}

	floatPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return fmt.Errorf("failed to parse price: %w", err)
	}

	params := email.Params{
		To:      subscription.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", floatPrice),
	}

	err = m.Services.Sender.Send(m.App.EmailConfig, params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
