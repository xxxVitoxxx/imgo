package storage

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	db *redis.Client
}

// NewRedis return a new redis instance
func NewRedis(user, password, address, db string) (*Redis, error) {
	client, err := connect(user, password, address, db)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}
	return &Redis{client}, nil
}

func connect(user, password, address, db string) (*redis.Client, error) {
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:6379/%s",
		user,
		password,
		address,
		db,
	))
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opt)

	return client, nil
}
