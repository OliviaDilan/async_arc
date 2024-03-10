package amqp

import (
	"context"

	"github.com/OliviaDilan/async_arc/pkg/contract"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher interface {
	Publish(ctx context.Context, message contract.Contract) error
	Close() error
}

type channelPublisher struct {
	ch *amqp.Channel
	queue amqp.Queue
}

func (p *channelPublisher) Close() error {
	return p.ch.Close()
}

func (p *channelPublisher) Publish(ctx context.Context, message contract.Contract) error {
	body, err := message.Marshal()
	if err != nil {
		return err
	}

	err = p.ch.PublishWithContext(ctx,
		"",     // exchange
		p.queue.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType: "text/plain",
			Body:        []byte(string(body)),
		})
	if err != nil {
		return err
	}

	return nil
}