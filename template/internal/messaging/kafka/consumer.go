package kafka

import (
	"context"
	"errors"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

// Handler processes one message. Returning an error skips the commit so
// the message is redelivered — swap in a retry/DLQ policy here as your
// reliability requirements grow.
type Handler func(ctx context.Context, msg kafka.Message) error

type Consumer struct {
	reader  *kafka.Reader
	handler Handler
	logger  *slog.Logger
}

func NewConsumer(brokers []string, groupID, topic string, handler Handler, logger *slog.Logger) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		GroupID: groupID,
		Topic:   topic,
	})
	return &Consumer{reader: reader, handler: handler, logger: logger}
}

// Run blocks, consuming messages until ctx is cancelled. Call it in its
// own goroutine from the composition root.
func (c *Consumer) Run(ctx context.Context) {
	for {
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return
			}
			c.logger.Error("kafka fetch failed", slog.Any("error", err), slog.String("topic", c.reader.Config().Topic))
			continue
		}

		if err := c.handler(ctx, msg); err != nil {
			c.logger.Error("kafka handler failed",
				slog.Any("error", err),
				slog.String("topic", msg.Topic),
				slog.Int("partition", msg.Partition),
				slog.Int64("offset", msg.Offset),
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			c.logger.Error("kafka commit failed", slog.Any("error", err))
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
