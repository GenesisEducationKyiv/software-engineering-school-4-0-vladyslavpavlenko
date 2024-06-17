package handlers

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"gopkg.in/gomail.v2"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// Fetcher interface defines an interface for fetching rates.
type Fetcher interface {
	Fetch() (string, error)
}

// Subscription interface defines methods to access models.Subscription data.
type Subscription interface {
	Create(string) error
	GetAll() ([]models.Subscription, error)
}

// NotifySubscribers handles sending currency update emails to all the subscribers.
func (m *Repository) NotifySubscribers() error {
	subscriptions, err := m.Subscription.GetAll()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, subscription := range subscriptions {
		log.Println("Adding to WaitGroup")
		wg.Add(1)
		go m.sendEmail(&wg, subscription)
	}
	wg.Wait()

	return nil
}

// sendEmail is a controller function to prepare and send emails
func (m *Repository) sendEmail(wg *sync.WaitGroup, subscription models.Subscription) {
	defer wg.Done()

	price, err := m.Fetcher.Fetch()
	if err != nil {
		log.Printf("Failed to retrieve rate: %v", err)
		return
	}

	floatPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		log.Printf("Failed to parse price: %v", err)
		return
	}

	params := email.Params{
		To:      subscription.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", floatPrice),
	}

	sender := &email.GomailSender{
		Dialer: gomail.NewDialer("smtp.gmail.com", 587, m.App.EmailConfig.Email, m.App.EmailConfig.Password),
	}

	email.SendEmail(sender, m.App.EmailConfig, params)
}
