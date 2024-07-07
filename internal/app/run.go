package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
	schedulerpkg "github.com/vladyslavpavlenko/genesis-api-project/pkg/scheduler"

	producerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/producer"

	consumerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/email/consumer"

	"github.com/robfig/cron/v3"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"
)

const (
	webPort  = 8080
	schedule = "0 10 * * *"
)

// scheduler is an interface for task scheduling.
type scheduler interface {
	Schedule(schedule string, task func()) (cron.EntryID, error)
	Start()
	Stop()
}

// consumer is an interface for event consumption.
type consumer interface {
	Consume(ctx context.Context)
}

// producer is an interface for event producing.
type producer interface {
	NewTopic(topic string, partitions int, replicationFactor int) error
	SetTopic(topic string)
	Produce(ctx context.Context, frequency time.Duration, topic string, partition int)
}

// Run is the application running process.
func Run(appConfig *config.AppConfig) error {
	appServices, err := setup(appConfig)
	if err != nil {
		return err
	}
	defer appServices.DBConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := schedulerpkg.NewCronScheduler()
	if err = scheduleEmails(s, appServices.Notifier); err != nil {
		return fmt.Errorf("failed to schedule emails: %w", err)
	}
	s.Start()
	defer s.Stop()

	kafkaURL := os.Getenv("KAFKA_URL")
	kafkaTopic := "emails-topic"

	kafkaProducer, err := producerpkg.NewKafkaProducer(kafkaURL, appServices.Outbox, appServices.DBConn)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}
	defer kafkaProducer.Writer.Close()

	err = kafkaProducer.NewTopic(kafkaTopic, 1, 1)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	kafkaProducer.SetTopic(kafkaTopic)
	go eventProducer(ctx, kafkaProducer, kafkaTopic, 1)

	kafkaGroupID := "emails-group"

	kafkaConsumer, err := consumerpkg.NewKafkaConsumer(
		kafkaURL,
		kafkaTopic,
		0,
		kafkaGroupID,
		appServices.Sender,
		appServices.DBConn)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	defer kafkaConsumer.Reader.Close()
	go eventConsumer(ctx, kafkaConsumer)

	log.Printf("Running on port %d", webPort)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	handleShutdown(srv, cancel)

	return nil
}

// handleShutdown handles a graceful shutdown of the application.
func handleShutdown(srv *http.Server, cancelFunc context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancelFunc() // Cancel context to shut down dispatcher
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		log.Println("Shutting down server...")
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("HTTP server shutdown failed: %v", err)
		}
		log.Println("Server has been stopped")
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// scheduleEmails sets up a mailing process.
func scheduleEmails(s scheduler, n *notifier.Notifier) error {
	_, err := s.Schedule(schedule, func() {
		err := n.Start()
		if err != nil {
			log.Printf("Error notifying subscribers: %v", err)
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule mailing task: %v", err)
	}

	return nil
}

// eventProducer runs an event dispatcher.
func eventProducer(ctx context.Context, producer producer, topic string, partition int) {
	producer.Produce(ctx, 10*time.Second, topic, partition)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down event producer...")
}

// eventConsumer runs an event dispatcher.
func eventConsumer(ctx context.Context, c consumer) {
	go c.Consume(ctx)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down event consumer...")
}
