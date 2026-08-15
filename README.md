# TTL Store

`ttlstore` is a concurrency-safe in-memory key/value store with optional expiration.

Pass a positive TTL to expire a value after a duration. A zero TTL stores a value until it is updated or deleted. The clock can be supplied for deterministic testing.

Run the example with `go run ./cmd/ttlstore`.
