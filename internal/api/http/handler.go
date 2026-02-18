package http

import (
	"bufio"
	"context"
	"rag-poc/internal/lang"
	"rag-poc/internal/rag"
	"rag-poc/pkg/helper"
	"strings"
	"time"

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
	this := "Ask"

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Minute)
	defer cancel()

	log := Logger(c)

	log.Info(this, zap.Any("headers", c.GetReqHeaders()), zap.Any("queries", c.Queries()))

	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing question")
	}

	lang := lang.Detect(q)

	sid := helper.GetSID(c)

	resp, err := h.RAG.Query(ctx, sid, q, lang)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	c.Set("X-Session-ID", sid)

	log.Info(this, zap.String("resp", helper.Cut(resp, 100)))

	return c.SendString(resp)
}

func (h *Handler) AskStream(c *fiber.Ctx) error {
	this := "AskStream"

	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Minute)
	defer cancel()

	log := Logger(c)

	log.Info(this, zap.Any("headers", c.GetReqHeaders()), zap.Any("queries", c.Queries()))

	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing question")
	}

	lang := lang.Detect(q)

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	sid := helper.GetSID(c)

	var resp strings.Builder
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {

		err := h.RAG.QueryStream(ctx, sid, q, lang, func(token string) {
			resp.WriteString(token)
			w.WriteString(token)
			w.Flush()
		})
		if err != nil {
			w.WriteString("event: error\ndata: " + err.Error() + "\n\n")
			w.Flush()
			return
		}
	})

	c.Set("X-Session-ID", sid)

	log.Info(this, zap.String("resp", helper.Cut(resp.String(), 100)))

	return nil
}
