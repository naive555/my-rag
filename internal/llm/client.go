package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type LlmClient struct {
	log     *zap.Logger
	baseURL string
	model   string
	http    *http.Client
}

func NewLlmClient(log *zap.Logger, baseURL, model string, timeout time.Duration) *LlmClient {
	if baseURL == "" {
		panic("llm: baseURL is required")
	}
	if model == "" {
		panic("llm: model is required")
	}

	return &LlmClient{
		log:     log,
		baseURL: baseURL,
		model:   model,
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

type GenerateRequest struct {
	Model   string         `json:"model"`
	Prompt  string         `json:"prompt"`
	Stream  bool           `json:"stream"`
	Options map[string]any `json:"options"`
}

type GenerateResponse struct {
	Response string `json:"response"`
}

type StreamChunk struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

func (c *LlmClient) GeneratePrompt(prompt string) (string, error) {
	context := "llm.GeneratePrompt"

	c.log.Info("request",
		zap.String("context", context),
		zap.String("prompt_preview", prompt[:min(len(prompt), 100)]),
	)

	reqBody := GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]any{
			"temperature": 0.2,
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	url := c.baseURL + "/api/generate"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama error: %s", string(body))
	}

	var result GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	c.log.Info("request",
		zap.String("context", context),
		zap.String("prompt_preview", result.Response[:min(len(result.Response), 100)]),
	)

	return result.Response, nil
}

func (c *LlmClient) StreamGeneratePrompt(prompt string, onToken func(string)) error {
	context := "llm.StreamGeneratePrompt"

	c.log.Info("request",
		zap.String("context", context),
		zap.String("prompt_preview", prompt[:min(len(prompt), 100)]),
	)

	reqBody := GenerateRequest{
		Model:  c.model,
		Prompt: prompt,
		Stream: true,
		Options: map[string]any{
			"temperature": 0.2,
		},
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	url := c.baseURL + "/api/generate"

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var result string
	for scanner.Scan() {
		var chunk StreamChunk
		if err := json.Unmarshal(scanner.Bytes(), &chunk); err != nil {
			continue
		}
		if chunk.Response != "" {
			result += chunk.Response
			onToken(chunk.Response)
		}
		if chunk.Done {
			c.log.Info("request",
				zap.String("context", context),
				zap.String("prompt_preview", result[:min(len(result), 100)]),
			)
			break
		}
	}

	return scanner.Err()
}
