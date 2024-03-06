package rabbitmq

import (
	"context"
	"testing"
)

func TestPublish(t *testing.T) {
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

	err = mq.Publish(context.Background(), queue, []byte("test"))
	if err != nil {
		t.Fatal(err)
	}
}
