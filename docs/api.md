# API Documentation

Base URL for local development:

```text
http://localhost:8080
```

## Health

```http
GET /health
```

Returns whether the process is alive.

Example response:

```json
{
  "status": "healthy"
}
```

## Readiness

```http
GET /ready
```

Returns whether required dependencies are reachable. The current readiness check pings PostgreSQL.

Example response:

```json
{
  "status": "ready"
}
```

## Create API Key

```http
POST /keys
Content-Type: application/json
```

Request:

```json
{
  "name": "billing-service"
}
```

Response:

```json
{
  "id": "key_1785400000000000000_abcd1234abcd1234",
  "name": "billing-service",
  "api_key": "ak_...",
  "key_prefix": "ak_123456789",
  "created_at": "2026-07-30T12:00:00Z"
}
```

The raw `api_key` is returned only once. The database stores only the SHA-256 hash and a short prefix.

## Validate API Key

```http
POST /auth/validate
Content-Type: application/json
```

Request:

```json
{
  "api_key": "ak_..."
}
```

Valid response:

```json
{
  "valid": true,
  "id": "key_...",
  "name": "billing-service",
  "key_prefix": "ak_123456789"
}
```

Invalid or revoked response:

```json
{
  "valid": false
}
```

## Revoke API Key

```http
POST /keys/{id}/revoke
```

Response:

```json
{
  "id": "key_...",
  "status": "revoked"
}
```

## Protected Endpoint

```http
GET /v1/protected
Authorization: Bearer ak_...
```

This endpoint demonstrates API key authentication plus per-key token bucket rate limiting.

Success:

```json
{
  "key_id": "key_...",
  "message": "authenticated request accepted"
}
```

If the token bucket is exhausted:

```http
HTTP/1.1 429 Too Many Requests
Retry-After: 1
```

## Metrics

```http
GET /metrics
```

Returns simple in-memory counters:

```json
{
  "requests_total": 10,
  "api_keys_created_total": 1,
  "validations_total": 2,
  "rate_limited_total": 1
}
```
