package storage

import (
	"context"
	"testing"
)

func TestHash(t *testing.T) {
	redis, err := NewRedis("default", "", "localhost", "0")
	if err != nil {
		t.Fatal(err)
	}
	defer redis.db.Close()

	err = redis.SetData(context.Background(), "test1", "1")
	if err != nil {
		t.Fatal(err)
	}

	value, err := redis.GetData(context.Background(), "test1")
	if err != nil {
		t.Fatal(err)
	}
	if value != "1" {
		t.Fatal("value is not 1")
	}

	redis.db.Del(context.Background(), "replicate")
}
