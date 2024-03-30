package amqp

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Client interface {
	NewPublisher(queue string) (Publisher, error)
	NewConsumer(queue string) (Consumer, error)
	Close() error
}

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

func (c *amqpClient) NewPublisher(queue string) (Publisher, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	return &channelPublisher{
		ch: ch,
		queue: q,
	}, nil
}



func (c *amqpClient) NewConsumer(queue string) (Consumer, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}


	return &queueConsumer{
		ch: ch,
		queue: q,
	}, nil

}

func (c *amqpClient) Close() error {
	return c.conn.Close()
}