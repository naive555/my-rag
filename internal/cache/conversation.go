package cache

import (
	"context"
	"time"
)

type Conversation struct {
	store ListStore
}

func NewConversation(store ListStore) *Conversation {
	return &Conversation{store: store}
}

func convKey(sid string) string {
	return "conv:" + sid
}

func summaryKey(sid string) string {
	return "conv:summary:" + sid
}

func (c *Conversation) AppendMessage(
	ctx context.Context,
	sid string,
	msg string,
	ttl time.Duration,
	max int64,
) error {
	key := convKey(sid)

	if err := c.store.RPush(ctx, key, msg); err != nil {
		return err
	}

	if max > 0 {
		if err := c.store.LTrim(ctx, key, -max, -1); err != nil {
			return err
		}
	}

	if ttl > 0 {
		if err := c.store.Expire(ctx, key, ttl); err != nil {
			return err
		}
	}

	return nil
}

func (c *Conversation) GetRecent(ctx context.Context, sid string,
	limit int64,
) ([]string, error) {
	return c.store.LRange(ctx, convKey(sid), -limit, -1)
}

func (c *Conversation) Clear(ctx context.Context, sid string) error {
	return c.store.Del(ctx, convKey(sid))
}
