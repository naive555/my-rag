package cache

import (
	"context"
	"time"
)

func (r *Redis) SetSummary(ctx context.Context, sessionID string, summary string, ttl time.Duration) error {
	key := "chat:summary:" + sessionID
	return r.Set(ctx, key, summary, ttl)
}

func (r *Redis) GetSummary(ctx context.Context, sessionID string) (string, error) {
	key := "chat:summary:" + sessionID
	return r.Get(ctx, key)
}
