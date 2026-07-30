# API Key Service

Production-oriented Go backend service for issuing, validating, revoking, and protecting API keys with PostgreSQL persistence, hashed key storage, Chi routing, authentication middleware, token bucket rate limiting, structured logging, metrics, and CI.

## Features

- Create API keys and return the raw key only once.
- Store only SHA-256 key hashes plus short key prefixes.
- Validate and revoke API keys through REST endpoints.
- Protect routes with `Authorization: Bearer <api_key>`.
- Enforce per-key in-memory token bucket rate limiting.
- Expose health, readiness, and metrics endpoints.
- Use PostgreSQL, SQL migrations, integration tests, and GitHub Actions CI.

## Tech Stack

- Go
- Chi
- PostgreSQL
- pgx
- log/slog
- GitHub Actions

## Run Locally

Create a PostgreSQL database and apply the migration:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/api_key_service?sslmode=disable"
psql "$DATABASE_URL" -f migrations/000001_create_api_keys_table.up.sql
```

Start the server:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/api_key_service?sslmode=disable"
export PORT=8080
go run ./cmd/server
```

Create a key:

```bash
curl -i -X POST http://localhost:8080/keys \
  -H "Content-Type: application/json" \
  -d '{"name":"billing-service"}'
```

Use the returned key:

```bash
curl -i http://localhost:8080/v1/protected \
  -H "Authorization: Bearer ak_replace_with_created_key"
```

Run tests:

```bash
go test ./...
```

## API

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/health` | Liveness check |
| `GET` | `/ready` | PostgreSQL readiness check |
| `POST` | `/keys` | Create an API key |
| `POST` | `/auth/validate` | Validate an API key |
| `POST` | `/keys/{id}/revoke` | Revoke an API key |
| `GET` | `/v1/protected` | Authenticated, rate-limited example route |
| `GET` | `/metrics` | In-memory service counters |

## Documentation

- [API reference](docs/api.md)
- [Architecture](docs/architecture.md)
- [Local development](docs/local-development.md)
- [Task list](docs/task-list.md)
