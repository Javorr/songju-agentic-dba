# Songju Agentic DBA

A learning project exploring AI-assisted database observability.

## What it does

Accepts HTTP POST requests with SQL query logs, masks sensitive data (PII), and publishes to NATS JetStream for asynchronous processing.

## What's implemented

- Go HTTP server with `/ingest` and `/health` endpoints
- PII masking for emails, SSNs, credit cards, numeric IDs, and UUIDs
- NATS JetStream integration (creates stream, publishes messages)
- Unit tests for PII masking
- Docker and Docker Compose configuration
- Benchmark definitions for PII masking function

## What's not implemented

- Python AI agents (the "Librarian", "Coder", "Executor" agents)
- Worker pool for concurrent log processing
- Performance testing at scale
- PostgreSQL metadata store
- Shadow database sandbox for query validation
- Consumer/worker to read from NATS

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Ingestor | Go 1.26 (net/http) |
| Message Broker | NATS JetStream |
| AI Agents | Python (not yet implemented) |
| Database | PostgreSQL (not yet implemented) |

## Architecture (planned)

When complete, the system will work like this:

```
App/DB logs ──▶ Go Ingestor (PII Mask) ──▶ NATS JetStream
                                                   │
                      ┌────────────────────────────┼────────────────────┐
                      ▼                            ▼                    ▼
               Librarian Agent              Coder Agent          Executor Agent
               (finds schemas)             (suggests SQL)       (validates safely)
```

### The Multi-Agent Workflow (planned)

1. **Agent A (Librarian)** — Receives a slow query, looks up relevant table schemas
2. **Agent B (Coder)** — Proposes optimized SQL or index suggestions
3. **Agent C (Executor)** — Validates in a sandbox with `EXPLAIN ANALYZE`

**Note:** This project targets **relational databases** (PostgreSQL, MySQL, etc.) since it works with SQL queries, table schemas, and `EXPLAIN ANALYZE`.

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

### Test It

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

The SQL gets masked automatically:
```sql
SELECT * FROM users WHERE email = "[REDACTED_EMAIL]"
```

## Running Tests

```bash
go test -v
go test -bench=.  # Run benchmarks
```

## Why I built this

To learn:
- Building services in Go
- NATS JetStream messaging
- AI agent system design
- Privacy-preserving architectures

## License

MIT License — see [LICENSE](LICENSE) for details.

---

*Learning project by Javier Martínez Segura.*
