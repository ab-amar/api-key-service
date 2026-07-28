# API Key Service Learning Task List

1. Start a minimal HTTP server with `net/http`
2. Move the health handler out of `main.go`
3. Add method checking to `/health`
4. Set explicit response headers and status code in `/health`
5. Handle `ListenAndServe()` errors properly
6. Run `gofmt` and learn standard Go formatting
7. Add a root endpoint to explain what the service is
8. Introduce a dedicated server setup function to keep `main.go` small
9. Add basic configuration for the server port
10. Read configuration from environment variables
11. Validate configuration at startup
12. Add graceful shutdown support
13. Add a `POST /keys` endpoint skeleton
14. Define request and response structs for API key creation
15. Generate secure API keys using Go standard library
16. Define the core API key model in `internal/model`
17. Move API key creation logic into `internal/service`
18. Add an in-memory repository implementation
19. Define a repository interface in `internal/repository`
20. Wire handler -> service -> repository dependencies
21. Store hashed API keys instead of raw keys
22. Add a `POST /auth/validate` endpoint
23. Add API key lookup by hash
24. Use Chi for key management routes
25. Add `POST /keys/{id}/revoke`
26. Add revoke status to the model
27. Add unit tests for the health handler
28. Add unit tests for key creation handler
29. Add unit tests for the service layer
30. Add unit tests for the repository layer
31. Improve error responses to be consistent
32. Introduce a small response-writing helper
33. Add structured logging with standard library basics
34. Add request logging middleware
35. Introduce middleware chaining
36. Add panic recovery middleware
37. Add request IDs
38. Add timeout handling at the server level
39. Add basic health vs readiness endpoint distinction
40. Add SQL migrations strategy
41. Add a PostgreSQL repository implementation
42. Add created-at, updated-at, and revoked-at fields
43. Add PostgreSQL connection configuration
44. Create a `pgxpool` connection pool in `main.go`
45. Wire the app to use `PostgresRepository`
46. Apply the initial migration to a local PostgreSQL instance
47. Manually test create, validate, and revoke against PostgreSQL
48. Add integration tests for HTTP endpoints
49. Add integration tests for repository behavior
50. Add authentication middleware using API keys
51. Implement in-memory token bucket rate limiting per API key
52. Return `429 Too Many Requests` for exhausted keys
53. Add basic metrics counters
54. Add GitHub Actions for formatting and tests
55. Add API documentation
56. Write a design doc for the project architecture
57. Prepare interview talking points for auth, hashing, and rate limiting
