# HTTP timeout middleware

This package provides a small `net/http` middleware that returns a gateway timeout when a handler exceeds its deadline while preserving responses that complete in time.
