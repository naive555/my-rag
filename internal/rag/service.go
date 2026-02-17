package rag

import (
	"context"
	"fmt"
	"strings"

	"rag-poc/internal/cache"
	"rag-poc/internal/classifier"
	"rag-poc/internal/llm"

	"go.uber.org/zap"
)

type Service struct {
	log   *zap.Logger
	llm   *llm.LlmClient
	redis *cache.Redis
}

func NewService(log *zap.Logger, llmClient *llm.LlmClient, redisClient *cache.Redis) *Service {
	if llmClient == nil {
		panic("rag: llmClient is required")
	}

	return &Service{
		log:   log,
		llm:   llmClient,
		redis: redisClient,
	}
}

func (s *Service) Query(ctx context.Context, q, lang string) (string, error) {
	if strings.TrimSpace(q) == "" {
		return "", fmt.Errorf("empty query")
	}

	res := classifier.Classify(q)
	policy := classifier.Resolve(res.Domain)

	var chunks []string
	if policy.UseRAG {
		chunks = s.retrieve(ctx, q, res.Domain, policy.TopK)
	}

	prompt := BuildPrompt(chunks, q, lang, policy.MaxSentences)

	return s.llm.GeneratePrompt(prompt)
}

func (s *Service) QueryStream(ctx context.Context, q, lang string, onToken func(string)) error {
	if strings.TrimSpace(q) == "" {
		return fmt.Errorf("empty query")
	}

	res := classifier.Classify(q)
	policy := classifier.Resolve(res.Domain)

	var chunks []string
	if policy.UseRAG {
		chunks = s.retrieve(ctx, q, res.Domain, policy.TopK)
	}

	prompt := BuildPrompt(chunks, q, lang, policy.MaxSentences)

	return s.llm.StreamGeneratePrompt(prompt, onToken)
}

func (s *Service) retrieve(ctx context.Context, q string, domain classifier.Domain, topK int) []string {
	// FUTURE:
	// - embed(q)
	// - vector search
	// - permission filter
	// - rerank
	// - return top chunks

	// For now: no context
	return []string{}
}
