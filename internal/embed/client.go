package embed

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type EmbedClient struct {
	baseURL string
	model   string
}

func NewEmbedClient(baseURL, model string) *EmbedClient {
	if baseURL == "" {
		panic("llm: baseURL is required")
	}
	if model == "" {
		panic("llm: model is required")
	}

	return &EmbedClient{
		baseURL: baseURL,
		model:   model,
	}
}

func (c *EmbedClient) Embed(text string) ([]float64, error) {
	reqBody := map[string]any{
		"model": c.model,
		"input": text,
	}

	b, _ := json.Marshal(reqBody)

	resp, err := http.Post(
		c.baseURL+"/api/embeddings",
		"application/json",
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var out struct {
		Embedding []float64 `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out.Embedding, nil
}
