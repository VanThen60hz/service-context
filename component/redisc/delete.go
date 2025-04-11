package redisc

import (
	"context"
)

func (r *RedisComponent) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}
