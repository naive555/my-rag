package rag

import (
	"context"
	"fmt"
	"strings"
	"time"

	"rag-poc/internal/cache"
	"rag-poc/internal/classifier"
	"rag-poc/internal/config"
	"rag-poc/internal/embed"
	"rag-poc/internal/llm"
	"rag-poc/internal/vector"

	"go.uber.org/zap"
)

type Retriever interface {
	Search(ctx context.Context, q string, domain classifier.Domain, topK int) ([]string, error)
}

type Service struct {
	cfg       *config.Config
	log       *zap.Logger
	llm       *llm.LlmClient
	conv      *cache.Conversation
	embed     *embed.Client
	vector    *vector.MongoStore
	retriever Retriever
}

func NewService(cfg *config.Config, log *zap.Logger, conv *cache.Conversation, embed *embed.Client, vector *vector.MongoStore, r Retriever, llmClient *llm.LlmClient) *Service {
	if llmClient == nil {
		panic("rag: llmClient is required")
	}

	return &Service{
		cfg:       cfg,
		log:       log,
		conv:      conv,
		embed:     embed,
		vector:    vector,
		retriever: r,
		llm:       llmClient,
	}
}

func (s *Service) Query(ctx context.Context, sid, q, lang string) (string, error) {
	if strings.TrimSpace(q) == "" {
		return "", fmt.Errorf("empty query")
	}

	res := classifier.Classify(q)
	policy := classifier.Resolve(res.Domain)

	var chunks []string
	if policy.UseRAG {
		chunks = s.retrieve(ctx, q, res.Domain, policy.TopK)
	}

	history, _ := s.conv.GetRecent(ctx, sid, 10)

	prompt := BuildPromptWithHistory(
		chunks,
		history,
		q,
		lang,
		policy.MaxSentences,
	)

	answer, err := s.llm.GeneratePrompt(prompt)
	if err != nil {
		return "", err
	}

	redisCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_ = s.conv.AppendMessage(redisCtx, sid, "U:"+q, time.Duration(s.cfg.ConvTtl)*time.Minute, int64(s.cfg.ConvMax))
	_ = s.conv.AppendMessage(redisCtx, sid, "A:"+answer, time.Duration(s.cfg.ConvTtl)*time.Minute, int64(s.cfg.ConvMax))

	return answer, nil
}

func (s *Service) QueryStream(
	ctx context.Context,
	sid, q, lang string,
	onToken func(string),
) error {

	if strings.TrimSpace(q) == "" {
		return fmt.Errorf("empty query")
	}

	res := classifier.Classify(q)
	policy := classifier.Resolve(res.Domain)

	var chunks []string
	if policy.UseRAG {
		chunks = s.retrieve(ctx, q, res.Domain, policy.TopK)
	}

	history, _ := s.conv.GetRecent(ctx, sid, 10)

	prompt := BuildPromptWithHistory(
		chunks,
		history,
		q,
		lang,
		policy.MaxSentences,
	)

	var sb strings.Builder

	err := s.llm.StreamGeneratePrompt(prompt, func(token string) {
		sb.WriteString(token)
		onToken(token)
	})
	if err != nil {
		return err
	}

	answer := sb.String()

	redisCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_ = s.conv.AppendMessage(redisCtx, sid, "U:"+q, time.Duration(s.cfg.ConvTtl)*time.Minute, int64(s.cfg.ConvMax))
	_ = s.conv.AppendMessage(redisCtx, sid, "A:"+answer, time.Duration(s.cfg.ConvTtl)*time.Minute, int64(s.cfg.ConvMax))

	return nil
}

func (s *Service) retrieve(
	ctx context.Context,
	q string,
	domain classifier.Domain,
	topK int,
) []string {

	if topK <= 0 {
		return nil
	}

	if s.vector != nil && s.embed != nil {
		vec, err := s.embed.Embed(ctx, q)
		if err == nil {
			hits, err := s.vector.Search(ctx, vec, string(domain), topK)
			if err == nil {
				out := make([]string, 0, len(hits))
				for _, h := range hits {
					out = append(out, h.Text)
				}
				return out
			}
		}
		s.log.Warn("vector search failed, fallback")
	}

	if s.retriever != nil {
		chunks, err := s.retriever.Search(ctx, q, domain, topK)
		if err != nil {
			s.log.Warn("retriever error", zap.Error(err))
			return nil
		}
		return chunks
	}

	return nil
}
