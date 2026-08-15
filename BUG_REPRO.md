# Cancelled Flag Lookups

## Behavior

Calling the flag service with an already cancelled context can still return a successful flag value.

## Reproduction

```sh
go test ./internal/flags -run '^TestServiceEnabledCancelledContextReturnsError$' -count=20
```

## Expected result

The call returns `context.Canceled`, returns `false`, and does not read from the store.

## Observed result on the base revision

The call returns the configured flag value with no error, even though its context was cancelled before the lookup.

## Verification

```sh
go test ./internal/flags -run '^TestServiceEnabledCancelledContext(ReturnsError|DoesNotPolluteCache)$' -count=20
go test ./...
```
