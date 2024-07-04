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

	"github.com/segmentio/kafka-go"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/gormoutbox"

	"github.com/robfig/cron/v3"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"
	schedulerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/scheduler"
)

const (
	webPort  = 8080
	schedule = "0 10 * * *"
)

type (
	scheduler interface {
		ScheduleTask(schedule string, task func()) (cron.EntryID, error)
		Start()
		Stop()
	}
)

// Run is the application running process.
func Run(appConfig *config.AppConfig) error {
	dbConn, err := setup(appConfig)
	if err != nil {
		return err
	}
	defer dbConn.Close()

	s := schedulerpkg.NewCronScheduler()
	if err = scheduleEmails(s); err != nil {
		return fmt.Errorf("failed to schedule emails: %w", err)
	}
	s.Start()
	defer s.Stop()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the event producer for Kafka
	outboxService, err := gormoutbox.NewOutbox(dbConn.DB)
	if err != nil {
		return fmt.Errorf("failed to create outbox: %w", err)
	}

	appConfig.Outbox = outboxService

	kafkaURL := os.Getenv("KAFKA_URL")
	kafkaTopic := "events-topic"
	kafkaGroupID := "email-sender"

	kafkaWriter := outbox.NewKafkaWriter(kafkaURL, kafkaTopic)
	defer kafkaWriter.Close()
	go eventProducer(ctx, outboxService, kafkaWriter)

	// Start the event consumer for Kafka
	kafkaReader := outbox.NewKafkaReader(kafkaURL, kafkaTopic, kafkaGroupID)
	defer kafkaReader.Close()
	go eventConsumer(ctx, kafkaReader)

	log.Printf("Running on port %d", webPort)
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Handle graceful shutdown
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
func scheduleEmails(s scheduler) error {
	_, err := s.ScheduleTask(schedule, func() {
		err := handlers.Repo.ProduceMailingEvents()
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
func eventProducer(ctx context.Context, o outbox.Outbox, w *kafka.Writer) {
	outbox.Worker(ctx, o, w)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down event producer...")
}

// eventProducer runs an event dispatcher.
func eventConsumer(ctx context.Context, r *kafka.Reader) {
	go outbox.ConsumeMessages(ctx, r)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	log.Println("Shutting down event consumer...")
}
