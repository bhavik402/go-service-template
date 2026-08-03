// Package kafka holds the Kafka producer and consumer. A consumer is
// just another entry point into the application, exactly like an HTTP
// handler: deserialize, validate, call a service, acknowledge.
package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

// Producer publishes messages to Kafka. Services depend on this directly
// for the example resource; if you want services to depend on an
// interface instead (e.g. to swap Kafka for SQS later), define one in
// internal/service and have Producer satisfy it.
type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:                   kafka.TCP(brokers...),
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, topic, key string, value []byte) error {
	return p.writer.WriteMessages(ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(key),
		Value: value,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
