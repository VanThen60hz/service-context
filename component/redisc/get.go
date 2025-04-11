package redisc

import (
	"context"
)

func (r *RedisComponent) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisComponent) MGet(ctx context.Context, keys ...string) ([]string, error) {
	vals, err := r.client.MGet(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	result := make([]string, len(vals))
	for i, val := range vals {
		if val == nil {
			result[i] = ""
		} else {
			result[i] = val.(string)
		}
	}

	return result, nil
}

func (r *RedisComponent) Exists(ctx context.Context, key string) (bool, error) {
	n, err := r.client.Exists(ctx, key).Result()
	return n > 0, err
}
