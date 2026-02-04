package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	Server     Server
	OllamaURL  string
	LLMModel   string
	EmbedModel string

	MongoURI string
	DBName   string
	ColName  string

	TopK int
}

type Server struct {
	Port        string
	Development bool
}

func Load() (*Config, error) {
	// Optional .env for local dev
	_ = godotenv.Load()

	cfg := &Config{
		OllamaURL:  getEnv("OLLAMA_URL", ""),
		LLMModel:   getEnv("LLM_MODEL", "smallthinker"),
		EmbedModel: getEnv("EMBED_MODEL", "nomic-embed-text"),

		MongoURI: getEnv("MONGO_URI", ""),
		DBName:   getEnv("DB_NAME", "rag"),
		ColName:  getEnv("COL_NAME", "chunks"),

		TopK: getEnvInt("TOP_K", 4),
	}

	viper.AddConfigPath("./config")

	env := getEnv("NODE_ENV", "")
	switch env {
	case "production":
		viper.SetConfigName("config-prod")
	case "dev":
		viper.SetConfigName("config-dev")
	default:
		viper.SetConfigName("config")
	}

	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
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
