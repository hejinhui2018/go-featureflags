# Record Store

This small HTTP service exposes records from an in-memory store.

Run it with:

```text
go run ./cmd/server
```

Read a record with `GET /v1/records/{key}`. A found record returns JSON with its key and value; a missing record returns a not-found error.
