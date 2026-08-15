# Revision Store

`revisionstore` is a concurrency-safe in-memory key/value store with optimistic revision checks.

Use revision zero when creating a key. Updates must supply the current positive revision and receive a record with the next revision.

Run the example with `go run ./cmd/revisionstore`.
