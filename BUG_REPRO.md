# Missing value is rejected without changing storage

## Observed behavior

`PUT /v1/records/{key}` with body `{}` returns `200` and writes an empty value. That makes a malformed update overwrite an existing record, and it also creates a record that should not exist.

## Expected behavior

The request must contain a string `value` field. When the field is missing, the service returns HTTP `400` with JSON error `invalid value`, leaves an existing value unchanged, and does not create a new record. Valid string values, including an explicitly empty string, continue to use the normal write response.

## Reproduction

From the repository root, run:

```text
go test ./internal/httpapi -run '^TestPutMissingValue$' -count=20
```

The regression test seeds a record, sends `{}`, checks the `400` response and error body, then verifies the old value is still present. It also checks that a missing-value write does not create a new record.
