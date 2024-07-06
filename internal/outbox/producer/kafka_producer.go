package producer

import (
	"context"
	"log"
	"strconv"
	"time"

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
}

// NewKafkaProducer initializes a new KafkaProducer.
func NewKafkaProducer(kafkaURL string, o Outbox, db dbConnection) (*KafkaProducer, error) {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	err := db.Migrate(&Offset{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate offset")
	}

	return &KafkaProducer{Writer: w, Outbox: o, db: db}, nil
}

// NewTopic creates a new kafka.TopicConfig if it does not exist.
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
			log.Println("Shutting down worker...")
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
		log.Printf("Failed to start transaction: %v", err)
		return
	}
	log.Println("Transaction started successfully")

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v, transaction rolled back", r)
		}
	}()

	// Retrieve the last processed offset
	lastOffset, err := p.db.GetLastOffset(topic, partition)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			log.Printf("Failed to fetch last offset: %v", err)
			return
		}

		lastOffset = Offset{
			Topic:     topic,
			Partition: partition,
			Offset:    0,
		}
	} else {
		log.Printf("Last offset fetched: %d", lastOffset.Offset)
	}

	// Fetch unpublished events based on the last offset
	events, err := p.db.FetchUnpublishedEvents(lastOffset.Offset)
	if err != nil {
		tx.Rollback()
		log.Printf("Failed to fetch unpublished events: %v", err)
		return
	}
	log.Printf("Fetched %d unpublished events", len(events))

	// Process each event
	for _, event := range events {
		key := []byte(strconv.Itoa(int(event.ID)))
		msg := &kafka.Message{
			Key:       key,
			Value:     []byte(event.Data),
			Partition: partition,
		}

		log.Printf("Preparing to send message with ID: %d", event.ID)

		if err = p.Writer.WriteMessages(ctx, *msg); err != nil {
			tx.Rollback()
			log.Printf("Failed to send Kafka message: %v", err)
			return
		}
		log.Printf("Message with ID %d sent successfully", event.ID)

		lastOffset.Offset = event.ID
		if err = p.db.UpdateOffset(&lastOffset); err != nil {
			tx.Rollback()
			log.Printf("Failed to update offset after sending message with ID %d: %v", event.ID, err)
			return
		}
		log.Printf("Offset updated successfully for message ID: %d", event.ID)
	}

	// Commit the transaction if all events are processed successfully
	if err = tx.Commit().Error; err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return
	}
	log.Println("Transaction committed and all events processed successfully")
}
