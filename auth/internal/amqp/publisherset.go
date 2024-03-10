package amqp

import (
	"context"

	"github.com/OliviaDilan/async_arc/auth/internal/user"
	"github.com/OliviaDilan/async_arc/pkg/amqp"
	"github.com/OliviaDilan/async_arc/pkg/contract/auth"
)

type PublisherSet struct {
	onRegisterPublisher amqp.Publisher
}

func NewPublisherSet(amqpClient amqp.Client) (*PublisherSet, error) {
	p, err := amqpClient.NewPublisher("auth.user_created")
	if err != nil {
		return nil, err
	}

	return &PublisherSet{
		onRegisterPublisher: p,
	}, nil
}

func (s *PublisherSet) UserCreatedV1(ctx context.Context, user user.User) error {
	message := auth.NewUserCreatedV1(user.Username)
	message.Marshal()
	message.Metadata.Marshal()

	return s.onRegisterPublisher.Publish(ctx, message)
}