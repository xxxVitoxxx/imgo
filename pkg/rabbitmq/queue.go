package rabbitmq

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

// DeclareReplicateQueue declare a queue for replicate
func (mq *MessageQueue) DeclareReplicateQueue() (amqp.Queue, error) {
	replicateQueue, err := mq.ch.QueueDeclare(
		"replicate", // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return amqp.Queue{}, fmt.Errorf("failed to declare replicate queue%w", err)
	}

	return replicateQueue, nil
}
