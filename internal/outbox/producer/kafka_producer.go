package producer

import (
	"context"
	"strconv"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"

	"github.com/pkg/errors"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"gorm.io/gorm"

	"github.com/segmentio/kafka-go"
)

type Outbox interface {
	AddEvent(data outbox.Data) error
}

type dbConnection interface {
	Migrate(models ...any) error
	BeginTransaction() (*gorm.DB, error)
	GetLastOffset(topic string, partition int) (Offset, error)
	FetchUnpublishedEvents(lastOffset uint) ([]outbox.Event, error)
	UpdateOffset(offset *Offset) error
}

type KafkaProducer struct {
	db     dbConnection
	Writer *kafka.Writer
	Outbox Outbox
	l      *logger.Logger
}

// NewKafkaProducer initializes a new KafkaProducer.
func NewKafkaProducer(kafkaURL string, o Outbox, db dbConnection, l *logger.Logger) (*KafkaProducer, error) {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	err := db.Migrate(&Offset{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate offset")
	}

	return &KafkaProducer{Writer: w, Outbox: o, db: db, l: l}, nil
}

// NewTopic creates a new kafka.TopicConfig if it doesn't exist.
func (p *KafkaProducer) NewTopic(topic string, partitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", p.Writer.Addr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}
	return nil
}

// SetTopic changes the topic of the current kafka.Writer
func (p *KafkaProducer) SetTopic(topic string) {
	p.Writer.Topic = topic
}

// Produce fetches for unpublished events, publishes them, and marks them as published.
func (p *KafkaProducer) Produce(ctx context.Context, frequency time.Duration, topic string, partition int) {
	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.l.Info("shutting down worker...")
			return
		case <-ticker.C:
			p.processEvents(ctx, topic, partition)
		}
	}
}

func (p *KafkaProducer) processEvents(ctx context.Context, topic string, partition int) {
	// Start a transaction
	tx, err := p.db.BeginTransaction()
	if err != nil {
		p.l.Error("failed to start transaction", zap.Error(err))
		return
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			p.l.Error("recovered from panic; transaction rolled back", zap.Any("recover", r))
		}
	}()

	// Retrieve the last processed offset
	lastOffset, err := p.db.GetLastOffset(topic, partition)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			p.l.Debug("failed to fetch last offset", zap.Error(err))
			return
		}

		lastOffset = Offset{
			Topic:     topic,
			Partition: partition,
			Offset:    0,
		}
	} else {
		p.l.Debug("last offset fetched", zap.Uint("offset", lastOffset.Offset))
	}

	// Fetch unpublished events based on the last offset
	events, err := p.db.FetchUnpublishedEvents(lastOffset.Offset)
	if err != nil {
		tx.Rollback()
		p.l.Error("failed to fetch unpublished events", zap.Error(err))
		return
	}

	if len(events) == 0 {
		tx.Rollback()
		return
	}

	// Process each event
	for _, event := range events {
		key := []byte(strconv.Itoa(int(event.ID)))
		msg := &kafka.Message{
			Key:       key,
			Value:     []byte(event.Data),
			Partition: partition,
		}

		p.l.Info("sending message to kafka", zap.Int("message_id", int(event.ID)))

		if err = p.Writer.WriteMessages(ctx, *msg); err != nil {
			tx.Rollback()
			p.l.Error("failed to send message", zap.Int("message_id", int(event.ID)), zap.Error(err))
			return
		}
		p.l.Debug("message sent", zap.Int("message_id", int(event.ID)))

		lastOffset.Offset = event.ID
		if err = p.db.UpdateOffset(&lastOffset); err != nil {
			tx.Rollback()
			p.l.Error("failed to update offset", zap.Error(err), zap.Int("message_id", int(event.ID)))
			return
		}
		p.l.Debug("offset updated", zap.Int("message_id", int(event.ID)))
	}

	// Commit the transaction if all events are processed successfully
	if err = tx.Commit().Error; err != nil {
		p.l.Warn("transaction failed", zap.Error(err))
		return
	}
	p.l.Debug("all kafka events processed", zap.String("topic", topic), zap.Int("partition", partition))
}
