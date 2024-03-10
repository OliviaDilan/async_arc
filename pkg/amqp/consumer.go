package amqp

import (
	"context"

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

	for {
		select {
		case d := <-msgs:
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