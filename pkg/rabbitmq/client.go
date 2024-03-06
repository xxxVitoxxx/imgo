package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageQueue struct {
	ch *amqp.Channel
}

// NewMessageQueue return a new message queue instance
func NewMessageQueue(user, password, address string) (*MessageQueue, error) {
	ch, err := connection(user, password, address)
	if err != nil {
		return nil, err
	}

	return &MessageQueue{ch}, nil
}

func connection(user, password, address string) (*amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s:5672/", user, password, address,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	return ch, nil
}
