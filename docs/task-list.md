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
13. Understand request lifecycle and `context.Context`
14. Add a `POST /keys` endpoint skeleton
15. Learn JSON decoding with `encoding/json`
16. Define request and response structs for API key creation
17. Discuss API key generation requirements
18. Generate secure API keys using Go standard library
19. Define the core API key model in `internal/model`
20. Move API key creation logic into `internal/service`
21. Design a service interface and discuss why interfaces matter
22. Add an in-memory repository implementation
23. Define a repository interface in `internal/repository`
24. Wire handler -> service -> repository dependencies
25. Store hashed API keys instead of raw keys
26. Explain hashing vs encryption for API keys
27. Add a `POST /auth/validate` endpoint
28. Add API key lookup by hash
29. Introduce Chi and route structure for key management
30. Add `POST /keys/{id}/revoke`
31. Add revoke status to the model
32. Add unit tests for the health handler
33. Learn table-driven tests in Go
34. Add unit tests for key creation handler
35. Add unit tests for the service layer
36. Add unit tests for the repository layer
37. Learn how to use `httptest`
38. Introduce `testify` after standard testing and `httptest`
39. Improve error responses to be consistent
40. Introduce a small response-writing helper
41. Add structured logging with standard library basics
42. Add request logging middleware
43. Introduce middleware chaining
44. Add panic recovery middleware
45. Add request IDs
46. Add timeout handling at the server level
47. Discuss server timeouts: read, write, idle, header
48. Add basic health vs readiness endpoint distinction
49. Add persistent storage design discussion
50. Introduce PostgreSQL concepts before implementation
51. Design the first database schema for API keys
52. Discuss primary keys, unique constraints, and indexes
53. Add SQL migrations strategy
54. Add a PostgreSQL repository implementation
55. Compare in-memory vs PostgreSQL repository behavior
56. Add created-at, updated-at, and revoked-at fields
57. Add PostgreSQL connection configuration
58. Create a `pgxpool` connection pool in `main.go`
59. Wire the app to use `PostgresRepository`
60. Apply the initial migration to a local PostgreSQL instance
61. Manually test create, validate, and revoke against PostgreSQL
62. Add integration tests for HTTP endpoints
63. Add integration tests for repository behavior
64. Add authentication middleware using API keys
65. Explain API key extraction from headers
66. Introduce rate limiting concepts before implementation
67. Compare fixed window vs sliding window vs token bucket
68. Implement in-memory token bucket rate limiting per API key
69. Return `429 Too Many Requests` for exhausted keys
70. Add observability discussion: logs, metrics, traces
71. Add basic metrics counters
72. Add GitHub Actions for formatting and tests
73. Add API documentation
74. Write a design doc for the project architecture
75. Prepare interview talking points for auth, hashing, and rate limiting
