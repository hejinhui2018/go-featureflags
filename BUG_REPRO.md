# Request body close reproduction

Run the focused test from this directory:

```text
go test ./... -run '^TestLimitMiddlewareClosesBodyWhenTooLarge$' -count=1
```

The test sends a body larger than the configured limit through the middleware. The expected behavior is HTTP 413, no downstream handler call, and a closed request body. The full suite also verifies that an in-limit body reaches the handler unchanged.
