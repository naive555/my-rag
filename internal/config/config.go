package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort    string
	OllamaURL  string
	LLMModel   string
	EmbedModel string

	MongoURI string
	DBName   string
	ColName  string

	TopK int
}

func Load() (*Config, error) {
	// Optional .env for local dev
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:    getEnv("APP_PORT", "3000"),
		OllamaURL:  getEnv("OLLAMA_URL", ""),
		LLMModel:   getEnv("LLM_MODEL", "smallthinker"),
		EmbedModel: getEnv("EMBED_MODEL", "nomic-embed-text"),

		MongoURI: getEnv("MONGO_URI", ""),
		DBName:   getEnv("DB_NAME", "rag"),
		ColName:  getEnv("COL_NAME", "chunks"),

		TopK: getEnvInt("TOP_K", 4),
	}

	// Required validation
	if cfg.OllamaURL == "" {
		return nil, fmt.Errorf("missing env: OLLAMA_URL")
	}
	// if cfg.MongoURI == "" {
	// 	return nil, fmt.Errorf("missing env: MONGO_URI")
	// }

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
