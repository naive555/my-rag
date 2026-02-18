package http

import (
	"context"
	"rag-poc/internal/lang"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type WSReq struct {
	Q   string `json:"q"`
	SID string `json:"sid"`
}

type WSResp struct {
	Type string `json:"type"`
	Data string `json:"data,omitempty"`
}

func (h *Handler) AskWS(c *websocket.Conn) {
	this := "AskWS"

	log := h.log.With(zap.String("handler", this))
	log.Info("ws connected")

	for {
		var req WSReq
		if err := c.ReadJSON(&req); err != nil {
			log.Info("ws closed")
			return
		}

		if req.Q == "" {
			_ = c.WriteJSON(WSResp{
				Type: "error",
				Data: "missing question",
			})
			continue
		}

		sid := req.SID
		if sid == "" {
			sid = uuid.NewString()
		}

		lang := lang.Detect(req.Q)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)

		err := h.RAG.QueryStream(
			ctx,
			sid,
			req.Q,
			lang,
			func(token string) {
				_ = c.WriteJSON(WSResp{
					Type: "token",
					Data: token,
				})
			},
		)

		cancel()

		if err != nil {
			_ = c.WriteJSON(WSResp{
				Type: "error",
				Data: err.Error(),
			})
			continue
		}

		_ = c.WriteJSON(WSResp{Type: "done"})
	}
}
