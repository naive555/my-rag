package http

import (
	"context"
	"rag-poc/internal/lang"

	"github.com/gofiber/contrib/websocket"
)

type Req struct {
	Question string `json:"question"`
}

func (h *Handler) AskWS(c *websocket.Conn) {
	var req Req
	if err := c.ReadJSON(&req); err != nil {
		return
	}

	ctx := context.Background()

	lang := lang.Detect(req.Question)

	err := h.RAG.QueryStream(ctx, req.Question, lang, func(token string) {
		if err := c.WriteMessage(websocket.TextMessage, []byte(token)); err != nil {
			return
		}
	})

	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("[ERROR] "+err.Error()))
		return
	}

	c.WriteMessage(websocket.TextMessage, []byte("[DONE]"))
}
