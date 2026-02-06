package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func LoggerMiddleware(log *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		reqID := c.Get("X-Request-ID")
		if reqID == "" {
			reqID = uuid.NewString()
		}

		c.Locals("logger", log.With(
			zap.String("request_id", reqID),
			zap.String("path", c.Path()),
		))
		return c.Next()
	}
}

func Logger(c *fiber.Ctx) *zap.Logger {
	if l, ok := c.Locals("logger").(*zap.Logger); ok {
		return l
	}
	return zap.NewNop()
}
