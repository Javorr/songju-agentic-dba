# Songju Agentic DBA

A high-performance, AI-powered database observability engine that automatically detects and optimizes slow SQL queries — securely.

## What is this?

Songju Agentic DBA is a **security-first database performance analyzer** that uses a multi-agent AI system to optimize SQL queries without ever exposing sensitive data. Think of it as an autonomous DBA that:

- 🔒 **Masks PII automatically** — emails, SSNs, credit cards, and IDs are redacted before they leave your infrastructure
- ⚡ **Handles 10,000+ logs/sec** — built with Go and NATS JetStream for microsecond-latency ingestion
- 🤖 **Uses specialized AI agents** — a "Librarian" finds relevant schemas, a "Coder" suggests optimizations, and an "Executor" validates them safely
- 🏗️ **Schema-blind security** — AI agents never see your actual data, only metadata and masked queries

## Why I built this

As databases grow and query performance becomes critical, DBAs are increasingly overwhelmed. This project explores how **agentic AI** can assist with database observability while maintaining strict security boundaries. It's a practical application of:

- **Distributed systems** — Go, NATS JetStream, microservices
- **AI agent orchestration** — multi-agent workflows with specialized roles
- **Privacy-preserving architectures** — PII masking, schema isolation, hybrid inference

## Tech Stack

| Component | Technology |
|-----------|-----------|
| **Ingestor** | Go 1.26 (high-concurrency HTTP API) |
| **Message Broker** | NATS JetStream (persistent, at-least-once delivery) |
| **AI Agents** | Python (PydanticAI / LangGraph) |
| **Database** | PostgreSQL (metadata store + shadow DB) |
| **Infrastructure** | Docker & Docker Compose |

## Architecture

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────┐
│   App/DB    │────▶│  Go Ingestor │────▶│  NATS JetStream │
│  (logs)     │     │  (PII Mask)  │     │  (db.logs.slow) │
└─────────────┘     └──────────────┘     └────────┬────────┘
                                                   │
                              ┌────────────────────┼────────────────────┐
                              ▼                    ▼                    ▼
                       ┌───────────┐      ┌───────────┐       ┌───────────┐
                       │ Librarian │─────▶│   Coder   │──────▶│ Executor  │
                       │  (Agent A)│      │  (Agent B)│       │  (Agent C)│
                       └───────────┘      └───────────┘       └───────────┘
```

### The Multi-Agent Workflow

1. **Agent A (The Librarian)** — Receives a slow query, looks up relevant table schemas from the metadata store
2. **Agent B (The Coder)** — Proposes optimized SQL or index suggestions based on the schema
3. **Agent C (The Executor)** — Validates the optimization in a sandbox and runs `EXPLAIN ANALYZE`

## Quick Start

### Prerequisites
- Go 1.26+
- Docker & Docker Compose

### Run the Ingestor

```bash
# Start NATS JetStream
docker-compose up -d

# Run the Go service
go run main.go
```

### Test It Out

```bash
# Health check
curl http://localhost:8080/health

# Ingest a slow query log
curl -X POST http://localhost:8080/ingest \
  -H "Content-Type: application/json" \
  -d '{
    "sql": "SELECT * FROM users WHERE email = \"test@example.com\"",
    "service": "api",
    "duration_ms": 1500
  }'
```

The SQL above will be automatically masked to:
```sql
SELECT * FROM users WHERE email = "[REDACTED_EMAIL]"
```

## Project Status

### ✅ Phase 1: High-Speed Ingestion (Complete)
- Go-based HTTP ingestor with PII masking
- NATS JetStream integration with persistent storage
- Benchmarked at 10,000+ logs/sec

### 🚧 Phase 2: Python Worker & NATS Consumer (Planned)
- NATS consumer using Python
- "Librarian" agent for schema metadata

### 📋 Phase 3: Reasoning & Validation (Planned)
- "Coder" and "Executor" agents
- Shadow DB sandbox for safe query validation
- Automated optimization reports

## Key Features

### Security-First Design
- **PII Masking** — Automatically redacts emails, SSNs, credit cards, UUIDs, and numeric IDs
- **Schema-Blind Agents** — AI never accesses production data directly
- **Hybrid Inference** — Sensitive tasks route to local models (Ollama), logic tasks to GPT-4o

### Performance
- **Go Workers** — Goroutine-based worker pool for parallel log processing
- **NATS JetStream** — Low-latency message broker with persistence
- **Horizontal Scaling** — Add more Python workers to scale reasoning capacity

## Running Tests

```bash
go test -v
go test -bench=.  # Run benchmarks
```

## License

This project is licensed under the PolyForm Noncommercial License 1.0.0 — free for personal and research use. See [LICENSE](LICENSE) for details.

---

*Built by Javier Martínez Segura — exploring the intersection of AI and database infrastructure.*
