package vector

import "context"

type Document struct {
	ID     string
	Domain string
	Text   string
	Vector []float32
	Meta   map[string]string
}

type Hit struct {
	ID     string
	Text   string
	Score  float64
	Domain string
}

type Store interface {
	Upsert(ctx context.Context, docs []Document) error
	Search(ctx context.Context, vec []float32, domain string, topK int) ([]Hit, error)
}
