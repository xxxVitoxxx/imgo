package storage

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

// a hash key for store replicate id and line user token
const replicateKey = "replicate"

// SetData store replicate id and line user token
func (r *Redis) SetData(ctx context.Context, key, value string) error {
	return r.db.HSet(ctx, replicateKey, key, value).Err()
}

// GetData get line user token by replicate id
func (r *Redis) GetData(ctx context.Context, key string) (string, error) {
	value, err := r.db.HGet(ctx, replicateKey, key).Result()
	if err == redis.Nil {
		return "", errors.New("key does not exist")
	}
	if err != nil {
		return "", err
	}

	return value, nil
}
