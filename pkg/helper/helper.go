package helper

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetSID(c *fiber.Ctx) string {
	sid := c.Get("X-Session-ID")
	if sid == "" {
		sid = c.Query("sid")
	}
	if sid == "" {
		sid = uuid.NewString()
	}
	return sid
}

func Cut(s string, n int) string {
	if len(s) <= n {
		return s
	}

	return s[:n] + "..."
}
