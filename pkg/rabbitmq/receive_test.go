package rabbitmq

import (
	"context"
	"testing"
)

func TestReceive(t *testing.T) {
	mq, err := NewMessageQueue("guest", "guest", "localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer mq.ch.Close()

	queue, err := mq.DeclareReplicateQueue()
	if err != nil {
		t.Fatal(err)
	}
	defer mq.ch.QueueDelete(queue.Name, true, false, true)

	body := "test"
	err = mq.Publish(context.Background(), queue, []byte(body))
	if err != nil {
		t.Fatal(err)
	}

	msgs, err := mq.Receive(queue)
	if err != nil {
		t.Fatal(err)
	}

	msg := <-msgs
	if string(msg.Body) != body {
		t.Fatal("body is not test")
	}
}
