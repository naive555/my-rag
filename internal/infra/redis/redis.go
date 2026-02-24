package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Redis struct {
	client *redis.Client
}

type Config struct {
	Addr     string
	Password string
	DB       int
}

func NewRedis(log *zap.Logger, cfg Config) (*Redis, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	log.Info("Redis connected!")

	return &Redis{client: rdb}, nil
}

func (r *Redis) Client() *redis.Client {
	return r.client
}
