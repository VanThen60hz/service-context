package redisc

import (
	"context"
	"time"
)

type Redis interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	MGet(ctx context.Context, keys ...string) ([]string, error)
	MSet(ctx context.Context, pairs map[string]interface{}) error
}
