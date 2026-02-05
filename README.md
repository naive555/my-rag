# my-rag

**RAG (Retrieval-Augmented Generation) Proof of Concept in Go**

This project demonstrates a clean, production-oriented architecture for building a RAG service in Go, with support for streaming responses, language handling, and pluggable LLM / vector backends. The focus is on clarity, extensibility, and real-world readiness.

---

## Features

- HTTP API built with Fiber
- RAG-style prompt orchestration
- Streaming and non-streaming LLM responses
- Ollama-compatible LLM client
- Embedding and vector search abstraction
- Redis cache support
- Clean hexagonal-style internal structure
- Ready for multi-language extension (EN / TH)

---

## Project Structure

```text
.
├── cmd
│   └── rag-server
│       └── main.go        # Application entrypoint
├── config
│   └── config.yaml        # Runtime configuration
├── internal
│   ├── api
│   │   └── http           # HTTP layer (router, handlers, middleware)
│   ├── cache              # Redis cache client
│   ├── classifier         # Language / domain classification
│   ├── config             # Config loader
│   ├── embed              # Embedding client
│   ├── llm                # LLM client (Ollama)
│   ├── rag                # RAG domain (policy, prompt, service)
│   └── vector             # Vector store implementation
├── Dockerfile
├── compose.yaml
├── go.mod
└── go.sum
```

---

## Requirements

- Go 1.25.6 or later
- Ollama
- Docker (optional)

---

## Environment Variables

- OLLAMA_URL=http://localhost:11434
- LLM_MODEL=qwen2.5:7b-instruct
- EMBED_MODEL=nomic-embed-text

---

## Running Locally

Start Ollama and pull the required models:

```text
ollama pull qwen2.5:7b-instruct
ollama pull nomic-embed-text
```

Run the server:

```text
go run ./cmd/rag-server
```

API Endpoints

Ask (non-streaming)

```text
GET /ask?q=your question
```

### Response:

- Plain text answer

Ask (streaming)

GET /ask-stream?q=your question

### Response:

- text/event-stream

---

## RAG Flow

1. Incoming question
2. Language or policy classification
3. Vector search (optional)
4. Prompt construction
5. LLM generation (streaming or non-streaming)

---

## Design Notes

- Core business logic lives in internal/rag
- HTTP layer is thin and stateless
- LLM and vector store implementations are replaceable
- Streaming is handled end-to-end

---

## Future Improvements

- Add real vector search integration
- Multi-language prompt policies
- Answer length and verbosity control
- Observability (metrics, tracing)
- Authentication and authorization middleware

---

## License

MIT

```

```
