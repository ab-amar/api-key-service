# Architecture

This project uses a layered backend structure:

```text
HTTP request
  -> middleware
  -> handler
  -> service
  -> repository
  -> PostgreSQL
```

## Layers

`cmd/server`

Application entry point. It loads configuration, validates startup requirements, opens the PostgreSQL connection pool, wires dependencies, starts the HTTP server, and handles graceful shutdown.

`internal/handler`

Owns HTTP concerns: route registration, JSON decoding, JSON responses, status codes, and URL parameters.

`internal/service`

Owns business rules: API key generation, API key hashing, validation, revocation behavior, and input validation.

`internal/repository`

Owns persistence. The service depends on an interface, while `PostgresRepository` implements that interface using `pgxpool`.

`internal/middleware`

Owns cross-cutting HTTP behavior: request IDs, request logging, panic recovery, API key authentication, and token bucket rate limiting.

`internal/metrics`

Owns concurrency-safe in-memory counters using `sync/atomic`.

## Security Model

Raw API keys are generated with `crypto/rand` and returned only during creation. The service stores only:

- `key_hash`: SHA-256 hash used for lookup.
- `key_prefix`: short visible prefix for debugging and audit-style display.

This is similar to password/token storage design: leaked database rows should not expose directly usable API keys.

## Rate Limiting

The protected endpoint uses an in-memory token bucket per API key.

Current settings:

- Capacity: 5 requests.
- Refill rate: 1 token per second.

This is intentionally process-local. In a multi-instance deployment, this should move to Redis or another shared store so all instances enforce a consistent limit.

## Operational Behavior

The service includes:

- `/health` for liveness.
- `/ready` for PostgreSQL readiness.
- Structured JSON logs with `log/slog`.
- Request IDs via `X-Request-ID`.
- Server read, write, idle, and graceful shutdown timeouts.
- GitHub Actions CI with formatting and tests.

## Tradeoffs

The project avoids a large framework and keeps the architecture explicit. Chi is used only for routing because path parameters and route composition become awkward with raw `net/http` as the API grows.

The in-memory rate limiter is simple and easy to reason about, but it is not distributed. That limitation is documented because it matters in real system design discussions.
