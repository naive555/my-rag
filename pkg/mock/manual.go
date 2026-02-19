package mock

import (
	"context"
	"rag-poc/internal/classifier"
)

type FakeRepo struct{}

func (FakeRepo) FindByDomain(
	ctx context.Context,
	d classifier.Domain,
	topK int,
) ([]string, error) {
	return []string{
		"<h1>SMS</h1> SMS is used for fast delivery and high open rate.",
	}, nil
}
