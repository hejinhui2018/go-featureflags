# Feature Flags Service

A small multi-tenant feature flag service written in Go. It exposes a JSON HTTP API and keeps frequently requested flag values in memory.

## Run

```bash
go test ./...
go run ./cmd/server
```

The server listens on `:8080`. Requests must include an `X-Tenant-ID` header:

```bash
curl -H "X-Tenant-ID: tenant-a" http://localhost:8080/v1/flags/checkout
```
