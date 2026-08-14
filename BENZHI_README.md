# Build and Run

This repository contains a multi-tenant feature flag HTTP service.

Standard commands:

```bash
go build ./...
go test ./...
go run ./cmd/server
```

The HTTP server listens on port 8080. Query a flag with `GET /v1/flags/{name}` and the `X-Tenant-ID` request header.
