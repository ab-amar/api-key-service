# Server Timeouts

Useful HTTP server timeouts include:

- `ReadHeaderTimeout`
- `ReadTimeout`
- `WriteTimeout`
- `IdleTimeout`

Why they matter:

- protect against slow clients
- reduce resource exhaustion
- improve resilience under bad network behavior
