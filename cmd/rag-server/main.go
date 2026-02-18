package main

import (
	"time"

	"rag-poc/internal/api/http"
	"rag-poc/internal/cache"
	"rag-poc/internal/config"
	"rag-poc/internal/llm"
	"rag-poc/internal/rag"
	"rag-poc/internal/redis"
	"rag-poc/pkg/logger"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log, err := logger.New(cfg.Env)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	log.Info("config loaded",
		zap.String("port", cfg.Server.Port),
	)

	app := fiber.New()

	app.Use(http.LoggerMiddleware(log))

	redisClient, err := redis.NewRedis(redis.Config{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPass,
		DB:       cfg.RedisDB,
	})
	if err != nil {
		panic(err)
	}

	store := cache.NewRedisListStore(redisClient.Client())
	conv := cache.NewConversation(store)

	llmClient := llm.NewLlmClient(
		log,
		cfg.OllamaURL,
		cfg.LLMModel,
		5*time.Minute,
	)

	ragSvc := rag.NewService(cfg, log, llmClient, conv)

	h := http.NewHandler(log, ragSvc)
	http.Register(app, h)

	log.Fatal("server stopped",
		zap.Error(app.Listen(":"+cfg.Server.Port)),
	)
}
