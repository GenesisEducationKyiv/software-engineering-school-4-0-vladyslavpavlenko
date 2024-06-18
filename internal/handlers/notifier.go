package handlers

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// Sender defines an interface for sending emails.
type Sender interface {
	Send(emailConfig email.Config, params email.Params) error
}

// Fetcher interface defines an interface for fetching rates.
type Fetcher interface {
	Fetch() (string, error)
}

// Subscriber interface defines methods to access models.Subscription data.
type Subscriber interface {
	AddSubscription(string) error
	GetSubscriptions() ([]models.Subscription, error)
}

// NotifySubscribers handles sending currency update emails to all the subscribers.
func (m *Repository) NotifySubscribers() error {
	subscriptions, err := m.Subscriber.GetSubscriptions()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, subscription := range subscriptions {
		log.Println("Adding to WaitGroup")
		wg.Add(1)
		go func() {
			err = m.sendEmail(&wg, subscription)
			if err != nil {
				log.Println(err)
			}
		}()
	}
	wg.Wait()

	return nil
}

// sendEmail is a controller function to prepare and send emails
func (m *Repository) sendEmail(wg *sync.WaitGroup, subscription models.Subscription) error {
	defer wg.Done()

	price, err := m.Fetcher.Fetch()
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

	err = m.Sender.Send(m.App.EmailConfig, params)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
