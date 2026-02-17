package cache

import (
	"context"
	"time"
)

func (r *Redis) AppendList(ctx context.Context, key string, val string, ttl time.Duration) error {
	if err := r.client.RPush(ctx, key, val).Err(); err != nil {
		return err
	}
	if ttl > 0 {
		r.client.Expire(ctx, key, ttl)
	}
	return nil
}

func (r *Redis) GetList(ctx context.Context, key string, limit int64) ([]string, error) {
	return r.client.LRange(ctx, key, -limit, -1).Result()
}

func (r *Redis) TrimList(ctx context.Context, key string, max int64) error {
	return r.client.LTrim(ctx, key, -max, -1).Err()
}
