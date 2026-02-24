package retriever

import (
	"context"
	"rag-poc/internal/classifier"
	"rag-poc/pkg/helper"
	"strings"

	"go.uber.org/zap"
)

type ManualRepo interface {
	FindByDomain(
		ctx context.Context,
		domain classifier.Domain,
		topK int,
	) ([]string, error)
}

type ManualRetriever struct {
	log  *zap.Logger
	repo ManualRepo
}

func NewManualRetriever(
	log *zap.Logger,
	r ManualRepo,
) *ManualRetriever {
	return &ManualRetriever{log: log, repo: r}
}

func (r *ManualRetriever) Search(
	ctx context.Context,
	q string,
	domain classifier.Domain,
	topK int,
) ([]string, error) {

	rows, err := r.repo.FindByDomain(ctx, domain, topK)
	if err != nil {
		return nil, err
	}

	q = strings.ToLower(q)
	words := strings.Fields(q)

	var hits []string

	for _, html := range rows {
		txt := strings.ToLower(helper.StripHTML(html))

		score := 0
		for _, w := range words {
			if len(w) < 3 {
				continue
			}
			if strings.Contains(txt, w) {
				score++
			}
		}

		if score > 0 {
			hits = append(hits, txt)
		}

		if len(hits) >= topK {
			break
		}
	}

	return hits, nil
}
