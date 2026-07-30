# Local Development

## Requirements

- Go installed.
- PostgreSQL installed locally or available remotely.
- `psql` available if you want to apply migrations from the terminal.

## Create Local Database

Example using `psql`:

```bash
createdb api_key_service
```

If `createdb` is unavailable, open `psql` and run:

```sql
CREATE DATABASE api_key_service;
```

## Apply Migration

From the project root:

```bash
psql "$DATABASE_URL" -f migrations/000001_create_api_keys_table.up.sql
```

Example `DATABASE_URL`:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/api_key_service?sslmode=disable"
```

## Run The Server

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/api_key_service?sslmode=disable"
export PORT=8080
go run ./cmd/server
```

## Smoke Test

Create an API key:

```bash
curl -i -X POST http://localhost:8080/keys \
  -H "Content-Type: application/json" \
  -d '{"name":"billing-service"}'
```

Validate the key:

```bash
curl -i -X POST http://localhost:8080/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"api_key":"ak_replace_with_created_key"}'
```

Call the protected endpoint:

```bash
curl -i http://localhost:8080/v1/protected \
  -H "Authorization: Bearer ak_replace_with_created_key"
```

Revoke the key:

```bash
curl -i -X POST http://localhost:8080/keys/key_replace_with_id/revoke
```

## Tests

Run all unit tests:

```bash
go test ./...
```

Run repository integration tests against PostgreSQL:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/api_key_service?sslmode=disable"
go test ./internal/repository -v
```
