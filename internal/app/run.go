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

	l.Info(fmt.Sprintf("running on port %d", webPort), zap.Int("port", webPort))
	srv := &http.Server{
		Addr:              fmt.Sprintf(":%d", webPort),
		Handler:           routes.Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdown(srv, cancel, l)

	return nil
}

// shutdown gracefully shuts down the application.
func shutdown(srv *http.Server, cancelFunc context.CancelFunc, l *logger.Logger) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		cancelFunc() // Cancel context to shut down dispatcher
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		l.Info("shutting down server...")
		if err := srv.Shutdown(ctx); err != nil {
			l.Error("failed to gracefully shutdown server", zap.Error(err))
		}
		l.Info("server has been shutdown")
	}()

	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		l.Fatal("failed to shutdown server", zap.Error(err))
	}
}

// scheduleEmails sets up a mailing process.
func scheduleEmails(s scheduler, n *notifier.Notifier, l *logger.Logger) error {
	_, err := s.Schedule(schedule, func() {
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
