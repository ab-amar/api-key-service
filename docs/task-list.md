# API Key Service Learning Task List

1. Initialize the Go module for the project
2. Start a minimal HTTP server with Chi
3. Add `/health` and `/ready` endpoints
4. Add method checking and explicit response headers
5. Introduce a dedicated server setup function to keep `main.go` small
6. Add basic configuration for the server port
7. Read configuration from environment variables

-----------------------------------------

8. Validate configuration at startup
9. Add graceful shutdown support
10. Add a root endpoint to explain what the service is
11. Define the core API key model in `internal/model`
12. Define a repository interface in `internal/repository`
13. Design the first database schema for API keys
14. Add SQL migrations strategy
15. Add created-at, updated-at, and revoked-at fields
16. Add PostgreSQL connection configuration
17. Create a `pgxpool` connection pool in `main.go`
18. Add a PostgreSQL repository implementation
19. Apply the initial migration to a local PostgreSQL instance
20. Move API key creation logic into `internal/service`
21. Define request and response structs for API key creation
22. Add a `POST /keys` endpoint skeleton
23. Generate secure API keys using Go standard library
24. Store hashed API keys instead of raw keys
25. Wire handler -> service -> repository dependencies
26. Add a `POST /auth/validate` endpoint
27. Add API key lookup by hash
28. Add `POST /keys/{id}/revoke`
29. Add unit tests for health/readiness handlers
30. Add unit tests for key creation handler
31. Add unit tests for the service layer
32. Add unit tests for the PostgreSQL repository layer
33. Improve error responses to be consistent
34. Introduce a small response-writing helper
35. Add structured logging with standard library basics
36. Add request logging middleware
37. Add middleware chaining
38. Add panic recovery middleware
39. Add request IDs
40. Add timeout handling at the server level
41. Manually test create, validate, and revoke against PostgreSQL
42. Add integration tests for HTTP endpoints
43. Add integration tests for repository behavior
44. Add authentication middleware using API keys
45. Implement in-memory token bucket rate limiting per API key
46. Return `429 Too Many Requests` for exhausted keys
47. Add basic metrics counters
48. Add GitHub Actions for formatting and tests
49. Add API documentation
50. Write a design doc for the project architecture
