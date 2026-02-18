package cache

import (
	"context"
	"time"
)

func (c *Conversation) SetSummary(
	ctx context.Context,
	sid string,
	summary string,
	ttl time.Duration,
) error {
	return c.store.Set(ctx, summaryKey(sid), summary, ttl)
}

func (c *Conversation) GetSummary(
	ctx context.Context,
	sid string,
) (string, error) {
	return c.store.Get(ctx, summaryKey(sid))
}
