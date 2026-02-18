package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisListStore struct {
	rdb *redis.Client
}

func NewRedisListStore(rdb *redis.Client) *RedisListStore {
	return &RedisListStore{rdb: rdb}
}

func (r *RedisListStore) Set(ctx context.Context, key string, val string, ttl time.Duration) error {
	return r.rdb.Set(ctx, key, val, ttl).Err()
}

func (r *RedisListStore) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *RedisListStore) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}

func (r *RedisListStore) RPush(ctx context.Context, key string, vals ...string) error {
	return r.rdb.RPush(ctx, key, vals).Err()
}

func (r *RedisListStore) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.rdb.LRange(ctx, key, start, stop).Result()
}

func (r *RedisListStore) LTrim(ctx context.Context, key string, start, stop int64) error {
	return r.rdb.LTrim(ctx, key, start, stop).Err()
}

func (r *RedisListStore) Expire(ctx context.Context, key string, ttl time.Duration) error {
	return r.rdb.Expire(ctx, key, ttl).Err()
}
