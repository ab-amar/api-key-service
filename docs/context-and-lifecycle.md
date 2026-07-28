# Context And Request Lifecycle

`context.Context` carries request-scoped values, deadlines, and cancellation signals through the call chain.

For an HTTP service, the request lifecycle usually means:

- request enters the router
- middleware runs
- handler executes
- downstream work uses the request context
- response is written

This matters because:

- timeouts should stop long-running work
- canceled client requests should not keep wasting server resources
- request-scoped metadata can be propagated cleanly
