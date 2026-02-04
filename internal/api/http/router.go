package http

import "github.com/gofiber/fiber/v2"

func Register(app *fiber.App, h *Handler) {
	app.Get("/ask", h.Ask)
	app.Get("/ask-stream", h.AskStream)
}
