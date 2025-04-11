package redisc

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	address := os.Getenv("REDIS_ADDRESS")
	if address == "" {
		address = "localhost:6379" // default address
	}

	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: os.Getenv("REDIS_PASSWORD"), // no password set
		DB:       0,                           // use default DB
	})

	return &RedisClient{client: client}
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Ping(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
