package rag

import (
	"context"
	"fmt"
	"strings"

	"rag-poc/internal/llm"
)

type Service struct {
	llm *llm.LlmClient
}

func NewService(llmClient *llm.LlmClient) *Service {
	if llmClient == nil {
		panic("rag: llmClient is required")
	}

	return &Service{
		llm: llmClient,
	}
}

func (s *Service) Query(ctx context.Context, q, lang string) (string, error) {
	if strings.TrimSpace(q) == "" {
		return "", fmt.Errorf("empty query")
	}

	chunks := s.retrieve(ctx, q)

	prompt := BuildPrompt(chunks, q, lang)

	return s.llm.GeneratePrompt(prompt)
}

func (s *Service) QueryStream(ctx context.Context, q, lang string, onToken func(string)) error {
	if strings.TrimSpace(q) == "" {
		return fmt.Errorf("empty query")
	}

	chunks := s.retrieve(ctx, q)

	prompt := BuildPrompt(chunks, q, lang)

	return s.llm.StreamGeneratePrompt(prompt, onToken)
}

func (s *Service) retrieve(ctx context.Context, q string) []string {
	// FUTURE:
	// - embed(q)
	// - vector search
	// - permission filter
	// - rerank
	// - return top chunks

	// For now: no context
	return []string{}
}
