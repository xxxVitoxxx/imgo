package rabbitmq

import "testing"

func TestDeclareQueue(t *testing.T) {
	mq, err := NewMessageQueue("guest", "guest", "localhost")
	if err != nil {
		t.Fatal(err)
	}
	defer mq.ch.Close()

	queue, err := mq.DeclareReplicateQueue()
	if err != nil {
		t.Fatal(err)
	}
	if queue.Name != "replicate" {
		t.Fatal("queue name is not replicate")
	}
}
