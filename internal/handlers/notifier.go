package handlers

import (
	"fmt"
	"log"
	"strconv"
	"sync"

	"gopkg.in/gomail.v2"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// Fetcher defines an interface for fetching rates.
type Fetcher interface {
	FetchRate(baseCode, targetCode string) (string, error)
}

// NotifySubscribers handles sending currency update emails to all the subscribers.
func (m *Repository) NotifySubscribers() error {
	subscriptions, err := m.App.Models.Subscription.GetSubscriptions()
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

	baseCode := subscription.BaseCurrency.Code
	targetCode := subscription.TargetCurrency.Code

	price, err := m.Fetcher.FetchRate(baseCode, targetCode)
	if err != nil {
		log.Printf("Failed to retrieve rate for %s to %s: %v", baseCode, targetCode, err)
		return
	}

	floatPrice, err := strconv.ParseFloat(price, 32)
	if err != nil {
		log.Printf("Failed to parse price: %v", err)
		return
	}

	params := email.Params{
		To:      subscription.User.Email,
		Subject: fmt.Sprintf("%s to %s Exchange Rate", baseCode, targetCode),
		Body:    fmt.Sprintf("The current exchange rate for %s to %s is %.2f.", baseCode, targetCode, floatPrice),
	}

	sender := &email.GomailSender{
		Dialer: gomail.NewDialer("smtp.gmail.com", 587, m.App.EmailConfig.Email, m.App.EmailConfig.Password),
	}

	email.SendEmail(sender, m.App.EmailConfig, params)
}
