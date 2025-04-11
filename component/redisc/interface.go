package redisc

import (
	"context"
	"time"
)

type Redis interface {
	// Set key-value pair with expiration
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get value by key
	Get(ctx context.Context, key string) (string, error)
	// Delete key
	Del(ctx context.Context, key string) error
	// Check if key exists
	Exists(ctx context.Context, key string) (bool, error)
	// Set key-value pair if key doesn't exist
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	// Get multiple values by keys
	MGet(ctx context.Context, keys ...string) ([]string, error)
	// Set multiple key-value pairs
	MSet(ctx context.Context, pairs map[string]interface{}) error
}

var _ Redis = (*RedisComponent)(nil)
