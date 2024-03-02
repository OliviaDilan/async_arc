package amqp

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	Publish(ctx context.Context, routingKey string, message any) error
	Subscribe(ctx context.Context, routingKey string, handler SubscribeHandlerFunc) error
	Close() error
}

type SubscribeHandlerFunc func(ctx context.Context, body []byte) error

type amqpClient struct {
	conn *amqp.Connection
}

func NewClient(connUrl string) (Client, error) {
	conn, err := amqp.Dial(connUrl)
	if err != nil {
		return nil, err
	}

	return &amqpClient{
		conn: conn,
	}, nil
}

func (c *amqpClient) Publish(ctx context.Context, routingKey string, message any) error {
	return nil
}

func (c *amqpClient) Subscribe(ctx context.Context, routingKey string, handler SubscribeHandlerFunc) error {
	return nil
}

func (c *amqpClient) Close() error {
	return c.conn.Close()
}