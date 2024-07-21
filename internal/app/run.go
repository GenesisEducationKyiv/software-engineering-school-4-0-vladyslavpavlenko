package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/notifier"
	schedulerpkg "github.com/vladyslavpavlenko/genesis-api-project/pkg/scheduler"

	producerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/outbox/producer"

	consumerpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/email/consumer"

	"github.com/robfig/cron/v3"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/routes"
)

const (
	apiPort         = 8080
	metricsPort     = 8081
	mailingSchedule = "0 10 * * *" // every day at 10 AM
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
func Run(app *config.Config, l *logger.Logger) error {
	svcs, err := setup(app, l)
	if err != nil {
		return err
	}
	defer svcs.DBConn.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	s := schedulerpkg.NewCronScheduler()
	if err = scheduleEmails(s, svcs.Notifier, l); err != nil {
		return fmt.Errorf("failed to schedule emails: %w", err)
	}
	s.Start()
	defer s.Stop()

	kafkaURL := os.Getenv("KAFKA_URL")
	kafkaTopic := "emails-topic"

	kafkaProducer, err := producerpkg.NewKafkaProducer(kafkaURL, svcs.Outbox, svcs.DBConn, l)
	if err != nil {
		return fmt.Errorf("failed to create kafka producer: %w", err)
	}
	defer kafkaProducer.Writer.Close()

	err = kafkaProducer.NewTopic(kafkaTopic, 1, 1)
	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}
	kafkaProducer.SetTopic(kafkaTopic)
	go eventProducer(ctx, kafkaProducer, kafkaTopic, 1, l)

	kafkaGroupID := "emails-group"

	kafkaConsumer, err := consumerpkg.NewKafkaConsumer(
		kafkaURL,
		kafkaTopic,
		0,
		kafkaGroupID,
		svcs.Sender,
		svcs.DBConn,
		l)
	if err != nil {
		return fmt.Errorf("failed to create kafka consumer: %w", err)
	}
	defer kafkaConsumer.Reader.Close()
	go eventConsumer(ctx, kafkaConsumer, l)

	apiServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", apiPort),
		Handler:           routes.API(svcs.Handlers),
		ReadHeaderTimeout: 5 * time.Second,
	}

	metricsServer := &http.Server{
		Addr:              fmt.Sprintf(":%d", metricsPort),
		Handler:           routes.Metrics(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdown(apiServer, metricsServer, cancel, l)

	return nil
}

// shutdown gracefully shuts down the application.
func shutdown(apiServer, metricsServer *http.Server, cancelFunc context.CancelFunc, l *logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancelFunc() // Cancel context to shut down dispatcher
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		l.Info("shutting down servers...")

		// Shutdown the API server
		if err := apiServer.Shutdown(ctx); err != nil {
			l.Error("failed to gracefully shutdown API server", zap.Error(err))
		} else {
			l.Info("API server has been shutdown")
		}

		// Shutdown the Metrics server
		if err := metricsServer.Shutdown(ctx); err != nil {
			l.Error("failed to gracefully shutdown metrics server", zap.Error(err))
		} else {
			l.Info("metrics server has been shutdown")
		}
	}()

	l.Info(fmt.Sprintf("running API server on port %d", apiPort), zap.Int("port", apiPort))
	go func() {
		if err := apiServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			l.Fatal("failed to shutdown API server", zap.Error(err))
		}
	}()

	l.Info(fmt.Sprintf("running metrics server on port %d", metricsPort), zap.Int("port", metricsPort))
	if err := metricsServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		l.Fatal("failed to shutdown metrics server", zap.Error(err))
	}
}

// scheduleEmails sets up a mailing process.
func scheduleEmails(s scheduler, n *notifier.Notifier, l *logger.Logger) error {
	_, err := s.Schedule(mailingSchedule, func() {
		err := n.Start()
		if err != nil {
			l.Error("error notifying subscribers", zap.Error(err))
		}
	})
	if err != nil {
		return fmt.Errorf("failed to schedule mailing task: %v", err)
	}

	return nil
}

// eventProducer runs an event dispatcher.
func eventProducer(ctx context.Context, producer producer, topic string, partition int, l *logger.Logger) {
	producer.Produce(ctx, 10*time.Second, topic, partition)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	l.Info("shutting down event producer...")
}

// eventConsumer runs an event dispatcher.
func eventConsumer(ctx context.Context, c consumer, l *logger.Logger) {
	go c.Consume(ctx)

	// Wait for context cancellation to handle graceful shutdown
	<-ctx.Done()
	l.Info("shutting down event consumer...")
}
