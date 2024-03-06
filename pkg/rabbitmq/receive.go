package rabbitmq

import (
	"github.com/rabbitmq/amqp091-go"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Receive receive message from queue
func (mq *MessageQueue) Receive(queue amqp.Queue) (<-chan amqp091.Delivery, error) {
	msgs, err := mq.ch.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}
