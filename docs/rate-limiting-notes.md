# Rate Limiting Notes

Rate limiting protects services from:

- abuse
- noisy clients
- accidental traffic spikes

For this project, in-memory per-key limiting is enough to learn:

- middleware design
- shared state
- concurrency
- rejection with `429 Too Many Requests`
