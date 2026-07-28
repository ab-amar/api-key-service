# API Key Service

Production-oriented Go backend project for learning API key management, secure key storage, authentication middleware, and rate limiting.

## Project Structure

- `cmd/server` - application entry point
- `internal/handler` - HTTP handlers
- `internal/service` - business logic
- `internal/repository` - persistence layer
- `internal/model` - core data models
- `internal/config` - configuration loading and validation
- `internal/middleware` - HTTP middleware
- `internal/metrics` - application metrics
- `docs` - learning notes and task list
- `scripts` - helper scripts
- `migrations` - SQL migrations

See [docs/task-list.md](/Users/amarbehera/go/api-key-service/docs/task-list.md) for the implementation roadmap.

Concept notes for the removed discussion/design topics live in:

- [docs/context-and-lifecycle.md](/Users/amarbehera/go/api-key-service/docs/context-and-lifecycle.md)
- [docs/api-key-generation.md](/Users/amarbehera/go/api-key-service/docs/api-key-generation.md)
- [docs/interfaces-and-dependencies.md](/Users/amarbehera/go/api-key-service/docs/interfaces-and-dependencies.md)
- [docs/hashing-vs-encryption.md](/Users/amarbehera/go/api-key-service/docs/hashing-vs-encryption.md)
- [docs/testing-notes.md](/Users/amarbehera/go/api-key-service/docs/testing-notes.md)
- [docs/server-timeouts.md](/Users/amarbehera/go/api-key-service/docs/server-timeouts.md)
- [docs/persistence-notes.md](/Users/amarbehera/go/api-key-service/docs/persistence-notes.md)
- [docs/postgres-concepts.md](/Users/amarbehera/go/api-key-service/docs/postgres-concepts.md)
- [docs/schema-and-indexes.md](/Users/amarbehera/go/api-key-service/docs/schema-and-indexes.md)
- [docs/repository-comparison.md](/Users/amarbehera/go/api-key-service/docs/repository-comparison.md)
- [docs/api-key-auth-notes.md](/Users/amarbehera/go/api-key-service/docs/api-key-auth-notes.md)
- [docs/rate-limiting-notes.md](/Users/amarbehera/go/api-key-service/docs/rate-limiting-notes.md)
- [docs/rate-limiting-algorithms.md](/Users/amarbehera/go/api-key-service/docs/rate-limiting-algorithms.md)
- [docs/observability.md](/Users/amarbehera/go/api-key-service/docs/observability.md)
