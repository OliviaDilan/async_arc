package amqp

import (
	"context"
	"log"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer interface {
	Consume(ctx context.Context, handler ConsumeHandlerFunc) error
	Close() error
}

type ConsumeHandlerFunc func(ctx context.Context, body []byte) error

type queueConsumer struct {
	ch *amqp.Channel
	queue amqp.Queue
}

func (c *queueConsumer) Consume(ctx context.Context, handler ConsumeHandlerFunc) error {
	msgs, err := c.ch.Consume(
		c.queue.Name, // queue
		uuid.New().String(),     // consumer
		true,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		select {
		case d, ok := <-msgs:
			if !ok {
				log.Println("channel closed")
				return nil
			}

			if err := handler(ctx, d.Body); err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (c *queueConsumer) Close() error {
	return c.ch.Close()
}