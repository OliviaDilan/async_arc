package amqp

import (
	"context"
	"sync"

	"github.com/OliviaDilan/async_arc/pkg/amqp"
)

type ConsumerSet struct {
	client    amqp.Client
	consumers []consumerWrapper
	mu        sync.Mutex
}

type consumerWrapper struct {
	amqp.Consumer
	handler amqp.ConsumeHandlerFunc
}

func NewConsumerSet(amqpClient amqp.Client) (*ConsumerSet) {
	return &ConsumerSet{
		client: amqpClient,
	}
}

func (s *ConsumerSet) StartConsumers(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, consumer := range s.consumers {
		go consumer.Consume(ctx, consumer.handler)
	}
}

func (s *ConsumerSet) Subscribe(topic string, handler amqp.ConsumeHandlerFunc) error {
	
	consumer, err := s.client.NewConsumer(topic)
	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.consumers = append(s.consumers, consumerWrapper{
		Consumer: consumer,
		handler:  handler,
	})

	return nil
}