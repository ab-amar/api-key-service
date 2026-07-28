# API Key Authentication Notes

Common extraction points for API keys:

- `Authorization` header
- custom header such as `X-API-Key`

For a small internal service, a custom header is usually simpler to start with.

Middleware should:

- extract the presented key
- hash it using the same scheme as storage
- look it up
- reject revoked or missing keys
