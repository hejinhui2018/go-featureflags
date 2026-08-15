# Whitespace in Successful Flag Responses

## Behavior

The service accepts leading or trailing spaces in the tenant header and flag path. The lookup uses the trimmed identifiers, but a successful JSON response echoes the untrimmed values.

## Reproduction

```sh
curl -H 'X-Tenant-ID:  tenant-a  ' 'http://localhost:8080/v1/flags/%20checkout%20'
```

## Expected result

The response identifies the queried objects with `"tenant":"tenant-a"` and `"flag":"checkout"`.

## Observed result on the base revision

The response contains the original whitespace, such as `"tenant":"  tenant-a  "` and `"flag":" checkout "`, even though the lookup uses the trimmed identifiers.

## Verification

```sh
go test ./internal/httpapi -run '^TestHandlerTrimsWhitespaceInResponse$' -count=20
go test ./...
```
