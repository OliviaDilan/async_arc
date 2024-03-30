package amqp

import (
	"context"

	"github.com/OliviaDilan/async_arc/task/internal/task"
	"github.com/OliviaDilan/async_arc/pkg/amqp"
	contract "github.com/OliviaDilan/async_arc/pkg/contract/task"
)

type PublisherSet struct {
	onTaskCreatedPublisher amqp.Publisher
	onTaskCompletedPublisher amqp.Publisher
	onTaskAssignedPublisher amqp.Publisher
}

func NewPublisherSet(amqpClient amqp.Client) (*PublisherSet, error) {
	taskCreatedPublisher, err := amqpClient.NewPublisher("task_created")
	if err != nil {
		return nil, err
	}

	taskAssignedPublisher, err := amqpClient.NewPublisher("task_assigned")
	if err != nil {
		return nil, err
	}

	taskCompletedPublisher, err := amqpClient.NewPublisher("task_completed")
	if err != nil {
		return nil, err
	}

	return &PublisherSet{
		onTaskCreatedPublisher: taskCreatedPublisher,
		onTaskCompletedPublisher: taskCompletedPublisher,
		onTaskAssignedPublisher: taskAssignedPublisher,
	}, nil
}

func (s *PublisherSet) TaskCreatedV1(ctx context.Context, task *task.Task) error {
	message := contract.NewTaskCreatedV1(task.ID)
	
	return s.onTaskCreatedPublisher.Publish(ctx, message)
}

func (s *PublisherSet) TaskAssignedV1(ctx context.Context, task *task.Task) error {
	message := contract.NewTaskAssignedV1(task.ID, task.Assignee)
	
	return s.onTaskAssignedPublisher.Publish(ctx, message)
}

func (s *PublisherSet) TaskCompletedV1(ctx context.Context, task *task.Task) error {
	message := contract.NewTaskCompletedV1(task.ID, task.Assignee)
	
	return s.onTaskCompletedPublisher.Publish(ctx, message)
}