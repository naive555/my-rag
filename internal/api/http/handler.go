package http

import (
	"bufio"
	"rag-poc/internal/rag"

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	RAG *rag.Service
}

func NewHandler(ragSvc *rag.Service) *Handler {
	return &Handler{RAG: ragSvc}
}

func (h *Handler) Ask(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing q param")
	}

	resp, err := h.RAG.Query(c.Context(), q)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.SendString(resp)
}

func (h *Handler) AskStream(c *fiber.Ctx) error {
	q := c.Query("q")
	if q == "" {
		return c.Status(400).SendString("missing q")
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		w.WriteString("Answer: ")
		_ = h.RAG.QueryStream(c.Context(), q, func(token string) {
			w.WriteString(token)
			w.Flush()
		})
	})

	return nil
}
