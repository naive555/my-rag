package http

import (
	"bufio"
	"rag-poc/internal/lang"
	"rag-poc/internal/rag"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Handler struct {
	log *zap.Logger
	RAG *rag.Service
}

func NewHandler(log *zap.Logger, ragSvc *rag.Service) *Handler {
	return &Handler{
		log: log, RAG: ragSvc}
}

func (h *Handler) Ask(c *fiber.Ctx) error {
	context := "Ask"

	log := Logger(c)
	log.Info(context)

	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing question")
	}

	lang := lang.Detect(q)

	resp, err := h.RAG.Query(c.Context(), q, lang)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	log.Info("context", zap.String("resp", resp[:min(len(resp), 100)]))

	return c.SendString(resp)
}

func (h *Handler) AskStream(c *fiber.Ctx) error {
	context := "AskStream"

	log := Logger(c)
	log.Info(context)

	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing question")
	}

	lang := lang.Detect(q)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		var resp string

		w.WriteString("Answer: ")
		_ = h.RAG.QueryStream(c.Context(), q, lang, func(token string) {
			resp += token
			w.WriteString(token)
			w.Flush()
		})

		log.Info("context", zap.String("resp", resp[:min(len(resp), 100)]))
	})

	return nil
}
