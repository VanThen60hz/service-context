package redisc

import (
	"context"
	"time"
)

func (r *RedisComponent) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisComponent) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.client.SetNX(ctx, key, value, expiration).Result()
}

func (r *RedisComponent) MSet(ctx context.Context, pairs map[string]interface{}) error {
	args := make([]interface{}, 0, len(pairs)*2)
	for k, v := range pairs {
		args = append(args, k, v)
	}
	return r.client.MSet(ctx, args...).Err()
}
