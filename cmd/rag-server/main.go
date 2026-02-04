package main

import (
	"log"
	"rag-poc/internal/api/http"
	"rag-poc/internal/config"
	"rag-poc/internal/llm"
	"rag-poc/internal/rag"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Config loaded, port:", cfg.AppPort)

	llmClient := llm.NewLlmClient(
		cfg.OllamaURL,
		cfg.LLMModel,
		5*time.Minute,
	)

	ragSvc := rag.NewService(llmClient)

	h := http.NewHandler(ragSvc)
	http.Register(app, h)

	log.Fatal(app.Listen(":" + cfg.AppPort))
}
